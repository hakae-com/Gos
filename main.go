package main

import (
	"fmt"
	"gos/gos"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gos.New()
	//add global middleware
	r.User(gos.Logger())

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	r.GET("/", func(c *gos.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	stu1 := &student{
		Name: "小米",
		Age:  12,
	}
	stu2 := &student{
		Name: "华为",
		Age:  32,
	}
	r.GET("/students", func(c *gos.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gos.H{
			"title":  "gos",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/date", func(c *gos.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gos.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	account := r.Group("/account")
	{
		account.POST("/register", func(c *gos.Context) {
			c.JSON(http.StatusOK, gos.H{"uid": c.PostForm("uid"), "name": c.PostForm("name")})
		})
		account.POST("/login", func(c *gos.Context) {
			c.JSON(http.StatusOK, gos.H{"phone": c.PostForm("phone"), "code": c.PostForm("code")})
		})
	}
	user := r.Group("/user")
	{
		user.GET("/info/:uid", func(c *gos.Context) {
			c.JSON(http.StatusOK, gos.H{"uid": c.Param("uid")})
		})
		user.GET("/assets/*filepath", func(c *gos.Context) {
			c.JSON(http.StatusOK, gos.H{"filepath": c.Param("filepath")})
		})
	}

	err := r.RUN(":9999")
	if err != nil {
		return
	}
}
