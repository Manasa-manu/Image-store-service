package storage

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserID   int // primary key
	UserName string
	Password string
}

type Image struct {
	ImageID      int
	AlbumID      int
	ImageName    string
}

type Album struct {
	AlbumID     int
	AlbumName   string
	CreatedAt   string
	UserID      int
}

type AlbumDetails struct{
	UserName string
	UserID int
	AlbumList[] Album
}
type Datastore interface {

	// User table actions
	AddUser(userName string, password string) error
	CheckIfUserExists(userName string,password string) (bool,error)
	GetUserID(userName string,password string) (int,error)
	GetUserName(userId int)(string,error)

	// Album table actions
	GetAlbums(userId int)([]Album,error)
}

type DB struct {
	*sql.DB
}


// connStr := "mysql://mysql:3306/image_store"
// connStr := "mysql:secret@/image_store"
// Connect to database
// const connStr = "mysql:secret@tcp(mysql:3306)/image_store"
const connStr = "test:password@tcp(localhost:3306)/image_store"

func NewDB() (*DB, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func InitDB() *sql.DB {
	log.Print("Start of InitDB")
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	if db == nil {
		log.Fatal(err)
		panic("db nil")
	}
	db.SetMaxIdleConns(5)
	log.Print("End of InitDB")
	return db
}

