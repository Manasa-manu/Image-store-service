package main

import (
	"image-store-service/storage"

	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"context"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	DataDir          = "photos"
	listen           = ":8080"

	// DB
	db,_ = storage.NewDB()
	Env 		     = &BeanAccess{db}

	// minio 
	MinioClientInstance *minio.Client = nil
	ctx = context.Background()
)
const (
	BucketName = "images"
	Location = "us-east-1"

)

type BeanAccess struct {
	db storage.Datastore
}
func initDB() {
	db := storage.InitDB()
	defer db.Close()
	storage.CreateAllTables(db)
}
func initMinio(){
        endpoint := "play.min.io"
        accessKeyID := "Q3AM3UQ867SPQQA43P2F"
        secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
        useSSL := true

        // Initialize minio client object.
		var err error
        MinioClientInstance, err = minio.New(endpoint, &minio.Options{
                Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
                Secure: useSSL,
        })
        if err != nil {
                log.Fatalln(err)
        }
}

func main() {
	/****Initialization***/
	initDB()
	/****End*/
	initMinio()
	
	r := gin.Default()
	r.LoadHTMLGlob("./static/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK,"index.html",nil)
	})

	r.POST("/user/register",UserRegisterHandler)

	r.GET("/register/success",func(c *gin.Context){
		c.String(http.StatusOK,"success!")
	})

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "resource not found")
	})

	r.GET("/login",func(c *gin.Context){
		c.HTML(http.StatusOK,"login.html",nil)
	})
	r.POST("/user/login",LoginHandler)

	r.GET("/user/:id",UserAlbumsHandler)
	
	r.GET("/user/:id/album/add",GetAlbumHandler)
	r.POST("/user/:id/album/add",AddAlbumHandler)
	
	r.GET("/user/:id/album/:albumid",GetImageHandler)
	r.GET("/user/:id/album/:albumid/download",DownloadAlbumHandler)
	r.GET("image/:imageid/download",DownloadImageHandler)
	r.POST("image/:imageid/delete",DeleteImageHandler) // DELETE Method 
	r.POST("user/:id/album/:albumid/delete",DeleteAlbumHandler) // DELETE Method 
	r.GET("user/:id/album/:albumid/image/add",AddImageFileFormHandler)
	r.POST("user/:id/album/:albumid/image/add",AddImageHandler)
	log.Fatal(r.Run(listen))
}
