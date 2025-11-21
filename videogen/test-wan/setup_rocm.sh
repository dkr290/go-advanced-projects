#!/bin/bash
# ROCm Setup Script for AMD GPUs (Ubuntu/Debian)

set -e

echo "=========================================="
echo "  AMD GPU (ROCm) Setup for Wan2.1 Server"
echo "=========================================="
echo ""

# Check if AMD GPU present
if ! lspci | grep -i amd | grep -i vga > /dev/null; then
    echo "❌ No AMD GPU detected!"
    echo "This script is for AMD GPUs only."
    exit 1
fi

echo "✓ AMD GPU detected"
lspci | grep -i amd | grep -i vga

echo ""
echo "Installing ROCm..."

# Detect Ubuntu version
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    VER=$VERSION_ID
else
    echo "Cannot detect OS version"
    exit 1
fi

echo "OS: $OS $VER"

# Install ROCm based on Ubuntu version
case "$VER" in
    "22.04")
        ROCM_VERSION="5.7"
        UBUNTU_CODENAME="jammy"
        ;;
    "20.04")
        ROCM_VERSION="5.7"
        UBUNTU_CODENAME="focal"
        ;;
    *)
        echo "⚠️  Unsupported Ubuntu version: $VER"
        echo "Trying with Ubuntu 22.04 packages..."
        UBUNTU_CODENAME="jammy"
        ;;
esac

# Download and install ROCm
echo "Downloading ROCm installer..."
wget -q https://repo.radeon.com/amdgpu-install/${ROCM_VERSION}/ubuntu/${UBUNTU_CODENAME}/amdgpu-install_${ROCM_VERSION}.50700-1_all.deb

echo "Installing ROCm..."
sudo apt install -y ./amdgpu-install_${ROCM_VERSION}.50700-1_all.deb
sudo amdgpu-install -y --usecase=rocm --no-dkms

# Add user to required groups
echo "Adding user to render and video groups..."
sudo usermod -a -G render,video $USER

echo ""
echo "Installing PyTorch with ROCm support..."

cd python_backend

# Activate virtual environment
if [ ! -d "venv" ]; then
    echo "Creating Python virtual environment..."
    python3 -m venv venv
fi

source venv/bin/activate

# Uninstall CUDA PyTorch if present
pip uninstall -y torch torchvision torchaudio 2>/dev/null || true

# Install ROCm PyTorch
echo "Installing PyTorch for ROCm 5.7..."
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/rocm5.7

# Install other dependencies
echo "Installing other dependencies..."
pip install -r requirements.txt

cd ..

# Update .env
echo ""
echo "Updating configuration..."

if [ -f ".env" ]; then
    # Update existing .env
    sed -i 's/ENABLE_GPU=.*/ENABLE_GPU=true/' .env
    
    if ! grep -q "GPU_BACKEND" .env; then
        echo "GPU_BACKEND=rocm" >> .env
    fi
else
    # Copy from example
    cp .env.example .env
    sed -i 's/ENABLE_GPU=.*/ENABLE_GPU=true/' .env
    echo "GPU_BACKEND=rocm" >> .env
fi

# Detect GPU architecture
echo ""
echo "Detecting AMD GPU architecture..."
GPU_NAME=$(rocm-smi --showproductname | grep "Card series" | awk '{print $NF}')

case "$GPU_NAME" in
    *"Navi 31"*|*"7900"*)
        echo "Detected: RDNA 3 (Navi 31)"
        echo 'export HSA_OVERRIDE_GFX_VERSION=11.0.0' >> ~/.bashrc
        export HSA_OVERRIDE_GFX_VERSION=11.0.0
        ;;
    *"Navi 21"*|*"6900"*|*"6800"*)
        echo "Detected: RDNA 2 (Navi 21)"
        echo 'export HSA_OVERRIDE_GFX_VERSION=10.3.0' >> ~/.bashrc
        export HSA_OVERRIDE_GFX_VERSION=10.3.0
        ;;
    *)
        echo "⚠️  Could not detect GPU architecture"
        echo "You may need to set HSA_OVERRIDE_GFX_VERSION manually"
        ;;
esac

echo ""
echo "=========================================="
echo "  ROCm Setup Complete!"
echo "=========================================="
echo ""
echo "IMPORTANT: You must log out and log back in for group changes to take effect!"
echo ""
echo "After logging back in, verify installation:"
echo "  rocm-smi"
echo "  cd python_backend && source venv/bin/activate"
echo "  python -c 'import torch; print(torch.cuda.is_available())'"
echo ""
echo "Then start the server normally:"
echo "  ./wan2-video-server"
echo ""
