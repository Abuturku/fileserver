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
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":"+flag.Lookup("Port").Value.String(), flag.Lookup("ServerCrt").Value.String(), flag.Lookup("ServerKey").Value.String(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	//io.WriteString(w, "This is an example server.\n")
	title := req.URL.Path[len("/website/"):]
    p := loadPage(title)
    t, _ := template.ParseFiles("index.html")
    t.Execute(w, p)
}

func loadPage(title string) *Page {
    filename := title + ".txt"
    body, _ := ioutil.ReadFile(filename)
    return &Page{Title: title, Body: body}
}

type Page struct {
    Title string
    Body  []byte
}

 
 
