package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	FirstName string
	LastName  string
}

func (u *User) String() string {
	return fmt.Sprintf("First name is: %s, Lat name is: %s", u.FirstName, u.LastName)
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	db.Create(&User{FirstName: "El", LastName: "Nino"})

	user := User{}
	db.First(&user)

	fmt.Println(&user)

	mux := http.NewServeMux()

	mux.HandleFunc("/", showForm)
	mux.HandleFunc("/save", saveForm)
	mux.HandleFunc("/time", getTime)

	fs := http.FileServer(http.Dir("./"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Listen on port 3000")

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}

func showForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./html/form.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, []int{1, 2, 3, 4, 5, 6, 7})
	if err != nil {
		log.Fatal(err)
	}
}

func saveForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	pilou := User{
		FirstName: r.Form.Get("first_name"),
	}

	fmt.Printf("first name is: %v", pilou)
}

func getTime(w http.ResponseWriter, r *http.Request) {
	tmplString := fmt.Sprintf("<span>%s<span>", time.Now())

	tmpl, err := template.New("time").Parse(tmplString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(time.Now())

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
