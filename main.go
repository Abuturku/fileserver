package main

import "flag"
import "landrive/server"

func main(){
	flag.String("Port", "1234", "The Port of the Server")
	flag.String("ServerKey","src/landrive/server/server.key","The Server Key File")
	flag.String("ServerCrt", "src/landrive/server/server.crt", "The Server crt File")
	flag.Parse()
	
	
	server.StartFileserver()
}