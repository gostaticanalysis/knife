#!/bin/bash

# Check for broken links in documentation files

set -e

echo "Checking documentation links..."

# Find all markdown files
md_files=$(find . -name "*.md" -not -path "./vendor/*")

broken_links=0

for file in $md_files; do
    echo "Checking $file..."
    
    # Extract URLs from markdown files
    urls=$(grep -oE 'https?://[^)]+' "$file" 2>/dev/null || true)
    
    for url in $urls; do
        # Check HTTP status
        status=$(curl -I -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null || echo "000")
        
        if [[ "$status" -ge 400 ]] || [[ "$status" == "000" ]]; then
            echo "❌ BROKEN: $url (HTTP $status) in $file"
            broken_links=$((broken_links + 1))
        else
            echo "✅ OK: $url (HTTP $status)"
        fi
    done
done

if [[ $broken_links -gt 0 ]]; then
    echo "Found $broken_links broken link(s)"
    exit 1
else
    echo "All links are working!"
fi