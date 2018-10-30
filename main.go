package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	"github.com/jezman/gorion/models"
)

var (
	db      *models.DB
	router  *gin.Engine
	err     error
	env     models.Datastore
	timeNow = time.Now().Local()
)

func initDB() (db *models.DB) {
	dsn := os.Getenv("BOLID_DSN")

	if db, err = models.OpenDB(dsn); err != nil {
		panic(err)
	}

	env = db
	return
}

func indexHandler(c *gin.Context) {
	render(c, gin.H{"title": "Главная страница"}, "index.html")
}

func companyHandler(c *gin.Context) {
	companies, err := env.Company()
	if err != nil {
		fmt.Println(err)
	}

	render(c, gin.H{
		"title":   "Список компаний",
		"payload": companies}, "companies.html")

}

func render(c *gin.Context, data gin.H, templateName string) {

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		c.XML(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}

}

func main() {
	initDB()

	router = gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", indexHandler)
	router.GET("/list/company", companyHandler)

	router.Run()
}
