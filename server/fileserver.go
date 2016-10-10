package server;

import (
	"io"
	"log"
	"net/http"
	"flag"
)

func StartFileserver() { 
	log.Println("Server Startet")
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":"+flag.Lookup("Port").Value.String(), flag.Lookup("ServerKey").Value.String(), flag.Lookup("ServerCrt").Value.String(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "This is an example server.\n")
}
