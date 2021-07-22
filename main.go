package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go_news_app_git/news"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var tpl *template.Template
var newsClient *news.Client

func init(){
	tpl = template.Must(template.ParseFiles("index.html"))
}

func Index(w http.ResponseWriter, req *http.Request){
	err := tpl.Execute(w, nil)
	if err != nil{
		log.Panicln(err)
	}
}

func Search(newApi *news.Client) http.HandlerFunc {
	return func (w http.ResponseWriter, req * http.Request){
		u, err := url.Parse(req.URL.String())
		if err != nil {
			log.Fatalln(err.Error(), http.StatusInternalServerError)
		}

		params := u.Query()

		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		fmt.Printf("The search query is %s", searchQuery)
		fmt.Printf("The page is %s", page)
	}
}

func main(){
	err := godotenv.Load()
	if err != nil{
		log.Fatalln("Error loading .env file: ", err)
	}

	port := os.Getenv("port")

	fs := http.FileServer(http.Dir("./assets/"))

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == ""{
		log.Fatalln("ENV: API key must be set")
	}

	myClient := &http.Client{
		Timeout:       10 * time.Second,
	}
	newsClient = news.NewClient(myClient, apiKey, 20)


	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.HandleFunc("/", Index)
	http.HandleFunc("/search", Search(newsClient))
	panic(http.ListenAndServe(":"+port, nil))
}