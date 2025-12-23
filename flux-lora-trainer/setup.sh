#!/bin/bash

# Flux LoRA Training - Quick Setup Script
# This script sets up the environment and prepares for training

set -e  # Exit on error

echo "================================"
echo "Flux LoRA Training - Quick Setup"
echo "================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo -e "${RED}Error: Python 3 is not installed${NC}"
    echo "Please install Python 3.10 or higher"
    exit 1
fi

echo -e "${GREEN}✓${NC} Python 3 found: $(python3 --version)"

# Check if we're in a virtual environment
if [[ -z "$VIRTUAL_ENV" ]]; then
    echo ""
    echo -e "${YELLOW}!${NC} Not in a virtual environment"
    echo "Creating virtual environment..."
    
    # Create venv
    python3 -m venv venv
    
    echo -e "${GREEN}✓${NC} Virtual environment created"
    echo ""
    echo "To activate the virtual environment, run:"
    echo "  source venv/bin/activate"
    echo ""
    echo "Then run this script again."
    exit 0
fi

echo -e "${GREEN}✓${NC} Virtual environment active"

# Check if requirements are installed
echo ""
echo "Checking dependencies..."

if ! python3 -c "import torch" &> /dev/null; then
    echo -e "${YELLOW}!${NC} PyTorch not found, installing dependencies..."
    echo ""
    
    # Install requirements
    pip install --upgrade pip
    pip install -r requirements.txt
    
    echo ""
    echo -e "${GREEN}✓${NC} Dependencies installed"
else
    echo -e "${GREEN}✓${NC} Dependencies already installed"
fi

# Check CUDA availability
echo ""
echo "Checking GPU..."
python3 << EOF
import torch
if torch.cuda.is_available():
    print(f"✓ CUDA available: {torch.cuda.get_device_name(0)}")
    print(f"  VRAM: {torch.cuda.get_device_properties(0).total_memory / 1024**3:.1f} GB")
else:
    print("⚠ No CUDA GPU detected - training will be very slow on CPU")
    print("  Make sure you have an NVIDIA GPU with CUDA support")
EOF

# Check Hugging Face CLI
echo ""
echo "Checking Hugging Face authentication..."

if command -v huggingface-cli &> /dev/null; then
    echo -e "${GREEN}✓${NC} Hugging Face CLI installed"
    
    # Check if logged in
    if huggingface-cli whoami &> /dev/null; then
        echo -e "${GREEN}✓${NC} Logged in to Hugging Face"
    else
        echo -e "${YELLOW}!${NC} Not logged in to Hugging Face"
        echo ""
        echo "To download Flux models, you need to:"
        echo "1. Create account at https://huggingface.co"
        echo "2. Accept license at https://huggingface.co/black-forest-labs/FLUX.1-dev"
        echo "3. Run: huggingface-cli login"
    fi
else
    echo -e "${YELLOW}!${NC} Hugging Face CLI not found"
    echo "Installing..."
    pip install huggingface-hub[cli]
fi

# Create directory structure
echo ""
echo "Setting up directories..."

mkdir -p dataset
mkdir -p outputs
mkdir -p cache

echo -e "${GREEN}✓${NC} Directories created"

# Check if dataset exists
if [ -z "$(ls -A dataset/)" ]; then
    echo ""
    echo -e "${YELLOW}!${NC} Dataset directory is empty"
    echo ""
    echo "Next steps:"
    echo "1. Create a subdirectory in dataset/ for your subject:"
    echo "   mkdir -p dataset/my_subject"
    echo ""
    echo "2. Add 10-50 images to dataset/my_subject/"
    echo ""
    echo "3. Generate captions:"
    echo "   python caption_images.py --dataset_path dataset/my_subject --trigger_word MYSUBJECT"
    echo ""
    echo "4. Validate dataset:"
    echo "   python prepare_dataset.py --dataset_path dataset/my_subject --trigger_word MYSUBJECT"
    echo ""
    echo "5. Start training:"
    echo "   python train_lora.py --config config.yaml"
else
    echo -e "${GREEN}✓${NC} Dataset directory exists"
    
    # Count subdirectories
    dataset_count=$(find dataset -mindepth 1 -maxdepth 1 -type d | wc -l)
    echo "  Found $dataset_count dataset folder(s)"
fi

echo ""
echo "================================"
echo -e "${GREEN}Setup Complete!${NC}"
echo "================================"
echo ""
echo "Quick Start Guide:"
echo ""
echo "1. Prepare your images:"
echo "   mkdir -p dataset/my_subject"
echo "   # Copy your images to dataset/my_subject/"
echo ""
echo "2. Generate captions:"
echo "   python caption_images.py --dataset_path dataset/my_subject --trigger_word MYSUBJECT"
echo ""
echo "3. Validate dataset:"
echo "   python prepare_dataset.py --dataset_path dataset/my_subject"
echo ""
echo "4. Edit config.yaml to customize training settings"
echo ""
echo "5. Start training:"
echo "   python train_lora.py --config config.yaml"
echo ""
echo "6. Test your LoRA:"
echo "   python test_lora.py --lora_path outputs/my_lora/final.safetensors --prompt 'MYSUBJECT in space'"
echo ""
echo "For detailed documentation, see README.md"
echo ""
