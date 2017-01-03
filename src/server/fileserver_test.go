package server

import (
	"encoding/csv"
	"flag"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	//"log"
	//"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func generateCookie() http.Cookie {
	cookieValue := hash([]string{"Andy", "a879518e72e3aa6d82126e52d6a641e66005d68b44a31ea5797d0e24f90fd759"})
	maxAge, _ := strconv.Atoi(flag.Lookup("T").Value.String())
	cookie := http.Cookie{Name: "Andy", Value: cookieValue, MaxAge: maxAge, Expires: time.Now().Add(15 * time.Minute)}
	return cookie
}

func init() {
	pathToFile := "./user_test.csv"
	if _, err := os.Stat(pathToFile); err == nil {
		os.Remove(pathToFile)
	}

	file, err := os.Create(pathToFile)

	writer := csv.NewWriter(file)
	defer file.Close()

	writer.Write([]string{"Andy", "a879518e72e3aa6d82126e52d6a641e66005d68b44a31ea5797d0e24f90fd759", "0912951feb016907a1b762c7f83de9b0"})
	writer.Flush()
	err = writer.Error()
	if err != nil {

	}

	os.Mkdir("test", 0777)

	flag.String("L", pathToFile, "Path to file, where usernames, passwords and salts are stored")
	flag.String("T", "900", "Session timeout given in seconds")
	flag.String("F", "test/", "Folder where all Userfiles are stored")

}

//Der Zugang soll durch Benutzernamen und Passwort geschützt werden. Positives Beispiel
func TestAccess(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	v := url.Values{}
	v.Add("username", "Andy")
	v.Add("password", "andy")
	req.Form = v

	rr := httptest.NewRecorder()

	loginHandler(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

//Der Zugang soll durch Benutzernamen und Passwort geschützt werden. Negatives Beispiel: FalschesPassword
func TestAccessWrongPassword(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	v := url.Values{}
	v.Add("username", "Andy")
	v.Add("password", "andy1")
	req.Form = v

	rr := httptest.NewRecorder()

	loginHandler(rr, req)

	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

//Der Zugang soll durch Benutzernamen und Passwort geschützt werden. Negatives Beispiel: User exestiert nicht
func TestAccessUserDoesntExist(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	v := url.Values{}
	v.Add("username", "Andy1")
	v.Add("password", "andy")
	req.Form = v

	rr := httptest.NewRecorder()

	loginHandler(rr, req)

	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

//Zur weiteren Identifikation des Nutzers soll ein Session-ID Cookie verwendet werden.
func TestValidCookie(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	v := url.Values{}
	v.Add("username", "Andy")
	v.Add("password", "andy")
	req.Form = v

	rr := httptest.NewRecorder()

	loginHandler(rr, req)

	cookie := generateCookie()

	req.AddCookie(&cookie)

	isValid, _, _ := checkCookie(rr, req)

	if !isValid {
		t.Errorf("Cookie check failed. Expected true got %v", isValid)
	}

}

func TestUnvalidCookie(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	v := url.Values{}
	v.Add("username", "Andy")
	v.Add("password", "andy")
	req.Form = v

	rr := httptest.NewRecorder()

	loginHandler(rr, req)
	cookieValue := hash([]string{"Andy1", "a879518e72e3aa6d82126e52d6a641e66005d68b44a31ea5797d0e24f90fd759"})
	maxAge, _ := strconv.Atoi(flag.Lookup("T").Value.String())
	cookie := http.Cookie{Name: "Andy", Value: cookieValue, MaxAge: maxAge, Expires: time.Now().Add(15 * time.Minute)}

	req.AddCookie(&cookie)

	isValid, _, _ := checkCookie(rr, req)

	if isValid {
		t.Errorf("Cookie check failed. Expected false got %v", isValid)
	}
}

// Neue Nutzer sollen selbst einen Zugang anlegen können.
func TestCreateValidUser(t *testing.T) {
	req, err := http.NewRequest("POST", "/register", nil)

	if err != nil {
		t.Fatal(err)
	}

	v := url.Values{}
	v.Add("username", "Niklas")
	v.Add("password", "niklas")
	v.Add("password2", "niklas")
	req.Form = v

	rr := httptest.NewRecorder()

	newUserHandler(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}
}

func TestCreateUserPwFalse(t *testing.T) {
	req, err := http.NewRequest("POST", "/register", nil)

	if err != nil {
		t.Fatal(err)
	}

	v := url.Values{}
	v.Add("username", "Niklas")
	v.Add("password", "niklas")
	v.Add("password2", "niklas1")
	req.Form = v

	rr := httptest.NewRecorder()

	newUserHandler(rr, req)

	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusMovedPermanently)
	}
}

func TestCreateUserNameFalse(t *testing.T) {
	req, err := http.NewRequest("POST", "/register", nil)

	if err != nil {
		t.Fatal(err)
	}

	v := url.Values{}
	v.Add("username", "Andy")
	v.Add("password", "niklas")
	v.Add("password2", "niklas")
	req.Form = v

	rr := httptest.NewRecorder()

	newUserHandler(rr, req)

	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusMovedPermanently)
	}
}

