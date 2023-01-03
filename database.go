package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Organism struct {
	ID    int64
	Name  string
	Image []byte
	Lat   float64
	Lng   float64
}

func connectToDatabase(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS organisms (
	id INTEGER PRIMARY KEY,
	name TEXT,
	image BLOB,
	lat REAL,
	lng REAL
	)`)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func insertIntoTable(db *sql.DB, name string, image []byte, lat float64, lng float64) (int64, error) {
	result, err := db.Exec(`INSERT INTO organisms (name, image, lat, lng) VALUES (?, ?, ?, ?)`, name, image, lat, lng)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func selectAllFromTable(db *sql.DB) ([]*Organism, error) {
	// Select all rows from the table
	rows, err := db.Query(`SELECT id, name, image, lat, lng FROM organisms`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organisms []*Organism
	for rows.Next() {
		var id int64
		var name string
		var image []byte
		var lat float64
		var lng float64
		if err := rows.Scan(&id, &name, &image, &lat, &lng); err != nil {
			return nil, err
		}
		organisms = append(organisms, &Organism{
			ID:    id,
			Name:  name,
			Image: image,
			Lat:   lat,
			Lng:   lng,
		})
	}

	return organisms, rows.Err()
}

func deleteFromTable(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM organisms WHERE id = ?`, id)
	return err
}
