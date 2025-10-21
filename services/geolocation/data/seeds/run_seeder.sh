#!/bin/bash

# TheraClosure Geolocation Service Data Seeder Setup & Runner
# Automatically installs dependencies and runs the data seeding

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_URL="http://localhost:3003/api/v1"

echo "🌟 TheraClosure Geolocation Data Seeder"
echo "======================================"

# Check if Python is available
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 is required but not found. Please install Python 3."
    exit 1
fi

echo "✅ Python 3 found"

# Create virtual environment if it doesn't exist
VENV_DIR="$SCRIPT_DIR/.venv"
if [ ! -d "$VENV_DIR" ]; then
    echo "📦 Creating Python virtual environment..."
    python3 -m venv "$VENV_DIR"
fi

# Activate virtual environment
echo "🔄 Activating virtual environment..."
source "$VENV_DIR/bin/activate"

# Install dependencies
echo "📥 Installing Python dependencies..."
pip install --upgrade pip
pip install -r "$SCRIPT_DIR/requirements.txt"

# Check if geolocation service is running
echo "🔍 Checking geolocation service health..."
if curl -s "$SERVICE_URL/../health" > /dev/null 2>&1; then
    echo "✅ Geolocation service is running"
else
    echo "❌ Geolocation service is not running!"
    echo "Please start the geolocation service first:"
    echo "  cd services/geolocation"
    echo "  go run cmd/main.go"
    echo ""
    echo "Or using Docker:"
    echo "  docker-compose up geolocation-service"
    exit 1
fi

# Parse command line arguments
SEED_ARGS=""
for arg in "$@"; do
    case $arg in
        --countries-only|--states-only|--cities-only|--reset)
            SEED_ARGS="$SEED_ARGS $arg"
            ;;
        --url=*)
            SERVICE_URL="${arg#*=}"
            SEED_ARGS="$SEED_ARGS --url $SERVICE_URL"
            ;;
        --help|-h)
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --countries-only    Seed only countries"
            echo "  --states-only       Seed only states/provinces"  
            echo "  --cities-only       Seed only cities"
            echo "  --reset            Clear all data before seeding"
            echo "  --url=URL          Geolocation service API URL (default: $SERVICE_URL)"
            echo "  --help, -h         Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                              # Seed all data"
            echo "  $0 --countries-only            # Seed only countries"
            echo "  $0 --reset                     # Clear and reseed all data"
            echo "  $0 --url=http://prod:3003/api/v1  # Use different service URL"
            exit 0
            ;;
    esac
done

# Run the seeder
echo ""
echo "🚀 Running comprehensive world data seeder..."
echo "Service URL: $SERVICE_URL"

python3 "$SCRIPT_DIR/world_data_seeder.py" --url="$SERVICE_URL" $SEED_ARGS

echo ""
echo "🎉 Data seeding completed!"
echo ""
echo "You can verify the data by making API calls:"
echo "  curl $SERVICE_URL/countries"
echo "  curl $SERVICE_URL/countries/{country-id}/states"
echo "  curl $SERVICE_URL/states/{state-id}/cities"