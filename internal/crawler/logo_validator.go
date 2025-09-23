package crawler

import (
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"sync"
	"time"

	"github.com/Tanmay-Thanvi/logo-crawler/internal/utils"
)

// LogoValidator handles concurrent logo validation
type LogoValidator struct {
	semaphore chan struct{}
}

// NewLogoValidator creates a new logo validator
func NewLogoValidator(maxConcurrent int) *LogoValidator {
	return &LogoValidator{
		semaphore: make(chan struct{}, maxConcurrent),
	}
}

// ValidateConcurrently validates multiple logo URLs concurrently
func (lv *LogoValidator) ValidateConcurrently(candidates []string) []LogoInfo {
	if len(candidates) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results := make(chan LogoInfo, len(candidates))
	var wg sync.WaitGroup

	for _, url := range candidates {
		wg.Add(1)
		go lv.validateSingleLogo(ctx, url, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var valid []LogoInfo
	for logo := range results {
		valid = append(valid, logo)
	}

	return valid
}

// validateSingleLogo validates a single logo URL
func (lv *LogoValidator) validateSingleLogo(ctx context.Context, url string, results chan<- LogoInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	select {
	case lv.semaphore <- struct{}{}:
		defer func() { <-lv.semaphore }()
	case <-ctx.Done():
		return
	}

	width, height := lv.getImageDimensionsWithContext(ctx, url)
	if width > 0 && height > 0 {
		results <- LogoInfo{
			URL:    url,
			Width:  width,
			Height: height,
			Valid:  true,
		}
	}
}

// getImageDimensionsWithContext gets image dimensions with context
func (lv *LogoValidator) getImageDimensionsWithContext(ctx context.Context, url string) (int, int) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, 0
	}

	resp, err := utils.Client.Do(req)
	if err != nil {
		return 0, 0
	}
	defer resp.Body.Close()

	img, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		return 0, 0
	}
	return img.Width, img.Height
}
