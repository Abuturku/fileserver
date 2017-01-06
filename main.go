//Authors: Andreas Schick (2792119), Linda Latreider (7743782), Niklas Nikisch (9364290)
package main

import (
	"flag"
	"landrive/server"
)

/*
Main Methode setzt alle Startparameter und startet anschließend des Server
*/
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