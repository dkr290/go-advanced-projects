#!/bin/bash
# Setup script for Python backend

echo "Setting up Python backend for Wan2.1 Video Generation Server"

# Check Python version
python_version=$(python3 --version 2>&1 | awk '{print $2}')
echo "Python version: $python_version"

# Create virtual environment
echo "Creating virtual environment..."
python3 -m venv venv

# Activate virtual environment
echo "Activating virtual environment..."
source venv/bin/activate

# Upgrade pip
echo "Upgrading pip..."
pip install --upgrade pip

# Install PyTorch with CUDA support (for GPU)
echo "Installing PyTorch with CUDA support..."
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118

# Install requirements
echo "Installing Python dependencies..."
pip install -r requirements.txt

# Create necessary directories
echo "Creating directories..."
mkdir -p ../uploads
mkdir -p ../outputs
mkdir -p ../models

echo ""
echo "Setup complete!"
echo ""
echo "To start the Python backend:"
echo "  1. Activate the virtual environment: source venv/bin/activate"
echo "  2. Set your Hugging Face token (optional): export HUGGING_FACE_HUB_TOKEN=your_token"
echo "  3. Run the server: python server.py"
echo ""
echo "For CPU-only installation, reinstall PyTorch without CUDA:"
echo "  pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cpu"
