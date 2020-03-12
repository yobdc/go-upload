package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type ServeMux struct {
}

func (server *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path == "/upload" {
		switch r.Method {
		case "GET":
			t, err := template.ParseFiles("./views/upload.ctpl")
			if err != nil {
				log.Fatal(err)
			}
			t.Execute(w, nil)
		case "POST":
			r.ParseMultipartForm(32 << 20)
			file, header, err := r.FormFile("uploadfile")
			if err != nil {
				log.Fatal("form file err: ", err)
				return
			}
			defer file.Close()
			f, err := os.OpenFile("./files/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal("open file err: ", err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
		}
	} else if r.URL.Path == "/upload/files" {
		had := http.StripPrefix("/upload/files", http.FileServer(http.Dir("files")))
		had.ServeHTTP(w, r)
	}
}

func main() {
	server := &ServeMux{}
	port := "9090"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	err := http.ListenAndServe(":"+port, server)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
