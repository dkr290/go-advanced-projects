# GPU Compatibility Summary

## Your Questions Answered

### Q1: Will it work with 4GB GPU?
**❌ No** - LTX-Video requires minimum 12GB VRAM

**Your options:**
1. ✅ **CPU Mode** (Accepted!) - Slow but works
2. ✅ Google Colab - Free T4 with 16GB
3. ✅ Switch to ModelScope model - Works on 4GB

### Q2: Will it work with AMD GPU?
**⚠️ Yes, but Linux only with manual setup**

**AMD GPU Support:**
- ✅ Linux (Ubuntu 20.04/22.04) - Run `./setup_rocm.sh`
- ❌ Windows - Use CPU mode instead
- ❌ macOS - AMD GPUs not supported

---

## Complete Compatibility Matrix

```
┌─────────────────────────────────────────────────────────┐
│              WILL IT WORK?                              │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Hardware              │  Works?  │  How?              │
│  ──────────────────────┼──────────┼─────────────────── │
│  NVIDIA GPU (12GB+)    │  ✅ YES  │  Default setup     │
│  NVIDIA GPU (4-8GB)    │  ❌ NO   │  Use CPU/cloud     │
│  AMD GPU (Linux)       │  ✅ YES  │  ./setup_rocm.sh   │
│  AMD GPU (Windows)     │  ❌ NO   │  Use CPU mode      │
│  AMD GPU (4GB)         │  ❌ NO   │  Use CPU mode      │
│  Intel GPU             │  ❌ NO   │  Use CPU mode      │
│  Apple M1/M2/M3        │  ⚠️ Exp  │  Use CPU mode      │
│  Any CPU               │  ✅ YES  │  ENABLE_GPU=false  │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## Recommended Configurations

### For CPU Mode (What you chose!)

```bash
# .env configuration
ENABLE_GPU=false
DEFAULT_WIDTH=256
DEFAULT_HEIGHT=256
MAX_FRAMES=16
MAX_CONCURRENT_REQUESTS=1
```

**Expected performance:**
- 256x256, 16 frames: ~10-15 minutes
- 512x512, 32 frames: ~30-45 minutes

**Minimum requirements:**
- 16GB RAM (32GB better)
- 8+ CPU cores recommended
- Patience! ☕

---

### For AMD GPU (Linux)

```bash
# 1. Run ROCm setup
chmod +x setup_rocm.sh
./setup_rocm.sh

# 2. Log out and back in (important!)

# 3. Verify ROCm
rocm-smi

# 4. Check PyTorch
cd python_backend
source venv/bin/activate
python -c "import torch; print(torch.cuda.is_available())"

# 5. Start server normally
./wan2-video-server
```

**Requirements:**
- Linux (Ubuntu 20.04/22.04 or similar)
- AMD RDNA 2/3 GPU (RX 6000/7000 series)
- 12GB+ VRAM
- ROCm 5.4+

**Supported AMD GPUs:**
- ✅ RX 7900 XTX/XT (24GB/20GB)
- ✅ RX 7800/7700 XT (16GB/12GB)
- ✅ RX 6950/6900/6800 XT (16GB)
- ✅ RX 6700 XT (12GB)
- ⚠️ RX 5700 XT (8GB - limited)

---

## Step-by-Step: Getting Started

### Option A: CPU Mode (Easiest)

```bash
# 1. Setup
./setup.sh

# 2. Configure for CPU
nano .env
# Set: ENABLE_GPU=false

# 3. Download model
./wan2-video-server download

# 4. Start Python backend (Terminal 1)
cd python_backend && source venv/bin/activate && python server.py

# 5. Start Go server (Terminal 2)
./wan2-video-server

# 6. Test with small video
curl -X POST http://localhost:8080/api/v1/generate/text-to-video \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A cute cat",
    "num_frames": 8,
    "width": 128,
    "height": 128,
    "num_inference_steps": 20
  }'
