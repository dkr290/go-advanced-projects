#!/bin/bash
# Quick verification script - see what's been created

echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║        WAN2.1 VIDEO SERVER - FILE VERIFICATION                   ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""

echo "📍 Current directory: $(pwd)"
echo ""

echo "📊 File Statistics:"
echo "────────────────────────────────────────────────────────────────────"
TOTAL_FILES=$(find . -type f | wc -l)
GO_FILES=$(find . -name "*.go" -type f | wc -l)
PY_FILES=$(find . -name "*.py" -type f | wc -l)
MD_FILES=$(find . -name "*.md" -type f | wc -l)
SH_FILES=$(find . -name "*.sh" -type f | wc -l)

echo "  Total files:        $TOTAL_FILES"
echo "  Go files (.go):     $GO_FILES"
echo "  Python files (.py): $PY_FILES"
echo "  Markdown docs:      $MD_FILES"
echo "  Shell scripts:      $SH_FILES"
echo ""

echo "📁 Directory Structure (2 levels):"
echo "────────────────────────────────────────────────────────────────────"
tree -L 2 -d 2>/dev/null || find . -maxdepth 2 -type d | grep -v ".git" | sort
echo ""

echo "📄 Key Files Present:"
echo "────────────────────────────────────────────────────────────────────"
FILES_TO_CHECK=(
    "main.go"
    "go.mod"
    "README.md"
    "QUICKSTART.md"
    "setup.sh"
    "Makefile"
    "Dockerfile"
    "docker-compose.yml"
    ".env.example"
    "python_backend/server.py"
    "python_backend/server_multimodel.py"
    "python_backend/requirements.txt"
    "cmd/root.go"
    "pkg/server/server.go"
    "docs/API.md"
)

for file in "${FILES_TO_CHECK[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✅ $file"
    else
        echo "  ❌ $file (missing)"
    fi
done
echo ""

echo "📁 Created Directories:"
echo "────────────────────────────────────────────────────────────────────"
for dir in cmd pkg python_backend docs examples scripts models uploads outputs; do
    if [ -d "$dir" ]; then
        COUNT=$(find "$dir" -type f 2>/dev/null | wc -l)
        echo "  ✅ $dir/ ($COUNT files)"
    else
        echo "  ❌ $dir/ (missing)"
    fi
done
echo ""

echo "💾 Estimated Disk Usage:"
echo "────────────────────────────────────────────────────────────────────"
du -sh . 2>/dev/null | awk '{print "  Current size: " $1}'
echo "  After setup:  ~500 MB (with dependencies)"
echo "  After model:  ~10-20 GB (with LTX-Video model)"
echo ""

echo "🎯 Quick Commands:"
echo "────────────────────────────────────────────────────────────────────"
echo "  View structure:     tree -L 3"
echo "  List all files:     find . -type f | sort"
echo "  Start setup:        ./setup.sh"
echo "  Read docs:          cat README.md"
echo "  Check GPU:          python scripts/check_amd_gpu.py"
echo ""

echo "✨ Everything is in: $(pwd)"
echo ""
