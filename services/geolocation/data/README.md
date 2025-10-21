# TheraClosure Geolocation Service - Comprehensive World Data Seeding

This directory contains a comprehensive world data seeding system for populating the Geolocation microservice with complete global geographic data using authoritative sources.

## 🌟 Features

- **🌍 Complete World Coverage**: All 250+ countries and territories using ISO 3166 standards
- **🏛️ All World Subdivisions**: States, provinces, regions for all countries using pycountry
- **🏙️ Major World Cities**: 40+ major cities across all continents with coordinates
- **🏗️ Authoritative Sources**: ISO standards, REST Countries API, pycountry library
- **🔄 Incremental Loading**: Supports partial updates and selective seeding
- **� Progress Tracking**: Real-time progress bars and comprehensive statistics
- **� High Performance**: Batch processing with rate limiting and error recovery

## 📁 Files Overview

```
data/seeds/
├── README.md                    # This comprehensive guide
├── world_data_seeder.py        # Comprehensive world data seeder (main script)
├── requirements.txt            # Python dependencies with geographic libraries
└── run_seeder.sh              # Shell script to setup and run seeder
```

## 🚀 Quick Start

### Automated Setup (Recommended)

The shell script automatically sets up the environment and runs the comprehensive world seeder.

```bash
# Make sure the geolocation service is running first
cd /path/to/geolocation-service
go run cmd/main.go

# In another terminal, run the comprehensive world seeder
cd data/seeds
./run_seeder.sh --full
```

### Manual Python Setup

```bash
# Install Python dependencies with geographic libraries
cd data/seeds
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

# Run the comprehensive world seeder
python3 world_data_seeder.py --full --url=http://localhost:3003/api/v1
```

## 📋 Usage Options

### Complete World Data Seeding
```bash
./run_seeder.sh --full                    # Complete world data (250+ countries + subdivisions + cities)
python3 world_data_seeder.py --full       # Direct Python execution
```

### Selective Data Seeding
```bash
./run_seeder.sh --countries-only          # All world countries only (~250 countries)
./run_seeder.sh --subdivisions-only       # All world subdivisions (states, provinces, regions)
./run_seeder.sh --cities-only             # Major world cities only (~40 cities)
```

### Reset and Reseed
```bash
./run_seeder.sh --reset --full            # Clear all data and reseed complete world data
```

### Custom Configuration
```bash
./run_seeder.sh --url=http://production:3003/api/v1  # Custom service URL
python3 world_data_seeder.py --batch-size=100        # Custom batch size
python3 world_data_seeder.py --rate-limit=0.05       # Faster processing
```

## 📊 Data Sources

