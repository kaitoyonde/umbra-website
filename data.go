package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// --- DATA STORE ---

var (
	siteData *SiteData
	dataMu   sync.RWMutex
	dataPath string
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
