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

// SelectBest selects the best logo based on resolution preferences
func (bls *BestLogoSelector) SelectBest(logos []LogoInfo, prefs config.Preferences) *LogoInfo {
	var best *LogoInfo
	for _, logo := range logos {
		if bls.meetsMinimumRequirements(logo, prefs) {
			if best == nil || bls.isBetterQuality(logo, *best) {
				best = &logo
			}
		}
	}
	return best
}

// meetsMinimumRequirements checks if logo meets minimum size requirements
func (bls *BestLogoSelector) meetsMinimumRequirements(logo LogoInfo, prefs config.Preferences) bool {
	return logo.Width >= prefs.Preferred.MinWidth && logo.Height >= prefs.Preferred.MinHeight
}

// isBetterQuality checks if logo1 is better quality than logo2
func (bls *BestLogoSelector) isBetterQuality(logo1, logo2 LogoInfo) bool {
	return logo1.Width*logo1.Height > logo2.Width*logo2.Height
}