### World Countries (250+ countries)
- **Primary**: [pycountry library](https://pypi.org/project/pycountry/) - Official ISO 3166-1 standards
- **Enhanced**: [REST Countries API](https://restcountries.com/v3.1/all) - Regions, currencies, capitals
- **Includes**: All UN recognized countries, ISO codes (2 & 3 letter), regions, currencies, capitals

### World Subdivisions (3000+ subdivisions)
- **Source**: pycountry ISO 3166-2 database
- **Coverage**: States, provinces, regions, departments, cantons for all countries
- **Examples**: US states (50), Canadian provinces (13), German states (16), etc.

### Major World Cities (40+ cities)
- **Coverage**: All continents with major metropolitan areas
- **Criteria**: Population centers, economic importance, geographic distribution
- **Data**: Coordinates, postal codes, population estimates
- **Examples**: 
  - **Americas**: New York, Los Angeles, Toronto, São Paulo, Buenos Aires
  - **Europe**: London, Paris, Berlin, Madrid, Rome, Moscow
  - **Asia**: Tokyo, Beijing, Mumbai, Seoul, Bangkok, Singapore
  - **Africa**: Cairo, Lagos, Johannesburg, Casablanca, Nairobi
  - **Oceania**: Sydney, Melbourne, Auckland

## 🔧 Configuration

### Environment Variables
```bash
export GEOLOCATION_SERVICE_URL="http://localhost:3003/api/v1"
export SEED_BATCH_SIZE="50"
export SEED_RATE_LIMIT="100ms"
```

### Service Requirements
The seeder requires the Geolocation service to be running with these endpoints:
- `POST /api/v1/bulk/countries` - Bulk country creation
- `POST /api/v1/bulk/states` - Bulk state creation  
- `POST /api/v1/bulk/cities` - Bulk city creation
- `GET /api/v1/countries` - List countries
- `GET /api/v1/countries/{id}/states` - List states by country
- `GET /health` - Service health check

## 🏗️ Service Integration

### Adding Bulk Endpoints

The seeder expects bulk endpoints for efficient data loading. Add these to your HTTP handlers:

```go
// In your HTTP server
func (h *GeolocationHandler) SetupBulkRoutes(router *gin.Engine) {
    api := router.Group("/api/v1")
    bulk := api.Group("/bulk")
    {
        bulk.POST("/countries", h.BulkCreateCountries)
        bulk.POST("/states", h.BulkCreateStates)
        bulk.POST("/cities", h.BulkCreateCities)
    }
}
```

### Bulk Handler Example
```go
func (h *GeolocationHandler) BulkCreateCountries(c *gin.Context) {
    var request struct {
        Countries []struct {
            Name     string `json:"name" binding:"required"`
            Code     string `json:"code" binding:"required"`
            Code2    string `json:"code2" binding:"required"`
            Region   string `json:"region"`
            Currency string `json:"currency"`
        } `json:"countries" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Process bulk creation...
    c.JSON(201, gin.H{"created": len(request.Countries)})
}
```

## 📈 Performance & Scaling

### Batch Processing
- **Countries**: 50 per batch (typical: 195 total)
- **States**: 100 per batch (typical: 63 total)  
- **Cities**: 25 per batch (prevents timeout)

### Rate Limiting
- 100ms delay between batches to prevent overwhelming the service
- Configurable via environment variables
- Automatic retry with exponential backoff

### Memory Usage
- Streaming JSON parsing for large datasets
- Minimal memory footprint (~10MB peak)
- Efficient batch processing

## 🧪 Testing & Verification

### Verify Seeded Data
```bash
# Check countries
curl http://localhost:3003/api/v1/countries | jq '.countries | length'

# Check US states  
curl http://localhost:3003/api/v1/countries | jq '.countries[] | select(.code2=="US") | .id'
curl http://localhost:3003/api/v1/countries/{us-id}/states | jq '.states | length'

# Check major cities
curl http://localhost:3003/api/v1/states/{state-id}/cities | jq '.cities | length'
```

### Data Quality Checks
```bash
# Countries should have ISO codes
curl http://localhost:3003/api/v1/countries | jq '.countries[] | select(.code2 | length != 2)'

# States should belong to countries
curl http://localhost:3003/api/v1/states | jq '.states[] | select(.country_id == null)'

# Cities should have coordinates  
curl http://localhost:3003/api/v1/cities | jq '.cities[] | select(.latitude == null or .longitude == null)'
```

## 🚨 Troubleshooting

### Common Issues

**Service Not Running**
```
❌ Geolocation service is not available!
Solution: Start the service first: go run cmd/main.go
```

**Import Errors (Python)**
```
ModuleNotFoundError: No module named 'requests'
Solution: pip install -r requirements.txt
```

**API Timeout**
```
Connection timeout during bulk insert
Solution: Reduce batch size or check service performance
```

**Duplicate Data**
```
Error: country with code 'US' already exists  
Solution: Use --reset flag to clear data first
```

### Debug Mode
```bash
# Enable detailed logging
export SEED_DEBUG=true
./run_seeder.sh --countries-only
```

### Manual Cleanup
```bash
# Connect to database and clear data
psql -h localhost -U postgres -d geolocation_db
DELETE FROM cities; DELETE FROM states; DELETE FROM countries;
```

## 🔄 Data Updates

### Regular Updates
- Countries: Updated when new nations are recognized
- States: Rarely change, but may need updates for territorial changes
- Cities: Updated quarterly for new major metropolitan areas

### Custom Data Sources
You can extend the seeder with additional data sources:

1. **GeoNames API**: For comprehensive city data
2. **World Bank API**: For country economic data  
3. **Natural Earth**: For geographic boundaries
4. **OpenStreetMap**: For detailed location data

### Adding New Countries/States
```python
# Add to fallback data in seed_data.py
fallback_countries = [
    {"name": "New Country", "code": "NCO", "code2": "NC", "region": "Region", "currency": "CUR"}
]
```

## 📝 License & Attribution

- **REST Countries API**: Open source, no API key required
- **Country/State Data**: Public domain administrative data
- **City Coordinates**: Various open sources (OpenStreetMap, GeoNames)

## 🤝 Contributing

To add support for more countries/regions:

1. Add state/province data to the appropriate seeder function
2. Add major cities with coordinates  
3. Update the CSV backup files
4. Test with a fresh database
5. Submit a pull request

---

**Need help?** Check the logs in `.venv/seed.log` or run with `--debug` flag for detailed output.