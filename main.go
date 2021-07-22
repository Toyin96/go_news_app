package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tpl *template.Template

func init(){
	tpl = template.Must(template.ParseFiles("index.html"))
}

func Index(w http.ResponseWriter, req *http.Request){
	err := tpl.Execute(w, nil)
	if err != nil{
		log.Panicln(err)
	}
}

func main(){
	router := mux.NewRouter()
	err := godotenv.Load()
	if err != nil{
		log.Fatalln("Error loading .env file: ", err)
	}

	port := os.Getenv("port")

	http.Handle("/", router)
	router.HandleFunc("/", Index)
	panic(http.ListenAndServe(":"+port, router))
}