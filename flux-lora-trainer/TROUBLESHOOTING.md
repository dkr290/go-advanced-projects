# Flux LoRA Training - Troubleshooting Guide

## Common Issues and Solutions

### Installation Issues

#### Issue: "torch not found" or CUDA errors
```bash
# Solution: Install PyTorch with CUDA support
pip uninstall torch torchvision
pip install torch torchvision --index-url https://download.pytorch.org/whl/cu121
```

#### Issue: "No module named 'transformers'"
```bash
# Solution: Install all requirements
pip install -r requirements.txt
```

#### Issue: "bitsandbytes is not supported on Windows"
```bash
# Solution for Windows:
# 1. Comment out use_8bit_adam in config.yaml
# 2. Or install Windows-compatible version:
pip install bitsandbytes-windows
```

---

### Hugging Face / Download Issues

#### Issue: "401 Unauthorized" when loading model
**Cause:** Not logged in or haven't accepted model license

**Solution:**
```bash
# 1. Login to Hugging Face
huggingface-cli login

# 2. Accept model license
# Visit: https://huggingface.co/black-forest-labs/FLUX.1-dev
# Click "Agree and access repository"
```

#### Issue: "Repository not found"
**Cause:** Wrong model name or no access

**Solution:**
```bash
# Verify you can access the model:
huggingface-cli whoami
# Should show your username

# Test download:
python -c "from huggingface_hub import login; login()"
```

#### Issue: Very slow download / timeout
**Solution:**
```bash
# Use resume download:
export HF_HUB_ENABLE_HF_TRANSFER=1
pip install hf_transfer

# Or use local mirror if available
```

---

### GPU / Memory Issues

#### Issue: "CUDA out of memory"
**Solutions (try in order):**

1. **Reduce batch size** (config.yaml):
   ```yaml
   train_batch_size: 1
   ```

2. **Increase gradient accumulation**:
   ```yaml
   gradient_accumulation_steps: 8  # or 16
   ```

3. **Enable all memory optimizations**:
   ```yaml
   gradient_checkpointing: true
   use_8bit_adam: true
   mixed_precision: "bf16"
   ```

4. **Reduce LoRA rank**:
   ```yaml
   lora_rank: 8  # instead of 16
   ```

5. **Reduce resolution**:
   ```yaml
   resolution: 768  # instead of 1024
   ```

6. **Clear CUDA cache** before training:
   ```python
   import torch
   torch.cuda.empty_cache()
   ```

7. **Close other programs** using GPU

8. **Use CPU offloading** (slower but works):
   ```python
   # In training script
   model.enable_model_cpu_offload()
   ```

#### Issue: "RuntimeError: CUDA error: device-side assert triggered"
**Solution:**
```bash
# Run with CUDA debugging:
CUDA_LAUNCH_BLOCKING=1 python train_lora.py --config config.yaml

# This will show the actual error
```

#### Issue: Training is very slow on GPU
**Causes and solutions:**

1. **Using float32 instead of bf16**:
   ```yaml
   mixed_precision: "bf16"  # Make sure this is set
   ```

2. **Too many dataloader workers**:
   ```yaml
   dataloader_num_workers: 0  # Try reducing
   ```

3. **Gradient checkpointing overhead**:
   ```yaml
   gradient_checkpointing: false  # If you have enough VRAM
   ```

---

### Dataset Issues

#### Issue: "No images found in dataset"
**Solution:**
```bash
# Check your directory structure:
ls -la dataset/my_subject/

# Should show .jpg, .png files
# Make sure images are directly in the folder, not in subfolders
```

#### Issue: "No caption file found for image.jpg"
**Solutions:**

1. **Auto-generate captions**:
   ```bash
   python caption_images.py \
       --dataset_path dataset/my_subject \
       --trigger_word MYWORD
   ```

2. **Manually create .txt files**:
   ```bash
   # For image_001.jpg, create image_001.txt
   echo "MYWORD wearing a hat" > dataset/my_subject/image_001.txt
   ```

#### Issue: "Image too small: 256x256"
**Solution:**
```bash
# Resize images to at least 512x512
# Using ImageMagick:
mogrify -resize 1024x1024^ -gravity center -extent 1024x1024 dataset/my_subject/*.jpg

# Or use a Python script to resize
```

#### Issue: "Corrupted image" errors
**Solution:**
```bash
# Find corrupted images:
python prepare_dataset.py --dataset_path dataset/my_subject

# Remove or replace corrupted files
# Test image integrity:
python -c "from PIL import Image; Image.open('image.jpg').verify()"
```

---

### Training Issues

#### Issue: Loss not decreasing / stuck at high value
**Possible causes:**

1. **Learning rate too low**:
   ```yaml
   learning_rate: 5.0e-4  # Try higher
   ```

2. **Captions don't match images**:
   - Review your captions
   - Regenerate with BLIP

3. **Dataset too small**:
   - Add more images (aim for 20+)

4. **Wrong trigger word**:
   - Make sure trigger word is in ALL captions
   ```bash
   python prepare_dataset.py --dataset_path dataset/my_subject --trigger_word MYWORD
   ```

#### Issue: Loss decreasing too fast / overfitting
**Symptoms:** Loss goes to nearly 0, generated images are identical to training images

**Solutions:**

1. **Reduce training steps**:
   ```yaml
   max_train_steps: 500  # Instead of 2000
   ```

2. **Lower learning rate**:
   ```yaml
   learning_rate: 5.0e-5  # Instead of 1e-4
   ```

3. **Add more variety** to dataset

