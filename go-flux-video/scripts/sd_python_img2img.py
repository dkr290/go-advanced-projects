import argparse
import json
import os
import sys
import time

import torch
from diffusers import (
    AutoPipelineForImage2Image,
    StableDiffusionImg2ImgPipeline,
    StableDiffusionXLImg2ImgPipeline,
    StableDiffusion3Img2ImgPipeline,
)
from diffusers.utils import logging as diffusers_logging
from PIL import Image


def save_image(image: Image.Image, output_path: str) -> None:
    """Save image as PNG, matching Go's png.Encode behavior."""
    os.makedirs(os.path.dirname(output_path) or ".", exist_ok=True)
    with open(output_path, "wb") as f:
        image.save(f, format="PNG", optimize=False)


def load_image(image_path: str, target_width: int = None, target_height: int = None) -> Image.Image:
    """Load and prepare input image for img2img."""
    if not os.path.exists(image_path):
        raise FileNotFoundError(f"Input image not found: {image_path}")
    
    image = Image.open(image_path)
    
    # Convert to RGB if needed
    if image.mode != "RGB":
        image = image.convert("RGB")
    
    # Resize if target dimensions specified
    if target_width and target_height:
        # Resize to target size
        image = image.resize((target_width, target_height), Image.Resampling.LANCZOS)
        print(f"  Resized image to {target_width}x{target_height}", file=sys.stderr)
    
    return image


def detect_model_type(model_id: str) -> str:
    """Detect Stable Diffusion model type from model ID."""
    model_lower = model_id.lower()
    
    if "stable-diffusion-3" in model_lower or "sd3" in model_lower:
        return "sd3"
    elif "sdxl" in model_lower or "xl" in model_lower:
        return "sdxl"
    elif "stable-diffusion-2" in model_lower or "sd-2" in model_lower:
        return "sd2"
    else:
        return "sd15"  # Default to SD 1.5


