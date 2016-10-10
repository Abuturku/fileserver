package fileserver;

import (
	"io"
	"log"
	"net/http"
)

func StartFileserver() { 
	configuration :=  GetConfig() 
	log.Println("Server Startet")
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":"+configuration.Port, configuration.ServerCrt, configuration.ServerKey, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "This is an example server.\n")
}
