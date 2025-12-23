# Complete Flux LoRA Training - Hybrid Workflow

## 🎯 Best of Both Worlds

This guide combines **our toolkit** (for preparation and testing) with **ai-toolkit** (for training) - the optimal workflow for Flux LoRA training.

---

## 📋 Why This Hybrid Approach?

### Our Toolkit Strengths
✅ Easy dataset preparation
✅ Automatic captioning with BLIP
✅ Dataset validation and analysis
✅ Flexible testing and inference
✅ Well-documented configuration

### ai-toolkit Strengths
✅ Battle-tested Flux training
✅ Optimized for performance
✅ Active development and support
✅ Proven results in production
✅ Built-in Flux-specific features

### Combined = Perfect Workflow 🎉

---

## 🚀 Complete Workflow: Start to Finish

### Phase 1: Environment Setup (15 minutes)

#### 1.1 Setup Our Toolkit
```bash
# Clone and setup our toolkit for data preparation
git clone [your-toolkit-repo]
cd flux-lora-training
bash setup.sh
source venv/bin/activate
```

#### 1.2 Setup ai-toolkit
```bash
# In a separate directory
cd ~
git clone https://github.com/ostris/ai-toolkit.git
cd ai-toolkit
python3.10 -m venv venv
source venv/bin/activate
pip install torch torchvision --index-url https://download.pytorch.org/whl/cu121
pip install -r requirements.txt
```

#### 1.3 Hugging Face Login
```bash
# Do this in both environments
huggingface-cli login
# Paste your token
# Accept FLUX.1-dev license at https://huggingface.co/black-forest-labs/FLUX.1-dev
```

---

### Phase 2: Dataset Preparation with Our Toolkit (30 minutes)

#### 2.1 Organize Images
```bash
cd flux-lora-training

# Create dataset directory
mkdir -p dataset/my_subject

# Copy your images (15-30 high quality images)
cp /path/to/your/photos/*.jpg dataset/my_subject/
```

#### 2.2 Auto-Generate Captions
```bash
# Activate our toolkit environment
source venv/bin/activate

# Generate captions with BLIP
python caption_images.py \
    --dataset_path dataset/my_subject \
    --trigger_word TOK \
    --model blip-base

# This creates a .txt file for each image
```

**Output:**
```
Found 25 images
Captioning...
✓ Saved to: dataset/my_subject/image_001.txt
✓ Saved to: dataset/my_subject/image_002.txt
...
✓ Captioning complete! Generated 25 captions
```

#### 2.3 Validate Dataset
```bash
# Check dataset quality
python prepare_dataset.py \
    --dataset_path dataset/my_subject \
    --trigger_word TOK
```

**Review the report:**
- ✅ All images valid
- ✅ All captions present
- ✅ Trigger word in all captions
- ✅ Resolution requirements met

**Fix any issues before proceeding!**

---

### Phase 3: Copy Dataset to ai-toolkit (2 minutes)

```bash
# Copy prepared dataset to ai-toolkit
cp -r dataset/my_subject ~/ai-toolkit/dataset/

# Verify
ls ~/ai-toolkit/dataset/my_subject/
# Should show: image files + txt files
```

---

### Phase 4: Configure Training (10 minutes)

#### 4.1 Create ai-toolkit Config

```bash
cd ~/ai-toolkit

# Create config file
nano config/my_training.yaml
```

#### 4.2 Configuration Template

Choose based on your GPU:

**For 12GB VRAM (RTX 3060, 4060):**
````yaml
---
job: extension
config:
  name: my_flux_lora_v1
  process:
    - type: sd_trainer
      training_folder: output
      device: cuda:0
      trigger_word: TOK
      
      network:
        type: lora
        linear: 16
        linear_alpha: 16
      
      save:
        dtype: float16
        save_every: 200
        max_step_saves_to_keep: 3
      
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
      
      train:
        batch_size: 1
        steps: 1000
        gradient_accumulation_steps: 4
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
        width: 1024
        height: 1024
        prompts:
          - "TOK person in a spacesuit"
          - "TOK in anime style"
        neg: "blurry, low quality"
        seed: 42
        guidance_scale: 3.5
        sample_steps: 20