def load_pipeline(args):
    print(f"Loading Stable Diffusion img2img model: {args.model}", file=sys.stderr, flush=True)
    
    start = time.time()
    diffusers_logging.set_verbosity_error()
    
    model_type = detect_model_type(args.model)
    print(f"✓ Detected model type: {model_type.upper()}", file=sys.stderr)
    
    # Determine dtype based on model
    dtype = torch.float16  # SD/SDXL use FP16
    
    # Load from single safetensors file (Civitai models)
    if args.safetensors and os.path.exists(args.safetensors):
        print(f"Loading from safetensors: {args.safetensors}", file=sys.stderr)
        
        if model_type == "sdxl":
            pipe = StableDiffusionXLImg2ImgPipeline.from_single_file(
                args.safetensors,
                torch_dtype=dtype,
                use_safetensors=True,
            )
            print("✓ Loaded SDXL img2img from single file", file=sys.stderr)
        elif model_type == "sd3":
            pipe = StableDiffusion3Img2ImgPipeline.from_single_file(
                args.safetensors,
                torch_dtype=dtype,
                use_safetensors=True,
            )
            print("✓ Loaded SD3 img2img from single file", file=sys.stderr)
        else:  # SD 1.5 / 2.1
            pipe = StableDiffusionImg2ImgPipeline.from_single_file(
                args.safetensors,
                torch_dtype=dtype,
                use_safetensors=True,
            )
            print(f"✓ Loaded {model_type.upper()} img2img from single file", file=sys.stderr)
    
    # Load from HuggingFace Hub
    else:
        try:
            pipe = AutoPipelineForImage2Image.from_pretrained(
                args.model,
                torch_dtype=dtype,
                variant="fp16",  # Use FP16 variant for memory savings
            )
            print("✓ Loaded from HuggingFace Hub with FP16", file=sys.stderr)
        except Exception as e:
            # Fallback without variant if fp16 not available
            print(f"⚠ FP16 variant not found, loading standard weights", file=sys.stderr)
            pipe = AutoPipelineForImage2Image.from_pretrained(
                args.model,
                torch_dtype=dtype,
            )
    
    # Tokenizer configuration
    if hasattr(pipe, "tokenizer") and pipe.tokenizer is not None:
        if hasattr(pipe.tokenizer, "model_max_length"):
            current_max_length = pipe.tokenizer.model_max_length
            if model_type == "sd15" or model_type == "sd2":
                if current_max_length != 77:
                    pipe.tokenizer.model_max_length = 77
                    print(f"✓ Tokenizer max length set to 77 (SD standard)", file=sys.stderr)
            else:
                print(f"✓ Tokenizer max length: {current_max_length}", file=sys.stderr)
        
        if hasattr(pipe.tokenizer, "add_prefix_space"):
            if pipe.tokenizer.add_prefix_space:
                pipe.tokenizer.add_prefix_space = False
                print("✓ Disabled 'add_prefix_space'", file=sys.stderr)
    
    # Memory optimizations
    if args.low_vram:
        if args.sequential_offload:
            pipe.enable_sequential_cpu_offload()
            print("✓ Sequential CPU offload enabled (ultra low VRAM mode)", file=sys.stderr)
        else:
            pipe.enable_model_cpu_offload()
            print("✓ Model CPU offload enabled (low VRAM mode)", file=sys.stderr)
        
        pipe.enable_attention_slicing("auto")
        print("✓ Attention slicing enabled", file=sys.stderr)
        
        # VAE optimizations for SDXL/SD
        if hasattr(pipe, 'vae'):
            pipe.vae.enable_slicing()
            pipe.vae.enable_tiling()
            print("✓ VAE slicing and tiling enabled", file=sys.stderr)
    else:
        pipe.to("cuda")
        print("✓ Full GPU mode", file=sys.stderr)
    
    # Torch compile for speed boost (2x faster)
    if args.compile:
        print("Compiling model for faster inference...", file=sys.stderr)
        if hasattr(pipe, 'unet'):
            pipe.unet = torch.compile(pipe.unet, mode="reduce-overhead", fullgraph=True)
            print("✓ UNet compiled (first generation will be slow)", file=sys.stderr)
        elif hasattr(pipe, 'transformer'):
            pipe.transformer = torch.compile(pipe.transformer, mode="reduce-overhead", fullgraph=True)
            print("✓ Transformer compiled (first generation will be slow)", file=sys.stderr)
    
    # xformers for memory efficiency
    try:
        pipe.enable_xformers_memory_efficient_attention()
        print("✓ xformers memory efficient attention enabled", file=sys.stderr)
    except Exception:
        pass
    
    # Load LoRA if specified
    if args.lora_file:
        print(f"Loading LoRA: {args.lora_file}", file=sys.stderr)
        try:
            lora_dir = os.path.dirname(args.lora_file)
            lora_filename = os.path.basename(args.lora_file)
            
            pipe.load_lora_weights(
                lora_dir,
                weight_name=lora_filename,
            )
            print("✓ LoRA loaded successfully", file=sys.stderr)
            
            # Disable safety checker if present
            if hasattr(pipe, "safety_checker"):
                pipe.safety_checker = None
                print("✓ Safety checker disabled", file=sys.stderr)
        except Exception as e:
            print(f"⚠ LoRA loading failed: {e}", file=sys.stderr)
            print("  Continuing without LoRA...", file=sys.stderr)
    
    # Disable safety checker if requested
    if args.disable_safety and hasattr(pipe, "safety_checker"):
        pipe.safety_checker = None
        print("✓ Safety checker disabled", file=sys.stderr)
    
    print(f"✓ Model loaded for img2img in {time.time() - start:.1f}s", file=sys.stderr)
    print(f"✓ Recommended resolution: {get_recommended_resolution(model_type)}", file=sys.stderr)
    print(f"✓ Strength: {args.strength} (0.0=keep original, 1.0=maximum change)", file=sys.stderr)
    
    return pipe


def get_recommended_resolution(model_type: str) -> str:
    """Get recommended resolution for model type."""
    resolutions = {
        "sd15": "512x512",
        "sd2": "768x768",
        "sdxl": "1024x1024",
        "sd3": "1024x1024",
    }
    return resolutions.get(model_type, "512x512")


