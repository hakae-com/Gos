package main

import (
	"gos/gos"
	"net/http"
)

func main() {
	r := gos.New()
	r.GET("/", func(c *gos.Context) {
		c.HTML(http.StatusOK, "<h1>Hello World</h1>")
	})
	r.POST("/register", func(c *gos.Context) {
		c.JSON(http.StatusOK, gos.H{"uid": c.PostForm("uid"), "name": c.PostForm("name")})
	})
	r.POST("/login", func(c *gos.Context) {
		c.JSON(http.StatusOK, gos.H{"phone": c.PostForm("phone"), "code": c.PostForm("code")})
	})
	r.GET("/userInfo/:uid", func(c *gos.Context) {
		c.JSON(http.StatusOK, gos.H{"uid": c.Param("uid")})
	})
	r.GET("/assets/*filepath", func(c *gos.Context) {
		c.JSON(http.StatusOK, gos.H{"filepath": c.Param("filepath")})
	})

	err := r.RUN(":9999")
	if err != nil {
		return
	}
}
