package app

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/Tanmay-Thanvi/logo-crawler/config"
	"github.com/Tanmay-Thanvi/logo-crawler/internal/crawler"
	"github.com/Tanmay-Thanvi/logo-crawler/internal/io"
	"github.com/Tanmay-Thanvi/logo-crawler/internal/output"
	"github.com/Tanmay-Thanvi/logo-crawler/internal/utils"
	"github.com/joho/godotenv"
)

// LogoCrawlerApp represents the main application
type LogoCrawlerApp struct {
	config     *AppConfig
	prefs      config.Preferences
	publishers []string
}

// AppConfig holds application configuration
type AppConfig struct {
	PublisherFilePath string
	ConfigFilePath    string
	MaxWorkers        int
	HTMLOutputPath    string
}

// NewLogoCrawlerApp creates a new application instance
func NewLogoCrawlerApp() *LogoCrawlerApp {
	return &LogoCrawlerApp{}
}

// Run executes the main application logic
func (app *LogoCrawlerApp) Run() {
	app.loadEnvironment()
	app.loadConfiguration()
	app.loadPublishers()
	app.displayStartupInfo()

	results, totalDuration := app.processPublishers()
	app.displayResults(results)
	app.generateHTMLReport(results, totalDuration)
}

// loadEnvironment loads environment variables and .env file
func (app *LogoCrawlerApp) loadEnvironment() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	app.config = &AppConfig{
		PublisherFilePath: os.Getenv("PUBLISHER_FILE_PATH"),
		ConfigFilePath:    os.Getenv("CONFIG_FILE_PATH"),
		MaxWorkers:        app.getMaxWorkers(),
		HTMLOutputPath:    app.getHTMLOutputPath(),
	}

	app.validateConfig()
}

// validateConfig validates required configuration
func (app *LogoCrawlerApp) validateConfig() {
	if app.config.PublisherFilePath == "" {
		log.Fatal("❌ Missing PUBLISHER_FILE_PATH env variable")
	}
	if app.config.ConfigFilePath == "" {
		log.Fatal("❌ Missing CONFIG_FILE_PATH env variable")
	}
}

// loadConfiguration loads the YAML configuration
func (app *LogoCrawlerApp) loadConfiguration() {
	app.prefs = config.LoadConfig(app.config.ConfigFilePath)
}

// loadPublishers reads publishers from file
func (app *LogoCrawlerApp) loadPublishers() {
	loader := utils.NewLoader("Reading publishers from file...")
	loader.Start()

	var err error
	app.publishers, err = io.ReadPublishers(app.config.PublisherFilePath)

	loader.Stop()

	if err != nil {
		log.Fatalf("Failed to read publishers: %v", err)
	}
	if len(app.publishers) == 0 {
		log.Fatal("❌ No publishers found in file")
	}

	fmt.Printf("✅ Loaded %d publishers\n", len(app.publishers))
}

// displayStartupInfo shows startup information
func (app *LogoCrawlerApp) displayStartupInfo() {
	fmt.Printf("🚀 Starting concurrent logo crawler with %d workers for %d publishers\n",
		app.config.MaxWorkers, len(app.publishers))
	fmt.Printf("⚡ Using %d CPU cores\n", runtime.NumCPU())
}

// processPublishers processes all publishers concurrently
func (app *LogoCrawlerApp) processPublishers() ([]crawler.PublisherResult, time.Duration) {
	fmt.Println("\n🔄 Starting logo crawling process...")

	// Create progress bar for overall progress
	progressBar := utils.NewProgressBar(len(app.publishers), "Processing publishers")

	start := time.Now()
	results := crawler.FetchPublishersConcurrently(app.publishers, app.prefs, app.config.MaxWorkers)
	totalDuration := time.Since(start)

	progressBar.Complete()

	fmt.Printf("\n📊 Results Summary:\n")
	fmt.Printf("⏱️  Total time: %v\n", totalDuration)
	fmt.Printf("📈 Average time per publisher: %v\n", totalDuration/time.Duration(len(app.publishers)))

	return results, totalDuration
}

// displayResults displays the processing results
func (app *LogoCrawlerApp) displayResults(results []crawler.PublisherResult) {
	stats := app.calculateStats(results)

	for _, result := range results {
		app.displayPublisherResult(result)
	}

	app.displayFinalStats(stats)
}