def main():
    # Diagnostic info
    print(f"PyTorch version: {torch.__version__}", file=sys.stderr)
    print(f"CUDA available: {torch.cuda.is_available()}", file=sys.stderr)
    if torch.cuda.is_available():
        print(f"CUDA device: {torch.cuda.get_device_name(0)}", file=sys.stderr)
        print(f"CUDA memory: {torch.cuda.get_device_properties(0).total_memory / 1e9:.2f} GB", file=sys.stderr)
    else:
        print("WARNING: CUDA not available, generation will be very slow!", file=sys.stderr)
    
    parser = argparse.ArgumentParser(description="Stable Diffusion Image-to-Image Generation")
    parser.add_argument("--model", required=True, help="HuggingFace model ID or path")
    parser.add_argument("--safetensors", default="", help="Path to safetensors file (Civitai models)")
    parser.add_argument("--negative-prompt", default="")
    parser.add_argument("--width", type=int, default=1024)
    parser.add_argument("--height", type=int, default=1024)
    parser.add_argument("--steps", type=int, default=30)
    parser.add_argument("--guidance-scale", type=float, default=7.5)
    parser.add_argument("--strength", type=float, default=0.75, 
                       help="Transformation strength (0.0-1.0). Higher = more creative, lower = closer to input")
    parser.add_argument("--seed", type=int, default=42)
    parser.add_argument("--output-dir", required=True)
    parser.add_argument("--lora-file", default="", help="LoRA safetensors file")
    parser.add_argument("--prompts-data", required=True, 
                       help='JSON string of prompt data with input_image paths')
    
    # Memory optimization flags
    parser.add_argument(
        "--low-vram",
        action="store_true",
        help="Enable CPU offload and attention slicing for low VRAM GPUs (8GB)",
    )
    parser.add_argument(
        "--sequential-offload",
        action="store_true",
        help="Enable sequential CPU offload for ultra low VRAM (4GB)",
    )
    parser.add_argument(
        "--compile",
        action="store_true",
        help="Compile model for 2x faster inference (requires PyTorch 2.0+)",
    )
    parser.add_argument(
        "--disable-safety",
        action="store_true",
        help="Disable safety checker",
    )
    
    args = parser.parse_args()
    
    try:
        pipe = load_pipeline(args)
        
        # Parse prompts data
        prompts_data = json.loads(args.prompts_data)
        all_results = []
        
        print(f"\n🎨 Starting img2img generation of {len(prompts_data)} images...\n", file=sys.stderr)
        
        for i, p_data in enumerate(prompts_data):
            prompt = p_data["prompt"]
            filename = p_data["filename"]
            seed = p_data.get("seed", args.seed)
            input_image_path = p_data["input_image"]
            output_path = os.path.join(args.output_dir, filename)
            
            # Load input image
            try:
                input_image = load_image(input_image_path, args.width, args.height)
                print(f"[{i+1}/{len(prompts_data)}] Loaded input: {os.path.basename(input_image_path)}", file=sys.stderr)
            except Exception as e:
                print(f"✗ Failed to load input image {input_image_path}: {e}", file=sys.stderr)
                all_results.append({
                    "status": "error",
                    "error": f"Failed to load input image: {str(e)}",
                    "prompt_index": i
                })
                continue
            
            # Generate
            generator = torch.Generator(device="cuda" if torch.cuda.is_available() else "cpu").manual_seed(seed)
            
            print(f"  Prompt: {prompt[:60]}...", file=sys.stderr)
            print(f"  Resolution: {args.width}x{args.height}, Steps: {args.steps}, Strength: {args.strength}, Seed: {seed}, CFG: {args.guidance_scale}", file=sys.stderr)
            gen_start = time.time()
            
            result = pipe(
                prompt=prompt,
                image=input_image,
                negative_prompt=args.negative_prompt,
                num_inference_steps=args.steps,
                guidance_scale=args.guidance_scale,
                strength=args.strength,
                generator=generator,
            )
            
            image = result.images[0]
            
            # Save image
            save_image(image, output_path)
            
            elapsed = time.time() - gen_start
            print(f"✓ Saved to {output_path} in {elapsed:.1f}s\n", file=sys.stderr)
            
            all_results.append({
                "status": "success",
                "output": output_path,
                "prompt_index": i
            })
        
        print(json.dumps({"all_status": "success", "generations": all_results}))
    
    except Exception as e:
        print(json.dumps({"status": "error", "error": str(e)}))
        import traceback
        traceback.print_exc(file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
