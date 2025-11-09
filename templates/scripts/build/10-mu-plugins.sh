#!/bin/bash
# Build MU plugins
# Installs PHP dependencies for the Firecrown mu-plugin

set -e

PLUGIN_DIR="wp-content/mu-plugins/firecrown"

if [ ! -d "$PLUGIN_DIR" ]; then
    echo "Warning: MU plugin directory not found: $PLUGIN_DIR"
    exit 0
fi

if [ ! -f "$PLUGIN_DIR/composer.json" ]; then
    echo "Warning: No composer.json found in $PLUGIN_DIR"
    exit 0
fi

echo "Building MU plugin..."
cd "$PLUGIN_DIR"

# Install composer dependencies
composer install --no-dev --prefer-dist --ignore-platform-reqs

echo "MU plugin build complete"
cd - > /dev/null
