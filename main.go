package main

import (
	"flag"
	"landrive/server"
)

func main(){
	flag.String("Port", "1234", "The port of the server")
	flag.String("ServerKey","server/server.key","Path to key file")
	flag.String("ServerCrt", "server/server.crt", "Path to certificate file")
	flag.Parse()
	
	
	server.StartFileserver()
}