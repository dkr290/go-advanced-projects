# Qwen-Image-Edit Integration

This project now supports the Qwen/Qwen-Image-Edit model for advanced image editing tasks. The Qwen-Image-Edit model is specialized for image-to-image transformations and editing.

## Features

- **Image-to-Image Editing**: Transform input images based on text prompts
- **Text-to-Image Generation**: Generate new images from text prompts
- **LoRA Support**: Use custom LoRA weights for specialized styles
- **Low VRAM Mode**: Optimized for GPUs with limited memory
- **Batch Processing**: Process multiple images in a single run

## Installation Requirements

Make sure you have the required Python packages:

```bash
pip install torch diffusers transformers pillow
```

## Input Image Location (for Image-to-Image Mode)

When using image-to-image mode (`--img2img` flag), you need to place your input images in the `./images/` directory.

### Directory Structure:
```
your-project/
├── images/           # ← PUT YOUR INPUT IMAGES HERE
│   ├── photo1.jpg
│   ├── photo2.png
│   └── photo3.webp
├── scripts/
│   └── qwen_img2img.py
├── main.go
└── config.json
```

### How to Set Up:

1. **Create the images directory:**
   ```bash
   mkdir -p images
   ```

2. **Copy your input images into it:**
   ```bash
   cp ~/Pictures/my_photo.jpg images/
   cp ~/Downloads/another_image.png images/
   ```

3. **Image Matching:** The program matches prompts with images in order:
   - First prompt in config → First image in directory (alphabetical order)
   - Second prompt → Second image
   - etc.

### Example:

**config.json:**
```json
{
  "prompts": [
    "make it look like a painting",
    "add sunset colors",
    "convert to black and white"
  ]
}
```

**./images/ directory:**
```
./images/
├── landscape.jpg    # Used with "make it look like a painting"
├── portrait.png     # Used with "add sunset colors"
└── cityscape.webp   # Used with "convert to black and white"
```

### Important Notes:

1. **File Formats Supported**: JPG, PNG, WebP, BMP, GIF, TIFF, etc. (anything PIL/Pillow supports)
2. **Order Matters**: Images are processed in alphabetical order by filename
3. **Image Count**: You need at least as many images as you have prompts in your config
4. **Directory Creation**: The program doesn't create the `images/` directory automatically
5. **Image Preparation**: Images are automatically converted to RGB format if needed

## Usage Examples

### 1. Basic Image-to-Image Editing

```bash
# Step 1: Create input images directory
mkdir -p images

# Step 2: Place your input images in the images directory
# Example: cp ~/Pictures/your_image.jpg images/

# Step 3: Run Qwen-Image-Edit with img2img mode
./gfluxgo \
  --config config_qwen_example.json \
  --hf-model "Qwen/Qwen-Image-Edit" \
  --use-qwen \
  --img2img \
  --strength 0.7 \
  --output ./qwen_output \
  --seed 42 \
  --resolution 1024x1024 \
  --steps 28 \
  --guidence_scale 3.5
```

### 2. Text-to-Image Generation

```bash
# Generate new images from text prompts
./gfluxgo \
  --config config_qwen_example.json \
  --hf-model "Qwen/Qwen-Image-Edit" \
  --use-qwen \
  --output ./qwen_output \
  --seed 42 \
  --resolution 1024x1024 \
  --steps 28 \
  --guidence_scale 3.5
```

### 3. With LoRA Weights

```bash
# Use custom LoRA weights for specific styles
./gfluxgo \
  --config config_qwen_example.json \
  --hf-model "Qwen/Qwen-Image-Edit" \
  --use-qwen \
  --img2img \
  --lora-url "https://huggingface.co/your-username/your-lora" \
  --output ./qwen_output \
  --strength 0.75
```

### 4. Low VRAM Mode

```bash
# For GPUs with limited memory (8GB or less)
./gfluxgo \
  --config config_qwen_example.json \
  --hf-model "Qwen/Qwen-Image-Edit" \
  --use-qwen \
  --img2img \
  --low-vram \
  --output ./qwen_output
```

## Configuration File Example

Create a JSON configuration file (e.g., `config_qwen_example.json`):

```json
{
  "style_suffix": "professional photography, 8k, highly detailed",
  "negative_prompt": "blurry, low quality, distorted, watermark, text",
  "prompts": [
    "a beautiful sunset over mountains",
    "a futuristic city at night",
    "a serene lake with autumn trees",
    "a magical forest with glowing mushrooms"
  ]
}
```

## Environment Variables

You can also use environment variables for configuration:

```bash
export HF_MODEL="Qwen/Qwen-Image-Edit"
export USE_QWEN="true"
export IMG2IMG="true"
export STRENGTH="0.7"
export OUTPUT="./qwen_output"
export SEED="42"
export RESOLUTION="1024x1024"
export STEPS="28"
export GUIDANCE_SCALE="3.5"

./gfluxgo --config config_qwen_example.json
```

## Command Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--use-qwen` | Use Qwen-Image-Edit model | `false` |
| `--hf-model` | HuggingFace model ID | `"Qwen/Qwen-Image-Edit"` |
| `--img2img` | Enable image-to-image mode | `false` |
| `--strength` | Transformation strength (0.0-1.0) | `0.75` |
| `--low-vram` | Enable low VRAM mode | `false` |
| `--output` | Output directory | `./output` |
| `--seed` | Random seed | `42` |
| `--resolution` | Image resolution | `1024x1024` |
| `--steps` | Inference steps | `28` |
| `--guidence_scale` | Guidance scale | `3.5` |

## Python Script

The Qwen-Image-Edit functionality is implemented in `scripts/qwen_img2img.py`. You can also run it directly:

```bash
python3 scripts/qwen_img2img.py \
  --model "Qwen/Qwen-Image-Edit" \
  --prompts-data '[{"prompt": "a cat wearing a hat", "filename": "cat_hat.png", "seed": 123, "input_image": "input_cat.png"}]' \
  --output-dir ./output \
  --width 1024 \
  --height 1024 \
  --steps 28 \
  --guidance-scale 3.5 \
  --strength 0.75 \
  --seed 42
```

## Notes

1. **Model Size**: Qwen-Image-Edit is a large model (~15GB). Ensure you have sufficient disk space and GPU memory.
2. **Performance**: The first run will download the model, which may take several minutes.
3. **Input Images**: For img2img mode, place your input images in the `./images` directory.
4. **Output Format**: All images are saved as PNG files in the specified output directory.

## Troubleshooting

### General Issues:
- **CUDA Out of Memory**: Enable `--low-vram` flag or reduce batch size
- **Slow Generation**: Reduce `--steps` or `--resolution`
- **Model Not Found**: Check your internet connection and HuggingFace access
- **Python Errors**: Ensure all required packages are installed

### Image-Related Issues:
- **"not enough input images for prompt"**: You have more prompts than images in `./images/` directory
- **Input images not found**: Make sure `./images/` directory exists and contains images
- **Image loading errors**: Check file permissions and format compatibility
- **Wrong image used for prompt**: Images are matched in alphabetical order - rename files if needed

### Quick Fixes:
1. **Create images directory**: `mkdir -p images`
2. **Add enough images**: Ensure you have at least as many images as prompts
3. **Check image formats**: Use common formats like JPG, PNG, WebP
4. **Verify file permissions**: Ensure images are readable

## License

This integration uses the Qwen-Image-Edit model under its respective license. Please review the model's license terms on HuggingFace.
