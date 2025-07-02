#!/bin/bash

# Test script to verify Go proxy publication workflow
# This simulates what happens in the GitHub Actions workflow

set -e

echo "=== Go Module Publication Test Script ==="
echo

# 1. Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

echo "✅ In git repository"

# 2. Check if go.mod exists
if [ ! -f "go.mod" ]; then
    echo "❌ go.mod not found"
    exit 1
fi

echo "✅ go.mod found"

# 3. Get module name from go.mod
MODULE_NAME=$(go list -m)
echo "📦 Module name: $MODULE_NAME"

# 4. Run go mod tidy
echo "🧹 Running go mod tidy..."
go mod tidy

# 5. Run tests (exclude examples directory to avoid main function conflicts)
echo "🧪 Running tests..."
go test -v $(go list ./... | grep -v examples)

# 6. Get latest tag (if any)
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [ -n "$LATEST_TAG" ]; then
    echo "🏷️  Latest tag: $LATEST_TAG"
    
    # 7. Test Go proxy publication (simulation)
    echo "🚀 Testing Go proxy publication..."
    echo "Would run: GOPROXY=proxy.golang.org go list -m $MODULE_NAME@$LATEST_TAG"
    
    # Actually test it (commented out to avoid spam)
    # GOPROXY=proxy.golang.org go list -m $MODULE_NAME@$LATEST_TAG
    
    echo "🌐 Module would be available at:"
    echo "   https://pkg.go.dev/$MODULE_NAME@$LATEST_TAG"
    echo "   https://pkg.go.dev/$MODULE_NAME"
else
    echo "⚠️  No tags found - would skip publication"
fi

echo
echo "✅ Publication test completed successfully!"
echo
echo "To create a new release:"
echo "  1. git commit -m 'feat: your changes'"
echo "  2. git tag v0.x.x"
echo "  3. git push origin v0.x.x"
echo "  4. GitHub Actions will handle the rest!"
