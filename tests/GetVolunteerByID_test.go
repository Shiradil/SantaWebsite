package tests

import (
	"SantaWeb/internal/db"
	"SantaWeb/internal/handlers"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetVolunteerByID(t *testing.T) {
	err := db.DbConnection()
	volunteer, err := handlers.GetVolunteerByID("65bdf00b869485d29a4c66e0")

	if err != nil {
		t.Failed()
		return
	}
	if volunteer.Name != "Adilzhan" {
		t.Failed()
		return
	}
}

func TestHomeHandler(t *testing.T) {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	logFile, _ := os.OpenFile("cmd/logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	handlers.InitLogger(log)

	server := httptest.NewServer(http.HandlerFunc(handlers.HomeHandler))
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Failed()
	}

}
