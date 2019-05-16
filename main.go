package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Articles struct which contains an array of artilce
type Articles struct {
	APIKey   string
	Articles []Article `json:"articles"`
}

// Article struct which contains a list of source, an author, a title, a description, an URL, an urlToImage, a publishedAt, a content
type Article struct {
	Source      Source `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	URLToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
}

// Source struct which contains an ID, a name
type Source struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func sortArticles(res []Article, sort string) {
	if sort == "title" && len(res) > 1 {
		for i := 0; i < len(res)-1; i++ {
			for j := i + 1; j < len(res); j++ {
				if res[i].Title > res[j].Title {
					res[i], res[j] = res[j], res[i]
				}
			}
		}
	} else if sort == "publishedat" && len(res) > 1 {
		layout := "2006-01-02T15:04:05Z"
		for i := 0; i < len(res)-1; i++ {
			for j := i + 1; j < len(res); j++ {
				t1, err := time.Parse(layout, res[i].PublishedAt)
				if err != nil {
					fmt.Println(err)
				}
				t2, err := time.Parse(layout, res[j].PublishedAt)
				if err != nil {
					fmt.Println(err)
				}
				if t2.Before(t1) {
					res[i], res[j] = res[j], res[i]
				}
			}
		}
	}
}

func getNews(w http.ResponseWriter, r *http.Request) {
	database, err := os.Open("database.json")
	if err != nil {
		fmt.Println(err)
	}
	defer database.Close()

	byteValue, _ := ioutil.ReadAll(database)

	var articles Articles
	var result []Article

	json.Unmarshal([]byte(byteValue), &articles)

	q := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sortBy")

	keyTable := articles.APIKey
	articlesTable := articles.Articles
	fmt.Println(keyTable)

	for i := 0; i < len(articlesTable); i++ {
		if strings.Contains(articlesTable[i].Title, q) == false {
			if strings.Contains(articlesTable[i].Content, q) == false {
				if strings.Contains(articlesTable[i].Description, q) != false {
					result = append(result, articlesTable[i])
				}
			} else {
				result = append(result, articlesTable[i])
			}
		} else {
			result = append(result, articlesTable[i])
		}
	}

	if strings.ToLower(sort) == "publishedat" || strings.ToLower(sort) == "title" {
		sortArticles(result, sort)
	}

	render.JSON(w, r, result)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/v2/everything", getNews)
	http.ListenAndServe(":8080", r)
}