4. **Use lower LoRA scale** during inference:
   ```bash
   python test_lora.py --lora_scale 0.6  # Instead of 1.0
   ```

#### Issue: Training crashes randomly
**Solutions:**

1. **Check GPU temperature**:
   ```bash
   nvidia-smi -l 1  # Monitor GPU
   ```

2. **Reduce batch size / increase accumulation**

3. **Enable gradient clipping**:
   ```yaml
   max_grad_norm: 1.0
   ```

4. **Check disk space**:
   ```bash
   df -h
   ```

---

### Inference Issues

#### Issue: "LoRA weights not loading"
**Solutions:**

1. **Check file path**:
   ```bash
   ls -lh outputs/my_lora/final.safetensors
   ```

2. **Use correct LoRA file**:
   ```bash
   # Use .safetensors file, not the directory
   python test_lora.py --lora_path outputs/my_lora/final.safetensors
   ```

3. **Check file isn't corrupted**:
   ```python
   import safetensors.torch
   state_dict = safetensors.torch.load_file("final.safetensors")
   print(f"Loaded {len(state_dict)} tensors")
   ```

#### Issue: Generated images don't show LoRA effect
**Solutions:**

1. **Increase LoRA scale**:
   ```bash
   python test_lora.py --lora_scale 1.5  # Instead of 1.0
   ```

2. **Use trigger word in prompt**:
   ```bash
   # Make sure to include your trigger word!
   python test_lora.py --prompt "MYWORD as an astronaut"
   ```

3. **Check if model trained enough**:
   ```bash
   # Try earlier checkpoint
   python test_lora.py --lora_path outputs/my_lora/checkpoint-600/pytorch_lora_weights.safetensors
   ```

#### Issue: "Low quality images" or artifacts
**Solutions:**

1. **Increase inference steps**:
   ```bash
   python test_lora.py --num_inference_steps 50  # Instead of 20
   ```

2. **Adjust guidance scale**:
   ```bash
   python test_lora.py --guidance_scale 7.5  # Try 5.0-10.0
   ```

3. **Use better prompt**:
   ```bash
   # More detailed prompt
   python test_lora.py --prompt "a professional photo of MYWORD, high quality, detailed, 4k"
   ```

4. **Add negative prompt**:
   ```bash
   python test_lora.py --negative_prompt "blurry, low quality, distorted, ugly, bad anatomy, watermark"
   ```

---

### Performance Issues

#### Issue: Captioning is very slow
**Solutions:**

1. **Use smaller model**:
   ```bash
   python caption_images.py --model blip-base  # Instead of blip-large
   ```

2. **Reduce batch size**:
   ```bash
   python caption_images.py --batch_size 1
   ```

3. **Use GPU if available** (automatic)

#### Issue: Training time estimation way too long
**Expected times (1000 steps):**
- RTX 3060 12GB: 2-4 hours
- RTX 4070 12GB: 1.5-2.5 hours  
- RTX 4090 24GB: 30-60 minutes

**If much slower:**
1. Check GPU is actually being used:
   ```bash
   nvidia-smi
   # Should show python process using GPU
   ```

2. Check CUDA version matches PyTorch:
   ```python
   import torch
   print(f"CUDA available: {torch.cuda.is_available()}")
   print(f"CUDA version: {torch.version.cuda}")
   ```

---

### Environment Issues

#### Issue: "ModuleNotFoundError" even after pip install
**Solution:**
```bash
# Make sure you're in the virtual environment
which python
# Should show: /path/to/venv/bin/python

# Reinstall in correct environment
source venv/bin/activate  # Linux/Mac
# or
venv\Scripts\activate  # Windows

pip install -r requirements.txt
```

#### Issue: Different results each time despite same seed
**Causes:**
- CUDA nondeterminism
- Different hardware

**Solution for reproducibility:**
```python
# Add to training script
import torch
import random
import numpy as np

def set_seed(seed):
    random.seed(seed)
    np.random.seed(seed)
    torch.manual_seed(seed)
    torch.cuda.manual_seed_all(seed)
    torch.backends.cudnn.deterministic = True
    torch.backends.cudnn.benchmark = False

set_seed(42)
```

---

## Getting Help

If you're still stuck:

1. **Check logs carefully** - error messages are usually helpful

2. **Run dataset validation**:
   ```bash
   python prepare_dataset.py --dataset_path dataset/my_subject
   ```

3. **Test with minimal config**:
   ```bash
   cp config_simple.yaml config_test.yaml
   # Edit to use your dataset
   python train_lora.py --config config_test.yaml
   ```

4. **Check system resources**:
   ```bash
   nvidia-smi  # GPU
   free -h     # RAM
   df -h       # Disk space
   ```

5. **Create minimal reproduction**:
   - Test with just 5 images
   - 100 training steps
   - Simple prompts

6. **Community resources**:
   - Hugging Face forums
   - GitHub issues for diffusers/PEFT
   - Reddit: r/StableDiffusion

---

## Debug Checklist

Before asking for help, verify:

- [ ] Python 3.10+ installed
- [ ] CUDA GPU with 12GB+ VRAM
- [ ] Virtual environment activated
- [ ] Requirements installed: `pip install -r requirements.txt`
- [ ] Logged in to Hugging Face: `huggingface-cli whoami`
- [ ] Model license accepted
- [ ] Dataset has 10+ valid images
- [ ] All images have captions (.txt files)
- [ ] Captions include trigger word
- [ ] Config file is valid YAML
- [ ] Enough disk space (50GB+)
- [ ] No other programs using GPU

Include this info when asking for help!
