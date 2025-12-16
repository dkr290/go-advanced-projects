# Stable Diffusion Image Generator

Go-based FLUX Stable Diffusion image generator with CUDA acceleration for GPU servers like RunPod.

## Prerequisites

### install Nvidia Container tookit if it is not installed

`````
````bash
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | \
  sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg

curl -fsSL https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
  sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#' | \
  sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list

sudo apt-get update
sudo apt-get install -y nvidia-container-toolkit


docker run --rm --gpus all nvidia/cuda:12.4.1-runtime-ubuntu22.04 nvidia-smi

`````

### Configure Docker to use the NVIDIA runtime and restart:

```bash
sudo nvidia-ctk runtime configure --runtime=docker
sudo systemctl restart docker
```

- Test GPU inside a container:

````bash
docker run --rm --gpus all nvidia/cuda:12.4.1-runtime-ubuntu22.04 nvidia-smi



## Build the image

```bash

docker buildx build -t go-gen .

```


## 📋 Requirements

- **GPU**: NVIDIA GPU with 16GB+ VRAM (RTX 3090/4090, A100, etc.)
- **CUDA**: Version 12.x drivers
- **Storage**: ~10GB for models (downloaded automatically on first run)

## 🎯 Features

- ✅ Automatic model downloading from HuggingFace
- ✅ CUDA-accelerated inference via stable-diffusion.cpp
- ✅ Batch processing of multiple prompts
- ✅ Configurable resolution, steps, and guidance
- ✅ LoRA support for fine-tuning
- ✅ Progress tracking and ETA display
- ✅ Docker and native binary deployment options


## 🔧 Configuration

Edit `character_config.json`:

```json
{
  "seed": 42,
  "output_dir": "./output",
  "resolution": [768, 1024],
  "num_inference_steps": 20,
  "guidance_scale": 3.5,
  "style_suffix": "high quality, detailed",
  "negative_prompt": "blurry, low quality",
  "prompts": [
    "your prompt here",
    "another prompt"
  ]
}
```

**Key Parameters:**
- `resolution`: Image dimensions [width, height]
- `num_inference_steps`: More steps = better quality but slower (15-30 recommended)
- `guidance_scale`: How closely to follow prompt (3.0-7.5 recommended)
- `seed`: Random seed for reproducibility

## 🐳 Docker Deployment (RunPod Ready)

Perfect for RunPod or any GPU cloud service:
```bash

# Run with your config
docker run --name test --gpus all \
  -v $(pwd)/character_config.json:/app/character_config.json:ro \
  -v $(pwd)/models:/app/models \
  -v $(pwd)/output:/app/output \
  -p 2022:22 \
  go-gen:latest \

```


## 🚀 RunPod Deployment


1. Push to Docker Hub:
   ```bash
   docker tag sd-generator:latest yourusername/sd-generator:latest
   docker push yourusername/sd-generator:latest
   ```


````
