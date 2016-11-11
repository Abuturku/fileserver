package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//Die Web-Seite soll nur per HTTPS erreichbar sein.
func TestPingPage(t *testing.T) {

}

//Der Zugang soll durch Benutzernamen und Passwort geschuëtzt werden.
func TestAcess(t *testing.T) {

}

//Die Passwoerter duërfen nicht im Klartext gespeichert werden.
func TestPassword(t *testing.T) {

}

//Zur weiteren Identifikation des Nutzers soll ein Session-ID Cookie ver- wendet werden.
func TestCookie(t *testing.T) {

}

// Neue Nutzer sollen selbst einen Zugang anlegen ko ̈nnen.
func TestCreatUser(t *testing.T) {

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
func TestDeletFile(t *testing.T) {

}

//Auch in diese Unterordner sollen sich Dateien laden lassen.
func TestCreateFolder(t *testing.T) {

}
 
