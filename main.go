package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	Id       uint16
	Title    string
	Anons    string
	FullText string
}

var posts = []Article{}
var showPost = Article{}

func connectToBD() *sql.DB {
	db, err := sql.Open("mysql", "sql11426417:iawK1tqJS6@tcp(sql11.freesqldatabase.com:3306)/sql11426417")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html")
	if err != nil {
		panic(err.Error())
	}

	err = t.ExecuteTemplate(w, "create", nil)

}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html")
	if err != nil {
		panic(err.Error())
	}

	res, err := connectToBD().Query("SELECT * FROM `articles` ")
	if err != nil {
		panic(err.Error())
	}

	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err.Error())
		}

		posts = append(posts, post)
	}
	err = t.ExecuteTemplate(w, "index", posts)

	defer connectToBD().Close()
}

func save_article(w http.ResponseWriter, r *http.Request) {

	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	insert, err := connectToBD().Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
	if err != nil {
		panic(err.Error())
	}
	defer connectToBD().Close()
	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html", "templates/header.html")
	if err != nil {
		panic(err.Error())
	}

	res, err := connectToBD().Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s' ", vars["id"]))
	if err != nil {
		panic(err.Error())
	}
	defer connectToBD().Close()

	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err.Error())
		}
		showPost = post
	}

	err = t.ExecuteTemplate(w, "show", showPost)

}

func delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	deleteArticle, err := connectToBD().Query(fmt.Sprintf("DELETE FROM `articles` WHERE `articles`.`id` = '%s'", vars["id"]))
	if err != nil {
		panic(err.Error())
	}

	deleteArticle.Close()
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/delete/{id:[0-9]+}", delete).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	http.Handle("/", rtr)
	http.ListenAndServe(":4614", nil)

}

func main() {
	handleFunc()
}
