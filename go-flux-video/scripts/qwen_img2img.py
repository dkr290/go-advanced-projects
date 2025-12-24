import argparse
import json
import os
import sys
import time

import torch
from diffusers import AutoPipelineForImage2Image
from diffusers.utils import logging as diffusers_logging
from PIL import Image


def save_image(image: Image.Image, output_path: str) -> None:
    """Save image as PNG, matching Go's png.Encode behavior."""
    # Ensure directory exists
    os.makedirs(os.path.dirname(output_path) or ".", exist_ok=True)

    # Save as PNG with same approach as Go
    with open(output_path, "wb") as f:
        image.save(f, format="PNG", optimize=False)


def load_image(image_path: str) -> Image.Image:
    """Load and prepare input image for img2img."""
    if not os.path.exists(image_path):
        raise FileNotFoundError(f"Input image not found: {image_path}")

    image = Image.open(image_path)
    # Convert to RGB if needed (removes alpha channel, converts grayscale, etc.)
    if image.mode != "RGB":
        image = image.convert("RGB")

    return image


def load_pipeline(args):
    """Load Qwen-Image-Edit pipeline with optional optimizations."""

    print(f"Loading Qwen-Image-Edit model: {args.model}", file=sys.stderr, flush=True)

    start = time.time()
    diffusers_logging.set_verbosity_error()

    # Load Qwen-Image-Edit pipeline
    # Note: Qwen-Image-Edit uses a different pipeline structure
    # We'll use AutoPipelineForImage2Image which should work with Qwen models
    pipe = AutoPipelineForImage2Image.from_pretrained(
        args.model,
        torch_dtype=torch.bfloat16 if torch.cuda.is_available() else torch.float32,
    )
    pipe.set_progress_bar_config(disable=None)
    print(f"✓ Loaded Qwen-Image-Edit model: {args.model}", file=sys.stderr, flush=True)

    # Memory optimizations
    if args.low_vram:
        pipe.enable_model_cpu_offload()
        pipe.enable_attention_slicing("auto")
        print("✓ Low VRAM mode enabled", file=sys.stderr)
    else:
        if torch.cuda.is_available():
            pipe.to("cuda")
            print("✓ Full GPU mode", file=sys.stderr)
        else:
            print("✓ CPU mode", file=sys.stderr)

    # Try to enable xformers for memory efficiency
    try:
        pipe.enable_xformers_memory_efficient_attention()
        print("✓ xformers enabled", file=sys.stderr)
    except Exception:
        print("⚠ xformers not available, using default attention", file=sys.stderr)

    # Load LoRA if specified (Qwen models may support LoRA)
    if args.lora_file:
        print(f"Loading LoRA: {args.lora_file}", file=sys.stderr)
        try:
            # Get the directory containing the lora file
            lora_dir = os.path.dirname(args.lora_file)
            lora_filename = os.path.basename(args.lora_file)

            # Load from local directory
            pipe.load_lora_weights(
                lora_dir,
                weight_name=lora_filename,
                adapter_name="custom_lora",
            )
            print("✓ LoRA loaded successfully", file=sys.stderr)
        except Exception as e:
            print(f"⚠ LoRA loading failed: {e}", file=sys.stderr)
            print("  Continuing without LoRA...", file=sys.stderr)

    print(f"✓ Pipeline loaded in {time.time() - start:.1f}s", file=sys.stderr)

    return pipe


