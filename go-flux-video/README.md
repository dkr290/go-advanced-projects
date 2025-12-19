Go-orchestrated FLUX image generator using Python diffusers with CUDA acceleration.

## 🎯 Features

- ✅ **Go orchestration** - All parameters controlled via CLI flags or environment variables
- ✅ **Automatic model downloading** - Downloads GGUF models and LoRA from HuggingFace
- ✅ **Batch processing** - Generate multiple images from a JSON prompt list
- ✅ **CUDA acceleration** - Full GPU support via PyTorch
- ✅ **LoRA support** - Load custom LoRA weights for fine-tuned generations
- ✅ **Smart filenames** - Images named using first 10 characters of the prompt
- ✅ **Docker ready** - Pre-configured for RunPod deployment

## 📋 Requirements

- **GPU**: NVIDIA GPU with 12GB+ VRAM (RTX 3080/3090/4090, A100, etc.)
- **CUDA**: Version 12.x drivers
- **Storage**: ~15GB for FLUX models

## 🚀 Quick Start

### Local Development

1. **Install dependencies:**

   ```bash
   # Go 1.23+
   go mod download

   # Python 3.10+
   pip install torch diffusers transformers accelerate safetensors pillow xformers
   ```

2. **Create config file:**

   ```bash
   cp character_config.json.example character_config.json
   ```

3. **Run:**
   ```bash
   go run main.go --config character_config.json
   ```

# Use default FLUX.1-dev model (no GGUF)

./go-flux-video --config config.json

# Use a different HuggingFace model

./go-flux-video --config config.json --hf-model "stabilityai/stable-diffusion-2"

# Use environment variable

export HF_MODEL="runwayml/stable-diffusion-v1-5"
./go-flux-video --config config.json

# Use with GGUF file

./go-flux-video --config config.json --gguf-model-url "https://example.com/flux.gguf"

````

### Docker (RunPod)

1. **Build the image:**

```bash
docker buildx build -t gfluxgo .
````

2. **Run container:**

   ```bash
   docker run --gpus all -d \
     -v $(pwd)/models:/app/models \
     -v $(pwd)/output:/app/output \
     -p 2022:22 \
     gfluxgo
   ```

3. **SSH into container and run:**
   ```bash
   ssh root@localhost -p 2022
   ./sd-generator --config /app/config.json \
   --lora-repo "Heartsync/FLUX-NSFW-uncensored" \
   --lora-url "https://huggingface.co/.../lora.safetensors"
   ```

## 🔧 Configuration

### CLI Flags

| Flag                | Default           | Description                          |
| ------------------- | ----------------- | ------------------------------------ |
| `--config`          | (required)        | Path to JSON prompt configuration    |
| `--model-url`       | FLUX.1-dev Q4_K_S | URL to GGUF model                    |
| `--lora-url`        | -                 | URL to LoRA safetensors file         |
| `--output`          | `./output`        | Output directory for images          |
| `--resolution`      | `1024x1024`       | Image resolution (WIDTHxHEIGHT)      |
| `--num_steps`       | `28`              | Number of inference steps            |
| `--guidence_scale`  | `7.0`             | Guidance scale (3.0-7.5 recommended) |
| `--seed`            | `42`              | Random seed for reproducibility      |
| `--model-down-path` | `./models`        | Model download directory             |
| `--lora-down-path`  | `./models/lora`   | LoRA download directory              |
| `--lora-repo`       |                   | Lora repo ID                         |
| `--debug`           | `false`           | Enable debug logging                 |

### Environment Variables

All flags can be overridden with environment variables:

| Variable         | Overrides          |
| ---------------- | ------------------ |
| `MODEL_URL`      | `--model-url`      |
| `LORA_URL`       | `--lora-url`       |
| `OUTPUT`         | `--output`         |
| `RESOLUTION`     | `--resolution`     |
| `STEPS`          | `--num_steps`      |
| `GUIDANCE_SCALE` | `--guidence_scale` |
| `SEED`           | `--seed`           |

### Prompt Configuration (JSON)

```json
{
  "style_suffix": "high quality, detailed, professional digital art, sharp focus",
  "negative_prompt": "blurry, low quality, distorted, deformed, ugly, bad anatomy",
  "prompts": [
    "a woman standing in front view, full body",
    "a woman portrait, face closeup"
  ]
}
```

| Field             | Description                                   |
| ----------------- | --------------------------------------------- |
| `style_suffix`    | Appended to every prompt for consistent style |
| `negative_prompt` | What to avoid in generation                   |
| `prompts`         | List of prompts to generate                   |

## 📤 Output

Generated images are saved to the output directory with filenames based on the prompt:

```
output/
├── 01_a-woman.png
|__ 02_a-woman.png

```

Format: `{index}_{first-10-chars-of-prompt}.png`

3. **SSH into pod** and run:
   ```bash
   cd /app
   ./sd-generator --config your_config.json --output /workspace/output
   ```

## 🛠️ NVIDIA Container Toolkit Setup (Host Machine)

If running locally with Docker, install NVIDIA Container Toolkit:

```bash
# Add NVIDIA repo
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | \
  sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg

curl -fsSL https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
  sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#' | \
  sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list

# Install
sudo apt-get update
sudo apt-get install -y nvidia-container-toolkit

# Configure Docker
sudo nvidia-ctk runtime configure --runtime=docker
sudo systemctl restart docker

# Test
docker run --rm --gpus all nvidia/cuda:12.4.1-runtime-ubuntu22.04 nvidia-smi
```
