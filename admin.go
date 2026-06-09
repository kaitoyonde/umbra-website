package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// --- ADMIN HANDLERS ---

func handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			renderAdmin(w, "login.html", struct{ Error string }{"Gagal membaca form."})
			return
		}
		if r.FormValue("password") == adminPass {
			b := make([]byte, 16)
			rand.Read(b)
			sessionID := hex.EncodeToString(b)
			sessionsMu.Lock()
			sessions[sessionID] = time.Now().Add(24 * time.Hour)
			sessionsMu.Unlock()
			http.SetCookie(w, &http.Cookie{
				Name: "admin", Value: sessionID, Path: "/",
				HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: 86400,
			})
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}
		renderAdmin(w, "login.html", struct{ Error string }{"Password salah."})
		return
	}
	renderAdmin(w, "login.html", struct{ Error string }{""})
}

func handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	if _, ok := adminSession(r); !ok {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	dataMu.RLock()
	d := *siteData
	dataMu.RUnlock()
	renderAdmin(w, "dashboard.html", struct {
		Site SiteData
	}{d})
}

func handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "admin", Value: "", Path: "/", HttpOnly: true, MaxAge: -1})
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func handleAdminDivision(w http.ResponseWriter, r *http.Request) {
	if _, ok := adminSession(r); !ok {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	slug := strings.TrimPrefix(r.URL.Path, "/admin/division/")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		dataMu.Lock()
		div := getDivisionLocked(slug)
		if div != nil {
			if desc := r.FormValue("description"); desc != "" {
				div.Description = desc
			}
			if bt := r.FormValue("banner_title"); bt != "" {
				div.BannerTitle = bt
			}
			if bs := r.FormValue("banner_subtitle"); bs != "" {
				div.BannerSubtitle = bs
			}
			if bi := r.FormValue("banner_image"); bi != "" {
				div.BannerImage = bi
			}
			if bid := r.FormValue("banner_image_id"); bid != "" {
				div.BannerImageID = bid
				for _, img := range siteData.Images {
					if img.ID == bid {
						div.BannerImage = img.URL
						break
					}
				}
			}
			if skills := r.FormValue("skills"); skills != "" {
				div.Skills = strings.Split(skills, "\n")
				for i := range div.Skills {
					div.Skills[i] = strings.TrimSpace(div.Skills[i])
				}
			}
			for ci := range div.PortfolioColumns {
				key := fmt.Sprintf("col_title_%d", ci)
				if t := r.FormValue(key); t != "" {
					div.PortfolioColumns[ci].Title = t
				}
			}
		}
		dataMu.Unlock()
		saveData()
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	div := getDivision(slug)
	if div == nil {
		http.NotFound(w, r)
		return
	}
	dataMu.RLock()
	site := *siteData
	dataMu.RUnlock()
	renderAdmin(w, "division.html", struct {
		Site     SiteData
		Division DivisionData
		Images   []ImageData
	}{site, *div, getAllImages()})
}

func handleAdminProjectAdd(w http.ResponseWriter, r *http.Request) {
	if _, ok := adminSession(r); !ok {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/admin/project/add/"), "/")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	slug, colIndex := parts[0], parts[1]

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		dataMu.Lock()
		div := getDivisionLocked(slug)
		ci := 0
		fmt.Sscanf(colIndex, "%d", &ci)
		if div != nil && ci >= 0 && ci < len(div.PortfolioColumns) {
			siteData.nextID++
			imageURL := r.FormValue("image_url")
			imageID := r.FormValue("image_id")
			if imageID != "" && imageURL == "" {
				for _, img := range siteData.Images {
					if img.ID == imageID {
						imageURL = img.URL
						break
					}
				}
			}
			p := ProjectData{
				ID:              fmt.Sprintf("%d", siteData.nextID),
				Year:            r.FormValue("year"),
				Name:            r.FormValue("name"),
				Subtitle:        r.FormValue("subtitle"),
				Client:          r.FormValue("client"),
				ImageURL:        imageURL,
				ImageID:         imageID,
				DescriptionMD:   r.FormValue("description_md"),
				TechDescription: r.FormValue("tech_description"),
				TechTable:       parseTechTable(r.FormValue("tech_table")),
				Gallery:         parseGallery(r.FormValue("gallery")),
			}
			div.PortfolioColumns[ci].Projects = append(div.PortfolioColumns[ci].Projects, p)
		}
		dataMu.Unlock()
		saveData()
		http.Redirect(w, r, "/admin/division/"+slug, http.StatusSeeOther)
		return
	}

	dataMu.RLock()
	site := *siteData
	dataMu.RUnlock()
	renderAdmin(w, "project-form.html", struct {
		Site         SiteData
		DivisionSlug string
		ColIndex     string
		Project      ProjectData
		Images       []ImageData
		IsNew        bool
	}{site, slug, colIndex, ProjectData{}, getAllImages(), true})
}

func handleAdminProjectEdit(w http.ResponseWriter, r *http.Request) {
	if _, ok := adminSession(r); !ok {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/admin/project/edit/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		dataMu.Lock()
		di, ci, pi := getProjectLocked(id)
		if di >= 0 {
			proj := &siteData.Divisions[di].PortfolioColumns[ci].Projects[pi]
			if y := r.FormValue("year"); y != "" {
				proj.Year = y
			}
			if n := r.FormValue("name"); n != "" {
				proj.Name = n
			}
			proj.Subtitle = r.FormValue("subtitle")
			if c := r.FormValue("client"); c != "" {
				proj.Client = c
			}
			proj.DescriptionMD = r.FormValue("description_md")
			proj.TechDescription = r.FormValue("tech_description")
			proj.TechTable = parseTechTable(r.FormValue("tech_table"))
			proj.Gallery = parseGallery(r.FormValue("gallery"))
			imageID := r.FormValue("image_id")
			if imageID != "" {
				proj.ImageID = imageID
				for _, img := range siteData.Images {
					if img.ID == imageID {
						proj.ImageURL = img.URL
						break
					}
				}
			}
			if u := r.FormValue("image_url"); u != "" {
				proj.ImageURL = u
			}
		}
		dataMu.Unlock()
		saveData()
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	di, ci, pi := getProject(id)
	if di < 0 {
		http.NotFound(w, r)
		return
	}
	dataMu.RLock()
	proj := siteData.Divisions[di].PortfolioColumns[ci].Projects[pi]
	slug := siteData.Divisions[di].Slug
	site := *siteData
	dataMu.RUnlock()
	renderAdmin(w, "project-form.html", struct {
		Site         SiteData
		DivisionSlug string
		ColIndex     string
		Project      ProjectData
		Images       []ImageData
		IsNew        bool
	}{site, slug, fmt.Sprint(ci), proj, getAllImages(), false})
}

func handleAdminProjectDelete(w http.ResponseWriter, r *http.Request) {
	if _, ok := adminSession(r); !ok {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/admin/project/delete/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	dataMu.Lock()
	di, ci, pi := getProjectLocked(id)
	if di >= 0 {
		col := &siteData.Divisions[di].PortfolioColumns[ci]
		col.Projects = append(col.Projects[:pi], col.Projects[pi+1:]...)
	}
	dataMu.Unlock()
	saveData()
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func handleAdminImages(w http.ResponseWriter, r *http.Request) {
	if _, ok := adminSession(r); !ok {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		r.ParseMultipartForm(10 << 20)
		file, header, err := r.FormFile("file")
		if err != nil {
			dataMu.RLock()
			site := *siteData
			dataMu.RUnlock()
			renderAdmin(w, "images.html", struct {
				Site   SiteData
				Images []ImageData
				Error  string
			}{site, getAllImages(), "Gagal membaca file: " + err.Error()})
			return
		}
		defer file.Close()

		ext := filepath.Ext(header.Filename)
		timestamp := fmt.Sprint(time.Now().UnixNano())
		filename := timestamp + ext

		dst, err := os.Create(filepath.Join("static", "uploads", filename))
		if err != nil {
			dataMu.RLock()
			site := *siteData
			dataMu.RUnlock()
			renderAdmin(w, "images.html", struct {
				Site   SiteData
				Images []ImageData
				Error  string
			}{site, getAllImages(), "Gagal menyimpan file."})
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			dataMu.RLock()
			site := *siteData
			dataMu.RUnlock()
			renderAdmin(w, "images.html", struct {
				Site   SiteData
				Images []ImageData
				Error  string
			}{site, getAllImages(), "Gagal menyimpan file: " + err.Error()})
			return
		}

		dataMu.Lock()
		siteData.nextID++
		img := ImageData{
			ID:          fmt.Sprintf("%d", siteData.nextID),
			Name:        r.FormValue("name"),
			Alt:         r.FormValue("alt"),
			Description: r.FormValue("description"),
			FileName:    filename,
			URL:         "/static/uploads/" + filename,
			UploadedAt:  time.Now().Format("2006-01-02 15:04"),
		}
		siteData.Images = append(siteData.Images, img)
		dataMu.Unlock()
		saveData()
		http.Redirect(w, r, "/admin/images", http.StatusSeeOther)
		return
	}

	dataMu.RLock()
	site := *siteData
	dataMu.RUnlock()
	renderAdmin(w, "images.html", struct {
		Site   SiteData
		Images []ImageData
		Error  string
	}{site, getAllImages(), ""})
}

func handleAdminImageDelete(w http.ResponseWriter, r *http.Request) {
	if _, ok := adminSession(r); !ok {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/admin/images/delete/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	dataMu.Lock()
	for i, img := range siteData.Images {
		if img.ID == id {
			os.Remove(filepath.Join("static", "uploads", img.FileName))
			siteData.Images = append(siteData.Images[:i], siteData.Images[i+1:]...)
			break
		}
	}
	dataMu.Unlock()
	saveData()
	http.Redirect(w, r, "/admin/images", http.StatusSeeOther)
}

func registerAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin", handleAdminLogin)
	mux.HandleFunc("/admin/dashboard", handleAdminDashboard)
	mux.HandleFunc("/admin/logout", handleAdminLogout)
	mux.HandleFunc("/admin/division/", handleAdminDivision)
	mux.HandleFunc("/admin/project/add/", handleAdminProjectAdd)
	mux.HandleFunc("/admin/project/edit/", handleAdminProjectEdit)
	mux.HandleFunc("/admin/project/delete/", handleAdminProjectDelete)
	mux.HandleFunc("/admin/images", handleAdminImages)
	mux.HandleFunc("/admin/images/delete/", handleAdminImageDelete)
}
