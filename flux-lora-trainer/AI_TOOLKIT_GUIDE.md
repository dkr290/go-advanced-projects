# Complete Guide: Flux LoRA Training with ai-toolkit

## 🎯 Overview

This guide shows you how to train a Flux LoRA using **ai-toolkit by ostris** - the easiest and most reliable method for Flux LoRA training.

**What is ai-toolkit?**
- Battle-tested LoRA training solution
- Specifically designed for Flux models
- Simple configuration, powerful results
- Active development and community support

**Time Required:** 2-4 hours total (mostly GPU time)

**Difficulty:** Easy to Moderate

---

## 📋 Prerequisites

### Hardware Requirements
- **GPU**: NVIDIA with 12GB+ VRAM
  - Minimum: RTX 3060 12GB
  - Recommended: RTX 4090 24GB
  - Budget: Rent cloud GPU (RunPod, Vast.ai)
- **RAM**: 16GB+ system RAM
- **Storage**: 50GB+ free space (SSD preferred)

### Software Requirements
- **OS**: Linux (Ubuntu 20.04+), Windows 10/11, or macOS
- **Python**: 3.10 or 3.11 (NOT 3.12 yet)
- **Git**: For cloning repositories
- **CUDA**: 11.8 or 12.1 (for NVIDIA GPUs)

### Account Requirements
- **Hugging Face account** (free)
- Access to FLUX.1-dev model

---

## 🚀 Step-by-Step Installation

### Step 1: System Preparation (5 minutes)

#### Check Your GPU
```bash
# Linux/WSL
nvidia-smi

# Should show your GPU and CUDA version
```

#### Install Python 3.10 (if needed)

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install python3.10 python3.10-venv python3.10-dev
```

**Windows:**
- Download from https://www.python.org/downloads/
- Choose Python 3.10.x
- Check "Add Python to PATH" during installation

**macOS:**
```bash
brew install python@3.10
```

---

### Step 2: Clone ai-toolkit (2 minutes)

```bash
# Navigate to where you want to install
cd ~  # Or C:\Users\YourName\ on Windows

# Clone the repository
git clone https://github.com/ostris/ai-toolkit.git
cd ai-toolkit

# Verify you're in the right place
ls -la  # Should see config/, scripts/, etc.
```

---

### Step 3: Create Virtual Environment (3 minutes)

```bash
# Create virtual environment
python3.10 -m venv venv

# Activate it
# Linux/Mac:
source venv/bin/activate

# Windows (Command Prompt):
venv\Scripts\activate.bat

# Windows (PowerShell):
venv\Scripts\Activate.ps1

# Your prompt should now show (venv)
```

---

### Step 4: Install Dependencies (10-15 minutes)

```bash
# Upgrade pip first
pip install --upgrade pip

# Install PyTorch with CUDA support
# For CUDA 12.1:
pip install torch torchvision --index-url https://download.pytorch.org/whl/cu121

# For CUDA 11.8:
pip install torch torchvision --index-url https://download.pytorch.org/whl/cu118

# Install ai-toolkit requirements
pip install -r requirements.txt

# Verify installation
python -c "import torch; print(f'PyTorch: {torch.__version__}'); print(f'CUDA: {torch.cuda.is_available()}')"
# Should show: CUDA: True
```

**Common Issues:**
- If CUDA shows False, reinstall PyTorch with correct CUDA version
- If out of memory during install, close other programs
- Windows users: Install Visual Studio Build Tools if needed

---

### Step 5: Hugging Face Setup (5 minutes)

```bash
# Install Hugging Face CLI
pip install huggingface-hub

# Login to Hugging Face
huggingface-cli login

