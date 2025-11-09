#!/bin/bash
# Build themes
# Compiles parent and child theme assets

set -e

# Build parent theme
PARENT_THEME="wp-content/themes/firecrown-parent"
if [ -d "$PARENT_THEME" ]; then
    echo "Building parent theme..."
    cd "$PARENT_THEME"

    if [ -f "package.json" ]; then
        echo "  - Installing npm dependencies..."
        npm install --legacy-peer-deps

        echo "  - Building assets..."
        npm run build
    fi

    if [ -f "composer.json" ]; then
        echo "  - Installing composer dependencies..."
        composer install --no-dev --ignore-platform-reqs
    fi

    echo "Parent theme build complete"
    cd - > /dev/null
else
    echo "Warning: Parent theme not found: $PARENT_THEME"
fi

# Build child theme
CHILD_THEME="wp-content/themes/firecrown-child"
if [ -d "$CHILD_THEME" ]; then
    echo ""
    echo "Building child theme..."
    cd "$CHILD_THEME"

    if [ -f "package.json" ]; then
        echo "  - Installing npm dependencies..."
        npm install --legacy-peer-deps

        echo "  - Building assets..."
        npm run build
    fi

    if [ -f "composer.json" ]; then
        echo "  - Installing composer dependencies..."
        composer install --no-dev --ignore-platform-reqs
    fi

    echo "Child theme build complete"
    cd - > /dev/null
else
    echo "Warning: Child theme not found: $CHILD_THEME"
fi
