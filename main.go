package main

import (
	"flag"
	"landrive/server"
)

func main(){
	flag.String("P", "443", "The port of the server")
	flag.String("K","server/server.key","Path to key file")
	flag.String("C", "server/server.crt", "Path to certificate file")
	flag.String("L", "server/users.csv", "Path to file, where usernames, passwords and salts are stored")
	flag.String("T", "900", "Session timeout given in seconds")
	flag.Parse()
	
	
	server.StartFileserver()
}