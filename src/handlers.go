package main

import (
	"strconv"
	"net/http"
	"log"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"os"
	"archive/zip"
	"io"
	"golang.org/x/crypto/bcrypt"
)

const tempLocation ="/tmp/image/"
// user/login
func LoginHandler(c *gin.Context){
	name:= c.PostForm("name");
	password:= c.PostForm("password");
	exists,_ := Env.db.CheckIfUserExists(name);
	if exists{
		hashFromDatabase ,_:= Env.db.GetPasswordHash(name);
		hashFromDBInByte := []byte(hashFromDatabase);
		// compare user entered password and password hash from DB
		if err := bcrypt.CompareHashAndPassword(hashFromDBInByte, []byte(password)); err != nil {
			// wrong password so redirect to login page
			c.Redirect(http.StatusFound,"/login")
		}

		userid,_ := Env.db.GetUserID(name);
		c.Redirect(http.StatusFound,"/user/"+strconv.Itoa(userid))
	}else{
		c.Redirect(http.StatusFound,"/login")
	}	
}
// /user/register
func UserRegisterHandler(c *gin.Context){
	name:= c.PostForm("name")
	password:= c.PostForm("password")
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal(err)
    }

	Env.db.AddUser(name,string(hash))
	c.Redirect(http.StatusFound,"/register/success")
}

//user/:userid -- lists all Albums
func UserAlbumsHandler(c *gin.Context){
	userId,_:= strconv.Atoi(c.Param("id"))
	albumList,_ := Env.db.GetAlbums(userId)
	userName,_:= Env.db.GetUserName(userId)
	c.HTML(http.StatusOK, "UserAlbumList.html", gin.H{"UserName":userName,"UserID":userId,"List":albumList})
}

func GetAlbumHandler(c *gin.Context){
	userId:=c.Param("id")
	c.HTML(http.StatusOK,"CreateAlbum.html",gin.H{"UserID":userId})
}

func GetImageHandler(c *gin.Context){
	albumId,_:=strconv.Atoi(c.Param("albumid"))
	userId:=c.Param("id")
	imageList,_ := Env.db.GetImages(albumId);
	c.HTML(http.StatusOK,"AlbumImageList.html",gin.H{"UserID":userId,"AlbumID":albumId,"List":imageList})
}

func AddAlbumHandler(c *gin.Context) {
    // Multipart form
    form, _ := c.MultipartForm()
    files := form.File["filename[]"]
	userId,_ := strconv.Atoi(c.Param("id"))
	albumName := c.PostForm("name")
	Env.db.AddAlbum(userId,albumName)
	albumId,_ := Env.db.GetAlbumID(userId,albumName)
    for _, file := range files {
	  filename := file.Filename
	  Env.db.AddImage(albumId,filename)
	  imageId_int,_:=Env.db.GetImageID(albumId,filename)
	  imageId :=strconv.Itoa(imageId_int)
	  key := tempLocation+imageId
      // Upload the file to specific dst.
	  c.SaveUploadedFile(file, key)
	  // Upload image
	  uploadImage(file,imageId)
    }
    c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
  }

func DownloadAlbumHandler(c *gin.Context){
	// userid:= c.Param("id")
	albumId,_:= strconv.Atoi(c.Param("albumid"))
	c.Writer.Header().Set("Content-type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=filename.zip")
	ar :=  zip.NewWriter(c.Writer)
	images,_ :=Env.db.GetImages(albumId)
	for _,image :=range images{
		// downloadImage(strconv.Itoa(image.ImageID),image.ImageName,ar)
		object, err := MinioClientInstance.GetObject(ctx, BucketName, strconv.Itoa(image.ImageID), minio.GetObjectOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
		f1, _ := ar.Create(image.ImageName)
		if _, err = io.Copy(f1, object); err != nil {
			fmt.Println(err)
			return
		}
	}
	defer ar.Close()
}


  func makeNewBucket(){
	// Make a new bucket called images.

	err := MinioClientInstance.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{Region: Location})
	if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			exists, errBucketExists := MinioClientInstance.BucketExists(ctx, BucketName)
			if errBucketExists == nil && exists {
					log.Printf("We already own %s\n", BucketName)
			} else {
					log.Fatalln(err)
			}
	} else {
			log.Printf("Successfully created %s\n", BucketName)
	}
  }
  func uploadImage(fileHeader *multipart.FileHeader,imageId string){
	makeNewBucket()
	file, err := os.Open(tempLocation+imageId)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	info, err := MinioClientInstance.PutObject(ctx, BucketName,imageId, file,fileHeader.Size, minio.PutObjectOptions{ContentType: fileHeader.Header["Content-Type"][0]})
	if err != nil {
			log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", fileHeader.Filename, info.Size)
	// remove temp data stored in server
	os.Remove(tempLocation+imageId)
  }
  func deleteImageFromMinio(imageId string){
	opts := minio.RemoveObjectOptions {
		}
	err := MinioClientInstance.RemoveObject(ctx, BucketName, imageId, opts)
	if err != nil {
		fmt.Println(err)
		return
	}
  }
  func DeleteImageHandler(c *gin.Context){
	// can implement this in a Task Engine, so that schedular will take care of it
	// Delete Images from bucket 
	// Delete Entry in DB
	imageId:=c.Param("imageid")
	deleteImageFromMinio(imageId)
	imageId_int,_ :=strconv.Atoi(imageId)
	Env.db.DeleteImage(imageId_int)
	c.String(http.StatusOK, "Deleted successfully")

  }
  func DownloadImageHandler(c *gin.Context){
	imageId:=c.Param("imageid")
	object, err := MinioClientInstance.GetObject(ctx, BucketName, imageId, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	imageId_int,_ :=strconv.Atoi(imageId)
	imageName,_ := Env.db.GetImageName(imageId_int)
	c.Header("Content-Disposition", "attachment; filename="+imageName )
    c.Header("Content-Type", "application/octet-stream")
	io.Copy(c.Writer, object) 
 }
 func DeleteAlbumHandler(c *gin.Context){
	// can implement this in a Task Engine, so that schedular will take care of it
	// Delete Images from bucket 
	// Delete Entry in DB
	albumId:=c.Param("albumid")
	albumId_int,_ :=strconv.Atoi(albumId)
	images,_ := Env.db.GetImages(albumId_int)
	for _,image :=range images {
		deleteImageFromMinio(strconv.Itoa(image.ImageID))
	}
	Env.db.DeleteAlbum(albumId_int)
	c.String(http.StatusOK, "Deleted successfully")
  }


func AddImageFileFormHandler(c *gin.Context){
	albumId,_:=strconv.Atoi(c.Param("albumid"))
	userId:=c.Param("id")
	c.HTML(http.StatusOK,"AddImageToAlbum.html",gin.H{"UserID":userId,"AlbumID":albumId})
}

func AddImageHandler(c *gin.Context) {
    // Multipart form
    form, _ := c.MultipartForm()
    files := form.File["filename[]"]
	albumId,_ := strconv.Atoi(c.Param("albumid"))
    for _, file := range files {
	  filename := file.Filename
	  Env.db.AddImage(albumId,filename)
	  imageId_int,_:=Env.db.GetImageID(albumId,filename)
	  imageId :=strconv.Itoa(imageId_int)
	  key := tempLocation+imageId
      // Upload the file to specific dst.
	  c.SaveUploadedFile(file, key)
	  // Upload image
	  uploadImage(file,imageId)
    }
    c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
  }
