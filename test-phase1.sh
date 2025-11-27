#!/bin/bash
set -e

echo "=== Phase 1 Integration Tests ==="

# Get the full path to the binary
MAAJISE_BIN="$(cd "$(dirname "$0")" && pwd)/maajise.exe"

# Use timestamp to create unique test directories
TIMESTAMP=$(date +%s%N | cut -b1-13)
TEST_DIR="phase1-tests-$TIMESTAMP"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

echo ""
echo "Test 1: Basic init with non-interactive flags"
"$MAAJISE_BIN" init --git-name "Test User" --git-email "test@example.com" --skip-beads test-project-1
cd test-project-1/test-project-1
if [ ! -d .git ]; then
    echo "FAIL: .git directory not created"
    exit 1
fi
cd ../../..
echo "✓ Test 1 passed"

echo ""
echo "Test 2: Convenience shorthand (flags before project name)"
"$MAAJISE_BIN" --git-name "Test User" --git-email "test@example.com" --skip-beads test-project-2
if [ ! -d "$TEST_DIR/test-project-2/$TEST_DIR/test-project-2/.git" ]; then
    # Try alternative path structure
    if [ ! -d test-project-2/test-project-2/.git ]; then
        echo "FAIL: convenience shorthand failed"
        exit 1
    fi
fi
echo "✓ Test 2 passed"

echo ""
echo "Test 3: Help command"
"$MAAJISE_BIN" help
"$MAAJISE_BIN" help init
echo "✓ Test 3 passed"

echo ""
echo "Test 4: Version command"
"$MAAJISE_BIN" version
echo "✓ Test 4 passed"

echo ""
echo "Test 5: Flags work (--skip-git)"
"$MAAJISE_BIN" init --skip-git --skip-beads test-project-3
if [ -d test-project-3/test-project-3/.git ]; then
    echo "FAIL: --skip-git didn't work"
    exit 1
fi
echo "✓ Test 5 passed"

echo ""
echo "Test 6: Global --help"
"$MAAJISE_BIN" --help
echo "✓ Test 6 passed"

echo ""
echo "Test 7: Global --version"
"$MAAJISE_BIN" --version
echo "✓ Test 7 passed"

echo ""
echo "=== All Phase 1 Integration Tests Passed ==="

# Cleanup test directory
cd ..
echo "Cleaning up test directory: $TEST_DIR"
rm -rf "$TEST_DIR" 2>/dev/null || true
