package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Image string `json:"image,omitempty"`
	Shop  string `json:"shop,omitempty"`
	Price string `json:"price,omitempty"`
}

var (
	itemsLinkRe = regexp.MustCompile(`/items/\d+($|[?#])`)
)

func buildSearchURL(lang, query, sort string, page int) (string, error) {
	if lang == "" {
		lang = "ja"
	}
	if sort == "" {
		sort = "popular"
	}
	if page < 1 {
		page = 1
	}
	// Use url.PathEscape to properly escape the query in the URL path
	escapedQuery := url.PathEscape(query)
	base := fmt.Sprintf("https://booth.pm/%s/search/%s", lang, escapedQuery)
	q := url.Values{}
	q.Set("sort", sort)
	q.Set("order", "desc")
	q.Set("page", fmt.Sprintf("%d", page))
	return base + "?" + q.Encode(), nil
}

func fetchDocument(targetURL string) (*goquery.Document, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.9,en;q=0.8,zh-TW;q=0.7")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return goquery.NewDocumentFromReader(resp.Body)
}

func absolute(base, href string) string {
	u, err := url.Parse(href)
	if err != nil {
		return href
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	if u.Host == "" {
		u.Host = "booth.pm"
	}
	if !strings.HasPrefix(u.Path, "/") {
		u.Path = "/" + u.Path
	}
	return u.String()
}

func extractItems(doc *goquery.Document) []Item {
	results := make([]Item, 0, 60)
	seen := make(map[string]struct{})

	// Strategy: find all anchors pointing to item detail pages and lift info from their surrounding card.
	doc.Find("a").Each(func(i int, a *goquery.Selection) {
		href, ok := a.Attr("href")
		if !ok || !itemsLinkRe.MatchString(href) {
			return
		}
		title := strings.TrimSpace(a.Text())
		// walk up to approximate card container - try different parent elements
		container := a.ParentsFiltered("li, article, div").First()
		if container.Length() == 0 {
			container = a.Parent()
		}
		// If still no container, use the document root for broader search
		if container.Length() == 0 {
			container = doc.Selection
		}

		// Image
		imgURL := ""
		if s := container.Find("img").First(); s.Length() > 0 {
			if v, ok := s.Attr("data-src"); ok {
				imgURL = v
			} else if v, ok := s.Attr("src"); ok {
				imgURL = v
			}
		}

		// Shop name
		shop := ""
		container.Find("a").Each(func(_ int, sa *goquery.Selection) {
			if sh, ok := sa.Attr("href"); ok && strings.Contains(sh, ".booth.pm") {
				shop = strings.TrimSpace(sa.Text())
			}
		})

		// Price (first text with currency symbol)
		price := ""
		container.Find("span,div,p").Each(func(_ int, sp *goquery.Selection) {
			if price != "" {
				return
			}
			text := strings.TrimSpace(sp.Text())
			if strings.Contains(text, "¥") || strings.Contains(strings.ToUpper(text), "JPY") {
				price = text
			}
		})

		abs := absolute("https://booth.pm", href)
		if _, exists := seen[abs]; exists {
			return
		}
		seen[abs] = struct{}{}

		if title == "" {
			// fallback to image alt or last segment of URL
			if s := container.Find("img[alt]").First(); s.Length() > 0 {
				title = strings.TrimSpace(s.AttrOr("alt", ""))
			}
			if title == "" {
				parts := strings.Split(strings.Trim(abs, "/"), "/")
				title = parts[len(parts)-1]
			}
		}

		results = append(results, Item{Title: title, URL: abs, Image: imgURL, Shop: shop, Price: price})
	})

	return results
}

func main() {
	query := flag.String("query", "vrchat 衣装", "Search keywords (e.g. 'vrchat 衣装' or 'vrchat 衣服')")
	sort := flag.String("sort", "popular", "Sort: popular or new")
	page := flag.Int("page", 1, "Page number (>=1)")
	lang := flag.String("lang", "ja", "Language: ja or zh-tw")
	flag.Parse()

	u, err := buildSearchURL(*lang, *query, *sort, *page)
	if err != nil {
		log.Fatalf("build url: %v", err)
	}

	doc, err := fetchDocument(u)
	if err != nil {
		log.Fatalf("fetch: %v", err)
	}

	items := extractItems(doc)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(items); err != nil {
		log.Fatalf("encode: %v", err)
	}
}
