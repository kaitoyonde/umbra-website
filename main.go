package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	loadEnv(".env")

	adminPass = os.Getenv("ADMIN_PASSWORD")
	if adminPass == "" {
		log.Fatal("ADMIN_PASSWORD environment variable is required")
	}

	waPhone = os.Getenv("WA_PHONE_NUMBER")
	if waPhone == "" {
		waPhone = "6281234567890"
		log.Println("WA_PHONE_NUMBER not set, using placeholder 6281234567890")
	}

	waPhoneSuara = os.Getenv("WA_PHONE_SUARA")
	if waPhoneSuara == "" {
		waPhoneSuara = waPhone
		log.Println("WA_PHONE_SUARA not set, falling back to WA_PHONE_NUMBER")
	}

	if err := loadData(); err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}
	os.MkdirAll("static/uploads", 0755)

	mux := http.NewServeMux()

	// Public routes
	registerPublicRoutes(mux)

	// Admin routes
	registerAdminRoutes(mux)

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		http.FileServer(http.Dir("static")).ServeHTTP(w, r)
	})))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
