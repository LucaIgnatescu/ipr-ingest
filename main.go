package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

const (
	CODES_FILE      = "codes.csv"
	MUSIC_DATA_FILE = "music_data.csv"
)

type Index map[string]string

func parse_file(fileName string) ([][]string, error) {
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

func setup_maps(lines [][]string, id_email_map Index, email_id_map Index) error {
	for _, line := range lines {
		if len(line) < 2 {
			return errors.New("Each line should have at least 2 entries")
		}

		id := line[0]
		email := line[1]
		// if _, exists := id_email_map[id]; exists {
		// 	fmt.Printf("Found duplicate id %v at row %v with email %v\n", id, i, email)
		// }
		id_email_map[id] = email
		email_id_map[email] = id
	}
	return nil
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

func upload_rows(rows [][]string, id_email_map Index, email_id_map Index) [][]string {
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
	return failed
}

func main() {
	id_email_map := make(Index)
	email_id_map := make(Index)

	lines, err := parse_file(CODES_FILE)
	if err != nil {
		log.Fatal(err)
	}

	err = setup_maps(lines, id_email_map, email_id_map)
	if err != nil {
		log.Fatal(err)
	}

	db, err := connect_db()
	if err != nil {
		log.Fatal(err)
	}

	err = create_tables(db)
	if err != nil {
		log.Fatal(err)
	}

	lines, err = parse_file(MUSIC_DATA_FILE)
	if err != nil {
		log.Fatal(err)
	}
	_ = upload_rows(lines, id_email_map, email_id_map)
}
