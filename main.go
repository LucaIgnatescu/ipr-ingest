package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"sync"
)

const (
	CODES_FILE      = "codes.csv"
	MUSIC_DATA_FILE = "music_data.csv"
)

type Index map[string]string

func parseFile(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}
	return records, nil
}

func setupMaps(lines [][]string, id_email_map Index, email_id_map Index) [][]string {
	ignored := make([][]string, 0)
	for _, line := range lines {
		id := line[0]
		email := line[1]
		_, id_exists := id_email_map[id]
		_, email_exists := email_id_map[email]
		if id_exists || email_exists {
			ignored = append(ignored, line)
			continue
		}
		id_email_map[id] = email
		email_id_map[email] = id
	}
	return ignored
}

func constructPreference(id string, row []string) MusicPreference {
	mp := MusicPreference{}
	v := reflect.ValueOf(&mp).Elem()
	hours, err := strconv.ParseUint(row[54], 10, 32)
	if err != nil {
		fmt.Println(row[54])
		mp.Hours = 0 // NOTE: The only error I noticed was a missing value. I will just assign to 0
	} else {
		mp.Hours = uint(hours)
	}
	mp.Importance = row[55]
	mp.Lyrics = row[56]
	genreStartRow := 57
	nGenres := 27
	for i := 0; i < nGenres; i++ {
		v.Field(4 + i).SetBool(len(row[genreStartRow+i]) > 0)
	}
	mp.Other = row[genreStartRow+nGenres]
	mp.ParticipandId = id
	return mp
}

func constructInstruments(id string, row []string) []Instrument {
	instruments := make([]Instrument, 0)
	nInstruments := 5
	startingRow := 33
	width := 4
	for i := 0; i < nInstruments; i++ {
		j := startingRow + width*i
		if len(row[j]) == 0 {
			break
		}
		instrument := Instrument{}
		instrument.Instrument = row[j]
		instrument.Experience = row[j+1]
		instrument.StillPlaying = row[j+2]
		instrument.Choice = row[j+3]
		instrument.ParticipandId = id
		instruments = append(instruments, instrument)
	}
	return instruments
}

func uploadRows(db *sql.DB, rows [][]string, id_email_map Index, email_id_map Index) [][]string {
	failed := make([][]string, 0)
	instruments := make([]Instrument, 0)
	preferences := make([]MusicPreference, 0)
	for _, row := range rows {
		id := row[9]
		if _, exists := id_email_map[id]; !exists {
			email := row[4]
			if id, exists = email_id_map[email]; !exists {
				failed = append(failed, row)
				continue
			}
		}
		mp := constructPreference(id, row)
		new_instruments := constructInstruments(id, row)
		instruments = append(instruments, new_instruments...)
		preferences = append(preferences, mp)
	}
	fmt.Printf("Parsed %v instruments\n", len(instruments))
	fmt.Printf("Parsed %v preferences\n", len(preferences))
	fmt.Printf("Failed to parse %v/%v rows\n", len(failed), len(rows))

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := insertMusicPreferences(db, preferences); err != nil {
			panic(err)
		}
		fmt.Println("Inserted music preferences")
	}()
	go func() {
		defer wg.Done()
		if err := insertInstruments(db, instruments); err != nil {
			panic(err)
		}
		fmt.Println("Inserted instruments")
	}()
	wg.Wait()
	return failed
}

func writeRows(rows [][]string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	return writer.WriteAll(rows)
}

func main() {
	id_email_map := make(Index)
	email_id_map := make(Index)

	lines, err := parseFile(CODES_FILE)
	if err != nil {
		log.Fatal(err)
	}

	unparsed := setupMaps(lines[1:], id_email_map, email_id_map)
	writeRows(unparsed, "code_unparsed.csv")

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}

	err = create_tables(db)
	if err != nil {
		log.Fatal(err)
	}
	lines, err = parseFile(MUSIC_DATA_FILE)
	if err != nil {
		log.Fatal(err)
	}
	unparsed = uploadRows(db, lines[2:], id_email_map, email_id_map)
	writeRows(unparsed, "music_data_unparsed.csv")
}
