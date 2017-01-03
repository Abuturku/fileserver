package main

import (
	"flag"
	"src/server"
)

func main(){
	flag.String("P", "1234", "The port of the server")
	flag.String("K","server/server.key","Path to key file")
	flag.String("C", "server/server.crt", "Path to certificate file")
	flag.String("L", "server/users.csv", "Path to file, where usernames, passwords and salts are stored")
	flag.String("T", "900", "Session timeout given in seconds")
	flag.String("F", "files/" ,"Folder where all Userfiles are stored")
	flag.Parse()
	
	
	server.StartFileserver()
}