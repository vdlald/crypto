package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/upload", Upload).Methods("POST")
	r.HandleFunc("/", handler)
	http.Handle("/", r)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")
	file, handler, err := r.FormFile("file")

	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Disposition", "attachment; filename="+handler.Filename)
	// w.Header().Set("Content-Length", string(handler.Size))

	f, err := os.OpenFile("temp/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	io.Copy(f, file)
	f.Close()

	data := Crypt(r.FormValue("action"), "temp/"+handler.Filename, r.FormValue("password"), r.FormValue("salt"))
	w.Write(data)
	err = os.Remove("temp/" + handler.Filename)
	if err != nil {
		log.Println(err)
		return
	}
}
