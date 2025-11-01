package main

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestBuildSearchURL(t *testing.T) {
	u, err := buildSearchURL("ja", "vrchat 衣装", "popular", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got, want := parsed.Host, "booth.pm"; got != want {
		t.Fatalf("host = %q, want %q", got, want)
	}
	if !strings.Contains(parsed.Path, "/ja/search/") {
		t.Fatalf("path should contain /ja/search/, got %q", parsed.Path)
	}
	// Check if the original URL contains the escaped query
	if !strings.Contains(u, url.PathEscape("vrchat 衣装")) {
		t.Fatalf("URL should contain escaped query, got %q", u)
	}
	// Alternatively, check if the parsed path contains the unescaped query
	if !strings.Contains(parsed.Path, "vrchat 衣装") {
		t.Fatalf("path should contain unescaped query, got %q", parsed.Path)
	}
	q := parsed.Query()
	if got, want := q.Get("sort"), "popular"; got != want {
		t.Fatalf("sort = %q, want %q", got, want)
	}
	if got, want := q.Get("order"), "desc"; got != want {
		t.Fatalf("order = %q, want %q", got, want)
	}
	if got, want := q.Get("page"), "1"; got != want {
		t.Fatalf("page = %q, want %q", got, want)
	}
}

func TestExtractItems_Simple(t *testing.T) {
	html := `
<!doctype html><html><body>
<ul>
  <li class="item-card">
    <a class="item-card__link" href="/ja/items/7587937">【しなの-Shinano-】Cat Hoodie</a>
    <img src="https://example.com/cat.jpg" alt="Cat Hoodie">
    <a href="https://micare.booth.pm/">Micare Sewing</a>
    <span>¥ 1,200</span>
  </li>
</ul>
</body></html>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("doc: %v", err)
	}
	items := extractItems(doc)
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	it := items[0]
	if it.Title == "" || !strings.Contains(it.Title, "Cat Hoodie") {
		t.Fatalf("unexpected title: %q", it.Title)
	}
	if !strings.HasPrefix(it.URL, "https://booth.pm/ja/items/") {
		t.Fatalf("unexpected url: %q", it.URL)
	}
	if it.Image == "" {
		t.Fatalf("image should not be empty")
	}
	if it.Shop == "" {
		t.Fatalf("shop should not be empty")
	}
	if it.Price == "" {
		t.Fatalf("price should not be empty")
	}

	// Ensure JSON encoding works as expected
	if _, err := json.Marshal(it); err != nil {
		t.Fatalf("json: %v", err)
	}
}

func TestExtractItems_Dedup(t *testing.T) {
	html := `
<!doctype html><html><body>
<div>
  <div>
    <a href="/ja/items/7600000">Item A</a>
    <a href="/ja/items/7600000">Item A (duplicate link)</a>
    <img src="https://example.com/a.jpg">
    <span>¥ 100</span>
  </div>
</div>
</body></html>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("doc: %v", err)
	}
	items := extractItems(doc)
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1 (dedup)", len(items))
	}
}
