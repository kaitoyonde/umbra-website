package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// --- TEMPLATE RENDERING ---

func render(w http.ResponseWriter, page string, data interface{}) {
	tmpl := template.New("base.html").Funcs(template.FuncMap{
		"imgByID": func(id string) ImageData {
			if id == "" {
				return ImageData{}
			}
			dataMu.RLock()
			defer dataMu.RUnlock()
			for _, img := range siteData.Images {
				if img.ID == id {
					return img
				}
			}
			return ImageData{}
		},
		"markdown": renderMarkdown,
	})
	var err error
	tmpl, err = tmpl.ParseFiles("templates/base.html", "templates/"+page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("template error (%s): %v", page, err)
	}
}

func renderAdmin(w http.ResponseWriter, page string, data interface{}) {
	tmpl := template.New("admin-base.html").Funcs(template.FuncMap{
		"imgByID": func(id string) ImageData {
			if id == "" {
				return ImageData{}
			}
			dataMu.RLock()
			defer dataMu.RUnlock()
			for _, img := range siteData.Images {
				if img.ID == id {
					return img
				}
			}
			return ImageData{}
		},
		"markdown": renderMarkdown,
		"add": func(a, b int) int {
			return a + b
		},
		"json": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(string(b))
		},
	})
	var err error
	tmpl, err = tmpl.ParseFiles("templates/admin-base.html", "templates/admin-"+page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("admin template error (%s): %v", page, err)
	}
}
