#!/bin/bash
# Post-start hook for Stax WordPress projects
# This script runs after DDEV containers start

echo "Running post-start hooks..."

# Wait for database to be ready
echo "Waiting for database..."
MAX_TRIES=30
TRIES=0
while ! mysql -h db -u db -pdb -e "SELECT 1" >/dev/null 2>&1; do
    TRIES=$((TRIES + 1))
    if [ $TRIES -ge $MAX_TRIES ]; then
        echo "Database failed to become ready after ${MAX_TRIES} attempts"
        exit 1
    fi
    sleep 1
done
echo "Database ready!"

# Verify WordPress installation
if wp core is-installed --allow-root 2>/dev/null; then
    echo "WordPress is installed"

    # Flush rewrite rules
    echo "Flushing rewrite rules..."
    wp rewrite flush --allow-root

    # Flush object cache if available
    if wp cache flush --allow-root 2>/dev/null; then
        echo "Cache flushed"
    fi

    # Verify core checksums
    echo "Verifying WordPress core..."
    wp core verify-checksums --allow-root || true

    echo "Post-start hooks complete!"
else
    echo "WordPress not yet installed - skipping post-start hooks"
fi
