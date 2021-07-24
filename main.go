package main

import (
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	"go_news_app/news"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var tpl *template.Template
var newsapi *news.Client

func init(){
	tpl = template.Must(template.ParseFiles("index.html"))
}

type Search struct {
	Query string
	NextPage int
	TotalPages int
	Results *news.Results
}

func (s *Search) IsLastPage() bool{
	return s.NextPage >= s.TotalPages
}

func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}
	return s.NextPage - 1
}

func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}

func Index(w http.ResponseWriter, req *http.Request){
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, nil)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	buf.WriteTo(w)

}

func SearchHandler(newsapi *news.Client) http.HandlerFunc {
	return func (w http.ResponseWriter, req * http.Request){
		u, err := url.Parse(req.URL.String())
		if err != nil {
			log.Fatalln(err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()

		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		nextPage, err := strconv.Atoi(page)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		search := &Search{
			Query: searchQuery,
			NextPage: nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results: results,
		}

		if ok := !search.IsLastPage(); ok{
			search.NextPage++
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		buf.WriteTo(w)
	}
}

func main(){
	err := godotenv.Load()
	if err != nil{
		fmt.Println("It is already in heroku config var")
	}

	port := os.Getenv("PORT")

	fs := http.FileServer(http.Dir("./assets/"))

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == ""{
		fmt.Println("It is already in heroku config var")
	}

	myClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	newsapi= news.NewClient(myClient, apiKey, 20)


	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.HandleFunc("/", Index)
	http.HandleFunc("/search", SearchHandler(newsapi))
	panic(http.ListenAndServe(":"+port, nil))
}