meta:
  name: my_flux_lora
  version: '1.0'
````

**For 24GB VRAM (RTX 4090, A5000):**
````yaml
---
job: extension
config:
  name: my_flux_lora_v1
  process:
    - type: sd_trainer
      training_folder: output
      device: cuda:0
      trigger_word: TOK
      
      network:
        type: lora
        linear: 32          # Higher rank
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
        steps: 1500
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
          - "TOK person as an astronaut on mars"
          - "TOK as a renaissance painting"
          - "portrait of TOK, professional photography"
        neg: "blurry, low quality"
        seed: 42
        guidance_scale: 3.5
        sample_steps: 28

meta:
  name: my_flux_lora
  version: '1.0'
````

---

### Phase 5: Training with ai-toolkit (1-3 hours)

```bash
cd ~/ai-toolkit
source venv/bin/activate

# Start training
python run.py config/my_training.yaml
```

**Monitor in another terminal:**
```bash
# GPU usage
watch -n 1 nvidia-smi

# Check samples as they generate
ls -lh output/my_flux_lora_v1/samples/
```

**What to expect:**
```
Loading model: black-forest-labs/FLUX.1-dev
Quantizing model...
Found 25 images in dataset/my_subject
Caching latents... [25/25]
Starting training...

Step 1/1000 - Loss: 0.234
Step 50/1000 - Loss: 0.156
Step 100/1000 - Loss: 0.098
Step 200/1000 - Loss: 0.067 [Checkpoint saved] [Sample generated]
Step 400/1000 - Loss: 0.045 [Checkpoint saved] [Sample generated]
...
Step 1000/1000 - Loss: 0.028 [Final model saved]

Training complete!
```

**Training time:**
- RTX 3060: 2-3 hours
- RTX 4070: 1-2 hours
- RTX 4090: 45-60 minutes

---

### Phase 6: Testing with Our Toolkit (5 minutes)

```bash
# Switch back to our toolkit
cd ~/flux-lora-training
source venv/bin/activate

# Test the final model
python test_lora.py \
    --lora_path ~/ai-toolkit/output/my_flux_lora_v1/my_flux_lora_v1.safetensors \
    --prompt "TOK person as an astronaut floating in space, detailed, high quality, 8k" \
    --negative_prompt "blurry, low quality, distorted, ugly" \
    --num_inference_steps 30 \
    --guidance_scale 3.5 \
    --lora_scale 1.0 \
    --width 1024 \
    --height 1024 \
    --output results/astronaut.png \
    --seed 42
```

**Try multiple prompts:**
```bash
# Batch testing script
cat > test_all.sh << 'EOF'
#!/bin/bash

LORA="~/ai-toolkit/output/my_flux_lora_v1/my_flux_lora_v1.safetensors"

python test_lora.py --lora_path $LORA \
    --prompt "TOK as a superhero, cape flowing" \
    --output results/superhero.png

python test_lora.py --lora_path $LORA \
    --prompt "TOK in anime style, vibrant colors" \
    --output results/anime.png

python test_lora.py --lora_path $LORA \
    --prompt "professional portrait of TOK, studio lighting" \
    --output results/portrait.png

python test_lora.py --lora_path $LORA \
    --prompt "TOK as a watercolor painting" \
    --output results/painting.png

python test_lora.py --lora_path $LORA \
    --prompt "TOK in cyberpunk style, neon city" \
    --output results/cyberpunk.png
EOF

chmod +x test_all.sh
./test_all.sh
```

**Check results:**
```bash
ls -lh results/
# Should show generated images
```

---

### Phase 7: Evaluation and Iteration

#### 7.1 Evaluate Results

**Good results checklist:**
- [ ] Subject is clearly recognizable
- [ ] Prompt is followed accurately
- [ ] Image quality is high
- [ ] Style variations work well
- [ ] No major artifacts

#### 7.2 If Results Are Not Satisfactory

