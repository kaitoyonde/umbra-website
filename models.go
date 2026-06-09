package main

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

type pageData struct {
	ImgVer       string
	WAPhone      string
	WAPhoneSuara string
}

func pageDataAll() pageData {
	return pageData{ImgVer: imgVer, WAPhone: waPhone, WAPhoneSuara: waPhoneSuara}
}


