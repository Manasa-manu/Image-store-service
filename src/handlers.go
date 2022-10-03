package main

import (
	"strconv"
	"net/http"
	"log"
	"github.com/gin-gonic/gin"
)


// user/login
func LoginHandler(c *gin.Context){
	name:= c.PostForm("name");
	password:= c.PostForm("password");
	log.Print("Name from:"+name)
	log.Print("Passwor:"+password)
	exists,_ := Env.db.CheckIfUserExists(name,password);
	if exists{
		userid,_ := Env.db.GetUserID(name,password);
		c.Redirect(http.StatusFound,"/user/"+strconv.Itoa(userid))
	}else{
		c.Redirect(http.StatusFound,"/login")
	}	
}
// /user/register
func UserRegisterHandler(c *gin.Context){
	name:= c.PostForm("name")
	password:= c.PostForm("password")

	log.Print("Register:"+name+" pass"+password);

	Env.db.AddUser(name,password)
	c.Redirect(http.StatusFound,"/register/success")
}

//user/:userid -- lists all Albums
func UserAlbumsHandler(c *gin.Context){
	userId,_:= strconv.Atoi(c.Param("id"))
	albumList,_ := Env.db.GetAlbums(userId)
	userName,_:= Env.db.GetUserName(userId)
	c.HTML(http.StatusOK, "UserAlbumList.html", gin.H{"UserName":userName,"UserID":userId,"List":albumList})
}

//user/:userid/album/add
func AddAlbumHandler(c *gin.Context){
	userId:=c.Param("id")
	c.HTML(http.StatusOK,"CreateAlbum.html",gin.H{"UserID":userId})
}
