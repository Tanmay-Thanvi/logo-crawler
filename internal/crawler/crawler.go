package crawler

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Tanmay-Thanvi/logo-crawler/config"
)

type LogoInfo struct {
	URL    string
	Width  int
	Height int
	Valid  bool
}

type PublisherResult struct {
	Publisher string
	Logos     []LogoInfo
	Best      *LogoInfo
	Error     error
	Duration  time.Duration
	Index     int // To preserve input order
}

// LogoCrawler orchestrates the logo crawling process
type LogoCrawler struct {
	extractor *LogoExtractor
	validator *LogoValidator
	processor *DomainProcessor
	selector  *BestLogoSelector
}

// NewLogoCrawler creates a new logo crawler
func NewLogoCrawler() *LogoCrawler {
	return &LogoCrawler{
		extractor: NewLogoExtractor(),
		validator: NewLogoValidator(10), // Max 10 concurrent validations
		processor: NewDomainProcessor(),
		selector:  NewBestLogoSelector(),
	}
}

// FetchPublisherLogos returns all valid logos and the best one
func (lc *LogoCrawler) FetchPublisherLogos(input string, prefs config.Preferences) ([]LogoInfo, *LogoInfo) {
	domain := lc.processor.DetectDomain(input)

	// Step 1: Extract candidates
	candidates := lc.extractor.ExtractCandidates(domain)

	// Step 2: Validate candidates concurrently
	valid := lc.validator.ValidateConcurrently(candidates)

	// Step 3: Select best logo
	best := lc.selector.SelectBest(valid, prefs)

	// Step 4: Sort logos with best logo first
	sortedLogos := lc.sortLogosWithBestFirst(valid, best)

	return sortedLogos, best
}

// FetchPublisherLogos is the public interface for backward compatibility
func FetchPublisherLogos(input string, prefs config.Preferences) ([]LogoInfo, *LogoInfo) {
	crawler := NewLogoCrawler()
	return crawler.FetchPublisherLogos(input, prefs)
}

// FetchPublishersConcurrently processes multiple publishers concurrently
func FetchPublishersConcurrently(publishers []string, prefs config.Preferences, maxWorkers int) []PublisherResult {
	if len(publishers) == 0 {
		return nil
	}

	// Create channels for work distribution
	type publisherTask struct {
		publisher string
		index     int
	}

	publisherChan := make(chan publisherTask, len(publishers))
	resultChan := make(chan PublisherResult, len(publishers))

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range publisherChan {
				start := time.Now()
				logos, best := FetchPublisherLogos(task.publisher, prefs)
				duration := time.Since(start)

				result := PublisherResult{
					Publisher: task.publisher,
					Logos:     logos,
					Best:      best,
					Duration:  duration,
					Index:     task.index,
				}

				// Handle any panics gracefully
				defer func() {
					if r := recover(); r != nil {
						result.Error = fmt.Errorf("panic occurred: %v", r)
						resultChan <- result
					}
				}()

				resultChan <- result
			}
		}()
	}

	// Send publishers to workers with their original index
	go func() {
		defer close(publisherChan)
		for index, publisher := range publishers {
			publisherChan <- publisherTask{
				publisher: publisher,
				index:     index,
			}
		}
	}()

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var results []PublisherResult
	for result := range resultChan {
		results = append(results, result)
	}

	// Sort results by original index to preserve input order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Index < results[j].Index
	})

	return results
}

// sortLogosWithBestFirst sorts logos with the best logo at the beginning
func (lc *LogoCrawler) sortLogosWithBestFirst(logos []LogoInfo, best *LogoInfo) []LogoInfo {
	if best == nil || len(logos) <= 1 {
		return logos
	}

	var sortedLogos []LogoInfo
	var otherLogos []LogoInfo

	// Separate best logo from others
	for _, logo := range logos {
		if logo.URL == best.URL {
			sortedLogos = append(sortedLogos, logo)
		} else {
			otherLogos = append(otherLogos, logo)
		}
	}

	// Add best logo first, then others
	sortedLogos = append(sortedLogos, otherLogos...)
	return sortedLogos
}
