package storage

import (
	"log"
)
func (db *DB) AddUser(userName string, password string) error {
	add_volume := `
	INSERT INTO User(
		UserName,
		Password
	) values(?, ?)
	`
	stmt, err := db.Prepare(add_volume)
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()
	_, err2 := stmt.Exec(userName, password)
	if err2 != nil {
		log.Print(err2)
		return err2
	}
	return nil
}

