package main

import (
	"net/http"
	"html/template"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var templates = template.Must(template.ParseFiles("./template/user_list.html"))

type User struct {
	Id 		int
	Email 	string
}

func main() {
	http.HandleFunc("/list/", listHandler)
	http.HandleFunc("/add/", addHandler)
	http.ListenAndServe(":8080", nil)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	email := r.FormValue("email")

	log.Println("id = " + id)
	log.Println("email = " + email)

	db, err := sql.Open("sqlite3", "./db/test.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO user(id, email) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id, email)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/list/", http.StatusFound)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/test.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows, err := db.Query("select id, email from user")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		user := User{}

		err := rows.Scan(&user.Id, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
		
		users = append(users, user)
	}

	log.Print("Users size = ")
	log.Println(len(users))

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	renderTemplate(w, "user_list", &users)
}

func renderTemplate(w http.ResponseWriter, tmpl string, users *[]User) {
	err := templates.ExecuteTemplate(w, tmpl+".html", users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
