package output

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/Tanmay-Thanvi/logo-crawler/internal/crawler"
)

// HTMLGenerator handles HTML report generation
type HTMLGenerator struct {
	outputPath string
}

// NewHTMLGenerator creates a new HTML generator
func NewHTMLGenerator(outputPath string) *HTMLGenerator {
	return &HTMLGenerator{
		outputPath: outputPath,
	}
}

// HTMLReport represents the data structure for HTML report
type HTMLReport struct {
	Title           string
	GeneratedAt     time.Time
	TotalPublishers int
	ValidPublishers int
	ErrorCount      int
	TotalLogos      int
	SuccessRate     float64
	TotalDuration   time.Duration
	AvgDuration     time.Duration
	Results         []crawler.PublisherResult
}

// GenerateReport generates an HTML report from the results
func (hg *HTMLGenerator) GenerateReport(results []crawler.PublisherResult, totalDuration time.Duration) error {
	stats := hg.calculateStats(results)

	report := HTMLReport{
		Title:           "Logo Crawler Report",
		GeneratedAt:     time.Now(),
		TotalPublishers: stats.TotalPublishers,
		ValidPublishers: stats.ValidPublishers,
		ErrorCount:      stats.ErrorCount,
		TotalLogos:      stats.TotalLogos,
		SuccessRate:     stats.SuccessRate,
		TotalDuration:   totalDuration,
		AvgDuration:     totalDuration / time.Duration(stats.TotalPublishers),
		Results:         results,
	}

	tmpl := hg.getHTMLTemplate()

	// Create output directory if it doesn't exist
	dir := filepath.Dir(hg.outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(hg.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create HTML file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, report); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
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
func (hg *HTMLGenerator) calculateStats(results []crawler.PublisherResult) Stats {
	stats := Stats{TotalPublishers: len(results)}

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

// getHTMLTemplate returns the HTML template
func (hg *HTMLGenerator) getHTMLTemplate() *template.Template {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 2.5em;
            font-weight: 300;
        }
        .header p {
            margin: 10px 0 0 0;
            opacity: 0.9;
        }
        .stats {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 20px;
            padding: 30px;
            background: #f8f9fa;
        }
        .stat-card {
            background: white;
            padding: 20px;
            border-radius: 8px;
            text-align: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .stat-number {
            font-size: 2em;
            font-weight: bold;
            color: #667eea;
            margin-bottom: 5px;
        }
        .stat-number.error {
            color: #d32f2f;
        }
        .stat-label {
            color: #666;
            font-size: 0.9em;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .results {
            padding: 30px;
        }
        .publisher {
            margin-bottom: 30px;
            border: 1px solid #e0e0e0;
            border-radius: 8px;
            overflow: hidden;
        }
        .publisher-header {
            background: #f8f9fa;
            padding: 15px 20px;
            border-bottom: 1px solid #e0e0e0;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .publisher-name {
            font-weight: bold;
            font-size: 1.1em;
            color: #333;
        }
        .publisher-duration {
            color: #666;
            font-size: 0.9em;
        }
        .publisher-error {
            background: #ffe6e6;
            color: #d32f2f;
            padding: 15px 20px;
        }
        .logos {
            padding: 20px;
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
            gap: 15px;
        }
        .logo-card {
            background: #f8f9fa;
            border-radius: 8px;
            border-left: 4px solid #ddd;
            overflow: hidden;
            transition: transform 0.2s ease, box-shadow 0.2s ease;
        }
        .logo-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        }
        .logo-card.best {
            border-left-color: #4caf50;
            background: #f1f8e9;
        }
        .logo-image-container {
            height: 120px;
            display: flex;
            align-items: center;
            justify-content: center;
            background: white;
            position: relative;
        }
        .logo-image {
            max-width: 100px;
            max-height: 100px;
            object-fit: contain;
            border-radius: 4px;
            transition: opacity 0.3s ease;
        }
        .logo-image.loading {
            opacity: 0.5;
        }
        .logo-image.error {
            display: none;
        }
        .logo-placeholder {
            display: none;
            color: #999;
            font-size: 0.8em;
            text-align: center;
        }
        .logo-placeholder.show {
            display: block;
        }
        .logo-info {
            padding: 12px;
        }
        .logo-url {
            font-size: 0.8em;
            color: #333;
            word-break: break-all;
            text-decoration: none;
            display: block;
            margin-bottom: 4px;
        }
        .logo-url:hover {
            color: #667eea;
            text-decoration: underline;
        }
        .logo-dimensions {
            font-size: 0.75em;
            color: #666;
            margin-bottom: 8px;
        }
        .best-badge {
            background: #4caf50;
            color: white;
            padding: 3px 8px;
            border-radius: 12px;
            font-size: 0.7em;
            font-weight: bold;
            display: inline-block;
        }
        .no-logos {
            text-align: center;
            color: #666;
            padding: 20px;
            font-style: italic;
        }
        .footer {
            background: #f8f9fa;
            padding: 20px;
            text-align: center;
            color: #666;
            border-top: 1px solid #e0e0e0;
        }
        
        /* Responsive Design */
        @media (max-width: 768px) {
            .stats {
                grid-template-columns: repeat(2, 1fr);
            }
            .logos {
                grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
            }
            .logo-image-container {
                height: 100px;
            }
            .logo-image {
                max-width: 80px;
                max-height: 80px;
            }
        }
        
        @media (max-width: 480px) {
            .stats {
                grid-template-columns: 1fr;
            }
            .logos {
                grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
            }
            .container {
                margin: 10px;
                border-radius: 4px;
            }
            .header {
                padding: 20px;
            }
            .header h1 {
                font-size: 2em;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 {{.Title}}</h1>
            <p>Generated on {{.GeneratedAt.Format "January 2, 2006 at 3:04 PM"}}</p>
        </div>
        
        <div class="stats">
            <div class="stat-card">
                <div class="stat-number">{{.TotalPublishers}}</div>
                <div class="stat-label">Total Publishers</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{.ValidPublishers}}</div>
                <div class="stat-label">With Logos</div>
            </div>
            <div class="stat-card">
                <div class="stat-number error">{{.ErrorCount}}</div>
                <div class="stat-label">Errors</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{.TotalLogos}}</div>
                <div class="stat-label">Total Logos</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{printf "%.1f" .SuccessRate}}%</div>
                <div class="stat-label">Success Rate</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{printf "%.3f" .TotalDuration.Seconds}}s</div>
                <div class="stat-label">Total Time</div>
            </div>
        </div>
        
        <div class="results">
            <h2>📊 Results</h2>
            {{range .Results}}
            <div class="publisher">
                {{if .Error}}
                <div class="publisher-error">
                    <strong>❌ {{.Publisher}}</strong> ({{.Duration}}) - ERROR: {{.Error}}
                </div>
                {{else}}
                <div class="publisher-header">
                    <div class="publisher-name">🔎 {{.Publisher}}</div>
                    <div class="publisher-duration">Processed in {{.Duration}}</div>
                </div>
                <div class="logos">
                    {{if .Logos}}
                        {{$bestURL := .Best.URL}}
                        {{range .Logos}}
                        <div class="logo-card {{if eq .URL $bestURL}}best{{end}}">
                            <div class="logo-image-container">
                                <img src="{{.URL}}" alt="Logo" class="logo-image" 
                                     onerror="this.classList.add('error'); this.nextElementSibling.classList.add('show');"
                                     onload="this.classList.remove('loading'); this.nextElementSibling.classList.remove('show');"
                                     onloadstart="this.classList.add('loading');">
                                <div class="logo-placeholder">
                                    🖼️ Image not available<br>
                                    <small>Click link to view</small>
                                </div>
                            </div>
                            <div class="logo-info">
                                <a href="{{.URL}}" target="_blank" class="logo-url">{{.URL}}</a>
                                <div class="logo-dimensions">{{.Width}}x{{.Height}} pixels</div>
                                {{if eq .URL $bestURL}}
                                <span class="best-badge">✅ BEST</span>
                                {{end}}
                            </div>
                        </div>
                        {{end}}
                    {{else}}
                    <div class="no-logos">❌ No valid logos found</div>
                    {{end}}
                </div>
                {{end}}
            </div>
            {{end}}
        </div>
        
        <div class="footer">
            <p>Generated by Logo Crawler - Enterprise Edition</p>
            <p><small>Note: Some images may not display due to CORS restrictions, but all links are clickable to view the logos directly.</small></p>
        </div>
    </div>
</body>
</html>`

	return template.Must(template.New("report").Parse(tmpl))
}
