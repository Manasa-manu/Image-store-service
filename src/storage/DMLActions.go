package storage

import (
	"log"
)
func (db *DB) AddUser(userName string, password string) error {
	add_user := `
	INSERT INTO User(
		UserName,
		Password
	) values(?, ?)
	`
	stmt, err := db.Prepare(add_user)
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

func (db *DB) CheckIfUserExists(userName string, password string) (bool,error) {
	var count int

	row := db.QueryRow("SELECT COUNT(*) FROM User WHERE UserName = ? AND Password = ?", userName,password)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	log.Print("count of similar mountpoints in db :", count)
	if count == 0 {
		return false, err
	} else {
		return true, nil
	}
}

func (db *DB) GetUserID(userName string,password string) (int, error) {
	stmt, err := db.Prepare("SELECT UserID FROM User WHERE UserName = ? AND Password = ? ")
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(userName,password);
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	defer rows.Close()
	var result int
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			log.Fatal(err)
			return -1, err
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return -1, err
	}
	return result, nil
}

func (db *DB) GetUserName(userId int) (string, error) {
	stmt, err := db.Prepare("SELECT UserName FROM User WHERE UserID = ? ")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId);
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer rows.Close()
	var result string
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return "", err
	}
	return result, nil
}

func (db *DB) GetAlbums(userId int) ([]Album, error) {
	stmt, err := db.Prepare("SELECT * FROM Album WHERE UserID = ?")
	if err != nil {
		log.Fatal(err)
		return []Album{}, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId);
	if err != nil {
		log.Fatal(err)
		return []Album{}, err
	}
	defer rows.Close()
	var result_arr []Album
	var result Album;
	for rows.Next() {
		err := rows.Scan(&result.AlbumID,&result.AlbumName,&result.CreatedAt,&result.UserID)
		if err != nil {
			log.Fatal(err)
			return []Album{}, err
		}
		result_arr =append(result_arr,result)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return []Album{}, err
	}
	return result_arr, nil
}
