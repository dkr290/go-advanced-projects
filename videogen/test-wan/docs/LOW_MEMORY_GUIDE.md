# 4GB GPU Survival Guide 🎮

## ⚠️ Reality Check

**The LTX-Video model is too large for 4GB GPUs.** Here are your options:

---

## Option 1: CPU Mode (Recommended for 4GB GPU)

Use CPU instead - it's slow but will work:

```bash
# Edit .env
ENABLE_GPU=false

# Or use the provided config
cp .env.low_memory .env
# Then edit and set:
ENABLE_GPU=false
```

**Expected Performance:**
- 256x256, 16 frames: ~10-20 minutes per video
- Not practical for production, but works for testing

---

## Option 2: Use Smaller Alternative Models

Instead of LTX-Video, use lighter models:

### Recommended: ModelScope Text-to-Video

```bash
# In .env
HUGGINGFACE_MODEL_ID=damo-vilab/text-to-video-ms-1.7b
```

This model works with 4GB VRAM!

### Or: ZeroScope v2

```bash
HUGGINGFACE_MODEL_ID=cerspense/zeroscope_v2_576w
```

**I can create configurations for these models if you want!**

---

## Option 3: Try Extreme Optimization (May Still Fail)

If you REALLY want to try LTX-Video on 4GB:

```bash
# Copy the low memory config
cp .env.low_memory .env

# Edit and ensure these settings:
ENABLE_GPU=true
LOW_MEMORY_MODE=true
MAX_FRAMES=8          # Very low!
DEFAULT_WIDTH=128     # Tiny!
DEFAULT_HEIGHT=128    # Tiny!
MAX_CONCURRENT_REQUESTS=1
```

**Modify Python backend to use float16 + more optimizations:**

Add to `python_backend/server.py` after line where pipeline loads:

```python
# Ultra low memory mode
if torch.cuda.is_available():
    torch.cuda.empty_cache()
    self.pipeline.to(torch.float16)  # Half precision
    self.pipeline.enable_model_cpu_offload()  # Offload to CPU when not in use
```

**This will likely still crash with 4GB!**

---

## Option 4: Cloud GPU (Best Solution)

Use a cloud provider with larger GPU:

### Free/Cheap Options:

1. **Google Colab** (Free with T4 16GB)
   ```python
   # In Colab notebook
   !git clone <your-repo>
   !cd wan2-video-server/python_backend && pip install -r requirements.txt
   !python server.py
   ```

2. **Kaggle** (Free with P100 16GB)
   - Similar to Colab
   - 30 hours/week free GPU

3. **RunPod** (~$0.20/hour for RTX 3090)
   - Pay as you go
   - Good for testing

4. **Vast.ai** (~$0.15/hour for various GPUs)
   - Cheapest option
   - Community GPUs

---

## Recommended Configuration for 4GB

Create a new file `docker-compose.low-memory.yml`:

```yaml
version: '3.8'

services:
  wan2-video-server:
    build: .
    container_name: wan2-video-server-low-mem
    ports:
      - "8080:8080"
      - "5000:5000"
    environment:
      - ENABLE_GPU=false  # Use CPU!
      - MODEL_CACHE_DIR=/app/models
      - MAX_FRAMES=16
      - DEFAULT_WIDTH=256
      - DEFAULT_HEIGHT=256
      - LOW_MEMORY_MODE=true
    volumes:
      - ./models:/app/models
      - ./outputs:/app/outputs
    restart: unless-stopped
```

---

## What WILL Work on 4GB GPU

### Stable Diffusion (Images, not video)
- Text-to-image: ✅ Works well
- Image-to-image: ✅ Works well

### Smaller Video Models
- **ModelScope 1.7B**: ✅ Works (256x256 videos)
- **ZeroScope v2**: ✅ Works (576x320 videos)
- **AnimateDiff**: ⚠️ Might work with optimizations

---

## My Recommendations (Priority Order)

### If you want to use THIS code:

1. ✅ **Use CPU mode** for testing/learning
   - Copy `.env.low_memory` → `.env`
   - Set `ENABLE_GPU=false`
   - Accept slow speeds

2. ✅ **Switch to smaller model** (I can help configure)
   - Use ModelScope or ZeroScope
   - Actually works on 4GB
   - Reasonable quality

3. ✅ **Use cloud GPU** for production
   - Google Colab (free)
   - Vast.ai (cheap)
   - Your code runs unchanged

### If you're flexible:

4. ✅ **Use Stable Diffusion** instead
   - Generate images, not videos
   - Works perfectly on 4GB
   - I can create an image server instead

---

## Quick Test for Your 4GB GPU

Try this to see what's possible:

```bash
# Install dependencies
pip install torch torchvision diffusers transformers

# Python test script
python3 << EOF
import torch
print(f"CUDA Available: {torch.cuda.is_available()}")
print(f"GPU Memory: {torch.cuda.get_device_properties(0).total_memory / 1024**3:.1f}GB")

# Try loading a small model
from diffusers import DiffusionPipeline
pipe = DiffusionPipeline.from_pretrained(
    "damo-vilab/text-to-video-ms-1.7b",
    torch_dtype=torch.float16
)
pipe = pipe.to("cuda")
print("✅ ModelScope loaded successfully!")
EOF
```

If this works → You can use ModelScope  
If this fails → Use CPU mode or cloud GPU

---

## Want me to create a version that works on 4GB?

I can create:

1. **Configuration for ModelScope** (text-to-video, 4GB compatible)
2. **Stable Diffusion image server** (images only, works great on 4GB)
3. **Video frame interpolation server** (different approach)

Just let me know which you prefer! 🚀

---

## Summary for 4GB GPU

| Option | Works? | Speed | Quality | Cost |
|--------|--------|-------|---------|------|
| LTX-Video on 4GB GPU | ❌ No | N/A | N/A | Free |
| CPU Mode | ✅ Yes | Very Slow | Good | Free |
| ModelScope on 4GB | ✅ Yes | Medium | OK | Free |
| Cloud GPU (Colab) | ✅ Yes | Fast | Excellent | Free |
| Cloud GPU (Paid) | ✅ Yes | Fast | Excellent | ~$0.20/hr |

**My recommendation: Use Google Colab (free T4 16GB) or switch to ModelScope model.**
