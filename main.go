package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var imgVer = fmt.Sprint(time.Now().Unix())

type pageData struct {
	ImgVer string
}

func pageDataAll() pageData {
	return pageData{ImgVer: imgVer}
}

func render(w http.ResponseWriter, page string) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/"+page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, pageDataAll())
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" { http.NotFound(w, r); return }
		render(w, "index.html")
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		render(w, "register.html")
	})

	subPages := map[string]string{
		"/UmbraCreativeSoftworks": "softworks.html",
		"/UmbraSuara":             "suara.html",
		"/Penumbra":               "penumbra.html",
	}
	for path, page := range subPages {
		path := path
		page := page
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != path { http.NotFound(w, r); return }
			render(w, page)
		})
		mux.HandleFunc(path+"/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, path, http.StatusMovedPermanently)
		})
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		http.FileServer(http.Dir("static")).ServeHTTP(w, r)
	})))

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
