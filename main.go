package main

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
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

func main() {
	id_email_map := make(Index)
	email_id_map := make(Index)

	lines, err := parse_file("codes.csv")
	if err != nil {
		log.Fatal(err)
	}
	if setup_maps(lines, id_email_map, email_id_map) != nil {
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

}
