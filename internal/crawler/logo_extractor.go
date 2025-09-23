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

	// Always try web scraping first to get more options
	htmlCandidates := le.extractFromHTML(baseURL)
	candidates = append(candidates, htmlCandidates...)

	// Always add common fallbacks
	candidates = append(candidates, le.getCommonFallbacks(domain)...)

	// Add Clearbit as a fallback (but not primary)
	candidates = append(candidates, le.getClearbitLogo(domain))

	return le.unique(candidates)
}

// extractFromHTML extracts logo candidates from HTML meta tags and links
func (le *LogoExtractor) extractFromHTML(baseURL string) []string {
	var allCandidates []string

	// Try multiple URL variations to get more logos
	urls := []string{baseURL}

	// Add www version if not already present
	if !strings.HasPrefix(baseURL, "https://www.") {
		wwwURL := strings.Replace(baseURL, "https://", "https://www.", 1)
		urls = append(urls, wwwURL)
	}

	// Try each URL variation
	for _, url := range urls {
		candidates := le.extractFromSingleURL(url)
		allCandidates = append(allCandidates, candidates...)
	}

	return le.unique(allCandidates)
}

// extractFromSingleURL extracts logos from a single URL
func (le *LogoExtractor) extractFromSingleURL(baseURL string) []string {
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

	// Extract from img tags with logo-related attributes
	candidates = append(candidates, le.extractImgTags(doc, base)...)

	return candidates
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

// extractImgTags extracts logo URLs from img tags with logo-related attributes
func (le *LogoExtractor) extractImgTags(doc *goquery.Document, base *url.URL) []string {
	var candidates []string
	domain := base.Hostname()

	// Look for img tags with logo-related attributes
	doc.Find("img").Each(func(i int, sel *goquery.Selection) {
		src, exists := sel.Attr("src")
		if !exists {
			return
		}

		// Check various attributes that might indicate a logo
		alt, _ := sel.Attr("alt")
		class, _ := sel.Attr("class")
		id, _ := sel.Attr("id")

		// Combine all attributes for checking
		combined := strings.ToLower(alt + " " + class + " " + id)

		// Skip if it's clearly not a domain logo
		if le.isUnrelatedLogo(combined, src, domain) {
			return
		}

		// Check if this looks like a domain logo
		if le.isDomainLogo(combined, src, domain) {
			candidates = append(candidates, le.resolveURL(base, src))
		}
	})

	return candidates
}

// isDomainLogo checks if the image is likely a domain-specific logo
func (le *LogoExtractor) isDomainLogo(combined, src, domain string) bool {
	// Check for domain-specific logo keywords
	domainLogoKeywords := []string{
		"logo", "brand", "header", "nav", "site-icon", "company",
		"main-logo", "brand-logo", "header-logo", "navigation-logo",
		"site-logo", "corporate-logo", "primary-logo",
	}

	hasLogoKeyword := false
	for _, keyword := range domainLogoKeywords {
		if strings.Contains(combined, keyword) {
			hasLogoKeyword = true
			break
		}
	}

	// If it has logo keywords, it's likely a domain logo
	if hasLogoKeyword {
		return true
	}

	// Check if it's in a logo-related path
	logoPaths := []string{
		"/logo", "/brand", "/assets/logo", "/images/logo", "/static/logo",
		"/img/logo", "/media/logo", "/uploads/logo",
	}

	for _, path := range logoPaths {
		if strings.Contains(src, path) {
			return true
		}
	}

	return false
}

// isUnrelatedLogo checks if the image is clearly not a domain logo
func (le *LogoExtractor) isUnrelatedLogo(combined, src, domain string) bool {
	// Partner/third-party logos
	partnerKeywords := []string{
		"partner", "sponsor", "advertisement", "ad", "banner",
		"pci", "dss", "iso", "certified", "award", "badge",
		"credit-card", "visa", "mastercard", "amex", "rupay",
		"bank", "payment", "security", "ssl", "trust",
		"social", "facebook", "twitter", "instagram", "linkedin",
		"youtube", "google", "apple", "microsoft", "amazon",
		"hero", "banner", "cover", "background", "splash",
		"testimonial", "review", "rating", "star",
	}

	for _, keyword := range partnerKeywords {
		if strings.Contains(combined, keyword) {
			return true
		}
	}

	// Check for partner/third-party domains in src
	thirdPartyDomains := []string{
		"cdn.", "assets.", "static.", "media.", "uploads.",
		"partner", "sponsor", "ad", "banner", "social",
	}

	for _, thirdParty := range thirdPartyDomains {
		if strings.Contains(src, thirdParty) && !strings.Contains(src, domain) {
			// If it's from a third-party domain and not the main domain
			return true
		}
	}

	// Check for advertisement-related paths
	adPaths := []string{
		"/ads/", "/advertisements/", "/banners/", "/promotions/",
		"/partners/", "/sponsors/", "/certifications/", "/awards/",
		"/testimonials/", "/reviews/", "/ratings/",
	}

	for _, path := range adPaths {
		if strings.Contains(src, path) {
			return true
		}
	}

	// Check for social media sharing images
	if strings.Contains(combined, "share") || strings.Contains(combined, "social") {
		return true
	}

	// Check for very large images (likely banners/hero images)
	// This will be handled by the best logo selector, but we can pre-filter obvious ones
	if strings.Contains(combined, "hero") || strings.Contains(combined, "banner") {
		return true
	}

	return false
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
