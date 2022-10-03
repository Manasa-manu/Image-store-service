package storage

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

)

func CreateAllTables(db *sql.DB) {
	// create table if not exists
	User_table := `
	CREATE TABLE if not exists User(
		UserID BIGINT NOT NULL AUTO_INCREMENT,
		UserName VARCHAR(255) NOT NULL,
		Password VARCHAR(255),
		PRIMARY KEY (UserID)
	);
	`
	_, err := db.Exec(User_table)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	Album_table := `
	CREATE TABLE if not exists Album(
		CreatedAt BIGINT NOT NULL,
		AlbumID BIGINT NOT NULL AUTO_INCREMENT,
		AlbumName VARCHAR(255) NOT NULL,
		UserID BIGINT,
		  PRIMARY key (AlbumID),
		FOREIGN KEY (UserID) 
			REFERENCES User(UserID) 
			ON DELETE CASCADE
	);	
	`
	_, err2 := db.Exec(Album_table)
	if err2 != nil {
		log.Fatal(err2)
		panic(err2)
	}

	Image_table := `
	CREATE TABLE if not exists Image(
		ImageID BIGINT NOT NULL AUTO_INCREMENT,
		ImageName TEXT NOT NULL,
		  AlbumID BIGINT,
		  PRIMARY KEY (ImageID),
		FOREIGN KEY (AlbumID) 
			REFERENCES Album(AlbumID) 
			ON DELETE CASCADE
	);	
	`
	_, err1 := db.Exec(Image_table)
	if err1 != nil {
		log.Fatal(err1)
		panic(err1)
	}
}



