package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//"bufio"
	"os"
	"strconv"
)

func StartFileserver() {
	log.Println("Server Startet")
	http.HandleFunc("/", index)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", newUserHandler)
	http.HandleFunc("/landrive", landrive)
	http.HandleFunc("/uploadFile", uploadFile)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("website"))))
	err := http.ListenAndServeTLS(":"+flag.Lookup("P").Value.String(), flag.Lookup("C").Value.String(), flag.Lookup("K").Value.String(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	if checkCookie(req) {
		http.Redirect(w, req, "/landrive", http.StatusMovedPermanently)
	} else {
		//io.WriteString(w, "This is an example server.\n")
		title := req.URL.Path[len("/"):]
		p, _ := loadPage(title)
		t, _ := template.ParseFiles("website/index.html")
		t.Execute(w, p)
	}
}

func landrive(w http.ResponseWriter, req *http.Request) {
	if checkCookie(req) {
		title := req.URL.Path[len("/"):]
		p, _ := loadPage(title)
		t, _ := template.ParseFiles("website/landrive.html")
		t.Execute(w, p)
	} else {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	}
}

func checkCookie(req *http.Request) bool {
	cookies := req.Cookies()

	for _, cookie := range cookies {
		cookieName := cookie.Name
		cookiePw := cookie.Value
		user := loadUser(cookieName)

		if cookiePw == hash([]string{user.name, user.password}) {
			return true
		}

	}
	return false
}

func uploadFile(w http.ResponseWriter, req *http.Request) {

	req.ParseMultipartForm(32 << 20)
	file, handler, err := req.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./file/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("User tried to log in")
	username := req.FormValue("username")
	password := req.FormValue("password")
	log.Println("User:", username, "Password:", password)
	user := loadUser(username)
	log.Println("Found user: ", user)
	authenticationSuccessful := authenticate(user, password)

	if authenticationSuccessful {
		loginUser(user, w, req)
	} else {
		http.Redirect(w, req, "?login=false", http.StatusMovedPermanently)
	}

}

func loginUser(user *user, w http.ResponseWriter, req *http.Request) {
	cookieValue := hash([]string{user.name, user.password})
	maxAge, _ := strconv.Atoi(flag.Lookup("T").Value.String())
	cookie := http.Cookie{Name: user.name, Value: cookieValue, MaxAge: maxAge}
	log.Println("Setting cookie")
	http.SetCookie(w, &cookie)
	log.Println("Redirecting to landrive")
	http.Redirect(w, req, "/landrive", http.StatusMovedPermanently)
}

func newUserHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("User tried to register")
	username := req.FormValue("username")
	password := req.FormValue("password")
	password2 := req.FormValue("password2")

	if password == password2 {
		user := loadUser(username)
		if user.name != "" {
			http.Redirect(w, req, "?register=userfalse", http.StatusMovedPermanently)
		} else {
			user := createUser(username, password)
			loginUser(&user, w, req)
		}
	} else {
		http.Redirect(w, req, "?register=pwfalse", http.StatusMovedPermanently)
	}

}

func createUser(username string, password string) user {
	salt := generateSalt()
	hashedPw := hash([]string{password, salt})
	log.Println("Creating user with parameters: ", username, hashedPw, salt)
	f, err := os.OpenFile(flag.Lookup("L").Value.String(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	writer := csv.NewWriter(f)
	defer f.Close()

	//writer.Write(username)
	//writer.Write(hashedPw)
	//writer.Write(salt)
	log.Println("Writing to csv")
	writer.Write([]string{username, hashedPw, salt})
	writer.Flush()
	err = writer.Error()
	if err != nil {
		log.Println(err)
	}
	
	createFolder(username)
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
	name     string
	password string
	salt     string
}

func loadUser(username string) *user {
	f, _ := os.Open(flag.Lookup("L").Value.String())
	r := csv.NewReader(f)
	defer f.Close()
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if record[0] == username {
			return &user{name: record[0], password: record[1], salt: record[2]}
		}
		//log.Println(record)
	}
	return &user{name: "", password: "", salt: ""}
}

func authenticate(user *user, password string) bool {
	//hasher := sha256.New()
	//hasher.Write([]byte(password))
	//hasher.Write([]byte(user.salt))

	hash := hash([]string{password, user.salt})

	if hash == user.password {
		return true
	}
	return false
}

func hash(strings []string) string {
	hasher := sha256.New()
	for _, value := range strings {
		hasher.Write([]byte(value))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func createFolder(path string){
	os.Mkdir("files/"+path, 0777)
}