// displayPublisherResult displays result for a single publisher
func (app *LogoCrawlerApp) displayPublisherResult(result crawler.PublisherResult) {
	if result.Error != nil {
		fmt.Printf("\n❌ Publisher: %s (processed in %v) - ERROR: %v\n",
			result.Publisher, result.Duration, result.Error)
		return
	}

	fmt.Printf("\n🔎 Publisher: %s (processed in %v)\n", result.Publisher, result.Duration)
	if len(result.Logos) == 0 {
		fmt.Println("❌ No valid logos found")
		return
	}

	for _, logo := range result.Logos {
		mark := ""
		if result.Best != nil && logo.URL == result.Best.URL {
			mark = " <- ✅ BEST"
		}
		fmt.Printf("   %s (%dx%d)%s\n", logo.URL, logo.Width, logo.Height, mark)
	}
}

// Stats holds processing statistics
type Stats struct {
	TotalPublishers int
	ValidPublishers int
	ErrorCount      int
	TotalLogos      int
	SuccessRate     float64
}

// calculateStats calculates processing statistics
func (app *LogoCrawlerApp) calculateStats(results []crawler.PublisherResult) Stats {
	stats := Stats{TotalPublishers: len(app.publishers)}

	for _, result := range results {
		if result.Error != nil {
			stats.ErrorCount++
			continue
		}

		stats.TotalLogos += len(result.Logos)
		if len(result.Logos) > 0 {
			stats.ValidPublishers++
		}
	}

	stats.SuccessRate = float64(stats.ValidPublishers) / float64(stats.TotalPublishers) * 100
	return stats
}

// displayFinalStats displays final processing statistics
func (app *LogoCrawlerApp) displayFinalStats(stats Stats) {
	fmt.Printf("\n📈 Final Stats:\n")
	fmt.Printf("   Total publishers: %d\n", stats.TotalPublishers)
	fmt.Printf("   Publishers with logos: %d\n", stats.ValidPublishers)
	fmt.Printf("   Publishers with errors: %d\n", stats.ErrorCount)
	fmt.Printf("   Total logos found: %d\n", stats.TotalLogos)
	fmt.Printf("   Success rate: %.1f%%\n", stats.SuccessRate)
}

// generateHTMLReport generates an HTML report
func (app *LogoCrawlerApp) generateHTMLReport(results []crawler.PublisherResult, totalDuration time.Duration) {
	if app.config.HTMLOutputPath == "" {
		return // Skip HTML generation if no output path specified
	}

	loader := utils.NewLoader("Generating HTML report...")
	loader.Start()

	generator := output.NewHTMLGenerator(app.config.HTMLOutputPath)
	if err := generator.GenerateReport(results, totalDuration); err != nil {
		loader.Stop()
		log.Printf("⚠️ Failed to generate HTML report: %v", err)
		return
	}

	loader.Stop()
	fmt.Printf("📄 HTML report generated: %s\n", app.config.HTMLOutputPath)

	// Open the report in the default browser
	if err := utils.OpenHTMLFile(app.config.HTMLOutputPath); err != nil {
		log.Printf("⚠️ Failed to open browser: %v", err)
		fmt.Printf("💡 You can manually open the report at: %s\n", app.config.HTMLOutputPath)
	} else {
		fmt.Printf("🌐 Opening report in default browser...\n")
	}
}

// getMaxWorkers determines the optimal number of workers
func (app *LogoCrawlerApp) getMaxWorkers() int {
	if maxWorkersStr := os.Getenv("MAX_WORKERS"); maxWorkersStr != "" {
		if maxWorkers, err := strconv.Atoi(maxWorkersStr); err == nil && maxWorkers > 0 {
			return maxWorkers
		}
	}

	// Default to number of CPU cores, but cap at 10
	maxWorkers := runtime.NumCPU()
	if maxWorkers > 10 {
		maxWorkers = 10
	}
	return maxWorkers
}

// getHTMLOutputPath gets the HTML output path from environment or uses default
func (app *LogoCrawlerApp) getHTMLOutputPath() string {
	if path := os.Getenv("HTML_OUTPUT_PATH"); path != "" {
		return path
	}
	// Default to reports directory with timestamp
	return fmt.Sprintf("reports/logo-crawler-report-%s.html", time.Now().Format("2006-01-02-15-04-05"))
}
