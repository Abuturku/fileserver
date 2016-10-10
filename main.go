package main

import "flag"
import "de/vorlesung/projekt/landrive/server"

func main(){
	flag.String("Port", "1234", "The Port of the Server")
	flag.String("ServerKey","src/de/vorlesung/projekt/landrive/server/server.key","The Server Key File")
	flag.String("ServerCrt", "src/de/vorlesung/projekt/landrive/server/server.crt", "The Server crt File")
	flag.Parse()
	
	
	server.StartFileserver()
}