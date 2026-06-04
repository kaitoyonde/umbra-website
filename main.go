package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var imgVer = fmt.Sprint(time.Now().Unix())

type pageData struct {
	ImgVer string
}

func pageDataAll() pageData {
	return pageData{ImgVer: imgVer}
}

func render(w http.ResponseWriter, page string, data interface{}) {
	tmpl := template.New("base.html").Funcs(template.FuncMap{
		"imgByID": func(id string) ImageData {
			if id == "" { return ImageData{} }
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
	tmpl.Execute(w, data)
}

func renderAdmin(w http.ResponseWriter, page string, data interface{}) {
	tmpl := template.New("admin-base.html").Funcs(template.FuncMap{
		"imgByID": func(id string) ImageData {
			if id == "" { return ImageData{} }
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
	tmpl, err = tmpl.ParseFiles("templates/admin-base.html", "templates/admin-"+page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

// --- DATA MODEL ---

type ProjectData struct {
	ID              string      `json:"id"`
	Year            string      `json:"year"`
	Name            string      `json:"name"`
	Subtitle        string      `json:"subtitle"`
	Client          string      `json:"client"`
	ImageURL        string      `json:"image_url"`
	ImageID         string      `json:"image_id"`
	DescriptionMD   string      `json:"description_md"`
	TechDescription string      `json:"tech_description"`
	TechTable       []TechRow   `json:"tech_table"`
	Gallery         []MediaItem `json:"gallery"`
}

type TechRow struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type MediaItem struct {
	Type string `json:"type"` // "image" or "video"
	URL  string `json:"url"`
	Alt  string `json:"alt"`
}

type PortfolioColumn struct {
	Title    string        `json:"title"`
	Projects []ProjectData `json:"projects"`
}

type DivisionData struct {
	Slug             string            `json:"slug"`
	Name             string            `json:"name"`
	Route            string            `json:"route"`
	BannerTitle      string            `json:"banner_title"`
	BannerSubtitle   string            `json:"banner_subtitle"`
	BannerImage      string            `json:"banner_image"`
	BannerImageID    string            `json:"banner_image_id"`
	Description      string            `json:"description"`
	PortfolioColumns []PortfolioColumn `json:"portfolio_columns"`
	Skills           []string          `json:"skills"`
}

type SiteData struct {
	Divisions []DivisionData `json:"divisions"`
	Images    []ImageData    `json:"images"`
	nextID    int
}

type ImageData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Alt         string `json:"alt"`
	Description string `json:"description"`
	FileName    string `json:"filename"`
	URL         string `json:"url"`
	UploadedAt  string `json:"uploaded_at"`
}

// --- DATA STORE ---

var (
	siteData   *SiteData
	dataMu     sync.RWMutex
	dataPath   string
)

func loadData() error {
	dataPath = filepath.Join("data", "data.json")
	b, err := os.ReadFile(dataPath)
	if err != nil {
		return err
	}
	siteData = &SiteData{}
	if err := json.Unmarshal(b, siteData); err != nil {
		return err
	}
	siteData.nextID = findMaxID(siteData) + 1
	return nil
}

func findMaxID(sd *SiteData) int {
	max := 0
	for _, d := range sd.Divisions {
		for _, c := range d.PortfolioColumns {
			for _, p := range c.Projects {
				var n int
				fmt.Sscanf(p.ID, "%d", &n)
				if n > max {
					max = n
				}
			}
		}
	}
	for _, img := range sd.Images {
		var n int
		fmt.Sscanf(img.ID, "%d", &n)
		if n > max {
			max = n
		}
	}
	return max
}

func saveData() error {
	dataMu.RLock()
	b, err := json.MarshalIndent(siteData, "", "  ")
	dataMu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(dataPath, b, 0644)
}

func getDivision(slug string) *DivisionData {
	dataMu.RLock()
	defer dataMu.RUnlock()
	return getDivisionLocked(slug)
}

func getDivisionLocked(slug string) *DivisionData {
	for i := range siteData.Divisions {
		if siteData.Divisions[i].Slug == slug {
			return &siteData.Divisions[i]
		}
	}
	return nil
}

func getProject(id string) (int, int, int) {
	dataMu.RLock()
	defer dataMu.RUnlock()
	return getProjectLocked(id)
}

func getProjectLocked(id string) (int, int, int) {
	for di, d := range siteData.Divisions {
		for ci, c := range d.PortfolioColumns {
			for pi, p := range c.Projects {
				if p.ID == id {
					return di, ci, pi
				}
			}
		}
	}
	return -1, -1, -1
}

func getAllImages() []ImageData {
	dataMu.RLock()
	defer dataMu.RUnlock()
	out := make([]ImageData, len(siteData.Images))
	copy(out, siteData.Images)
	return out
}

func getProjectByID(id string) *ProjectData {
	dataMu.RLock()
	defer dataMu.RUnlock()
	for _, d := range siteData.Divisions {
		for _, c := range d.PortfolioColumns {
			for i := range c.Projects {
				if c.Projects[i].ID == id {
					return &c.Projects[i]
				}
			}
		}
	}
	return nil
}

func renderMarkdown(md string) template.HTML {
	if md == "" { return "" }
	text := md
	// Escape literal HTML
	text = template.HTMLEscapeString(text)
	// Code blocks
	reCodeBlock := regexp.MustCompile("(?s)```(\\w*)\n(.*?)```\n?")
	text = reCodeBlock.ReplaceAllString(text, "<pre><code>$2</code></pre>\n")
	// Headers
	reH := regexp.MustCompile("(?m)^(#{1,6})\\s+(.+)$")
	text = reH.ReplaceAllStringFunc(text, func(m string) string {
		sm := reH.FindStringSubmatch(m)
		if len(sm) < 3 { return m }
		level := len(sm[1])
		return fmt.Sprintf("<h%d>%s</h%d>", level, sm[2], level)
	})
	// Unordered lists
	reUL := regexp.MustCompile(`(?m)^\s*[-*]\s+(.+)$`)
	text = reUL.ReplaceAllString(text, "<li>$1</li>")
	// Ordered lists
	reOL := regexp.MustCompile(`(?m)^\s*\d+\.\s+(.+)$`)
	text = reOL.ReplaceAllString(text, "<li>$1</li>")
	text = strings.ReplaceAll(text, "<li>", "</ul><ul><li>") + "</ul>"
	text = strings.Replace(text, "</ul><ul><li>", "<li>", 1)
	text = strings.TrimSuffix(text, "</ul>")
	// Bold
	reBold := regexp.MustCompile(`\*\*(.+?)\*\*`)
	text = reBold.ReplaceAllString(text, "<strong>$1</strong>")
	// Italic
	reItalic := regexp.MustCompile(`\*(.+?)\*`)
	text = reItalic.ReplaceAllString(text, "<em>$1</em>")
	// Inline code
	reCode := regexp.MustCompile("`(.+?)`")
	text = reCode.ReplaceAllString(text, "<code>$1</code>")
	// Links
	reLink := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = reLink.ReplaceAllString(text, `<a href="$2">$1</a>`)
	// Paragraphs
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "<h") && !strings.HasPrefix(text, "<ul") && !strings.HasPrefix(text, "<pre") && !strings.HasPrefix(text, "<li") {
		para := strings.Split(text, "\n\n")
		for i := range para {
			p := strings.TrimSpace(para[i])
			if p == "" { continue }
			if strings.HasPrefix(p, "<") { continue }
			para[i] = "<p>" + strings.ReplaceAll(p, "\n", "<br/>") + "</p>"
		}
		text = strings.Join(para, "\n")
	}
	return template.HTML(text)
}

func parseTechTable(s string) []TechRow {
	var rows []TechRow
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" { continue }
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			rows = append(rows, TechRow{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])})
		}
	}
	return rows
}

func parseGallery(s string) []MediaItem {
	var items []MediaItem
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" { continue }
		parts := strings.Split(line, "|")
		if len(parts) >= 2 {
			item := MediaItem{Type: strings.TrimSpace(parts[0]), URL: strings.TrimSpace(parts[1])}
			if len(parts) >= 3 { item.Alt = strings.TrimSpace(parts[2]) }
			items = append(items, item)
		}
	}
	return items
}

// --- ADMIN SESSION ---

func adminSession(r *http.Request) (string, bool) {
	c, err := r.Cookie("admin")
	if err != nil || c.Value == "" {
		return "", false
	}
	return c.Value, true
}

// --- ROUTES ---

func main() {
	if err := loadData(); err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}
	os.MkdirAll("static/uploads", 0755)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" { http.NotFound(w, r); return }
		render(w, "index.html", pageDataAll())
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		success := r.URL.Query().Get("success") == "true"
		render(w, "register.html", struct {
			pageData
			Success bool
		}{pageDataAll(), success})
	})

	mux.HandleFunc("/do_inquiry", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/register?success=true", http.StatusSeeOther)
	})

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
			if r.URL.Path != path { http.NotFound(w, r); return }
			div := getDivision(slug)
			if div == nil { http.NotFound(w, r); return }
			render(w, "division.html", struct {
				pageData
				Division DivisionData
			}{pageDataAll(), *div})
		})
		mux.HandleFunc(path+"/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, path, http.StatusMovedPermanently)
		})
	}

	// Project detail page
	mux.HandleFunc("/project/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/project/")
		if id == "" { http.NotFound(w, r); return }
		proj := getProjectByID(id)
		if proj == nil { http.NotFound(w, r); return }
		render(w, "project.html", struct {
			pageData
			Project ProjectData
		}{pageDataAll(), *proj})
	})

	// Admin routes
	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()
			if r.FormValue("password") == "umbra2024" {
				http.SetCookie(w, &http.Cookie{
					Name: "admin", Value: "authenticated", Path: "/",
					HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: 86400,
				})
				http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
				return
			}
			renderAdmin(w, "login.html", struct{ Error string }{"Password salah."})
			return
		}
		renderAdmin(w, "login.html", struct{ Error string }{""})
	})

	mux.HandleFunc("/admin/dashboard", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := adminSession(r); !ok { http.Redirect(w, r, "/admin", http.StatusSeeOther); return }
		dataMu.RLock()
		d := *siteData
		dataMu.RUnlock()
		renderAdmin(w, "dashboard.html", struct {
			Site SiteData
		}{d})
	})

	mux.HandleFunc("/admin/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "admin", Value: "", Path: "/", HttpOnly: true, MaxAge: -1})
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	})

	mux.HandleFunc("/admin/division/", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := adminSession(r); !ok { http.Redirect(w, r, "/admin", http.StatusSeeOther); return }
		slug := strings.TrimPrefix(r.URL.Path, "/admin/division/")
		if slug == "" { http.NotFound(w, r); return }

		if r.Method == "POST" {
			r.ParseForm()
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
		if div == nil { http.NotFound(w, r); return }
		renderAdmin(w, "division.html", struct {
			Division DivisionData
			Images   []ImageData
		}{*div, getAllImages()})
	})

	mux.HandleFunc("/admin/project/add/", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := adminSession(r); !ok { http.Redirect(w, r, "/admin", http.StatusSeeOther); return }
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/admin/project/add/"), "/")
		if len(parts) != 2 { http.NotFound(w, r); return }
		slug, colIndex := parts[0], parts[1]

		if r.Method == "POST" {
			r.ParseForm()
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

		renderAdmin(w, "project-form.html", struct {
			DivisionSlug string
			ColIndex     string
			Project      ProjectData
			Images       []ImageData
			IsNew        bool
		}{slug, colIndex, ProjectData{}, getAllImages(), true})
	})

	mux.HandleFunc("/admin/project/edit/", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := adminSession(r); !ok { http.Redirect(w, r, "/admin", http.StatusSeeOther); return }
		id := strings.TrimPrefix(r.URL.Path, "/admin/project/edit/")
		if id == "" { http.NotFound(w, r); return }

		if r.Method == "POST" {
			r.ParseForm()
			dataMu.Lock()
			di, ci, pi := getProjectLocked(id)
			if di >= 0 {
				proj := &siteData.Divisions[di].PortfolioColumns[ci].Projects[pi]
				if y := r.FormValue("year"); y != "" { proj.Year = y }
				if n := r.FormValue("name"); n != "" { proj.Name = n }
				proj.Subtitle = r.FormValue("subtitle")
				if c := r.FormValue("client"); c != "" { proj.Client = c }
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
				if u := r.FormValue("image_url"); u != "" { proj.ImageURL = u }
			}
			dataMu.Unlock()
			saveData()
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}

		di, ci, pi := getProject(id)
		if di < 0 { http.NotFound(w, r); return }
		dataMu.RLock()
		proj := siteData.Divisions[di].PortfolioColumns[ci].Projects[pi]
		dataMu.RUnlock()
		renderAdmin(w, "project-form.html", struct {
			DivisionSlug string
			ColIndex     string
			Project      ProjectData
			Images       []ImageData
			IsNew        bool
		}{siteData.Divisions[di].Slug, fmt.Sprint(ci), proj, getAllImages(), false})
	})

	mux.HandleFunc("/admin/project/delete/", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := adminSession(r); !ok { http.Redirect(w, r, "/admin", http.StatusSeeOther); return }
		id := strings.TrimPrefix(r.URL.Path, "/admin/project/delete/")
		if id == "" { http.NotFound(w, r); return }

		dataMu.Lock()
		di, ci, pi := getProjectLocked(id)
		if di >= 0 {
			col := &siteData.Divisions[di].PortfolioColumns[ci]
			col.Projects = append(col.Projects[:pi], col.Projects[pi+1:]...)
		}
		dataMu.Unlock()
		saveData()
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	})

	// Image manager
	mux.HandleFunc("/admin/images", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := adminSession(r); !ok { http.Redirect(w, r, "/admin", http.StatusSeeOther); return }

		if r.Method == "POST" {
			r.ParseMultipartForm(10 << 20)
			file, header, err := r.FormFile("file")
			if err != nil {
				renderAdmin(w, "images.html", struct {
					Images []ImageData
					Error  string
				}{getAllImages(), "Gagal membaca file: " + err.Error()})
				return
			}
			defer file.Close()

			ext := filepath.Ext(header.Filename)
			timestamp := fmt.Sprint(time.Now().UnixNano())
			filename := timestamp + ext

			dst, err := os.Create(filepath.Join("static", "uploads", filename))
			if err != nil {
				renderAdmin(w, "images.html", struct {
					Images []ImageData
					Error  string
				}{getAllImages(), "Gagal menyimpan file."})
				return
			}
			defer dst.Close()
			io.Copy(dst, file)

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

		renderAdmin(w, "images.html", struct {
			Images []ImageData
			Error  string
		}{getAllImages(), ""})
	})

	mux.HandleFunc("/admin/images/delete/", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := adminSession(r); !ok { http.Redirect(w, r, "/admin", http.StatusSeeOther); return }
		id := strings.TrimPrefix(r.URL.Path, "/admin/images/delete/")
		if id == "" { http.NotFound(w, r); return }

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
	})

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		http.FileServer(http.Dir("static")).ServeHTTP(w, r)
	})))

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
