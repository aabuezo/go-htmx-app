package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB
var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func initDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err = sql.Open(
		os.Getenv("DB_CONN"),
		os.Getenv("DB_USERNAME")+":"+
			os.Getenv("DB_PASSWORD")+"@tcp(127.0.0.1:"+
			os.Getenv("DB_PORT")+")/"+
			os.Getenv("DB_NAME")+"?parseTime=True&loc=Local")

	if err != nil {
		log.Fatal(err)
	}

	// check database connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	initDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	//gRouter := mux.NewRouter()
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/static/css/styles.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/css/styles.css")
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	err := tmpl.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}