```

---

### Option B: AMD GPU (Linux Only)

```bash
# 1. Check your GPU
python scripts/check_amd_gpu.py

# 2. Install ROCm
./setup_rocm.sh

# 3. Reboot (important!)
sudo reboot

# 4. Verify
rocm-smi
rocminfo | grep "Name"

# 5. Test PyTorch
cd python_backend && source venv/bin/activate
python -c "import torch; print(f'GPU: {torch.cuda.is_available()}')"

# 6. Configure
nano .env
# Set: ENABLE_GPU=true
# Add: GPU_BACKEND=rocm

# 7. Start normally
./wan2-video-server
```

---

## Troubleshooting Guide

### CPU Mode Issues

**Problem: Too slow**
```bash
# Reduce complexity
DEFAULT_WIDTH=128
DEFAULT_HEIGHT=128
MAX_FRAMES=8
num_inference_steps=10
```

**Problem: Out of memory (RAM)**
```bash
# Close other applications
# Reduce concurrent requests
MAX_CONCURRENT_REQUESTS=1
```

---

### AMD GPU Issues

**Problem: ROCm not detected**
```bash
# Check installation
rocm-smi

# If not found, reinstall
sudo apt remove rocm-*
./setup_rocm.sh
```

**Problem: PyTorch not using GPU**
```bash
# Set architecture override
export HSA_OVERRIDE_GFX_VERSION=11.0.0  # RDNA 3
# or
export HSA_OVERRIDE_GFX_VERSION=10.3.0  # RDNA 2

# Add to ~/.bashrc for persistence
```

**Problem: Out of VRAM**
```bash
# In .env
LOW_MEMORY_MODE=true
MAX_FRAMES=32
DEFAULT_WIDTH=256
DEFAULT_HEIGHT=256
```

---

## Performance Expectations

### CPU Mode
| Processor | Resolution | Frames | Time |
|-----------|------------|--------|------|
| i5/Ryzen 5 | 128x128 | 8 | ~5 min |
| i5/Ryzen 5 | 256x256 | 16 | ~15 min |
| i7/Ryzen 7 | 256x256 | 32 | ~25 min |
| i9/Ryzen 9 | 512x512 | 32 | ~40 min |

### AMD GPU (ROCm)
| GPU | Resolution | Frames | Time |
|-----|------------|--------|------|
| RX 6700 XT | 256x256 | 32 | ~60s |
| RX 6900 XT | 512x512 | 64 | ~90s |
| RX 7900 XTX | 512x512 | 64 | ~60s |
| RX 7900 XTX | 768x768 | 128 | ~180s |

---

## Files Created for GPU Support

✅ `docs/AMD_GPU_GUIDE.md` - Complete AMD GPU guide
✅ `docs/LOW_MEMORY_GUIDE.md` - 4GB GPU survival guide
✅ `setup_rocm.sh` - Automated ROCm installation
✅ `scripts/check_amd_gpu.py` - AMD GPU detection
✅ `GPU_COMPATIBILITY.txt` - Quick reference
✅ `GPU_REQUIREMENTS.txt` - Memory requirements
✅ `.env.low_memory` - Low memory configuration

---

## Quick Commands

```bash
# Check what GPU you have
lspci | grep -i vga

# NVIDIA
nvidia-smi

# AMD (if ROCm installed)
rocm-smi

# Check your configuration
python scripts/check_amd_gpu.py

# Test CPU mode
ENABLE_GPU=false ./wan2-video-server

# Test AMD GPU
./setup_rocm.sh && ./wan2-video-server
```

---

## Summary

✅ **CPU Mode**: Works on any hardware, slow but reliable
✅ **NVIDIA GPU**: Best support, works out of box
⚠️ **AMD GPU**: Linux only, requires ROCm setup
❌ **4GB GPU**: Not enough for LTX-Video
❌ **Windows + AMD**: Use CPU mode instead

**Your best option with current hardware:**
→ CPU mode (accepted!)
→ Or use Google Colab (free 16GB GPU)

Happy generating! 🎥
