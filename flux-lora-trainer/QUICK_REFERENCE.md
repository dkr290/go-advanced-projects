# Flux LoRA Training - Quick Reference

## 🚀 Common Commands

### Setup
```bash
bash setup.sh                    # Initial setup
source venv/bin/activate         # Activate environment
pip install -r requirements.txt  # Install dependencies
huggingface-cli login           # Login to Hugging Face
```

### Dataset Preparation
```bash
# Interactive setup
python quick_start.py

# Generate captions
python caption_images.py \
    --dataset_path dataset/my_subject \
    --trigger_word MYWORD \
    --model blip-base

# Validate dataset
python prepare_dataset.py \
    --dataset_path dataset/my_subject \
    --trigger_word MYWORD
```

### Training
```bash
# Using framework (after completing implementation)
python train_lora.py --config config.yaml

# Or use ai-toolkit (recommended)
git clone https://github.com/ostris/ai-toolkit
cd ai-toolkit
# ... follow their guide
```

### Testing
```bash
# Basic test
python test_lora.py \
    --lora_path outputs/my_lora/final.safetensors \
    --prompt "MYWORD as an astronaut" \
    --output result.png

# Advanced test
python test_lora.py \
    --lora_path outputs/my_lora/final.safetensors \
    --prompt "MYWORD in a cyberpunk city" \
    --negative_prompt "blurry, low quality" \
    --num_inference_steps 50 \
    --guidance_scale 7.5 \
    --lora_scale 0.8 \
    --seed 42
```

## ⚙️ Key Configuration Parameters

### Memory Settings (for low VRAM)
```yaml
train_batch_size: 1
gradient_accumulation_steps: 8
gradient_checkpointing: true
use_8bit_adam: true
mixed_precision: "bf16"
lora_rank: 8
```

### Quality Settings (for better results)
```yaml
max_train_steps: 2000
learning_rate: 1.0e-4
lora_rank: 32
resolution: 1024
```

### Speed Settings (for faster training)
```yaml
max_train_steps: 500
train_batch_size: 4
gradient_accumulation_steps: 1
resolution: 768
```

## 📊 GPU Memory Requirements

| Config | VRAM | Training Time (1000 steps) |
|--------|------|---------------------------|
| Minimal (rank 8, batch 1) | 10GB | 4-6 hours |
| Standard (rank 16, batch 1) | 12GB | 2-4 hours |
| Optimal (rank 32, batch 2) | 20GB | 1-2 hours |
| Maximum (rank 64, batch 4) | 24GB+ | 30-60 min |

## 🎯 Recommended Settings by Dataset Size

### Small Dataset (10-15 images)
```yaml
max_train_steps: 500
learning_rate: 1.0e-4
lora_rank: 16
```

### Medium Dataset (20-30 images)
```yaml
max_train_steps: 1000
learning_rate: 1.0e-4
lora_rank: 16
```

### Large Dataset (40+ images)
```yaml
max_train_steps: 2000
learning_rate: 5.0e-5
lora_rank: 32
```

## 🐛 Quick Troubleshooting

### Out of Memory
```yaml
# Add to config.yaml:
train_batch_size: 1
gradient_accumulation_steps: 16
gradient_checkpointing: true
use_8bit_adam: true
lora_rank: 8
```

### Loss Not Decreasing
```yaml
# Try:
learning_rate: 5.0e-4  # Increase
max_train_steps: 2000  # Train longer
# Check captions have trigger word
```

### Overfitting (loss too low)
```yaml
# Try:
max_train_steps: 500   # Train less
learning_rate: 5.0e-5  # Decrease
```

### Low Quality Output
```bash
# Test with:
--num_inference_steps 50
--guidance_scale 7.5
--lora_scale 1.0
# Add detailed prompt
# Add negative prompt
```

## 📁 Directory Structure
```
flux-lora-training/
├── dataset/
│   └── my_subject/
│       ├── image_001.jpg
│       ├── image_001.txt
│       └── ...
├── outputs/
│   └── my_lora/
│       ├── checkpoint-200/
│       └── final.safetensors
├── cache/              # Model downloads
├── venv/              # Virtual environment
├── config.yaml        # Main config
├── train_lora.py      # Training script
├── test_lora.py       # Testing script
└── README.md          # Documentation
```

## 📋 Checklist Before Training

- [ ] Python 3.10+ installed
- [ ] NVIDIA GPU with 12GB+ VRAM
- [ ] CUDA installed and working
- [ ] Virtual environment created and activated
- [ ] Requirements installed
- [ ] Logged in to Hugging Face
- [ ] Flux model license accepted
- [ ] 10+ training images collected
- [ ] Images in dataset directory
- [ ] Captions generated (.txt files)
- [ ] Dataset validated (no errors)
- [ ] Config file edited
- [ ] Trigger word set
- [ ] Enough disk space (50GB+)

## 🔗 Important Links

- Flux Model: https://huggingface.co/black-forest-labs/FLUX.1-dev
- Diffusers Docs: https://huggingface.co/docs/diffusers
- PEFT Docs: https://huggingface.co/docs/peft
- ai-toolkit: https://github.com/ostris/ai-toolkit
- SimpleTuner: https://github.com/bghira/SimpleTuner

## 💡 Pro Tips

1. **Start Small**: Test with 10 images and 500 steps first
2. **Quality > Quantity**: Better to have 15 great images than 50 mediocre ones
3. **Monitor Loss**: Should decrease smoothly, typical final loss: 0.02-0.05
4. **Test Checkpoints**: Test every checkpoint to catch overfitting
5. **Use Trigger Word**: Always include in prompts and captions
6. **Adjust LoRA Scale**: Try 0.6-1.2 range during inference
7. **Backup**: Save your best checkpoints externally
8. **Experiment**: Try different learning rates and ranks

## 🎓 Learning Path

1. **Day 1**: Setup environment, prepare 10 images
2. **Day 2**: Generate captions, validate dataset
3. **Day 3**: Train first model (500 steps)
4. **Day 4**: Test and iterate
5. **Day 5**: Full training with tuned parameters

## ⌨️ Keyboard Shortcuts (during training)

- `Ctrl+C`: Stop training (saves checkpoint)
- Monitor: `nvidia-smi -l 1` (watch GPU)
- Logs: `tensorboard --logdir outputs/my_lora/logs`

## 📞 Getting Help

1. Check TROUBLESHOOTING.md
2. Review error messages carefully
3. Run prepare_dataset.py to validate
4. Check GitHub issues for similar problems
5. Ask in Hugging Face forums

## 🎉 Success Indicators

✅ Loss decreasing smoothly
✅ GPU utilization 90%+
✅ Checkpoints saving
✅ Validation shows trigger word effect
✅ Generated images look like subject
✅ Effect adjustable with lora_scale

---

**Quick Start**: `python quick_start.py`
**Full Guide**: See README.md
**Problems?**: See TROUBLESHOOTING.md
