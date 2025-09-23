# Logo Crawler - Enterprise Edition 🚀

A high-performance, enterprise-grade Go application that fetches logos for multiple publishers concurrently. Built with clean architecture principles, comprehensive error handling, and optimized for maximum performance.

## 🎯 Performance Improvements

### Before (Sequential)
- Publishers processed one by one
- Logo validation done sequentially
- No rate limiting
- Basic error handling

### After (Concurrent)
- **Concurrent publisher processing** with worker pool pattern
- **Concurrent logo validation** with semaphore-controlled concurrency
- **Enhanced error handling** with panic recovery
- **Performance metrics** and timing information
- **Configurable concurrency** via environment variables

## 🚀 Key Features

### Architecture & Design
- **Clean Architecture**: Separation of concerns with layered design
- **SOLID Principles**: Single responsibility, dependency injection
- **Enterprise Patterns**: Factory, Strategy, and Observer patterns
- **Modular Design**: Pluggable components for extensibility

### Performance & Concurrency
- **Worker Pool Pattern**: Configurable number of workers (default: CPU cores, max: 10)
- **Concurrent Logo Validation**: Up to 10 concurrent image dimension checks per publisher
- **Connection Pooling**: Optimized HTTP client with connection reuse
- **Context-based Cancellation**: Proper timeout and cancellation handling

### Reliability & Monitoring
- **Error Resilience**: Graceful error handling and panic recovery
- **Performance Metrics**: Detailed timing and success rate statistics
- **Comprehensive Logging**: Structured logging for debugging
- **Health Monitoring**: Built-in health checks and monitoring

## 📊 Expected Performance Gains

For a typical workload with 8 publishers:
- **Sequential**: ~40-60 seconds
- **Concurrent**: ~8-15 seconds
- **Speedup**: 3-5x faster

## 🛠️ Usage

### Environment Variables

```bash
# Required
export PUBLISHER_FILE_PATH="publishers.txt"
export CONFIG_FILE_PATH="config/config.yaml"

# Optional
export MAX_WORKERS="5"  # Default: CPU cores (max 10)
export HTML_OUTPUT_PATH="reports/logo-report.html"  # HTML report output path
```

### Running the Application

```bash
# Build
go build -o logo-crawler main.go

# Run
./logo-crawler
```

### Example Output

```
🚀 Starting concurrent logo crawler with 4 workers for 8 publishers
⚡ Using 8 CPU cores

📊 Results Summary:
⏱️  Total time: 12.5s
📈 Average time per publisher: 1.56s

🔎 Publisher: amazon.com (processed in 1.2s)
   https://logo.clearbit.com/amazon.com (512x512) <- ✅ BEST
   https://amazon.com/favicon.ico (32x32)

🔎 Publisher: google.com (processed in 0.8s)
   https://logo.clearbit.com/google.com (512x512) <- ✅ BEST
   https://google.com/favicon.ico (32x32)

📈 Final Stats:
   Total publishers: 8
   Publishers with logos: 7
   Publishers with errors: 1
   Total logos found: 15
   Success rate: 87.5%

📄 HTML report generated: reports/logo-crawler-report-2024-09-23-18-50-00.html
```

### HTML Report Features

The application generates a beautiful HTML report with:
- **Interactive Dashboard**: Statistics cards with key metrics
- **Visual Logo Display**: Thumbnail previews of all found logos
- **Best Logo Highlighting**: Clear indication of the best logo for each publisher
- **Responsive Design**: Works on desktop and mobile devices
- **Error Handling**: Clear display of any processing errors
- **Performance Metrics**: Detailed timing and success rate information

## 🔧 Configuration

### config.yaml
```yaml
preferred:
  min_width: 120
  min_height: 120
```

### publishers.txt
```
amazon.com
google.com
hotstar.com
jiohotstar.com
pubmatic.com
deepintent.com
honeywell.com
phonepe.com
```

## 🏗️ Architecture

### Clean Architecture Layers

1. **Application Layer**: `main.go` - Entry point and application orchestration
2. **Business Logic Layer**: `internal/app/` - Core application logic and coordination
3. **Service Layer**: `internal/crawler/` - Domain-specific services and components
4. **Infrastructure Layer**: `internal/io/`, `internal/utils/`, `config/` - External concerns

### Key Components

#### Core Services
- **LogoCrawler**: Main orchestrator coordinating all components
- **LogoExtractor**: Handles logo candidate extraction from multiple sources
- **LogoValidator**: Manages concurrent logo validation with semaphore control
- **DomainProcessor**: Domain detection and normalization
- **BestLogoSelector**: Logo quality assessment and selection

#### Concurrency Patterns
1. **Worker Pool Pattern**: For processing multiple publishers
2. **Semaphore Pattern**: For limiting concurrent logo validations
3. **Context-based Cancellation**: For timeout handling
4. **Channel-based Communication**: For result collection

### Design Patterns Used
- **Factory Pattern**: Component creation and initialization
- **Strategy Pattern**: Different logo extraction strategies
- **Observer Pattern**: Result collection and notification
- **Dependency Injection**: Loose coupling between components

## 🎛️ Tuning Parameters

- **MAX_WORKERS**: Number of concurrent publisher processors
- **Semaphore Size**: Currently set to 10 concurrent logo validations
- **HTTP Timeout**: 8 seconds per request
- **Validation Timeout**: 30 seconds per publisher

## 🔍 Monitoring

The application provides detailed metrics:
- Total processing time
- Per-publisher processing time
- Success/failure rates
- Total logos found
- Error details

## 🚨 Error Handling

- Graceful panic recovery in worker goroutines
- Context-based timeout handling
- Detailed error reporting
- Continues processing even if individual publishers fail

## 📈 Performance Tips

1. **Adjust MAX_WORKERS**: Start with CPU cores, increase if I/O bound
2. **Network Conditions**: Performance varies with network latency
3. **Target Server Load**: Some servers may be slower than others

## 🔄 Migration from Sequential Version

The concurrent version is backward compatible:
- Same configuration files
- Same input format
- Same output format
- Just much faster! ⚡
