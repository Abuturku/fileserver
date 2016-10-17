package server

import (
	"fmt"
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"crypto/rand"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	//"bufio"
	"os"
	"strconv"
)

func StartFileserver() {
	log.Println("Server Startet")
	http.HandleFunc("/", index)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", newUserHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("website"))))
	err := http.ListenAndServeTLS(":"+flag.Lookup("P").Value.String(), flag.Lookup("C").Value.String(), flag.Lookup("K").Value.String(), nil)
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
	user := loadUser(username)
	log.Println("Found user: ", user)
	authenticationSuccessful := authenticate(user, password)
	
	if authenticationSuccessful {
		loginUser(user, w, req)
	}else{
		http.Redirect(w, req, "?login=false", 301)
	}
	
}

func loginUser(user *user, w http.ResponseWriter, req *http.Request){
		cookieValue := hash([]string{user.name, user.password})
		maxAge, _ := strconv.Atoi(flag.Lookup("T").Value.String())
		cookie := http.Cookie{Name: user.name, Value: cookieValue, MaxAge: maxAge}
		log.Println("Setting cookie")
		http.SetCookie(w, &cookie)
		log.Println("Redirecting to landrive")
		http.Redirect(w, req, "/landrive", 301)
}

func newUserHandler(w http.ResponseWriter, req *http.Request){
	log.Println("User tried to register")
	username := req.FormValue("username")
	password := req.FormValue("password")
	password2 := req.FormValue("password2")
	
	if password == password2 {
		user := loadUser(username)
		if user.name != "" {
			http.Redirect(w, req, "?register=userfalse", 301)
		}else{
			user := createUser(username, password)
			loginUser(&user, w, req)
		}
	}else{
		http.Redirect(w, req, "?register=pwfalse", 301)
	}
	
}


func createUser(username string, password string) user{
	salt := generateSalt();
	hashedPw := hash([]string{password, salt})
	log.Println("Creating user with parameters: ", username, hashedPw, salt)
	f, _ := os.OpenFile(flag.Lookup("L").Value.String(), os.O_APPEND, 0644)
	writer := csv.NewWriter(f)
	defer f.Close()
	
	
	//writer.Write(username)
	//writer.Write(hashedPw)
	//writer.Write(salt)
	log.Println("Writing to csv")
	writer.Write([]string{username, hashedPw, salt})
	writer.Flush()
	return user{name: username, password: hashedPw, salt: salt}
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


func generateSalt() string {
	saltSize := 16
	buf := make([]byte, saltSize)
	_, err := rand.Read(buf)

	if err != nil {
		fmt.Printf("random read failed: %v", err)
	}
	
	return hex.EncodeToString(buf)
}

type user struct {
    name string
    password string
    salt string
}

func loadUser(username string) *user{
	f, _ := os.Open(flag.Lookup("L").Value.String())
	r := csv.NewReader(f)
	defer f.Close()
	for {
		record, err := r.Read()
		if err == io.EOF{
			break
		}
		if record[0] == username{
			return &user{name: record[0], password: record[1], salt: record[2]}
		}
		//log.Println(record)
	}
	return &user{name: "", password: "", salt: ""}
}

func authenticate(user *user, password string) bool{
	//hasher := sha256.New()
	//hasher.Write([]byte(password))
	//hasher.Write([]byte(user.salt))
	
	hash := hash([]string{password, user.salt})
	
	if hash == user.password {
		return true
	}
	return false
}

func hash(strings []string) string{
	hasher := sha256.New()
	for _, value := range strings {
		hasher.Write([]byte(value))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}