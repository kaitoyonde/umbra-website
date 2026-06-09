package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// --- CONFIGURATION ---

var (
	imgVer       = fmt.Sprint(time.Now().Unix())
	adminPass    string
	waPhone      string
	waPhoneSuara string
)

func loadEnv(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return // .env file is optional
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		if len(val) >= 2 && (val[0] == '"' && val[len(val)-1] == '"' || val[0] == '\'' && val[len(val)-1] == '\'') {
			val = val[1 : len(val)-1]
		}
		if key != "" {
			os.Setenv(key, val)
		}
	}
}
