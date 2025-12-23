# Flux LoRA Training Guide

A comprehensive guide and toolkit for training custom LoRA models on Flux (Flux-dev) using your own images.

## 📋 Table of Contents

- [Overview](#overview)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Detailed Guide](#detailed-guide)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)

## 🎯 Overview

This toolkit allows you to train a custom LoRA adapter for Flux models using your own images. LoRA (Low-Rank Adaptation) is an efficient fine-tuning technique that requires:
- Less VRAM (8-24GB depending on settings)
- Fewer training images (10-50 images recommended)
- Shorter training time (30 minutes to a few hours)

**Complexity Level**: Moderate - requires GPU with at least 12GB VRAM and basic Python knowledge.

## 🔧 Requirements

### Hardware
- **GPU**: NVIDIA GPU with minimum 12GB VRAM (16GB+ recommended)
  - RTX 3060 12GB (minimum)
  - RTX 4090 24GB (optimal)
- **RAM**: 16GB+ system RAM
- **Storage**: 50GB+ free space

### Software
- Python 3.10+
- CUDA 11.8+ or 12.1+
- Git

## 📦 Installation

### Step 1: Clone and Setup

```bash
# Create project directory
mkdir flux-lora-training
cd flux-lora-training

# Create virtual environment
python -m venv venv

# Activate virtual environment
# On Linux/Mac:
source venv/bin/activate
# On Windows:
# venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt
```

### Step 2: Hugging Face Token

You'll need a Hugging Face account and token to download Flux models:

1. Create account at https://huggingface.co
2. Go to Settings → Access Tokens
3. Create a token with read permissions
4. Accept the license at https://huggingface.co/black-forest-labs/FLUX.1-dev

```bash
# Login to Hugging Face
huggingface-cli login
```

## 🚀 Quick Start

### 1. Prepare Your Images

```bash
# Create dataset directory
mkdir -p dataset/my_subject

# Add your images (10-50 images recommended)
# Supported formats: .jpg, .jpeg, .png, .webp
# Copy your images to: dataset/my_subject/
```

### 2. Configure Training

Edit `config.yaml`:

```yaml
# Minimum required settings
project_name: "my_lora_model"
trigger_word: "MYSUBJECT"  # Unique identifier for your subject
```

### 3. Caption Your Images (Optional but Recommended)

```bash
# Auto-caption using BLIP
python caption_images.py --dataset_path dataset/my_subject
```

Or manually create `.txt` files with the same name as each image.

### 4. Start Training

```bash
# Simple training
python train_lora.py --config config.yaml

# Advanced training with custom parameters
python train_lora.py \
    --config config.yaml \
    --dataset_path dataset/my_subject \
    --output_dir outputs/my_lora \
    --max_train_steps 1000 \
    --learning_rate 1e-4
```

### 5. Test Your LoRA

```bash
python test_lora.py \
    --lora_path outputs/my_lora/final.safetensors \
    --prompt "a photo of MYSUBJECT in a forest" \
    --output test_output.png
```

## 📚 Detailed Guide

### Dataset Preparation

#### Image Guidelines
- **Quantity**: 10-50 images (20-30 is sweet spot)
- **Quality**: High resolution (1024x1024 or higher)
- **Variety**: Different angles, lighting, backgrounds
- **Consistency**: All images should feature the same subject/style
- **Format**: JPG, PNG, or WEBP

#### Image Organization

```
dataset/
└── my_subject/
    ├── image_001.jpg
    ├── image_001.txt  # Caption file
    ├── image_002.jpg
    ├── image_002.txt
    └── ...
```

#### Captioning Best Practices

Each `.txt` file should describe the image and include your trigger word:

**Good captions:**
```
a photo of MYSUBJECT wearing a red shirt, smiling, outdoors
MYSUBJECT standing in front of a brick wall, professional lighting
close-up portrait of MYSUBJECT with blue eyes
```

**Bad captions:**
```
image123  # Too vague
a person  # No trigger word
MYSUBJECT  # Too short, no context
```

### Training Configuration

The `config.yaml` file controls all training parameters:

#### Basic Settings
```yaml
# Project identification
project_name: "my_custom_lora"
trigger_word: "MYSUBJECT"  # Use uppercase, unique word

# Model selection
model_name: "black-forest-labs/FLUX.1-dev"  # or FLUX.1-schnell

# Dataset
dataset_path: "dataset/my_subject"
output_dir: "outputs/my_custom_lora"
```

#### Training Parameters

```yaml
# Training steps (adjust based on dataset size)
max_train_steps: 1000  # 500-2000 typical range
# Rule of thumb: 50-100 steps per image

# Learning rate (critical parameter)
learning_rate: 1e-4  # Start here
# Too high: unstable training, artifacts
# Too low: slow convergence, underfitting

# Batch size (GPU memory dependent)
train_batch_size: 1  # 1 for 12GB VRAM, 2-4 for 24GB

# LoRA rank (model capacity)
lora_rank: 16  # 8, 16, 32, or 64
# Higher = more capacity but more VRAM/time
```

#### Memory Optimization

```yaml
# Enable for lower VRAM usage
gradient_checkpointing: true
gradient_accumulation_steps: 4
mixed_precision: "bf16"  # or "fp16"

# Use 8-bit Adam optimizer
use_8bit_adam: true
```

### Training Process

#### Monitoring Training

The script will output:
```
Step 100/1000 - Loss: 0.0523
Step 200/1000 - Loss: 0.0312
Step 300/1000 - Loss: 0.0198
...
```

**Loss interpretation:**
- Loss decreasing = good learning
- Loss < 0.05 = generally well-trained
- Loss stuck/increasing = potential issues

#### Checkpoints

Checkpoints are saved every N steps:
```
outputs/my_lora/
├── checkpoint-200/
├── checkpoint-400/
├── checkpoint-600/
├── ...
└── final.safetensors
```

You can test intermediate checkpoints to prevent overfitting.

#### Training Time Estimates

- **12GB GPU (RTX 3060)**: ~2-4 hours for 1000 steps
- **16GB GPU (RTX 4060 Ti)**: ~1.5-3 hours for 1000 steps
- **24GB GPU (RTX 4090)**: ~30-60 minutes for 1000 steps

### Testing and Inference

```bash
# Basic test
python test_lora.py \
    --lora_path outputs/my_lora/final.safetensors \
    --prompt "MYSUBJECT in a cyberpunk city" \
    --output result.png

# Advanced options
python test_lora.py \
    --lora_path outputs/my_lora/final.safetensors \
    --prompt "a photo of MYSUBJECT as an astronaut" \
    --negative_prompt "blurry, low quality" \
    --num_inference_steps 30 \
    --guidance_scale 7.5 \
    --lora_scale 0.8 \
    --seed 42 \
    --output astronaut.png
```

**LoRA Scale Parameter:**
- `0.0` = no LoRA effect (base model)
- `0.5-0.8` = subtle effect
- `1.0` = full LoRA effect (default)
- `1.2+` = exaggerated effect

## ⚙️ Configuration

### config.yaml Reference

See `config.yaml` for all available options with detailed comments.

Key parameters to tune:
- `max_train_steps`: More steps = more training (but risk overfitting)
- `learning_rate`: Controls how fast the model learns
- `lora_rank`: Higher = more capacity but slower
- `save_steps`: How often to save checkpoints

## 🐛 Troubleshooting

### Common Issues

#### CUDA Out of Memory
```
RuntimeError: CUDA out of memory
```

**Solutions:**
1. Reduce batch size: `train_batch_size: 1`
2. Enable gradient checkpointing: `gradient_checkpointing: true`
3. Increase gradient accumulation: `gradient_accumulation_steps: 8`
4. Reduce LoRA rank: `lora_rank: 8`
5. Use 8-bit optimizer: `use_8bit_adam: true`

#### Model Not Learning (High Loss)
```
Step 1000/1000 - Loss: 0.2534 (stuck)
```

**Solutions:**
1. Increase learning rate: try `2e-4` or `5e-4`
2. Check captions include trigger word
3. Increase training steps
4. Verify image quality and variety

#### Overfitting (Perfect Training, Poor Generation)
- Model memorizes training images exactly
- Generated images look identical to training data

**Solutions:**
1. Reduce training steps
2. Decrease learning rate
3. Add more variety to dataset
4. Use lower LoRA scale during inference (0.6-0.8)

#### Permission/Download Errors
```
401 Unauthorized
```

**Solution:**
1. Accept model license on Hugging Face
2. Login: `huggingface-cli login`
3. Verify token has read permissions

### Performance Optimization

#### Faster Training
1. Use FLUX.1-schnell (faster but lower quality)
2. Reduce `num_inference_steps` during testing
3. Use mixed precision: `mixed_precision: "bf16"`
4. Increase batch size if VRAM allows

#### Better Quality
1. Use more high-quality training images
2. Increase LoRA rank to 32 or 64
3. Train for more steps (monitor validation)
4. Use detailed, accurate captions
5. Ensure trigger word is consistent

## 📖 Additional Resources

- [Flux Model Card](https://huggingface.co/black-forest-labs/FLUX.1-dev)
- [LoRA Paper](https://arxiv.org/abs/2106.09685)
- [Diffusers Documentation](https://huggingface.co/docs/diffusers)

## 🎓 Tips for Best Results

1. **Start Small**: Test with 10-15 images first
2. **Quality > Quantity**: Better to have 15 great images than 50 mediocre ones
3. **Monitor Checkpoints**: Test every checkpoint to catch overfitting early
4. **Experiment**: Try different learning rates and LoRA ranks
5. **Backup**: Save your best checkpoints externally

## 📝 License

This training toolkit is provided as-is. Make sure to comply with:
- Flux model license (non-commercial for FLUX.1-dev)
- Your training data rights
- Generated content usage rights

---

**Need help?** Check the troubleshooting section or review the detailed configuration in `config.yaml`.
