import argparse
import json
import os
import sys
import time

import torch
from diffusers import (
    AutoPipelineForImage2Image,
    FluxPipeline,
    FluxTransformer2DModel,
    GGUFQuantizationConfig,
)
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

    print(f"Loading model for img2img: {args.model}", file=sys.stderr, flush=True)
    if args.gguf:
        print(f"Using GGUF: {args.gguf}", file=sys.stderr, flush=True)

    start = time.time()
    diffusers_logging.set_verbosity_error()

    # Load pipeline for image-to-image
    if args.gguf and os.path.exists(args.gguf):
        transformer = FluxTransformer2DModel.from_single_file(
            args.gguf,
            quantization_config=GGUFQuantizationConfig(compute_dtype=torch.bfloat16),
            torch_dtype=torch.bfloat16,
        )

        pipe = FluxPipeline.from_pretrained(
            args.model,
            transformer=transformer,
            torch_dtype=torch.bfloat16,
        )

    else:
        pipe = AutoPipelineForImage2Image.from_pretrained(
            args.model,
            torch_dtype=torch.bfloat16,
        )

        print("✓ Using default FULL model for img2img", file=sys.stderr, flush=True)

    if hasattr(pipe, "tokenizer"):
        # Check and set model_max_length to 512
        if hasattr(pipe.tokenizer, "model_max_length"):
            current_max_length = pipe.tokenizer.model_max_length
            if current_max_length != 512:
                pipe.tokenizer.model_max_length = 512
                print(
                    f"✓ Tokenizer max length updated from {current_max_length} to 512.",
                    file=sys.stderr,
                )
            else:
                print("✓ Tokenizer max length is already 512.", file=sys.stderr)
        else:
            print(
                "⚠ Warning: Tokenizer does not have 'model_max_length' attribute. Cannot enforce 512 tokens.",
                file=sys.stderr,
            )

        # Check and disable add_prefix_space
        if hasattr(pipe.tokenizer, "add_prefix_space"):
            if pipe.tokenizer.add_prefix_space:
                pipe.tokenizer.add_prefix_space = False
                print(
                    "✓ Disabled 'add_prefix_space' to avoid warnings.", file=sys.stderr
                )
            else:
                print("✓ 'add_prefix_space' is already disabled.", file=sys.stderr)
        else:
            print(
                "⚠ Warning: Tokenizer does not have 'add_prefix_space' attribute.",
                file=sys.stderr,
            )
    else:
        print(
            "⚠ Warning: Pipeline does not have a tokenizer. Max length and prefix space checks skipped.",
            file=sys.stderr,
        )

    # Memory optimizations
    if args.low_vram:
        pipe.enable_model_cpu_offload()
        pipe.enable_attention_slicing("auto")
        print("✓ Low VRAM mode enabled", file=sys.stderr)
    else:
        pipe.to("cuda")
        print("✓ Full GPU mode", file=sys.stderr)

    try:
        pipe.enable_xformers_memory_efficient_attention()
        print("✓ xformers enabled", file=sys.stderr)
    except Exception:
        pass

    # Load LoRA if specified
    if args.lora_file:
        print(f"Loading LoRA: {args.lora_file}", file=sys.stderr)
        try:
            # Get the directory containing the lora file
            lora_dir = os.path.dirname(args.lora_file)
            lora_filename = os.path.basename(args.lora_file)

            # Load from local directory instead of HuggingFace
            pipe.load_lora_weights(
                lora_dir,
                weight_name=lora_filename,
                adapter_name="uncensored",
            )
            if hasattr(pipe, "safety_checker"):
                pipe.safety_checker = None
                print("✓ Safety checker disabled", file=sys.stderr)
        except Exception as e:
            print(f"⚠ LoRA loading failed: {e}", file=sys.stderr)
            print("  Continuing without LoRA...", file=sys.stderr)

    print(f"✓ Model loaded for img2img in {time.time() - start:.1f}s", file=sys.stderr)

    # Cache the pipeline
    return pipe


def main():
    # Add diagnostic checks here
    print(f"PyTorch version: {torch.__version__}", file=sys.stderr)
    print(f"CUDA available: {torch.cuda.is_available()}", file=sys.stderr)
    if torch.cuda.is_available():
        print(f"CUDA device count: {torch.cuda.device_count()}", file=sys.stderr)
        for i in range(torch.cuda.device_count()):
            print(f"  Device {i}: {torch.cuda.get_device_name(i)}", file=sys.stderr)
            print(
                f"    Memory allocated: {torch.cuda.memory_allocated(i)/1e9:.2f} GB",
                file=sys.stderr,
            )
            print(
                f"    Memory reserved: {torch.cuda.memory_reserved(i)/1e9:.2f} GB",
                file=sys.stderr,
            )
    else:
        print("WARNING: CUDA not available.", file=sys.stderr)
    
    parser = argparse.ArgumentParser()
    parser.add_argument("--model", required=True, help="HuggingFace model ID")
    parser.add_argument("--gguf", default="", help="Path to GGUF file (optional)")
    parser.add_argument("--negative-prompt", default="")
    parser.add_argument("--width", type=int, default=1024)
    parser.add_argument("--height", type=int, default=1024)
    parser.add_argument("--steps", type=int, default=28)
    parser.add_argument("--guidance-scale", type=float, default=3.5)
    parser.add_argument("--strength", type=float, default=0.75, help="Transformation strength (0.0-1.0). Higher = more creative, lower = closer to input")
    parser.add_argument("--seed", type=int, default=42)
    parser.add_argument("--output-dir", required=True)
    parser.add_argument("--lora-file", default="", help="safetensors file")
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
                print(f"✗ Failed to load input image {input_image_path}: {e}", file=sys.stderr)
                all_results.append(
                    {"status": "error", "error": f"Failed to load input image: {str(e)}", "prompt_index": i}
                )
                continue
            
            output_path = os.path.join(args.output_dir, filename)
            
            # Generate
            generator = torch.Generator().manual_seed(seed)

            print(f"Generating img2img: {prompt[:60]}...", file=sys.stderr)
            print(
                f"  Input: {os.path.basename(input_image_path)}, Size: {args.width}x{args.height}, Steps: {args.steps}, Strength: {args.strength}, Seed: {seed}",
                file=sys.stderr,
            )
            gen_start = time.time()

            result = pipe(
                prompt=prompt,
                image=input_image,
                negative_prompt=args.negative_prompt,
                width=args.width,
                height=args.height,
                num_inference_steps=args.steps,
                guidance_scale=args.guidance_scale,
                strength=args.strength,
                generator=generator,
            )

            image = result.images[0]

            # Save image (matching Go's approach)
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
