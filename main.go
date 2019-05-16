package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	_ "github.com/lib/pq"
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

const (
	HOST       = "localhost"
	PORT       = 5432
	DBUSER     = "postgres"
	DBPASSWORD = "password"
	DBNAME     = "postgres"
)

func database() {
	database, err := os.Open("database.json")
	if err != nil {
		fmt.Println(err)
	}
	defer database.Close()

	byteValue, _ := ioutil.ReadAll(database)

	var articles Articles

	json.Unmarshal([]byte(byteValue), &articles)

	keyTable := articles.APIKey
	articlesTable := articles.Articles

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, DBUSER, DBPASSWORD, DBNAME)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		fmt.Println("ERROR :", err)
	}

	for i := 0; i < len(articlesTable); i++ {
		sqlStatement := "INSERT INTO Articles VALUES (" + strconv.Itoa(i+1) + `, '{"id":` + strconv.Itoa(articlesTable[i].Source.ID) + `,"name":"` + articlesTable[i].Source.Name + `"}', ` + `'` + articlesTable[i].Author + `', '` + articlesTable[i].Title + `', '` + articlesTable[i].Description + `', '` + articlesTable[i].URL + `', '` + articlesTable[i].URLToImage + `', '` + articlesTable[i].PublishedAt + `', '` + articlesTable[i].Content + `')`

		_, err := db.Exec(sqlStatement)
		if err != nil {
			fmt.Println(err)
		}
	}

	sqlStatement := "INSERT INTO Key VALUES (1,'" + keyTable + "')"

	_, err = db.Exec(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}

}

func getNews(w http.ResponseWriter, r *http.Request) {
	database()
	q := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sortBy")

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, DBUSER, DBPASSWORD, DBNAME)

	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		fmt.Println("ERROR :", err)
	}

	var sqlStatement string
	if strings.ToLower(sort) == "title" {
		sqlStatement = "SELECT * FROM Articles WHERE title LIKE '%" + q + "%' OR decription LIKE '%" + q + "%' OR content LIKE '%" + q + "%' ORDER BY title"
	} else if strings.ToLower(sort) == "publishedat" {
		sqlStatement = "SELECT * FROM Articles WHERE title LIKE '%" + q + "%' OR decription LIKE '%" + q + "%' OR content LIKE '%" + q + "%' ORDER BY publishedAt"
	} else {
		sqlStatement = "SELECT * FROM Articles WHERE title LIKE '%" + q + "%' OR decription LIKE '%" + q + "%' OR content LIKE '%" + q + "%'"
	}

	result, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println("ERROR :", err)
	}

	var out []Article

	for result.Next() {
		var (
			id          int
			source      string
			author      string
			title       string
			decription  string
			url         string
			urlToImage  string
			publishedAt string
			content     string
			tempSource  Source
			temp        Article
		)
		err = result.Scan(&id, &source, &author, &title, &decription, &url, &urlToImage, &publishedAt, &content)
		if err != nil {
			fmt.Println(err)
		}

		err = json.Unmarshal([]byte(source), &tempSource)
		if err != nil {
			fmt.Println(err)
		}

		temp.Source = tempSource
		temp.Author = author
		temp.Title = title
		temp.Description = decription
		temp.URL = url
		temp.URLToImage = urlToImage
		temp.PublishedAt = publishedAt
		temp.Content = content
		out = append(out, temp)
	}
	render.JSON(w, r, out)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/v2/everything", getNews)
	http.ListenAndServe(":8080", r)
}
