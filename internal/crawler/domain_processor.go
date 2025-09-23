package crawler

import (
	"regexp"
	"strings"

	"github.com/Tanmay-Thanvi/logo-crawler/config"
)

// DomainProcessor handles domain detection and normalization
type DomainProcessor struct{}

// NewDomainProcessor creates a new domain processor
func NewDomainProcessor() *DomainProcessor {
	return &DomainProcessor{}
}

// DetectDomain normalizes input into a domain string
func (dp *DomainProcessor) DetectDomain(input string) string {
	domainRe := regexp.MustCompile(`[\w\.-]+\.[a-z]{2,}`)
	if domainRe.MatchString(input) {
		return domainRe.FindString(input)
	}
	// Fallback: assume it's a name -> append .com
	return strings.ToLower(strings.ReplaceAll(input, " ", "")) + ".com"
}

// BestLogoSelector selects the best logo based on preferences
type BestLogoSelector struct{}

// NewBestLogoSelector creates a new best logo selector
func NewBestLogoSelector() *BestLogoSelector {
	return &BestLogoSelector{}
}

// SelectBest selects the best logo using intelligent scoring
func (bls *BestLogoSelector) SelectBest(logos []LogoInfo, prefs config.Preferences) *LogoInfo {
	if len(logos) == 0 {
		return nil
	}

	var best *LogoInfo
	bestScore := -1

	for _, logo := range logos {
		score := bls.calculateLogoScore(logo, prefs)
		if score > bestScore {
			bestScore = score
			best = &logo
		}
	}

	return best
}

// calculateLogoScore calculates an intelligent score for logo selection
func (bls *BestLogoSelector) calculateLogoScore(logo LogoInfo, prefs config.Preferences) int {
	score := 0
	url := strings.ToLower(logo.URL)

	// Base score for meeting minimum requirements
	if logo.Width >= prefs.Preferred.MinWidth && logo.Height >= prefs.Preferred.MinHeight {
		score += 10
	} else {
		// Penalty for not meeting minimum requirements
		score -= 20
	}

	// Bonus for Clearbit logos (usually high quality)
	if strings.Contains(url, "logo.clearbit.com") {
		score += 15
	}

	// Bonus for favicon.ico (official icon)
	if strings.Contains(url, "favicon.ico") {
		score += 12
	}

	// Bonus for apple-touch-icon (high quality)
	if strings.Contains(url, "apple-touch-icon") {
		score += 10
	}

	// Bonus for SVG logos (scalable)
	if strings.Contains(url, ".svg") {
		score += 8
	}

	// Penalty for dashboard/cover images (usually large)
	if bls.isDashboardImage(logo, url) {
		score -= 30
	}

	// Penalty for social media images (og:image, twitter:image)
	if bls.isSocialMediaImage(url) {
		score -= 25
	}

	// Penalty for partner/third-party logos
	if bls.isPartnerLogo(url) {
		score -= 40
	}

	// Penalty for advertisement/promotional content
	if bls.isAdvertisement(url) {
		score -= 35
	}

	// Bonus for square logos (better for branding)
	if logo.Width == logo.Height {
		score += 5
	}

	// Bonus for reasonable aspect ratio (not too wide/tall)
	aspectRatio := float64(logo.Width) / float64(logo.Height)
	if aspectRatio >= 0.5 && aspectRatio <= 2.0 {
		score += 3
	}

	// Size-based scoring (prefer medium-sized logos)
	area := logo.Width * logo.Height
	if area >= 10000 && area <= 100000 { // 100x100 to 316x316 pixels
		score += 8
	} else if area >= 1000 && area < 10000 { // 32x32 to 100x100 pixels
		score += 5
	} else if area > 100000 { // Very large images
		score -= 10
	}

	// Bonus for PNG format (good quality)
	if strings.Contains(url, ".png") {
		score += 3
	}

	// Penalty for very small images
	if logo.Width < 32 || logo.Height < 32 {
		score -= 15
	}

	return score
}

// isDashboardImage checks if the logo is likely a dashboard/cover image
func (bls *BestLogoSelector) isDashboardImage(logo LogoInfo, url string) bool {
	// Very large images are likely dashboard/cover images
	if logo.Width > 800 || logo.Height > 600 {
		return true
	}

	// Check for dashboard-related keywords in URL
	dashboardKeywords := []string{
		"dashboard", "cover", "hero", "banner", "header-bg",
		"background", "splash", "landing", "homepage", "og-image",
		"social", "twitter", "facebook", "linkedin",
	}

	for _, keyword := range dashboardKeywords {
		if strings.Contains(url, keyword) {
			return true
		}
	}

	return false
}

// isSocialMediaImage checks if the logo is a social media image
func (bls *BestLogoSelector) isSocialMediaImage(url string) bool {
	socialKeywords := []string{
		"og-image", "twitter-image", "facebook-image", "social-image",
		"meta-image", "share-image", "preview-image",
	}

	for _, keyword := range socialKeywords {
		if strings.Contains(url, keyword) {
			return true
		}
	}

	return false
}

// isPartnerLogo checks if the logo is from a partner/third-party
func (bls *BestLogoSelector) isPartnerLogo(url string) bool {
	partnerKeywords := []string{
		"pci", "dss", "iso", "certified", "award", "badge",
		"credit-card", "visa", "mastercard", "amex", "rupay",
		"bank", "payment", "security", "ssl", "trust",
		"partner", "sponsor", "collaboration", "alliance",
	}

	for _, keyword := range partnerKeywords {
		if strings.Contains(url, keyword) {
			return true
		}
	}

	return false
}

// isAdvertisement checks if the logo is an advertisement/promotional content
func (bls *BestLogoSelector) isAdvertisement(url string) bool {
	adKeywords := []string{
		"advertisement", "ad", "promotion", "banner", "campaign",
		"offer", "deal", "discount", "sale", "limited-time",
		"testimonial", "review", "rating", "feedback",
		"hero", "cover", "background", "splash",
	}

	for _, keyword := range adKeywords {
		if strings.Contains(url, keyword) {
			return true
		}
	}

	return false
}
