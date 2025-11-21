#!/bin/bash
# Quick start script for the entire system

echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║         WAN2.1 VIDEO GENERATOR - FULL SYSTEM START             ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""

# Function to start a component in a new terminal
start_in_terminal() {
    local title=$1
    local command=$2
    local dir=$3
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        osascript -e "tell application \"Terminal\" to do script \"cd '$dir' && echo '$title' && $command\""
    elif command -v gnome-terminal &> /dev/null; then
        # Linux with gnome-terminal
        gnome-terminal --title="$title" -- bash -c "cd '$dir' && echo '$title' && $command; exec bash"
    elif command -v xterm &> /dev/null; then
        # Linux with xterm
        xterm -T "$title" -e "cd '$dir' && echo '$title' && $command; bash" &
    else
        echo "⚠️  Please manually run in separate terminal:"
        echo "   cd $dir && $command"
        echo ""
        read -p "Press Enter when ready to continue..."
    fi
}

# Get absolute paths
MAIN_DIR=$(cd ../.. && pwd)
WEB_DIR=$(pwd)

echo "📍 Main API Directory: $MAIN_DIR"
echo "📍 Web UI Directory: $WEB_DIR"
echo ""

# Check if main project exists
if [ ! -f "$MAIN_DIR/main.go" ]; then
    echo "❌ Main project not found at $MAIN_DIR"
    echo "Make sure you're in videogen/web/"
    exit 1
fi

echo "🚀 Starting all components..."
echo ""

# Start Python backend
echo "1️⃣  Starting Python Backend..."
start_in_terminal "Python Backend (Port 5000)" \
    "source venv/bin/activate && python server.py" \
    "$MAIN_DIR/python_backend"
sleep 2

# Start Go API server
echo "2️⃣  Starting Go API Server..."
start_in_terminal "Go API Server (Port 8080)" \
    "./wan2-video-server" \
    "$MAIN_DIR"
sleep 2

# Start Web UI
echo "3️⃣  Starting Web UI..."
start_in_terminal "Web UI (Port 3000)" \
    "go run main.go" \
    "$WEB_DIR"
sleep 2

echo ""
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║                    ALL COMPONENTS STARTED!                       ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""
echo "🌐 Access points:"
echo "   Web UI:        http://localhost:3000  ← Start here!"
echo "   API Server:    http://localhost:8080"
echo "   Python Backend: http://localhost:5000"
echo ""
echo "📖 Quick links:"
echo "   • Home:         http://localhost:3000/"
echo "   • Text to Video: http://localhost:3000/text-to-video"
echo "   • Gallery:      http://localhost:3000/gallery"
echo ""
echo "⚡ Tips:"
echo "   • Use Ctrl+C in each terminal to stop"
echo "   • Check logs if something doesn't work"
echo "   • Make sure ports 3000, 5000, 8080 are free"
echo ""
