package main

import (
	"encoding/csv"
	"errors"
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

func constructPreference(row []string) (MusicPreference, error) {
	mp := MusicPreference{}
	v := reflect.ValueOf(mp)
	hours, err := strconv.ParseUint(row[54], 10, 32)
	mp.Hours = uint(hours)
	if err != nil {
		return mp, err
	}
	mp.Importance = row[55]
	mp.Lyrics = row[56]
	genreStartRow := 57
	nGenres := 27
	for i := 0; i < nGenres; i++ {
		if v.Field(3 + i).CanSet() {
			v.Field(3 + i).SetBool(len(row[genreStartRow+i]) > 0)
		}
	}
	mp.Other = row[genreStartRow+nGenres]
	return mp, nil
}

func constructInstruments(row []string) ([]Instrument, error) {
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
		experience, err := strconv.ParseUint(row[j+1], 10, 32)
		if err != nil {
			return nil, err
		}
		instrument.Experience = uint(experience)
		instrument.Still_Playing = row[j+2]
		instrument.Choice = row[j+3]
		instruments = append(instruments, instrument)
	}
	return instruments, nil
}

func upload_rows(rows [][]string, id_email_map Index, email_id_map Index) [][]string {
	failed := make([][]string, 0)
	instruments := make([]Instrument, 100)
	preferences := make([]MusicPreference, 100)
	for _, row := range rows {
		id := row[9]
		if _, exists := id_email_map[id]; !exists {
			email := row[4]
			if id, exists = email_id_map[email]; !exists {
				failed = append(failed, row)
				continue
			}
			mp, err := constructPreference(row)
			if err != nil {
				failed = append(failed, row)
				continue
			}
			new_instruments, err := constructInstruments(row)
			if err != nil {
				failed = append(failed, row)
				continue
			}
			instruments = append(instruments, new_instruments...)
			preferences = append(preferences, mp)
		}
	}
	return nil
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
	upload_rows(lines, id_email_map, email_id_map)
}
