package server

import (
	"encoding/csv"
	"flag"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

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

	flag.String("L", pathToFile, "Path to file, where usernames, passwords and salts are stored")
	flag.String("T", "900", "Session timeout given in seconds")
	flag.String("F", "files/", "Folder where all Userfiles are stored")
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

//Zur weiteren Identifikation des Nutzers soll ein Session-ID Cookie ver- wendet werden.
func TestCookie(t *testing.T) {
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

// Neue Nutzer sollen selbst einen Zugang anlegen ko ̈nnen.
func TestCreateUser(t *testing.T) {

}

// Ein Nutzer soll sein Passwort  aendern koennen.
func TestChangePassword(t *testing.T) {

}

// Es soll mo ̈glich sein, Dateien ”hochzuladen“
func TestSaveFile(t *testing.T) {

}

// Es soll mo ̈glich sein, Dateien ”herunterzuladen“
func TestDownloadFile(t *testing.T) {

}

// Es soll mo ̈glich sein, Dateien ”herunterzuladen“ ueber wget
func TestDownloadFileWGET(t *testing.T) {

}

//Es mo ̈glich sein, Dateien zu lo ̈schen.
func TestDeleteFile(t *testing.T) {

}

//Auch in diese Unterordner sollen sich Dateien laden lassen.
func TestCreateFolder(t *testing.T) {

}
