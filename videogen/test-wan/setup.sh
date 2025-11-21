#!/bin/bash

# Wan2.1 Video Server - Quick Setup Script
# This script sets up both Go and Python environments

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}"
echo "=========================================="
echo "  Wan2.1 Video Server - Setup"
echo "=========================================="
echo -e "${NC}"

# Check if Go is installed
echo -e "${YELLOW}Checking dependencies...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}Go is not installed. Please install Go 1.21 or higher.${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Go $(go version | awk '{print $3}')${NC}"

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo -e "${RED}Python3 is not installed. Please install Python 3.9 or higher.${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Python $(python3 --version | awk '{print $2}')${NC}"

# Check if CUDA is available (optional)
if command -v nvidia-smi &> /dev/null; then
    echo -e "${GREEN}✓ NVIDIA GPU detected${NC}"
    nvidia-smi --query-gpu=name,memory.total --format=csv,noheader | head -1
else
    echo -e "${YELLOW}⚠ No NVIDIA GPU detected. Will use CPU (slower).${NC}"
fi

echo ""
echo -e "${BLUE}Setting up Go environment...${NC}"

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download
go mod tidy
echo -e "${GREEN}✓ Go dependencies installed${NC}"

# Build the application
echo "Building application..."
go build -o wan2-video-server main.go
echo -e "${GREEN}✓ Application built successfully${NC}"

echo ""
echo -e "${BLUE}Setting up Python backend...${NC}"

cd python_backend

# Create virtual environment
if [ ! -d "venv" ]; then
    echo "Creating Python virtual environment..."
    python3 -m venv venv
    echo -e "${GREEN}✓ Virtual environment created${NC}"
else
    echo -e "${YELLOW}Virtual environment already exists${NC}"
fi

# Activate virtual environment
source venv/bin/activate

# Upgrade pip
echo "Upgrading pip..."
pip install --upgrade pip > /dev/null 2>&1

# Install PyTorch
echo "Installing PyTorch..."
if command -v nvidia-smi &> /dev/null; then
    echo "  Installing with CUDA support..."
    pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118
else
    echo "  Installing CPU version..."
    pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cpu
fi
echo -e "${GREEN}✓ PyTorch installed${NC}"

# Install other dependencies
echo "Installing Python dependencies..."
pip install -r requirements.txt
echo -e "${GREEN}✓ Python dependencies installed${NC}"

cd ..

# Create necessary directories
echo ""
echo -e "${BLUE}Creating directories...${NC}"
mkdir -p uploads outputs models
echo -e "${GREEN}✓ Directories created${NC}"

# Setup configuration
if [ ! -f ".env" ]; then
    echo ""
    echo -e "${BLUE}Setting up configuration...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ Configuration file created (.env)${NC}"
    echo -e "${YELLOW}⚠ Please edit .env file with your settings${NC}"
else
    echo -e "${YELLOW}.env file already exists, skipping...${NC}"
fi

# Make scripts executable
chmod +x examples/api_examples.sh

echo ""
echo -e "${GREEN}"
echo "=========================================="
echo "  Setup Complete!"
echo "=========================================="
echo -e "${NC}"

echo ""
echo -e "${BLUE}Next steps:${NC}"
echo ""
echo "1. Edit configuration (optional):"
echo -e "   ${YELLOW}nano .env${NC}"
echo ""
echo "2. Download the model (required):"
echo -e "   ${YELLOW}./wan2-video-server download${NC}"
echo ""
echo "3. Start the Python backend (Terminal 1):"
echo -e "   ${YELLOW}cd python_backend${NC}"
echo -e "   ${YELLOW}source venv/bin/activate${NC}"
echo -e "   ${YELLOW}python server.py${NC}"
echo ""
echo "4. Start the Go server (Terminal 2):"
echo -e "   ${YELLOW}./wan2-video-server${NC}"
echo ""
echo "5. Test the API:"
echo -e "   ${YELLOW}./examples/api_examples.sh${NC}"
echo ""
echo -e "${GREEN}Enjoy generating videos with Wan2.1!${NC}"
echo ""
