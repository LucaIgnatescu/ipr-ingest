package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const DATABASE = "machine2.db"

type Instrument struct {
	Instrument    string
	Experience    uint
	Still_Playing string
	Choice        string
}

type MusicPreference struct {
	Hours            uint
	Importance       string
	Lyrics           string
	AlternativeIndie bool
	Blues            bool
	Classical        bool
	Country          bool
	Dance            bool
	Disco            bool
	Electronic       bool
	EDM              bool
	Folk             bool
	Funk             bool
	Gospel           bool
	House            bool
	Jazz             bool
	Latin            bool
	Metal            bool
	Opera            bool
	Pop              bool
	Punk             bool
	RapHiphop        bool
	RnBSoul          bool
	Reggae           bool
	Religious        bool
	Rock             bool
	Soundtrack       bool
	Swing            bool
	Techno           bool
	World            bool
	Other            string
}

const QUERY_CREATE_LISTENING = `
CREATE TABLE music_preference (
    participantId TEXT PRIMARY KEY,
    Hours INTEGER,
    Importance INTEGER,
    Lyrics INTEGER,
    AlternativeIndie INTEGER,
    Blues INTEGER,
    Classical INTEGER,
    Country INTEGER,
    Dance INTEGER,
    Disco INTEGER,
    Electronic INTEGER,
    EDM INTEGER,
    Folk INTEGER,
    Funk INTEGER,
    Gospel INTEGER,
    House INTEGER,
    Jazz INTEGER,
    Latin INTEGER,
    Metal INTEGER,
    Opera INTEGER,
    Pop INTEGER,
    Punk INTEGER,
    RapHiphop INTEGER,
    RnBSoul INTEGER,
    Reggae INTEGER,
    Religious INTEGER,
    Rock INTEGER,
    Soundtrack INTEGER,
    Swing INTEGER,
    Techno INTEGER,
    World INTEGER,
    Other INTEGER
);
  `
const QUERY_CREATE_INSTRUMENT = `
CREATE TABLE instrument (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    participantId TEXT,
    Instrument TEXT,
    Experience INTEGER,
    Still_Playing TEXT,
    Choice TEXT
);
  `

const QUERY_RESET_TABLES = `
  DROP TABLE IF EXISTS instrument;
  DROP TABLE IF EXISTS music_preference;
  `

func connect_db() (*sql.DB, error) {
	fmt.Println("Connecting to database...")
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected")
	return db, nil
}

func create_tables(db *sql.DB) error {
	if db == nil {
		return errors.New("Database connection is nil")
	}
	_, err := db.Exec(QUERY_RESET_TABLES)
	if err != nil {
		return err
	}
	_, err = db.Exec(QUERY_CREATE_LISTENING)
	if err != nil {
		fmt.Println("Error creating listening table")
		return err
	}
	_, err = db.Exec(QUERY_CREATE_INSTRUMENT)
	if err != nil {
		fmt.Println("Error creating instrument table")
		return err
	}
	fmt.Println("Created tables")
	return nil
}
