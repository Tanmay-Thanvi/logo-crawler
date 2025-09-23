package crawler

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Tanmay-Thanvi/logo-crawler/internal/utils"
)

// LogoExtractor handles logo extraction from various sources
type LogoExtractor struct{}

// NewLogoExtractor creates a new logo extractor
func NewLogoExtractor() *LogoExtractor {
	return &LogoExtractor{}
}

// ExtractCandidates extracts logo candidates from HTML and common paths
func (le *LogoExtractor) ExtractCandidates(domain string) []string {
	baseURL := "https://" + domain

	var candidates []string
	candidates = append(candidates, le.extractFromHTML(baseURL)...)
	candidates = append(candidates, le.getCommonFallbacks(domain)...)
	candidates = append(candidates, le.getClearbitLogo(domain))

	return le.unique(candidates)
}

// extractFromHTML extracts logo candidates from HTML meta tags and links
func (le *LogoExtractor) extractFromHTML(baseURL string) []string {
	resp, err := utils.Client.Get(baseURL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil
	}

	var candidates []string
	base := resp.Request.URL

	// Extract from meta tags
	candidates = append(candidates, le.extractMetaTags(doc, base)...)

	// Extract from link tags
	candidates = append(candidates, le.extractLinkTags(doc, base)...)

	return le.unique(candidates)
}

// extractMetaTags extracts logo URLs from meta tags
func (le *LogoExtractor) extractMetaTags(doc *goquery.Document, base *url.URL) []string {
	var candidates []string
	metaProps := []string{"og:image", "twitter:image", "og:image:url"}

	for _, prop := range metaProps {
		// Check property attribute
		if content, exists := doc.Find("meta[property='" + prop + "']").Attr("content"); exists {
			candidates = append(candidates, le.resolveURL(base, content))
		}
		// Check name attribute
		if content, exists := doc.Find("meta[name='" + prop + "']").Attr("content"); exists {
			candidates = append(candidates, le.resolveURL(base, content))
		}
	}

	return candidates
}

// extractLinkTags extracts logo URLs from link tags
func (le *LogoExtractor) extractLinkTags(doc *goquery.Document, base *url.URL) []string {
	var candidates []string

	doc.Find("link[rel]").Each(func(i int, sel *goquery.Selection) {
		rel, _ := sel.Attr("rel")
		href, _ := sel.Attr("href")
		if strings.Contains(strings.ToLower(rel), "icon") && href != "" {
			candidates = append(candidates, le.resolveURL(base, href))
		}
	})

	return candidates
}

// getCommonFallbacks returns common logo/icon paths for a domain
func (le *LogoExtractor) getCommonFallbacks(domain string) []string {
	base := "https://" + domain
	paths := []string{
		"/favicon.ico",
		"/favicon.png",
		"/favicon.svg",
		"/apple-touch-icon.png",
		"/apple-touch-icon-precomposed.png",
		"/logo.png",
		"/assets/logo.png",
		"/images/logo.png",
	}

	var urls []string
	for _, path := range paths {
		urls = append(urls, base+path)
	}
	return urls
}

// getClearbitLogo returns the Clearbit logo API URL for the domain
func (le *LogoExtractor) getClearbitLogo(domain string) string {
	return "https://logo.clearbit.com/" + domain
}

// resolveURL resolves a relative URL against a base URL
func (le *LogoExtractor) resolveURL(base *url.URL, href string) string {
	u, err := base.Parse(href)
	if err != nil {
		return href
	}
	return u.String()
}

// unique removes duplicate URLs from the slice
func (le *LogoExtractor) unique(list []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, v := range list {
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
