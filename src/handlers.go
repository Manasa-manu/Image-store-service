package main

import (
		
	"net/http"
	"log"
	"github.com/gin-gonic/gin"
)
var m map[string]string = make(map[string]string)
func LoginHandler(c *gin.Context){
	name:= c.PostForm("name");
	// password:= c.PostForm("password");
	log.Print("Name from:"+name)
	_, ok :=m[name]
	if ok{
		c.Redirect(http.StatusFound,"/user/1234")
	}else{
		c.Redirect(http.StatusFound,"/usd")
	}	
}

func UserRegisterHandler(c *gin.Context){
	name:= c.PostForm("name")
	password:= c.PostForm("password")

	Env.db.AddUser(name,password)
	m[name]=password
	c.Redirect(http.StatusFound,"/register/success")
}