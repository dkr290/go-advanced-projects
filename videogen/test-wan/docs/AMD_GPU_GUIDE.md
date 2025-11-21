# AMD GPU (ROCm) Support Guide

## Current Status

The default setup uses **CUDA (NVIDIA only)**. For AMD GPUs, you need **ROCm**.

---

## AMD GPU Compatibility

### ✅ Supported AMD GPUs (with ROCm 5.7+)

- **RDNA 2/3 Series**
  - RX 7900 XTX/XT ✅
  - RX 7800 XT ✅
  - RX 7700 XT ✅
  - RX 6950 XT ✅
  - RX 6900 XT ✅
  - RX 6800 XT ✅
  - RX 6700 XT ✅

- **Older Cards (Limited Support)**
  - RX 5700 XT ⚠️
  - Vega 64/56 ⚠️

### ❌ Not Supported
- RX 500 series and older
- APUs (Ryzen integrated graphics)

---

## Option 1: ROCm Setup (Linux Only)

### Prerequisites
- **Linux Only** (Ubuntu 20.04/22.04, RHEL 8+)
- AMD GPU with 12GB+ VRAM
- ROCm 5.7 or newer

### Installation Steps

#### 1. Install ROCm

```bash
# Ubuntu 22.04
wget https://repo.radeon.com/amdgpu-install/5.7/ubuntu/jammy/amdgpu-install_5.7.50700-1_all.deb
sudo apt install ./amdgpu-install_5.7.50700-1_all.deb
sudo amdgpu-install --usecase=rocm

# Add user to groups
sudo usermod -a -G render,video $USER
newgrp render

# Verify installation
rocm-smi
```

#### 2. Install PyTorch for ROCm

```bash
cd python_backend
source venv/bin/activate

# Uninstall CUDA PyTorch
pip uninstall torch torchvision torchaudio

# Install ROCm PyTorch
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/rocm5.7

# Verify
python -c "import torch; print(f'ROCm available: {torch.cuda.is_available()}'); print(f'Device: {torch.cuda.get_device_name(0)}')"
```

#### 3. Update Environment Configuration

```bash
# In .env
ENABLE_GPU=true
GPU_DEVICE_ID=0
GPU_BACKEND=rocm  # New variable

# Add to your shell
export HSA_OVERRIDE_GFX_VERSION=11.0.0  # For RDNA 3
# or
export HSA_OVERRIDE_GFX_VERSION=10.3.0  # For RDNA 2
```

#### 4. Modify Python Backend

See modifications below for ROCm compatibility.

---

## Option 2: CPU Mode (Works Now!)

**Already configured!** Just use CPU mode:

```bash
# In .env
ENABLE_GPU=false
```

This works on ANY system (AMD, NVIDIA, Intel, Apple Silicon via Rosetta).

**Performance:**
- Slow but functional
- 256x256, 16 frames: ~10-20 minutes
- Good for testing/development

---

## Option 3: DirectML (Windows + AMD)

For Windows with AMD GPU, use DirectML:

### Installation

```bash
cd python_backend

# Install torch-directml
pip install torch-directml

# Install other dependencies
pip install diffusers transformers accelerate
```

### Modify server.py for DirectML

I'll create a DirectML version below.

---

## Option 4: Cloud GPU (Easiest!)

Use cloud providers with NVIDIA GPUs:
- **Google Colab**: Free T4 (16GB)
- **Vast.ai**: Rent NVIDIA GPUs cheaply
- Your code works unchanged!

---

## Modified Files for AMD Support

I'll create these files now...

