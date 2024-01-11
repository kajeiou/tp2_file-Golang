package api_mode

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	logFile *os.File
	logger  *log.Logger
)

func init() {
	currentTime := time.Now()
	logFileName := fmt.Sprintf("app_%s.log", currentTime.Format("2006-01-02"))

	var err error
	logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Impossible d'ouvrir le fichier de log:", err)
	}

	logger = log.New(logFile, "DICO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogToFile(context, message string) {
	logMessage := fmt.Sprintf("[%s] %s", context, message)
	logger.Println(logMessage)
}

func LogAndRespond(w http.ResponseWriter, r *http.Request, message string, status int) {
	LogToFile(r.URL.Path, fmt.Sprintf("Requête reçue : %s %s", r.Method, r.URL.Path))
	w.WriteHeader(status)
	fmt.Fprintln(w, message)
}
