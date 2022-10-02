package main

import (
	"image-store-service/storage"

	"log"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	DataDir          = "photos"
	listen           = ":8080"
	db,_ = storage.NewDB()
	Env 		     = &BeanAccess{db}
)

type BeanAccess struct {
	db storage.Datastore
}
func initDB() {
	db := storage.InitDB()
	defer db.Close()
	storage.CreateAllTables(db)
}

func main() {
	/****Initialization***/
	initDB()
	/****End*/
	/*****Init DB******/
	// db, err := storage.NewDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	/*****End**********/
	// env := &BeanAccess{db}
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	r.POST("/user/register",UserRegisterHandler)

	r.GET("/register/success",func(c *gin.Context){
		c.String(http.StatusOK,"success!")
	})

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "resource not found")
	})

	r.GET("/login",func(c *gin.Context){
		c.File("login.html")
	})
	r.POST("/user/login",LoginHandler)

	r.GET("/user/:id",func(c *gin.Context){
		c.String(http.StatusOK,"Hello, user!")
	})
	
	log.Fatal(r.Run(listen))
}
