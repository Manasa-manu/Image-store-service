package storage

import (
	"log"
	"time"
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

func (db *DB) CheckIfUserExists(userName string) (bool,error) {
	var count int

	row := db.QueryRow("SELECT COUNT(*) FROM User WHERE UserName = ?", userName)
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

func (db *DB) GetUserID(userName string) (int, error) {
	stmt, err := db.Prepare("SELECT UserID FROM User WHERE UserName = ?")
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(userName);
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
	stmt, err := db.Prepare("SELECT UserName FROM User WHERE UserName = ? ")
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

func (db *DB) GetPasswordHash(userName string) (string, error) {
	stmt, err := db.Prepare("SELECT Password FROM User WHERE UserID = ? ")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer stmt.Close()
	rows, err := stmt.Query(userName);
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

func (db *DB) GetAlbumName(albumId int) (string, error) {
	stmt, err := db.Prepare("SELECT AlbumName FROM Album WHERE AlbumID = ? ")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer stmt.Close()
	rows, err := stmt.Query(albumId);
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
	stmt, err := db.Prepare("SELECT AlbumID,AlbumName,CreatedAt,UserID FROM Album WHERE UserID = ?")
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

func (db *DB) GetAlbumID(userId int,albumName string) (int, error) {
	stmt, err := db.Prepare("SELECT AlbumID FROM Album WHERE UserID = ? AND AlbumName = ? ")
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId,albumName);
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

func (db *DB) AddAlbum(userId int, albumName string)error{
	add_Album := `
	INSERT INTO Album(
		AlbumName,
		CreatedAt,
		UserID
	) values(?, ?, ?)
	`
	stmt, err := db.Prepare(add_Album)
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()
	_, err2 := stmt.Exec(albumName, time.Now().Unix(),userId)
	if err2 != nil {
		log.Print(err2)
		return err2
	}
	return nil
}

func (db *DB) DeleteAlbum(albumId int)(int,error){
	stmt := "DELETE FROM Album where AlbumID = ?"
	res, err := db.Exec(stmt, albumId)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		// panic(err)
		log.Print(err)
		return 0, err
	}
	return int(count),nil
}

func (db *DB) GetImageName(imageId int) (string, error) {
	stmt, err := db.Prepare("SELECT ImageName FROM Image WHERE ImageID = ? ")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer stmt.Close()
	rows, err := stmt.Query(imageId);
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

func (db *DB) GetImages(albumId int) ([]Image, error) {
	stmt, err := db.Prepare("SELECT ImageID,ImageName,AlbumID FROM Image WHERE AlbumID = ?")
	if err != nil {
		log.Fatal(err)
		return []Image{}, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(albumId);
	if err != nil {
		log.Fatal(err)
		return []Image{}, err
	}
	defer rows.Close()
	var result_arr []Image
	var result Image;
	for rows.Next() {
		err := rows.Scan(&result.ImageID,&result.ImageName,&result.AlbumID)
		if err != nil {
			log.Fatal(err)
			return []Image{}, err
		}
		result_arr =append(result_arr,result)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return []Image{}, err
	}
	return result_arr, nil
}

func (db *DB) GetImageID(albumId int,imageName string) (int, error) {
	stmt, err := db.Prepare("SELECT ImageID FROM Image WHERE AlbumID = ? AND ImageName = ? ")
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(albumId,imageName);
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

func (db *DB) AddImage(albumId int, imageName string)error{
	add_Image := `
	INSERT INTO Image(
		AlbumID,
		ImageName
	) values(?, ?)
	`
	stmt, err := db.Prepare(add_Image)
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmt.Close()
	_, err2 := stmt.Exec(albumId,imageName)
	if err2 != nil {
		log.Print(err2)
		return err2
	}
	return nil
}

func (db *DB) DeleteImage(imageId int) (int,error){
	stmt := "DELETE FROM Image where ImageID = ?"
	res, err := db.Exec(stmt, imageId)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		// panic(err)
		log.Print(err)
		return 0, err
	}
	return int(count),nil
}

