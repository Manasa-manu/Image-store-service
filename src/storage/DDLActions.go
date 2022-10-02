package storage

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

)

func CreateAllTables(db *sql.DB) {
	// create table if not exists
	User_table := `
	CREATE TABLE IF NOT EXISTS User(
		UserID INTEGER NOT NULL UNIQUE PRIMARY KEY,
		UserName TEXT NOT NULL UNIQUE,
		Password TEXT NOT NULL
	);
	`
	_, err := db.Exec(User_table)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	Album_table := `
	CREATE TABLE IF NOT EXISTS Album(
		CreatedAt BIGINT NONT NULL,
		AlbumID BIGINT NOT NULL UNIQUE PRIMARY KEY,
		AlbumName TEXT NOT NULL,
		FOREIGN KEY (UserID) 
			REFERENCES USER(VolumeID) 
			ON DELETE CASCADE
	);
	`
	_, err2 := db.Exec(Album_table)
	if err2 != nil {
		log.Fatal(err2)
		panic(err2)
	}

	Image_table := `
	CREATE TABLE IF NOT EXISTS Image(
		ImageID BIGINT NOT NULL UNIQUE PRIMARY KEY,
		ImageName TEXT NOT NULL,
		FOREIGN KEY (UserID) 
			REFERENCES USER(VolumeID) 
			ON DELETE CASCADE
	);
	`
	_, err1 := db.Exec(Image_table)
	if err1 != nil {
		log.Fatal(err1)
		panic(err1)
	}
}
