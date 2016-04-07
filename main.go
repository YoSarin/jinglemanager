package main

import (
	"github.com/martin-reznik/logger"
	"net/http"
)

func main() {
	log := logger.NewLog(func(line *logger.LogLine) { line.Print() })
	log.LogSeverity[logger.DEBUG] = true
	defer log.Close()

	http.Handle("/", Index{log})

	http.Handle("/css/", Static{log})
	http.Handle("/js/", Static{log})
	http.Handle("/images/", Static{log})

	log.Info("Server is up and running, open 'http://localhost:8080' in your browser")

	http.ListenAndServe(":8080", nil)
}