# Paste your token (get from https://huggingface.co/settings/tokens)
# Choose: "Add token as git credential? (Y/n)" → Y
```

**Get Access to Flux:**
1. Go to https://huggingface.co/black-forest-labs/FLUX.1-dev
2. Click "Agree and access repository"
3. Wait for approval (usually instant)

**Verify access:**
```bash
huggingface-cli whoami
# Should show your username
```

---

## 📂 Prepare Your Dataset

### Step 6: Organize Your Images (15-30 minutes)

#### Create Dataset Directory
```bash
# In ai-toolkit directory
mkdir -p dataset/my_subject
```

#### Image Requirements
- **Quantity**: 15-30 images (sweet spot: 20-25)
- **Quality**: High resolution (1024x1024 minimum, 2048x2048+ ideal)
- **Format**: JPG, PNG, WebP
- **Content**: 
  - All same subject (person, character, style, object)
  - Variety in poses, angles, lighting
  - Avoid duplicates or near-duplicates
  - Clean, uncluttered backgrounds preferred

#### Copy Your Images
```bash
# Copy your prepared images
cp /path/to/your/images/*.jpg dataset/my_subject/
cp /path/to/your/images/*.png dataset/my_subject/

# Verify
ls dataset/my_subject/
# Should show your images
```

---

### Step 7: Caption Your Images (10-20 minutes)

You have three options:

#### Option A: Use Our Auto-Captioning Tool (Recommended)

```bash
# Go back to our flux-lora-training toolkit directory
cd /path/to/flux-lora-training

# Activate its environment
source venv/bin/activate

# Generate captions
python caption_images.py \
    --dataset_path /path/to/ai-toolkit/dataset/my_subject \
    --trigger_word TOK \
    --model blip-base

# This creates .txt files for each image
```

#### Option B: Manual Captioning

Create a `.txt` file for each image with the same name:

**Example:**
- Image: `image_001.jpg`
- Caption file: `image_001.txt`
- Content: `TOK person wearing a red shirt, smiling, outdoors`

**Caption Guidelines:**
```
Good captions:
✅ "TOK woman with long brown hair, professional headshot"
✅ "TOK dog running in a park, sunny day, grass"
✅ "TOK in cyberpunk style, neon lights, detailed"

Bad captions:
❌ "TOK" (too short, no context)
❌ "a person" (no trigger word)
❌ "image123.jpg" (not descriptive)
```

#### Option C: Use ai-toolkit's Built-in Captioning

```bash
# ai-toolkit can generate captions during training
# See config below - set "caption_ext": "auto"
```

**Trigger Word Guidelines:**
- Use a unique, uncommon word (not "man", "dog", "style")
- Examples: `TOK`, `SKSDOG`, `OHWX`, `XYZ123`
- Use consistently in ALL captions
- UPPERCASE by convention

---

### Step 8: Validate Your Dataset

```bash
# Quick check
ls dataset/my_subject/*.txt | wc -l
ls dataset/my_subject/*.jpg | wc -l
# Numbers should match (one caption per image)

# Verify a caption
cat dataset/my_subject/image_001.txt
# Should show: "TOK [description]"

# Or use our validation tool
cd /path/to/flux-lora-training
python prepare_dataset.py \
    --dataset_path /path/to/ai-toolkit/dataset/my_subject \
    --trigger_word TOK
```

---

## ⚙️ Configure Training

### Step 9: Create Training Configuration (10 minutes)

```bash
# Back to ai-toolkit directory
cd /path/to/ai-toolkit

# Create your config file
nano config/my_lora_training.yaml
# Or use any text editor
```

#### Basic Configuration (12GB VRAM)

````yaml
---
job: extension
config:
  # Model configuration
  name: my_flux_lora_v1
  process:
    - type: sd_trainer
      training_folder: output
      
      # Device settings
      device: cuda:0
      
      # Trigger word
      trigger_word: TOK
      
      # Network settings (LoRA)
      network:
        type: lora
        linear: 16          # LoRA rank
        linear_alpha: 16    # LoRA alpha
      
      # Save settings
      save:
        dtype: float16
        save_every: 200     # Save checkpoint every N steps
        max_step_saves_to_keep: 3
      
      # Dataset settings
      datasets:
        - folder_path: dataset/my_subject
          caption_ext: txt
          caption_dropout_rate: 0.05
          shuffle_tokens: false
          cache_latents_to_disk: true
          resolution: 
            - 512
            - 768
            - 1024
      
      # Training parameters
      train:
        batch_size: 1
        steps: 1000         # Total training steps
        gradient_accumulation_steps: 1
        train_unet: true
        train_text_encoder: false
        
        # Learning rate
        lr: 0.0001          # 1e-4
        
        # Optimizer
        optimizer: adamw8bit
        lr_scheduler: constant
        
        # Training settings
        max_denoising_steps: 1000
        
        # Memory optimization
        gradient_checkpointing: true
        noise_scheduler: flowmatch  # For Flux
        
        # dtype
        dtype: bf16
      
      # Model to train on
      model:
        name_or_path: black-forest-labs/FLUX.1-dev
        is_flux: true
        quantize: true      # Saves VRAM
      
      # Sample generation during training (optional)
      sample:
        sampler: flowmatch
        sample_every: 200   # Generate sample every N steps
        width: 1024
        height: 1024
        prompts:
          - "TOK person in a superhero costume"
          - "TOK in anime style"
          - "portrait of TOK, professional photo"
        neg: "blurry, low quality"
        seed: 42
        walk_seed: true
        guidance_scale: 4.0
        sample_steps: 20

meta:
  name: my_flux_lora
  version: '1.0'
````

#### Advanced Configuration (24GB VRAM)

````yaml
---
job: extension
config:
  name: my_flux_lora_v1_advanced
  process:
    - type: sd_trainer
      training_folder: output
      device: cuda:0
      trigger_word: TOK
      
      network:
        type: lora
        linear: 32          # Higher rank for more capacity
        linear_alpha: 32
      
      save:
        dtype: float16
        save_every: 250
        max_step_saves_to_keep: 5
      
      datasets:
        - folder_path: dataset/my_subject
          caption_ext: txt
          caption_dropout_rate: 0.1
          shuffle_tokens: false
          cache_latents_to_disk: true
          resolution: 
            - 512
            - 768
            - 1024
            - 1536
      
      train:
        batch_size: 2       # Larger batch
        steps: 1500         # More steps
        gradient_accumulation_steps: 2
        train_unet: true
        train_text_encoder: false
        lr: 0.0001
        optimizer: adamw8bit
        lr_scheduler: constant
        max_denoising_steps: 1000
        gradient_checkpointing: true
        noise_scheduler: flowmatch
        dtype: bf16
      
      model:
        name_or_path: black-forest-labs/FLUX.1-dev
        is_flux: true
        quantize: true
      
      sample:
        sampler: flowmatch
        sample_every: 250
        width: 1024
        height: 1024
        prompts:
          - "TOK person wearing a spacesuit on mars"
          - "TOK as a renaissance painting"
          - "close up portrait of TOK, dramatic lighting"
          - "TOK in cyberpunk style, neon city background"
        neg: "blurry, low quality, distorted"
        seed: 42
        walk_seed: true
        guidance_scale: 3.5
        sample_steps: 28

meta:
  name: my_flux_lora_advanced
  version: '1.0'
````

#### Configuration for Limited VRAM (8-10GB)

````yaml
---
job: extension
config:
  name: my_flux_lora_v1_lowvram
  process:
    - type: sd_trainer
      training_folder: output
      device: cuda:0
      trigger_word: TOK
      
      network:
        type: lora
        linear: 8           # Lower rank
        linear_alpha: 8
      
      save:
        dtype: float16
        save_every: 200
        max_step_saves_to_keep: 2
      
      datasets:
        - folder_path: dataset/my_subject
          caption_ext: txt
          caption_dropout_rate: 0.05
          shuffle_tokens: false
          cache_latents_to_disk: true
          resolution: 
            - 512
            - 768           # Lower resolutions only
      
      train:
        batch_size: 1
        steps: 800
        gradient_accumulation_steps: 4  # Simulate larger batch
        train_unet: true
        train_text_encoder: false
        lr: 0.0001
        optimizer: adamw8bit
        lr_scheduler: constant
        max_denoising_steps: 1000
        gradient_checkpointing: true
        noise_scheduler: flowmatch
        dtype: bf16
      
      model:
        name_or_path: black-forest-labs/FLUX.1-dev
        is_flux: true
        quantize: true
      
      sample:
        sampler: flowmatch
        sample_every: 200
        width: 768          # Lower resolution samples
        height: 768
        prompts:
          - "TOK person smiling"
        neg: "blurry, low quality"
        seed: 42
        guidance_scale: 4.0
        sample_steps: 20

meta:
  name: my_flux_lora_lowvram
  version: '1.0'
````

**Save your config as:** `config/my_lora_training.yaml`

---

## 🎓 Training Configuration Explained

### Key Parameters to Understand

#### LoRA Settings
```yaml
linear: 16        # LoRA rank (4, 8, 16, 32, 64)
linear_alpha: 16  # Usually same as rank
```
- **Lower rank (8)**: Less VRAM, faster, less capacity
- **Higher rank (32)**: More VRAM, slower, more capacity
- **Sweet spot**: 16 for most cases

#### Training Steps
```yaml
steps: 1000
```
- **Rule of thumb**: 40-80 steps per training image
- 15 images → 600-1200 steps
- 30 images → 1200-2400 steps

#### Learning Rate
```yaml
lr: 0.0001  # 1e-4
```
- **Too high** (>5e-4): Unstable, artifacts
- **Too low** (<5e-5): Slow learning, underfitting
- **Start with**: 1e-4

#### Batch Size
```yaml
batch_size: 1
gradient_accumulation_steps: 1
```
- **Effective batch** = batch_size × gradient_accumulation_steps
- Use gradient_accumulation for larger effective batch without VRAM

#### Resolution
```yaml
resolution: 
  - 512
  - 768
  - 1024
```
- Train on multiple resolutions for flexibility
- Higher = better quality but more VRAM

---

## 🏋️ Start Training

### Step 10: Launch Training (2 minutes setup, then wait)

```bash
# Make sure you're in ai-toolkit directory
cd /path/to/ai-toolkit

# Activate virtual environment
source venv/bin/activate  # or venv\Scripts\activate on Windows

# Start training
python run.py config/my_lora_training.yaml

# Or with custom output name
python run.py config/my_lora_training.yaml
```

**What You'll See:**
```
Loading model: black-forest-labs/FLUX.1-dev
Quantizing model...
Loading dataset from dataset/my_subject
Found 20 images
Caching latents...
Starting training...

Step 1/1000 - Loss: 0.234
Step 10/1000 - Loss: 0.198
Step 50/1000 - Loss: 0.142
Step 100/1000 - Loss: 0.089
Step 200/1000 - Loss: 0.045 - Saved checkpoint
...
```

### Monitor Training

#### Watch GPU Usage
```bash
# In another terminal
watch -n 1 nvidia-smi

# Should show:
# - GPU utilization: 90-100%
# - Memory usage: 10-20GB (depending on config)
# - Temperature: <85°C ideal
```

#### Check Loss Values
- **Loss decreasing smoothly** = good learning
- **Loss stuck/increasing** = problem (see troubleshooting)
- **Target final loss**: 0.02-0.06 (varies)

#### Sample Images
```bash
# Check sample images during training
ls output/my_flux_lora_v1/samples/

# Images generated at step 200, 400, 600, etc.
```

### Training Time Estimates

| GPU | 1000 Steps | 2000 Steps |
|-----|-----------|-----------|
| RTX 3060 12GB | 3-4 hours | 6-8 hours |
| RTX 4070 12GB | 2-3 hours | 4-6 hours |
| RTX 4080 16GB | 1.5-2 hours | 3-4 hours |
| RTX 4090 24GB | 45-60 min | 1.5-2 hours |

**Tip**: You can safely leave training running overnight or while you work.

---

## ✅ Testing Your LoRA

### Step 11: Find Your Trained Model

```bash
# Your LoRA is saved in:
ls output/my_flux_lora_v1/

# Look for:
# - my_flux_lora_v1.safetensors (final model)
# - my_flux_lora_v1_000000200.safetensors (checkpoint at step 200)
# - my_flux_lora_v1_000000400.safetensors (checkpoint at step 400)
# etc.
```

---

### Step 12: Test with Our Toolkit

```bash
# Go to our flux-lora-training directory
cd /path/to/flux-lora-training

# Activate environment
source venv/bin/activate

# Test the final model
python test_lora.py \
    --lora_path /path/to/ai-toolkit/output/my_flux_lora_v1/my_flux_lora_v1.safetensors \
    --prompt "TOK person as an astronaut in space, detailed, high quality" \
    --negative_prompt "blurry, low quality, distorted" \
    --num_inference_steps 30 \
    --guidance_scale 3.5 \
    --lora_scale 1.0 \
    --output result_astronaut.png \
    --seed 42
```

**Test Different Prompts:**
```bash
# Superhero
python test_lora.py \
    --lora_path /path/to/ai-toolkit/output/my_flux_lora_v1/my_flux_lora_v1.safetensors \
    --prompt "TOK as a superhero, cape flowing, city background" \
    --output superhero.png

# Anime style
python test_lora.py \
    --lora_path /path/to/ai-toolkit/output/my_flux_lora_v1/my_flux_lora_v1.safetensors \
    --prompt "TOK in anime style, vibrant colors, detailed eyes" \
    --output anime.png

# Portrait
python test_lora.py \
    --lora_path /path/to/ai-toolkit/output/my_flux_lora_v1/my_flux_lora_v1.safetensors \
    --prompt "professional portrait of TOK, studio lighting, 4k" \
    --output portrait.png
```

---

### Step 13: Test with ai-toolkit (Alternative)

ai-toolkit has built-in inference:

Create test config `config/test_lora.yaml`:
````yaml
---
job: extension
config:
  name: test_my_lora
  process:
    - type: sd_trainer
      training_folder: output/test_results
      device: cuda:0
      
      network:
        type: lora
        linear: 16
        linear_alpha: 16
      
      model:
        name_or_path: black-forest-labs/FLUX.1-dev
        is_flux: true
      
      sample:
        sampler: flowmatch
        width: 1024
        height: 1024
        prompts:
          - "TOK person in a spacesuit on mars"
          - "TOK as a watercolor painting"
          - "portrait of TOK, professional photography"
        neg: "blurry, low quality"
        seed: 42
        guidance_scale: 3.5
        sample_steps: 28
        
      # Load your trained LoRA
      network_path: output/my_flux_lora_v1/my_flux_lora_v1.safetensors

meta:
  name: test_lora
  version: '1.0'
````

Run:
```bash
python run.py config/test_lora.yaml
```

---

## 🎨 Adjusting LoRA Strength

### LoRA Scale Parameter

```bash
# Subtle effect (0.5-0.7)
python test_lora.py --lora_scale 0.6 --prompt "TOK in a garden"

# Normal effect (0.8-1.0)
python test_lora.py --lora_scale 1.0 --prompt "TOK in a garden"

# Strong effect (1.2-1.5)
python test_lora.py --lora_scale 1.3 --prompt "TOK in a garden"
```

**When to adjust:**
- Effect too weak → Increase to 1.2-1.5
- Effect too strong/artifacts → Decrease to 0.6-0.8
- Natural results → Keep at 1.0

---

## 🔧 Troubleshooting Training

### Issue: CUDA Out of Memory

**Error:**
```
RuntimeError: CUDA out of memory
```

**Solutions:**

1. **Reduce batch size** (config):
```yaml
batch_size: 1
gradient_accumulation_steps: 4
```

2. **Lower LoRA rank**:
```yaml
linear: 8
linear_alpha: 8
```

3. **Reduce resolution**:
```yaml
resolution:
  - 512
  - 768
```

4. **Disable samples during training**:
```yaml
# Comment out or remove sample: section
```

5. **Close other programs** using GPU

---

### Issue: Loss Not Decreasing

**Symptoms:**
```
Step 500/1000 - Loss: 0.234
Step 600/1000 - Loss: 0.231
Step 700/1000 - Loss: 0.229
```

**Solutions:**

1. **Increase learning rate**:
```yaml
lr: 0.0002  # or 0.0005
```

2. **Check captions**:
```bash
# Verify all captions have trigger word
grep -L "TOK" dataset/my_subject/*.txt
# Should return nothing
```

3. **Train longer**:
```yaml
steps: 2000
```

4. **Check dataset quality**:
- Are images high quality?
- Enough variety?
- Consistent subject?

---

### Issue: Overfitting

**Symptoms:**
- Loss goes to nearly 0
- Generated images look exactly like training images
- No variation in outputs

**Solutions:**

1. **Reduce training steps**:
```yaml
steps: 500  # Train less
```

2. **Lower learning rate**:
```yaml
lr: 0.00005  # 5e-5
```

3. **Add more variety** to dataset

4. **Use earlier checkpoint**:
```bash
# Test checkpoint from step 400 instead of final
python test_lora.py --lora_path output/my_flux_lora_v1/my_flux_lora_v1_000000400.safetensors
```

5. **Reduce LoRA scale during inference**:
```bash
python test_lora.py --lora_scale 0.7
```

---

### Issue: Poor Quality Outputs

**Solutions:**

1. **Better training data**:
- Use higher resolution images
- More variety in dataset
- Remove low quality images

2. **Increase LoRA rank**:
```yaml
linear: 32
linear_alpha: 32
```

3. **Train longer**:
```yaml
steps: 2000
```

4. **Better prompts**:
```bash
# Detailed, specific prompts
python test_lora.py --prompt "professional studio portrait of TOK, dramatic lighting, high detail, 8k quality"
```

5. **Adjust inference parameters**:
```bash
python test_lora.py \
    --num_inference_steps 40 \
    --guidance_scale 4.0 \
    --lora_scale 1.0
```

---

### Issue: Training Crashes

**Solutions:**

1. **Check disk space**:
```bash
df -h
# Need at least 20GB free
```

2. **Monitor GPU temperature**:
```bash
nvidia-smi -l 1
# Should be <85°C
```

3. **Reduce settings** (see "CUDA Out of Memory" above)

4. **Check logs**:
```bash
# Look for error messages in terminal output
```

---

## 📊 Interpreting Results

### Good Training Signs ✅
- Loss decreases smoothly
- Loss reaches 0.02-0.06
- GPU utilization 90%+
- Sample images show progression
- Generated images feature your subject
- Outputs vary based on prompts

### Warning Signs ⚠️
- Loss stuck >0.1
- GPU utilization <50%
- No improvement in samples
- Out of memory errors
- Very fast convergence (overfitting)

### Bad Training Signs ❌
- Loss increases
- Training crashes repeatedly
- Sample images are blank/corrupted
- No subject recognition in outputs
- All outputs look identical

---

## 🎓 Advanced Tips

### Training Multiple LoRAs

```bash
# Train LoRA 1
python run.py config/character1_lora.yaml

# Train LoRA 2
python run.py config/character2_lora.yaml

# Use different trigger words: TOK1, TOK2, etc.
```

### Resume Training

```bash
# If training was interrupted, it auto-resumes from last checkpoint
python run.py config/my_lora_training.yaml

# ai-toolkit automatically detects existing checkpoints
```

### Experiment with Parameters

Create variations:
```yaml
# Fast iteration (testing)
steps: 300
batch_size: 2

# Quality training (final)
steps: 2000
linear: 32
```

### Dataset Augmentation

ai-toolkit supports automatic augmentation:
```yaml
datasets:
  - folder_path: dataset/my_subject
    caption_ext: txt
    shuffle_tokens: false
    random_crop: true       # Random cropping
    buckets:                # Multi-aspect ratio
      - [512, 768]
      - [768, 512]
      - [1024, 1024]
```

---

## 📤 Sharing Your LoRA

### Upload to Hugging Face

```bash
# Install huggingface-hub
pip install huggingface-hub

# Login
huggingface-cli login

# Create repository on HF website first
# Then upload
python -c "
from huggingface_hub import HfApi
api = HfApi()
api.upload_file(
    path_or_fileobj='output/my_flux_lora_v1/my_flux_lora_v1.safetensors',
    path_in_repo='my_flux_lora_v1.safetensors',
    repo_id='your-username/your-lora-name',
    repo_type='model',
)
"
```

### Upload to CivitAI

1. Go to https://civitai.com
2. Create account
3. Click "Upload Model"
4. Upload your `.safetensors` file
5. Add sample images and description

---

## 📚 Complete Workflow Summary

```
1. Install ai-toolkit ✅
   └─> git clone, create venv, install requirements

2. Prepare Dataset ✅
   └─> 15-30 images, high quality, consistent subject

3. Caption Images ✅
   └─> Use our tool or manual, include trigger word

4. Create Config ✅
   └─> Copy template, adjust for your VRAM

5. Start Training ✅
   └─> python run.py config/your_config.yaml

6. Monitor Progress ✅
   └─> Watch loss, check samples, monitor GPU

7. Test Results ✅
   └─> Use test_lora.py with different prompts

8. Iterate if Needed ✅
   └─> Adjust config, retrain if unsatisfied

9. Share (Optional) ✅
   └─> Upload to HuggingFace or CivitAI
```

---

## 🎯 Quick Reference Commands

```bash
# Setup
git clone https://github.com/ostris/ai-toolkit.git
cd ai-toolkit
python3.10 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
huggingface-cli login

# Prepare
mkdir -p dataset/my_subject
cp /path/to/images/* dataset/my_subject/

# Caption (using our tool)
cd /path/to/flux-lora-training
python caption_images.py --dataset_path /path/to/ai-toolkit/dataset/my_subject --trigger_word TOK

# Train
cd /path/to/ai-toolkit
python run.py config/my_lora_training.yaml

# Monitor
nvidia-smi -l 1

# Test (using our tool)
cd /path/to/flux-lora-training
python test_lora.py \
    --lora_path /path/to/ai-toolkit/output/my_lora/my_lora.safetensors \
    --prompt "TOK as an astronaut" \
    --output test.png
```

---

## 🌟 Success Checklist

Before starting training:
- [ ] GPU has 12GB+ VRAM
- [ ] Python 3.10 installed
- [ ] ai-toolkit cloned and dependencies installed
- [ ] Logged into Hugging Face
- [ ] FLUX.1-dev access granted
- [ ] 15-30 high quality images collected
- [ ] All images in dataset folder
- [ ] All images have caption files
- [ ] All captions include trigger word
- [ ] Config file created and customized
- [ ] At least 30GB free disk space
- [ ] Other GPU programs closed

During training:
- [ ] Loss is decreasing
- [ ] GPU utilization >80%
- [ ] No out of memory errors
- [ ] Sample images show progress

After training:
- [ ] Test with multiple prompts
- [ ] Adjust LoRA scale if needed
- [ ] Compare checkpoints
- [ ] Save best version

---

## 🎉 Congratulations!

You now know how to:
✅ Install and setup ai-toolkit
✅ Prepare high-quality datasets
✅ Configure Flux LoRA training
✅ Train and monitor LoRA models
✅ Test and iterate on results
✅ Troubleshoot common issues

**You're ready to train professional-quality Flux LoRAs!** 🚀

---

## 📖 Additional Resources

- **ai-toolkit GitHub**: https://github.com/ostris/ai-toolkit
- **ai-toolkit Examples**: Check `config/examples/` folder
- **Flux Model**: https://huggingface.co/black-forest-labs/FLUX.1-dev
- **Our Toolkit**: Use for dataset prep and testing
- **Community**: 
  - r/StableDiffusion on Reddit
  - Hugging Face Discord
  - CivitAI Community

**Happy training!** 🎨✨
