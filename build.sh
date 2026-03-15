#!/bin/bash

echo "🚀 Building portfolio site..."

# Install dependencies
go mod tidy

# Generate static site
go run main.go

echo "✅ Site generated successfully!"
echo "📁 Files are in the 'docs' directory"
echo "🌐 Open docs/index.html in your browser to preview"

# # Optional: Start local server for preview
# if command -v python3 &> /dev/null; then
#     echo "🔧 Starting local server at http://localhost:8000"
#     cd docs && python3 -m http.server 8000
# fi