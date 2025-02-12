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
	Experience    string
	StillPlaying  string
	Choice        string
	ParticipandId string
}

type MusicPreference struct {
	ParticipandId    string
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
    participantId TEXT,
    hours INTEGER,
    importance TEXT,
    lyrics TEXT,
    alternative_indie INTEGER,
    blues INTEGER,
    classical INTEGER,
    country INTEGER,
    dance INTEGER,
    disco INTEGER,
    electronic INTEGER,
    edm INTEGER,
    folk INTEGER,
    funk INTEGER,
    gospel INTEGER,
    house INTEGER,
    jazz INTEGER,
    latin INTEGER,
    metal INTEGER,
    opera INTEGER,
    pop INTEGER,
    punk INTEGER,
    rapHiphop INTEGER,
    rnBSoul INTEGER,
    reggae INTEGER,
    religious INTEGER,
    rock INTEGER,
    soundtrack INTEGER,
    swing INTEGER,
    techno INTEGER,
    world INTEGER,
    other TEXT
);
  `
const QUERY_CREATE_INSTRUMENT = `
CREATE TABLE instrument (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    participantId TEXT,
    instrument TEXT,
    experience TEXT,
    still_playing TEXT,
    choice TEXT
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

func insertMusicPreferences(db *sql.DB, prefs []MusicPreference) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO music_preference (
			participantId,
			hours,
			importance,
			lyrics,
			alternative_indie,
			blues,
			classical,
			country,
			dance,
			disco,
			electronic,
			edm,
			folk,
			funk,
			gospel,
			house,
			jazz,
			latin,
			metal,
			opera,
			pop,
			punk,
			rapHiphop,
			rnBSoul,
			reggae,
			religious,
			rock,
			soundtrack,
			swing,
			techno,
			world,
			other
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, pref := range prefs {
		_, err := stmt.Exec(
			pref.ParticipandId,    // participantId TEXT
			pref.Hours,            // hours INTEGER
			pref.Importance,       // importance TEXT
			pref.Lyrics,           // lyrics TEXT
			pref.AlternativeIndie, // alternative_indie INTEGER (bool converted automatically)
			pref.Blues,            // blues INTEGER
			pref.Classical,        // classical INTEGER
			pref.Country,          // country INTEGER
			pref.Dance,            // dance INTEGER
			pref.Disco,            // disco INTEGER
			pref.Electronic,       // electronic INTEGER
			pref.EDM,              // eDM INTEGER
			pref.Folk,             // folk INTEGER
			pref.Funk,             // funk INTEGER
			pref.Gospel,           // gospel INTEGER
			pref.House,            // house INTEGER
			pref.Jazz,             // jazz INTEGER
			pref.Latin,            // latin INTEGER
			pref.Metal,            // metal INTEGER
			pref.Opera,            // opera INTEGER
			pref.Pop,              // pop INTEGER
			pref.Punk,             // punk INTEGER
			pref.RapHiphop,        // rapHiphop INTEGER
			pref.RnBSoul,          // rnBSoul INTEGER
			pref.Reggae,           // reggae INTEGER
			pref.Religious,        // religious INTEGER
			pref.Rock,             // rock INTEGER
			pref.Soundtrack,       // soundtrack INTEGER
			pref.Swing,            // swing INTEGER
			pref.Techno,           // techno INTEGER
			pref.World,            // world INTEGER
			pref.Other,            // other TEXT
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func insertInstruments(db *sql.DB, instruments []Instrument) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO instrument (
			participantId,
			instrument,
			experience,
			still_playing,
			choice
		) VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, inst := range instruments {
		_, err = stmt.Exec(
			inst.ParticipandId, // participantId TEXT
			inst.Instrument,    // instrument TEXT
			inst.Experience,    // experience TEXT
			inst.StillPlaying,  // still_playing TEXT
			inst.Choice,        // choice TEXT
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
