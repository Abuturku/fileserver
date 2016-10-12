package server

import (
	"flag"
	"html/template"
	//"io"
	"log"
	"net/http"
	"io/ioutil"
)

func StartFileserver() {
	log.Println("Server Startet")
	http.HandleFunc("/", index)
	http.HandleFunc("/login", loginHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("website"))))
	err := http.ListenAndServeTLS(":"+flag.Lookup("Port").Value.String(), flag.Lookup("ServerCrt").Value.String(), flag.Lookup("ServerKey").Value.String(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	//io.WriteString(w, "This is an example server.\n")
	title := req.URL.Path[len("/"):]
    p, _ := loadPage(title)
    t, _ := template.ParseFiles("website/index.html")
    t.Execute(w, p)
}

func loginHandler(w http.ResponseWriter, req *http.Request){
	log.Println("User tried to log in")
	username := req.FormValue("username")
	password := req.FormValue("password")
	log.Println("User:", username, "Password:", password)
	
	
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

type Page struct {
    Title string
    Body  []byte
}