**Problem: Subject not well learned**
```bash
# Solution 1: Train longer
# Edit config: steps: 1500 or 2000
# Retrain

# Solution 2: Test earlier checkpoint
python test_lora.py \
    --lora_path ~/ai-toolkit/output/my_flux_lora_v1/my_flux_lora_v1_000000800.safetensors \
    --prompt "TOK person smiling" \
    --output test_checkpoint.png
```

**Problem: Overfitted (looks exactly like training images)**
```bash
# Solution 1: Use earlier checkpoint (step 400-600)

# Solution 2: Lower LoRA scale
python test_lora.py --lora_scale 0.6 --prompt "TOK in a garden"

# Solution 3: Retrain with fewer steps
# Edit config: steps: 600
```

**Problem: Low quality outputs**
```bash
# Solution 1: Increase inference steps
python test_lora.py --num_inference_steps 50 --prompt "TOK portrait"

# Solution 2: Better prompts
python test_lora.py --prompt "professional photograph of TOK, studio lighting, high detail, 8k quality, masterpiece"

# Solution 3: Retrain with higher rank
# Edit config: linear: 32, linear_alpha: 32
```

---

## 🎓 Advanced Workflows

### Workflow A: Quick Iteration

```bash
# 1. Prepare dataset (our toolkit)
python caption_images.py --dataset_path dataset/test1 --trigger_word TOK1

# 2. Quick training (ai-toolkit, 500 steps)
python run.py config/quick_test.yaml

# 3. Fast testing (our toolkit)
python test_lora.py --lora_path output/test1/test1.safetensors --prompt "TOK1 smiling"

# 4. Adjust and repeat
```

### Workflow B: Production Quality

```bash
# 1. Prepare dataset carefully (our toolkit)
python prepare_dataset.py --dataset_path dataset/production
# Fix all issues

# 2. Train with high settings (ai-toolkit)
# Config: rank 32, steps 2000, batch 2
python run.py config/production.yaml

# 3. Test all checkpoints (our toolkit)
for ckpt in output/prod/*.safetensors; do
    python test_lora.py --lora_path $ckpt --prompt "TOK portrait" --output test_$(basename $ckpt).png
done

# 4. Select best checkpoint
```

### Workflow C: Multiple Characters

```bash
# 1. Prepare separate datasets (our toolkit)
python caption_images.py --dataset_path dataset/char1 --trigger_word CHAR1
python caption_images.py --dataset_path dataset/char2 --trigger_word CHAR2

# 2. Train separately (ai-toolkit)
python run.py config/char1_training.yaml
python run.py config/char2_training.yaml

# 3. Test combinations (our toolkit)
python test_lora.py \
    --lora_path output/char1/char1.safetensors \
    --prompt "CHAR1 and CHAR2 having coffee together"
```

---

## 📊 Comparison: Tools Overview

| Feature | Our Toolkit | ai-toolkit | Combined |
|---------|------------|-----------|----------|
| Dataset Prep | ✅ Excellent | ⚠️ Manual | ✅ Best |
| Auto-Captioning | ✅ Built-in | ❌ None | ✅ Best |
| Dataset Validation | ✅ Comprehensive | ❌ Basic | ✅ Best |
| Flux Training | ⚠️ Framework | ✅ Production | ✅ Best |
| Configuration | ✅ Well-documented | ✅ Examples | ✅ Best |
| Testing/Inference | ✅ Flexible | ⚠️ Basic | ✅ Best |
| Documentation | ✅ Extensive | ⚠️ GitHub only | ✅ Best |

---

## 🔧 Tool-Specific Commands Reference

### Our Toolkit Commands
```bash
# Activate environment
cd ~/flux-lora-training
source venv/bin/activate

# Caption images
python caption_images.py \
    --dataset_path dataset/my_subject \
    --trigger_word TOK \
    --model blip-base \
    --batch_size 4

# Validate dataset
python prepare_dataset.py \
    --dataset_path dataset/my_subject \
    --trigger_word TOK

# Test LoRA
python test_lora.py \
    --lora_path /path/to/model.safetensors \
    --prompt "your prompt here" \
    --output result.png \
    --lora_scale 1.0 \
    --num_inference_steps 30 \
    --guidance_scale 3.5
```

