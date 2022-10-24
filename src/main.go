package main

import (
	"image-store-service/storage"
	"os"
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
    endpoint, endpointOk := os.LookupEnv("MINIO_ENDPOINT");
	accessKeyId, accessKeyIdOk := os.LookupEnv("MINIO_ACCESS_KEY_ID");
	secretAccessKey, secretAccessKeyOk := os.LookupEnv("MINIO_SECRET_ACCESS_KEY");
    useSSL := false

	var err error
	if(endpointOk && accessKeyIdOk && secretAccessKeyOk) {

		// Initialize minio client object.
		MinioClientInstance, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatal("Missing minio connection info, check environment variables")
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
