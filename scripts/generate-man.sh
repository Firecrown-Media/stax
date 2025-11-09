#!/bin/bash
# Generate man page for stax

set -e

# Get version from git or default
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
DATE=$(date +"%B %Y")

# Create output directory
mkdir -p dist/man

# Method 1: Use cobra to generate man page
if [ -f "./stax" ]; then
    echo "Generating man page with Cobra..."
    ./stax man -o dist/man/
fi

# Method 2: Use template if available
if [ -f "docs/stax.1.template" ]; then
    echo "Generating man page from template..."
    sed -e "s/{{.Version}}/$VERSION/g" \
        -e "s/{{.Date}}/$DATE/g" \
        docs/stax.1.template > dist/man/stax.1
fi

echo "Man page generated: dist/man/stax.1"
echo ""
echo "Preview with: man dist/man/stax.1"
echo "Install with: sudo cp dist/man/stax.1 /usr/local/share/man/man1/"
echo "Then run: man stax"
