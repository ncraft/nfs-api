package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	util "github.com/ncraft/machinery/pkg/base"
	"github.com/ncraft/nfs-api/pkg/nfs"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/shares", addShare)

	server := http.Server{
		Handler: router,
	}

	listenerPath := util.MustBeSet(os.Getenv("LISTENER_PATH"))

	listener, err := net.Listen("unix", listenerPath)
	if err != nil {
		panic(err)
	}

	log.Printf("Serving at %s [%s]", listenerPath, "unix")

	log.Fatalln(server.Serve(listener))
}

func addShare(w http.ResponseWriter, r *http.Request) {
	writeJsonResponse := func(shareResponse *nfs.ShareResponse) {
		_ = json.NewEncoder(w).Encode(*shareResponse)
	}

	errorHandler := func(status int, msg string) {
		log.Printf("error: %s", msg)

		w.WriteHeader(status)

		writeJsonResponse(&nfs.ShareResponse{
			Status:  status,
			Message: msg,
		})
	}

	shareRequest, err := nfs.JsonDecode(r.Body)
	if err != nil {
		errorHandler(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Adding nfs share: %+v", shareRequest)

	if err := nfs.Add(shareRequest); err != nil {
		errorHandler(http.StatusInternalServerError, err.Error())
		return
	}

	writeJsonResponse(&nfs.ShareResponse{
		Status:  http.StatusOK,
		Message: "SHARED ADDED",
	})
}