// Es soll möglich sein, Dateien ”hochzuladen“
func TestSaveFile(t *testing.T) {

	postData := `--xxx Content-Disposition: form-data; name="user_test.csv"; filename="user_test.csv" Content-Type: application/octet-stream Content-Transfer-Encoding: binary binary data --xxx--`
	req, err := http.NewRequest("POST", "/uploadFile", ioutil.NopCloser(strings.NewReader(postData)))

	if err != nil {
		t.Fatal(err)
	}

	cookie := generateCookie()
	req.AddCookie(&cookie)
	rr := httptest.NewRecorder()
	req.PostForm = url.Values{}
	req.PostForm.Add("path", "")
	req.Header.Set("Content-Type", `multipart/form-data; boundary="xxx"`)
	req.ParseMultipartForm(32 << 20)

	uploadFileHandler(rr, req)

}

// Es soll möglich sein, Dateien ”herunterzuladen“
func TestDownloadFile(t *testing.T) {
	//req, err := http.NewRequest("POST", "/download", nil)

	//if err != nil {
	//	t.Fatal(err)
	//}

	//cookie := generateCookie()
	//req.AddCookie(&cookie)

	//v := url.Values{}
	//v.Add("path", "test/Andy/user_test.csv")
	//req.Form = v

	//rr := httptest.NewRecorder()

	//downloadHandler(rr, req)

	//if status := rr.Code; status != http.StatusOK {
	//	t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	//}
}

// Es soll möglich sein, Dateien ”herunterzuladen“ über wget
func TestDownloadFileWGET(t *testing.T) {
}

//Auch in diese Unterordner sollen sich Dateien laden lassen.
func TestCreateFolder(t *testing.T) {
	req, err := http.NewRequest("POST", "/newFolder", nil)

	if err != nil {
		t.Fatal(err)
	}

	cookie := generateCookie()
	req.AddCookie(&cookie)

	v := url.Values{}
	v.Add("path", "")
	v.Add("newFolderName", "testFolder")
	req.Form = v

	rr := httptest.NewRecorder()

	createFolderHandler(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}
}

//Es möglich sein, Ordner zu löschen.
func TestDeleteFolder(t *testing.T) {
	req, err := http.NewRequest("POST", "/delete", nil)

	if err != nil {
		t.Fatal(err)
	}

	cookie := generateCookie()
	req.AddCookie(&cookie)

	v := url.Values{}
	v.Add("path", "testFolder")
	req.Form = v

	rr := httptest.NewRecorder()

	deleteHandler(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}
}

// Ein Nutzer soll sein Passwort ändern können.
func TestChangePasswordValid(t *testing.T) {
	req, err := http.NewRequest("POST", "/changePw", nil)

	if err != nil {
		t.Fatal(err)
	}

	v := url.Values{}
	v.Add("oldPassword", "andy")
	v.Add("newPassword", "niklas")
	v.Add("newPassword2", "niklas")
	req.Form = v

	cookie := generateCookie()
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()

	changePasswordHandler(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}
}