### ai-toolkit Commands
```bash
# Activate environment
cd ~/ai-toolkit
source venv/bin/activate

# Train
python run.py config/your_config.yaml

# Resume training (automatic)
python run.py config/your_config.yaml
# Detects existing checkpoints automatically

# Monitor
nvidia-smi -l 1
```

---

## 📝 Project Organization

Recommended directory structure:

```
~/projects/
│
├── flux-lora-training/          # Our toolkit
│   ├── venv/
│   ├── dataset/
│   │   ├── project1/
│   │   ├── project2/
│   │   └── project3/
│   ├── results/                 # Test outputs
│   └── caption_images.py
│
├── ai-toolkit/                  # Training toolkit
│   ├── venv/
│   ├── config/
│   │   ├── project1.yaml
│   │   ├── project2.yaml
│   │   └── project3.yaml
│   ├── dataset/                 # Copied from our toolkit
│   │   ├── project1/
│   │   ├── project2/
│   │   └── project3/
│   └── output/                  # Trained models
│       ├── project1/
│       ├── project2/
│       └── project3/
│
└── final_models/                # Your best models
    ├── project1_v1.safetensors
    ├── project1_v2.safetensors
    └── project2_v1.safetensors
```

---

## 🎯 Quick Start Script

Save this as `train_flux_lora.sh`:

```bash
#!/bin/bash

# Flux LoRA Training - Complete Workflow Script
# Usage: ./train_flux_lora.sh project_name trigger_word

PROJECT=$1
TRIGGER=$2

if [ -z "$PROJECT" ] || [ -z "$TRIGGER" ]; then
    echo "Usage: ./train_flux_lora.sh project_name TRIGGER_WORD"
    exit 1
fi

echo "🚀 Starting Flux LoRA Training: $PROJECT"
echo "   Trigger word: $TRIGGER"
echo ""

# Step 1: Prepare dataset
echo "📸 Step 1: Preparing dataset..."
cd ~/flux-lora-training
source venv/bin/activate
python caption_images.py --dataset_path dataset/$PROJECT --trigger_word $TRIGGER
python prepare_dataset.py --dataset_path dataset/$PROJECT --trigger_word $TRIGGER

read -p "Is dataset ready? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Fix dataset issues and run again"
    exit 1
fi

# Step 2: Copy to ai-toolkit
echo "📂 Step 2: Copying dataset to ai-toolkit..."
cp -r dataset/$PROJECT ~/ai-toolkit/dataset/

# Step 3: Start training
echo "🏋️ Step 3: Starting training..."
cd ~/ai-toolkit
source venv/bin/activate
python run.py config/${PROJECT}_training.yaml

# Step 4: Test results
echo "🎨 Step 4: Testing results..."
cd ~/flux-lora-training
source venv/bin/activate
python test_lora.py \
    --lora_path ~/ai-toolkit/output/${PROJECT}/${PROJECT}.safetensors \
    --prompt "$TRIGGER person in various styles" \
    --output results/${PROJECT}_test.png

echo "✅ Complete! Check results/ folder for output"
```

Make executable:
```bash
chmod +x train_flux_lora.sh
```

Use it:
```bash
./train_flux_lora.sh my_character TOK
```

---

## 🎉 Success!

You now have a **complete, professional workflow** for Flux LoRA training:

✅ **Easy dataset preparation** (our toolkit)
✅ **Automatic captioning** (our toolkit)
✅ **Quality validation** (our toolkit)
✅ **Professional training** (ai-toolkit)
✅ **Flexible testing** (our toolkit)
✅ **Comprehensive documentation** (both)

**This is the optimal setup for Flux LoRA training!** 🚀✨

---

## 📚 What to Read Next

1. **AI_TOOLKIT_GUIDE.md** - Detailed ai-toolkit usage
2. **QUICK_REFERENCE.md** - Command cheat sheet
3. **TROUBLESHOOTING.md** - Solutions to issues
4. **README.md** - Complete toolkit guide

**Happy training!** 🎨
