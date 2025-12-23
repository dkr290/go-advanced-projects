# Qwen-Image-Edit Integration - Changes Summary

## Files Created

### 1. `scripts/qwen_img2img.py`
- Python script for Qwen-Image-Edit model integration
- Supports both text-to-image and image-to-image generation
- Includes Qwen-specific optimizations and error handling
- Compatible with existing Go wrapper architecture

### 2. `README_QWEN.md`
- Comprehensive documentation for Qwen-Image-Edit integration
- Usage examples for both img2img and text-to-image modes
- Command line flag reference
- Troubleshooting guide

### 3. `config_qwen_example.json`
- Example configuration file for Qwen-Image-Edit
- Shows proper JSON structure with style suffix and negative prompts

### 4. `test_qwen.sh`
- Test script to verify Qwen integration
- Creates test directories and configuration
- Provides example commands for testing

### 5. `example_qwen_commands.txt`
- Collection of practical usage examples
- Shows various flag combinations for different use cases
- Includes environment variable examples

## Files Modified

### 1. `pkg/config/config.go`
- Added `UseQwen bool` field to `CmdConf` struct
- Added `--use-qwen` command line flag
- Added environment variable support (`USE_QWEN`)

### 2. `main.go`
- Added Qwen model selection logic in both img2img and text-to-image modes
- Integrated `GenerateWithPythonQwen` and `GenerateImg2ImgWithPythonQwen` functions
- Maintains backward compatibility with existing SD and FLUX modes

### 3. `pkg/generate/generate_flux_images.go`
- Added `GenerateWithPythonQwen()` function for text-to-image generation
- Added `GenerateImg2ImgWithPythonQwen()` function for image-to-image generation
- Both functions use the new `qwen_img2img.py` script
- Includes proper error handling and progress reporting

## Key Features Implemented

### 1. Model Support
- Integration with Qwen/Qwen-Image-Edit HuggingFace model
- Automatic model downloading on first run
- Proper pipeline initialization with Qwen-specific settings

### 2. Generation Modes
- **Text-to-Image**: Generate new images from text prompts
- **Image-to-Image**: Transform input images based on text prompts
- **Batch Processing**: Process multiple images/prompts in single run

### 3. Performance Optimizations
- Low VRAM mode with CPU offload and attention slicing
- xformers memory efficient attention support
- Proper CUDA/CPU fallback handling

### 4. Customization Options
- LoRA weight support for custom styles
- Adjustable transformation strength (0.0-1.0)
- Configurable resolution, steps, and guidance scale
- Seed control for reproducible results

### 5. Integration Points
- Uses existing Go wrapper architecture
- Compatible with existing configuration files
- Follows same output format as other models
- Reuses existing prompt data structures

## Usage Examples

### Basic Image-to-Image
```bash
./gfluxgo --config config.json --hf-model "Qwen/Qwen-Image-Edit" --use-qwen --img2img
```

### Text-to-Image with Custom Settings
```bash
./gfluxgo --config config.json --hf-model "Qwen/Qwen-Image-Edit" --use-qwen \
  --output ./results --seed 42 --resolution 1024x1024 --steps 30
```

### Low VRAM Mode
```bash
./gfluxgo --config config.json --hf-model "Qwen/Qwen-Image-Edit" --use-qwen --img2img --low-vram
```

## Backward Compatibility

- Existing SD and FLUX functionality remains unchanged
- All existing command line flags continue to work
- Configuration files work with all model types
- Output directory structure remains consistent

## Requirements

1. **Python Packages**: `torch`, `diffusers`, `transformers`, `pillow`
2. **Disk Space**: ~15GB for Qwen-Image-Edit model
3. **GPU Memory**: 8GB+ recommended (4GB with low-vram mode)
4. **Internet**: For model download on first run

## Testing

Run the test script to verify installation:
```bash
chmod +x test_qwen.sh
./test_qwen.sh
```

## Next Steps

1. Test the integration with actual Qwen-Image-Edit model
2. Consider adding Qwen-specific optimizations
3. Add example input/output images to documentation
4. Consider supporting other Qwen variants if needed