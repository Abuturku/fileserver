package server

import (
	//"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func StartFileserver() {
	log.Println("Server Startet")
	http.HandleFunc("/", index)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", newUserHandler)
	http.HandleFunc("/landrive", landrive)
	http.HandleFunc("/uploadFile", uploadFileHandler)
	http.HandleFunc("/getFolderStruct", folderStructHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/newFolder", createFolderHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/download", downloadHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("website"))))
	err := http.ListenAndServeTLS(":"+flag.Lookup("P").Value.String(), flag.Lookup("C").Value.String(), flag.Lookup("K").Value.String(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	cookiecheck, _, _ := checkCookie(w, req)
	if cookiecheck {
		http.Redirect(w, req, "/landrive", http.StatusMovedPermanently)
	} else {
		t, _ := template.ParseFiles("website/index.html")
		t.Execute(w, nil)
	}
}

func landrive(w http.ResponseWriter, req *http.Request) {
	cookiecheck, _, _ := checkCookie(w, req)
	if cookiecheck {
		t, err := template.ParseFiles("website/landrive.html")
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(string(b))

		t.Execute(w, nil)
	} else {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	}
}

func folderStructHandler(w http.ResponseWriter, req *http.Request) {
	cookiecheck, user, _ := checkCookie(w, req)
	if cookiecheck {
		log.Println("Loading FolderStruct for user " + user.name)
		folders := getFolderStruct(user.name)
		js, err := json.Marshal(folders)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	}
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	cookiecheck, user, cookie := checkCookie(w, req)
	if cookiecheck {
		log.Println("Logout: " + user.name)
		cookie.Expires = time.Now().Add(-1)
		cookie.Value = ""
		http.SetCookie(w, &cookie)
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	} else {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	}
}

func createFolderHandler(w http.ResponseWriter, req *http.Request) {
	cookiecheck, user, _ := checkCookie(w, req)
	if cookiecheck {

		path := req.FormValue("path")
		log.Println("Create Folder: " + path)
		createFolder(user.name + "/" + path)
	} else {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	}
}

func deleteHandler(w http.ResponseWriter, req *http.Request) {
	cookiecheck, user, _ := checkCookie(w, req)
	if cookiecheck {

		path := req.FormValue("path")
		log.Println("delete: " + path)
		os.RemoveAll(flag.Lookup("F").Value.String() + user.name + "/" + path)
	} else {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	}
}

func downloadHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Entering downloadHandler")
	cookiecheck, user, _ := checkCookie(w, req)
	if cookiecheck {
		path := req.FormValue("path")
		stringarray := strings.Split(path, "/")
		log.Println("Download File: " + path)
 
		w.Header().Set("Content-Disposition", "attachment; filename=\""+stringarray[len(stringarray)-1]+"\"")

		http.ServeFile(w, req, flag.Lookup("F").Value.String() + user.name + "/"+path)
 
	} else {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	}
}

func checkCookie(w http.ResponseWriter, req *http.Request) (bool, user, http.Cookie) {
	cookies := req.Cookies()

	for _, cookie := range cookies {
		cookieName := cookie.Name
		cookiePw := cookie.Value
		user := loadUser(cookieName)

		if cookiePw == hash([]string{user.name, user.password}) {
			//following 3 lines are responsible for setting the cookie expiration date 15 minutes to the future
			expiration := time.Now().Add(15 * time.Minute)
			cookie.Expires = expiration
			http.SetCookie(w, cookie)
			return true, *user, *cookie
		}

	}
	return false, user{}, http.Cookie{}
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

func createFolder(path string) {
	os.Mkdir(flag.Lookup("F").Value.String()+path, 0777)
}

type Folder struct {
	Name    string
	Files   []File
	Folders []Folder
}

type File struct {
	Name string
	Date time.Time
	Size int64
}

func getFolderStruct(path string) Folder {
	//log.Println(path)
	index := strings.Index(path, "/")
	var name string
	if index > 0 {
		name = path[index+1:]
	} else {
		name = path
	}
	files := make([]File, 0)
	folders := make([]Folder, 0)
	//log.Println(name + ": ")
	fileinfos, _ := ioutil.ReadDir(flag.Lookup("F").Value.String() + "/" + path)

	for _, file := range fileinfos {
		if file.IsDir() {
			folders = append(folders, getFolderStruct(path+"/"+file.Name()))
		} else {
			//log.Println(file.Name())
			fileStruct := File{Name: file.Name(), Date: file.ModTime(), Size: file.Size()}
			files = append(files, fileStruct)
		}
	}
	folder := Folder{Name: name, Files: files, Folders: folders}
	return folder
}

func uploadFileHandler(w http.ResponseWriter, req *http.Request) {
	cookiecheck, user, _ := checkCookie(w, req)
	if cookiecheck {
		log.Println("Request to upload a file was made from user " + user.name)

		defer http.Redirect(w, req, "/landrive", http.StatusMovedPermanently)

		//Fileupload orientiert nach https://www.socketloop.com/tutorials/golang-upload-file
		file, header, err := req.FormFile("uploadFile")
		folderPath := req.FormValue("path")
		log.Println("Path where file is going to be placed: " + folderPath)
		if err != nil {
			log.Println(w, err)
			return
		}

		defer file.Close()

		filepath := flag.Lookup("F").Value.String() + user.name + "/" + folderPath + "/" + header.Filename

		out, err := os.Create(filepath)
		if err != nil {
			log.Println(w, "Unable to create the file for writing. Check your write access privilege. Path: "+filepath)
			return
		}

		defer out.Close()

		// write the content from POST to the file
		_, err = io.Copy(out, file)
		if err != nil {
			log.Println(w, err)
		}

		log.Println(w, "File uploaded successfully : ")
		log.Println(w, header.Filename)
	}
}
