package main

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

// --- MARKDOWN RENDERER ---

func renderMarkdown(md string) template.HTML {
	if md == "" {
		return ""
	}
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
		if len(sm) < 3 {
			return m
		}
		level := len(sm[1])
		return fmt.Sprintf("<h%d>%s</h%d>", level, sm[2], level)
	})
	// Unordered lists
	reUL := regexp.MustCompile(`(?m)^\s*[-*]\s+(.+)$`)
	text = reUL.ReplaceAllString(text, "<!--ULITEM-->$1")
	reWrapUL := regexp.MustCompile(`(?:<!--ULITEM-->[^\n]*\n?)+`)
	text = reWrapUL.ReplaceAllStringFunc(text, func(m string) string {
		items := strings.TrimSpace(strings.ReplaceAll(m, "<!--ULITEM-->", ""))
		return "<ul><li>" + strings.ReplaceAll(items, "\n", "</li><li>") + "</li></ul>"
	})
	// Ordered lists
	reOL := regexp.MustCompile(`(?m)^\s*\d+\.\s+(.+)$`)
	text = reOL.ReplaceAllString(text, "<!--OLITEM-->$1")
	reWrapOL := regexp.MustCompile(`(?:<!--OLITEM-->[^\n]*\n?)+`)
	text = reWrapOL.ReplaceAllStringFunc(text, func(m string) string {
		items := strings.TrimSpace(strings.ReplaceAll(m, "<!--OLITEM-->", ""))
		return "<ol><li>" + strings.ReplaceAll(items, "\n", "</li><li>") + "</li></ol>"
	})
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
	if !strings.HasPrefix(text, "<h") && !strings.HasPrefix(text, "<ul") && !strings.HasPrefix(text, "<ol") && !strings.HasPrefix(text, "<pre") && !strings.HasPrefix(text, "<li") {
		para := strings.Split(text, "\n\n")
		for i := range para {
			p := strings.TrimSpace(para[i])
			if p == "" {
				continue
			}
			if strings.HasPrefix(p, "<") {
				continue
			}
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
		if line == "" {
			continue
		}
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
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 2 {
			item := MediaItem{Type: strings.TrimSpace(parts[0]), URL: strings.TrimSpace(parts[1])}
			if len(parts) >= 3 {
				item.Alt = strings.TrimSpace(parts[2])
			}
			items = append(items, item)
		}
	}
	return items
}