def main():
    # Add diagnostic checks
    print(f"PyTorch version: {torch.__version__}", file=sys.stderr)
    print(f"CUDA available: {torch.cuda.is_available()}", file=sys.stderr)
    if torch.cuda.is_available():
        print(f"CUDA device count: {torch.cuda.device_count()}", file=sys.stderr)
        for i in range(torch.cuda.device_count()):
            print(f"  Device {i}: {torch.cuda.get_device_name(i)}", file=sys.stderr)

    parser = argparse.ArgumentParser(
        description="Qwen-Image-Edit image-to-image generation"
    )
    parser.add_argument(
        "--model",
        required=True,
        help="HuggingFace model ID (e.g., Qwen/Qwen-Image-Edit)",
    )
    parser.add_argument(
        "--negative-prompt", default="", help="Negative prompt for generation"
    )
    parser.add_argument("--width", type=int, default=1024, help="Output image width")
    parser.add_argument("--height", type=int, default=1024, help="Output image height")
    parser.add_argument(
        "--steps", type=int, default=40, help="Number of inference steps"
    )
    parser.add_argument(
        "--guidance-scale", type=float, default=1.0, help="Guidance scale (CFG)"
    )
    parser.add_argument(
        "--strength",
        type=float,
        default=0.75,
        help="Transformation strength (0.0-1.0). Higher = more creative, lower = closer to input",
    )
    parser.add_argument("--seed", type=int, default=42, help="Random seed")
    parser.add_argument(
        "--output-dir", required=True, help="Output directory for generated images"
    )
    parser.add_argument(
        "--lora-file", default="", help="Path to LoRA safetensors file (optional)"
    )
    parser.add_argument(
        "--num-images", default=1, type=int, help="Number of images per prompt"
    )

    # New argument for Qwen-specific settings
    parser.add_argument(
        "--max-tokens",
        type=int,
        default=2048,
        help="Maximum tokens for Qwen model (default: 2048)",
    )

    # New argument to accept multiple prompts and their data for img2img
    parser.add_argument(
        "--prompts-data",
        required=True,
        help='JSON string of prompt data, e.g., \'[{"prompt": "a dog", "filename": "dog.png", "seed": 123, "input_image": "input.png"}]\'',
    )

    # Performance flag
    parser.add_argument(
        "--low-vram",
        action="store_true",
        help="Enable CPU offload and attention slicing for low VRAM GPUs",
    )

    args = parser.parse_args()

    try:
        pipe = load_pipeline(args)

        # Parse the incoming JSON string for prompts data
        prompts_data = json.loads(args.prompts_data)
        all_results = []

        for i, p_data in enumerate(prompts_data):
            prompt = p_data["prompt"]
            filename = p_data["filename"]
            seed = p_data["seed"]
            input_image_path = p_data["input_image"]

            # Load input image
            try:
                input_image = load_image(input_image_path)
            except Exception as e:
                print(
                    f"✗ Failed to load input image {input_image_path}: {e}",
                    file=sys.stderr,
                )
                all_results.append(
                    {
                        "status": "error",
                        "error": f"Failed to load input image: {str(e)}",
                        "prompt_index": i,
                    }
                )
                continue

            output_path = os.path.join(args.output_dir, filename)

            # Generate with Qwen-Image-Edit
            generator = torch.Generator().manual_seed(seed)

            print(f"Generating with Qwen-Image-Edit: {prompt[:60]}...", file=sys.stderr)
            print(
                f"  Input: {os.path.basename(input_image_path)}, Size: {args.width}x{args.height}, Steps: {args.steps}, Strength: {args.strength}, Seed: {seed}",
                file=sys.stderr,
            )
            gen_start = time.time()

            # Qwen-Image-Edit specific parameters
            result = pipe(
                prompt=prompt,
                image=input_image,
                negative_prompt=args.negative_prompt,
                width=args.width,
                height=args.height,
                num_inference_steps=args.steps,
                guidance_scale=args.guidance_scale,
                true_cfg_scale=4.0,
                strength=args.strength,
                generator=generator,
                num_images_per_prompt=1,
                # Qwen-specific parameters if needed
                # max_length=args.max_tokens,
            )

            image = result.images[0]

            # Save image
            save_image(image, output_path)

            elapsed = time.time() - gen_start
            print(f"✓ Saved to {output_path} in {elapsed:.1f}s", file=sys.stderr)
            all_results.append(
                {"status": "success", "output": output_path, "prompt_index": i}
            )

        print(json.dumps({"all_status": "success", "generations": all_results}))

    except Exception as e:
        print(json.dumps({"status": "error", "error": str(e)}))
        sys.exit(1)


if __name__ == "__main__":
    main()

