import argparse
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


def save_image(image: Image.Image, output_path: str) -> None:
    """Save image as PNG, matching Go's png.Encode behavior."""
    # Ensure directory exists
    os.makedirs(os.path.dirname(output_path) or ".", exist_ok=True)

    # Save as PNG with same approach as Go
    with open(output_path, "wb") as f:
        image.save(f, format="PNG", optimize=False)


def main():
    # Always safe optimizations
    torch.backends.cudnn.benchmark = True  # Works on all CUDA GPUs
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

    print(f"Loading model: {args.model}", file=sys.stderr)
    if args.gguf:
        print(f"Using GGUF: {args.gguf}", file=sys.stderr)
    start = time.time()
    diffusers_logging.set_verbosity_error()
    try:
        # Load pipeline
        if args.gguf and os.path.exists(args.gguf):
            transformer = FluxTransformer2DModel.from_single_file(
                args.gguf,
                quantization_config=GGUFQuantizationConfig(
                    compute_dtype=torch.bfloat16
                ),
                torch_dtype=torch.bfloat16,
            )

            pipe = FluxPipeline.from_pretrained(
                args.model,
                transformer=transformer,
                torch_dtype=torch.bfloat16,
            )
        else:
            pipe = AutoPipelineForText2Image.from_pretrained(
                args.model,
                torch_dtype=torch.float16,
            )

        # Memory optimizations
        if args.low_vram:
            pipe.enable_model_cpu_offload()
            pipe.enable_attention_slicing("auto")
            print(
                "✓ Low VRAM mode: CPU offload + attention slicing enabled",
                file=sys.stderr,
            )
        else:
            # Full GPU mode - fastest on high VRAM GPUs
            pipe.to("cuda")
            print("✓ Full GPU mode: maximum performance", file=sys.stderr)

        try:
            pipe.enable_xformers_memory_efficient_attention()
            print("✓ xformers enabled", file=sys.stderr)
        except Exception:
            pass

        # Load LoRA if specified
        if args.lora:
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
                    print("✓ Safety checker disabled")
                else:
                    print("✓ LoRA loaded (no safety checker found)")

            except Exception as e:
                print(f"⚠ LoRA loading failed: {e}", file=sys.stderr)
                print("  Continuing without LoRA...", file=sys.stderr)

        print(f"✓ Model loaded in {time.time() - start:.1f}s", file=sys.stderr)
        if torch.cuda.is_available():
            print(
                f"  VRAM used: {torch.cuda.memory_allocated() / 1024**3:.1f}GB",
                file=sys.stderr,
            )

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
