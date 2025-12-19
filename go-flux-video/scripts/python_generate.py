import argparse
import hashlib
import json
import os
import sys
import time

import torch
from diffusers import (
    AutoPipelineForText2Image,
    FluxPipeline,
    FluxTransformer2DModel,
    GGUFQuantizationConfig,
)
from diffusers.utils import logging as diffusers_logging
from PIL import Image

_PIPELINE_CACHE = {}


def get_cache_key(args):
    """Create a unique cache key based on model configuration"""
    key_parts = [
        args.model,
        args.gguf or "",
        args.lora or "",
        args.lora_file or "",
        str(args.low_vram),
    ]
    return hashlib.md5(":".join(key_parts).encode()).hexdigest()


def save_image(image: Image.Image, output_path: str) -> None:
    """Save image as PNG, matching Go's png.Encode behavior."""
    # Ensure directory exists
    os.makedirs(os.path.dirname(output_path) or ".", exist_ok=True)

    # Save as PNG with same approach as Go
    with open(output_path, "wb") as f:
        image.save(f, format="PNG", optimize=False)


def load_pipeline(args):
    """Load pipeline with caching"""
    cache_key = get_cache_key(args)

    # Return cached pipeline if available
    if cache_key in _PIPELINE_CACHE:
        print("✓ Using cached pipeline", file=sys.stderr)
        return _PIPELINE_CACHE[cache_key]

    print(f"Loading model: {args.model}", file=sys.stderr, flush=True)
    if args.gguf:
        print(f"Using GGUF: {args.gguf}", file=sys.stderr, flush=True)

    start = time.time()
    diffusers_logging.set_verbosity_error()

    # Load pipeline (same as before)
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
        if hasattr(pipe, "tokenizer") and hasattr(pipe.tokenizer, "model_max_length"):
            pipe.tokenizer.model_max_length = 512
            print("✓ Updated tokenizer max length to 512 for T5", file=sys.stderr)

        if hasattr(pipe.tokenizer, "add_prefix_space"):
            pipe.tokenizer.add_prefix_space = False
            print("✓ Disabled add_prefix_space to avoid warning", file=sys.stderr)

    else:
        pipe = AutoPipelineForText2Image.from_pretrained(
            args.model,
            torch_dtype=torch.float16,
        )
        if hasattr(pipe, "tokenizer"):
            if hasattr(pipe.tokenizer, "model_max_length"):
                pipe.tokenizer.model_max_length = 512
            print(
                "✓ Loaded FluxPipeline, set tokenizer to 512 tokens",
                file=sys.stderr,
                flush=True,
            )
        print("✓ Using  using default model", file=sys.stderr, flush=True)

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
    if args.lora and args.lora_file:
        print(f"Loading LoRA: {args.lora}", file=sys.stderr)
        print(f"  File: {args.lora_file}", file=sys.stderr)
        try:
            pipe.load_lora_weights(
                args.lora,
                weight_name=os.path.basename(args.lora_file),
                adapter_name="uncensored",
            )
            if hasattr(pipe, "safety_checker"):
                pipe.safety_checker = None
                print("✓ Safety checker disabled", file=sys.stderr)
        except Exception as e:
            print(f"⚠ LoRA loading failed: {e}", file=sys.stderr)
            print("  Continuing without LoRA...", file=sys.stderr)

    print(f"✓ Model loaded in {time.time() - start:.1f}s", file=sys.stderr)

    # Cache the pipeline
    _PIPELINE_CACHE[cache_key] = pipe
    return pipe


def main():

    parser = argparse.ArgumentParser()
    parser.add_argument("--model", required=True, help="HuggingFace model ID")
    parser.add_argument("--gguf", default="", help="Path to GGUF file (optional)")
    parser.add_argument("--prompt", required=True)
    parser.add_argument("--negative-prompt", default="")
    parser.add_argument("--width", type=int, default=1024)
    parser.add_argument("--height", type=int, default=1024)
    parser.add_argument("--steps", type=int, default=28)
    parser.add_argument("--guidance-scale", type=float, default=3.5)
    parser.add_argument("--seed", type=int, default=42)
    parser.add_argument("--output", required=True)
    parser.add_argument("--lora", default="", help="LoRA HuggingFace repo ID")
    parser.add_argument("--lora-file", default="", help="safesensoes file")
    # Performance flag
    parser.add_argument(
        "--low-vram",
        action="store_true",
        help="Enable CPU offload and attention slicing for low VRAM GPUs",
    )

    args = parser.parse_args()
    try:
        pipe = load_pipeline(args)

        # Generate
        generator = torch.Generator().manual_seed(args.seed)

        print(f"Generating: {args.prompt[:60]}...", file=sys.stderr)
        print(
            f"  Size: {args.width}x{args.height}, Steps: {args.steps}, Seed: {args.seed}",
            file=sys.stderr,
        )
        gen_start = time.time()

        result = pipe(
            prompt=args.prompt,
            negative_prompt=args.negative_prompt or None,
            width=args.width,
            height=args.height,
            num_inference_steps=args.steps,
            guidance_scale=args.guidance_scale,
            generator=generator,
        )

        image = result.images[0]

        # Save image (matching Go's approach)
        save_image(image, args.output)

        elapsed = time.time() - gen_start
        print(f"✓ Saved to {args.output} in {elapsed:.1f}s", file=sys.stderr)
        print(json.dumps({"status": "success", "output": args.output}))

    except Exception as e:
        print(json.dumps({"status": "error", "error": str(e)}))
        sys.exit(1)


if __name__ == "__main__":
    main()
