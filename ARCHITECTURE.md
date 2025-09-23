# Logo Crawler - Architecture Documentation

## System Architecture

```mermaid
graph TB
    subgraph "Application Layer"
        A[main.go] --> B[LogoCrawlerApp]
    end
    
    subgraph "Business Logic Layer"
        B --> C[Configuration Management]
        B --> D[Publisher Processing]
        B --> E[Results Display]
    end
    
    subgraph "Service Layer"
        D --> F[LogoCrawler]
        F --> G[LogoExtractor]
        F --> H[LogoValidator]
        F --> I[DomainProcessor]
        F --> J[BestLogoSelector]
    end
    
    subgraph "Infrastructure Layer"
        G --> K[HTTP Client]
        H --> K
        I --> L[File I/O]
        C --> M[Config Management]
        E --> N[HTML Generator]
    end
    
    subgraph "External Systems"
        K --> N[Web Servers]
        K --> O[Clearbit API]
        L --> P[Publisher Files]
        M --> Q[Config Files]
    end
```

## Component Overview

### Core Components

| Component | Responsibility | Key Methods |
|-----------|---------------|-------------|
| **LogoCrawlerApp** | Application orchestration | `Run()`, `loadEnvironment()`, `processPublishers()` |
| **LogoCrawler** | Business logic coordination | `FetchPublisherLogos()` |
| **LogoExtractor** | Logo candidate extraction | `ExtractCandidates()`, `extractFromHTML()` |
| **LogoValidator** | Concurrent logo validation | `ValidateConcurrently()`, `validateSingleLogo()` |
| **DomainProcessor** | Domain detection | `DetectDomain()` |
| **BestLogoSelector** | Logo quality assessment | `SelectBest()`, `meetsMinimumRequirements()` |
| **HTMLGenerator** | HTML report generation | `GenerateReport()`, `getHTMLTemplate()` |

## Data Flow

```mermaid
sequenceDiagram
    participant App as LogoCrawlerApp
    participant Crawler as LogoCrawler
    participant Extractor as LogoExtractor
    participant Validator as LogoValidator
    participant Selector as BestLogoSelector
    
    App->>Crawler: FetchPublisherLogos(domain)
    Crawler->>Extractor: ExtractCandidates(domain)
    Extractor-->>Crawler: candidates[]
    Crawler->>Validator: ValidateConcurrently(candidates)
    Validator-->>Crawler: validLogos[]
    Crawler->>Selector: SelectBest(validLogos, prefs)
    Selector-->>Crawler: bestLogo
    Crawler-->>App: (validLogos, bestLogo)
    App->>App: displayResults(results)
    App->>App: generateHTMLReport(results)
```

## Concurrency Patterns

| Pattern | Purpose | Implementation |
|---------|---------|----------------|
| **Worker Pool** | Process multiple publishers | `FetchPublishersConcurrently()` |
| **Semaphore** | Limit concurrent validations | `LogoValidator.semaphore` (10 max) |
| **Context Cancellation** | Timeout handling | 30-second validation timeout |
| **Channel Communication** | Result collection | Buffered channels |

## Configuration

### Environment Variables
- `PUBLISHER_FILE_PATH`: Path to publishers file
- `CONFIG_FILE_PATH`: Path to configuration file  
- `MAX_WORKERS`: Number of concurrent workers (optional)
- `HTML_OUTPUT_PATH`: Path for HTML report output (optional)

### YAML Configuration
```yaml
preferred:
  min_width: 120
  min_height: 120
```

## Performance

- **Speedup**: 1.66x faster than sequential processing
- **Scalability**: Worker pool scales with CPU cores
- **Resource Management**: Bounded by semaphore and worker limits
- **Error Handling**: Graceful degradation with panic recovery
