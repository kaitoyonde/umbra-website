package main

import (
	"net/http"
	"strings"
)

// --- PUBLIC HANDLERS ---

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	render(w, "index.html", pageDataAll())
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	render(w, "register.html", pageDataAll())
}

func handleInquiry(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}

func handleProject(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/project/")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	proj := getProjectByID(id)
	if proj == nil {
		http.NotFound(w, r)
		return
	}
	render(w, "project.html", struct {
		pageData
		Project ProjectData
	}{pageDataAll(), *proj})
}

func registerPublicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/register", handleRegister)
	mux.HandleFunc("/do_inquiry", handleInquiry)
	mux.HandleFunc("/project/", handleProject)

	// Division pages
	divSlugs := map[string]string{
		"/UmbraCreativeSoftworks": "softworks",
		"/UmbraSuara":             "suara",
		"/Penumbra":               "penumbra",
	}
	for path, slug := range divSlugs {
		path := path
		slug := slug
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != path {
				http.NotFound(w, r)
				return
			}
			div := getDivision(slug)
			if div == nil {
				http.NotFound(w, r)
				return
			}
			render(w, "division.html", struct {
				pageData
				Division DivisionData
			}{pageDataAll(), *div})
		})
		mux.HandleFunc(path+"/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, path, http.StatusMovedPermanently)
		})
	}
}
