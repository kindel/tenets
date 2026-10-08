package tenets_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestPromptKeepsEssayAndNonEssayLinksApart(t *testing.T) {
	prompt, err := os.ReadFile("prompt.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(prompt)
	if !strings.Contains(text, "https://kindel.com/essays/tenets/") {
		t.Fatal("tenets essay should use kindel.com/essays")
	}
	if !strings.Contains(text, "https://kindel.com/essays/load-bearing-words/") {
		t.Fatal("load-bearing words is an essay")
	}
	const engineer = "https://blog.kindel.com/2026/07/25/principal-engineer-tenets-unless-you-know-better-ones/"
	if !strings.Contains(text, engineer) {
		t.Fatal("the principal engineer post is not an essay and stays on the blog")
	}
	assertNoEssayPermalinks(t)
}

func assertNoEssayPermalinks(t *testing.T) {
	t.Helper()
	raw, err := os.ReadFile("data/essay_slugs.json")
	if err != nil {
		t.Fatal(err)
	}
	var snap struct {
		BySlug map[string]string `json:"by_slug"`
		ByID   map[string]string `json:"by_id"`
	}
	if err := json.Unmarshal(raw, &snap); err != nil {
		t.Fatal(err)
	}
	dated := regexp.MustCompile(`https?://(?:www\.)?blog\.kindel\.com/\d{4}/\d{2}/\d{2}/([a-z0-9]+(?:-[a-z0-9]+)*)/?`)
	byP := regexp.MustCompile(`https?://(?:www\.)?blog\.kindel\.com/(?:index\.php)?\?[^"'\s>]*\bp=(\d+)`)
	skip := map[string]bool{".git": true, "essay_slugs.json": true, "check-essay-links.js": true}
	var hits []string
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if skip[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if skip[info.Name()] {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".html", ".md", ".js", ".json", ".go", ".css", ".txt":
		default:
			return nil
		}
		text, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		body := string(text)
		for _, match := range dated.FindAllStringSubmatch(body, -1) {
			if _, ok := snap.BySlug[strings.ToLower(match[1])]; ok {
				hits = append(hits, path+": "+match[0])
			}
		}
		for _, match := range byP.FindAllStringSubmatch(body, -1) {
			id := strings.TrimLeft(match[1], "0")
			if id == "" {
				id = "0"
			}
			if _, ok := snap.ByID[id]; ok {
				hits = append(hits, path+": "+match[0])
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) > 0 {
		t.Fatalf("essay permalinks still on the blog:\n%s", strings.Join(hits, "\n"))
	}
}
