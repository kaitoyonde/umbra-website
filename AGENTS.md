# UmbraSoftwork — Go MPA Company Profile

## Build & Run

```
go build -o ./tmp/main.exe .   # build
./tmp/main.exe                  # run on :8080
air                             # hot-reload (watches .go, .json, .html, .css)
```

No test files exist. No CI/CD.

## Architecture

Single-file Go app (`main.go`, ~700 lines, stdlib only). No external dependencies. No database — all content lives in `data/data.json` (read/written with `sync.RWMutex`). Custom regex-based markdown renderer (no goldmark/blackfriday).

Templates loaded from disk on every request — edit `.html` or `.css` and refresh, no restart needed.

## Key Directories

| Path | Purpose |
|------|---------|
| `main.go` | All routes, models, handlers, data access, markdown, uploads |
| `templates/` | 11 Go `html/template` files: `base.html`, `index.html`, `division.html`, `project.html`, `register.html`, `admin-*.html` |
| `static/` | `style.css`, `LOGO.svg`, `uploads/` (user-uploaded images) |
| `data/data.json` | Entire site content: divisions, projects, images, skills |
| `.agents/skills/` | OpenCode skill definitions (loaded automatically by the system prompt when a task matches their description) |

## Routes

| Path | Description |
|------|-------------|
| `/` | Home page |
| `/register` | Inquiry form |
| `/UmbraCreativeSoftworks` | Division page |
| `/UmbraSuara` | Division page |
| `/Penumbra` | Division page |
| `/project/{id}` | Project detail page |
| `/admin` | Admin login |
| `/admin/dashboard` | CMS dashboard |
| `/admin/division/{slug}` | Edit division projects |
| `/admin/images` | Image manager (upload/delete) |

## Gotchas

- **Admin password**: hardcoded `umbra2024` in `main.go:412`
- **Port**: hardcoded `:8080`, no env-var override
- **Image upload limit**: 10 MB (`r.ParseMultipartForm(10 << 20)`)
- **Static files**: served with `Cache-Control: no-cache, must-revalidate`
- **All UI text in Indonesian** — templates and data content are in `id` locale
- **No `.gitignore`** — `tmp/main.exe` is tracked in git
- **Deadlock caveat**: `sync.RWMutex` is not reentrant — never call `getDivision()`/`getProject()` while holding `dataMu.Lock()`; use `getDivisionLocked()`/`getProjectLocked()` instead
