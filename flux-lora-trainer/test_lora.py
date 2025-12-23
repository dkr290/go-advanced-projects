#!/usr/bin/env python3
"""
Flux LoRA Testing and Inference Script

Test your trained LoRA model by generating images with custom prompts.
"""

import argparse
import logging
from pathlib import Path

import torch
from PIL import Image
from diffusers import DiffusionPipeline, FlowMatchEulerDiscreteScheduler
from diffusers.utils import load_image
import safetensors.torch

# Setup logging
logging.basicConfig(
    format="%(asctime)s - %(levelname)s - %(message)s",
    level=logging.INFO,
)
logger = logging.getLogger(__name__)


def load_flux_with_lora(
    model_name: str,
    lora_path: str,
    lora_scale: float = 1.0,
    device: str = "cuda",
    dtype: torch.dtype = torch.bfloat16,
):
    """
    Load Flux model with LoRA weights.
    
    Args:
        model_name: Base Flux model name
        lora_path: Path to LoRA weights (.safetensors)
        lora_scale: LoRA strength (0.0-2.0)
        device: Device to load on
        dtype: Data type for model
        
    Returns:
        Loaded pipeline
    """
    logger.info(f"Loading Flux model: {model_name}")
    
    # Load base pipeline
    pipe = DiffusionPipeline.from_pretrained(
        model_name,
        torch_dtype=dtype,
    )
    
    # Load LoRA weights
    logger.info(f"Loading LoRA weights from: {lora_path}")
    
    try:
        pipe.load_lora_weights(lora_path)
        pipe.fuse_lora(lora_scale=lora_scale)
        logger.info(f"LoRA loaded with scale: {lora_scale}")
    except Exception as e:
        logger.warning(f"Could not load LoRA using built-in method: {e}")
        logger.warning("Attempting manual LoRA loading...")
        
        # Manual loading fallback
        lora_state_dict = safetensors.torch.load_file(lora_path)
        # Note: Manual injection would go here - implementation depends on Flux architecture
        logger.warning("Manual LoRA loading not fully implemented in this example")
    
    # Move to device
    pipe = pipe.to(device)
    
    # Enable memory optimizations
    if device == "cuda":
        try:
            pipe.enable_model_cpu_offload()
            logger.info("Enabled CPU offload for memory efficiency")
        except:
            pass
        
        try:
            pipe.enable_attention_slicing()
            logger.info("Enabled attention slicing")
        except:
            pass
    
    return pipe


def generate_image(
    pipe: DiffusionPipeline,
    prompt: str,
    negative_prompt: str = "",
    num_inference_steps: int = 30,
    guidance_scale: float = 7.5,
    width: int = 1024,
    height: int = 1024,
    seed: int = None,
    **kwargs,
) -> Image.Image:
    """
    Generate an image using the pipeline.
    
    Args:
        pipe: Diffusion pipeline
        prompt: Text prompt
        negative_prompt: Negative prompt (what to avoid)
        num_inference_steps: Number of denoising steps
        guidance_scale: Classifier-free guidance scale
        width: Image width
        height: Image height
        seed: Random seed for reproducibility
        **kwargs: Additional arguments for pipeline
        
    Returns:
        Generated PIL Image
    """
    logger.info(f"Generating image with prompt: {prompt}")
    
    # Set seed if provided
    generator = None
    if seed is not None:
        generator = torch.Generator(device=pipe.device).manual_seed(seed)
        logger.info(f"Using seed: {seed}")
    
    # Generate
    with torch.no_grad():
        image = pipe(
            prompt=prompt,
            negative_prompt=negative_prompt,
            num_inference_steps=num_inference_steps,
            guidance_scale=guidance_scale,
            width=width,
            height=height,
            generator=generator,
            **kwargs,
        ).images[0]
    
    return image


def main():
    parser = argparse.ArgumentParser(
        description="Test Flux LoRA model by generating images"
    )
    
    # Model arguments
    parser.add_argument(
        "--model_name",
        type=str,
        default="black-forest-labs/FLUX.1-dev",
        help="Base Flux model name",
    )
    parser.add_argument(
        "--lora_path",
        type=str,
        required=True,
        help="Path to trained LoRA weights (.safetensors)",
    )
    parser.add_argument(
        "--lora_scale",
        type=float,
        default=1.0,
        help="LoRA strength (0.0-2.0, default: 1.0)",
    )
    
    # Generation arguments
    parser.add_argument(
        "--prompt",
        type=str,
        required=True,
        help="Text prompt for generation",
    )
    parser.add_argument(
        "--negative_prompt",
        type=str,
        default="blurry, low quality, distorted, ugly, bad anatomy",
        help="Negative prompt (what to avoid)",
    )
    parser.add_argument(
        "--num_inference_steps",
        type=int,
        default=30,
        help="Number of denoising steps (20-50)",
    )
    parser.add_argument(
        "--guidance_scale",
        type=float,
        default=7.5,
        help="Classifier-free guidance scale (3.0-15.0)",
    )
    parser.add_argument(
        "--width",
        type=int,
        default=1024,
        help="Image width",
    )
    parser.add_argument(
        "--height",
        type=int,
        default=1024,
        help="Image height",
    )
    parser.add_argument(
        "--seed",
        type=int,
        default=None,
        help="Random seed for reproducibility",
    )
    
    # Output arguments
    parser.add_argument(
        "--output",
        type=str,
        default="output.png",
        help="Output image path",
    )
    parser.add_argument(
        "--num_images",
        type=int,
        default=1,
        help="Number of images to generate",
    )
    
    # Device arguments
    parser.add_argument(
        "--device",
        type=str,
        default="cuda",
        help="Device to use (cuda/cpu)",
    )
    parser.add_argument(
        "--dtype",
        type=str,
        default="bfloat16",
        choices=["float32", "float16", "bfloat16"],
        help="Model data type",
    )
    
    args = parser.parse_args()
    
    # Convert dtype string to torch dtype
    dtype_map = {
        "float32": torch.float32,
        "float16": torch.float16,
        "bfloat16": torch.bfloat16,
    }
    dtype = dtype_map[args.dtype]
    
    # Check LoRA path exists
    lora_path = Path(args.lora_path)
    if not lora_path.exists():
        logger.error(f"LoRA path does not exist: {lora_path}")
        return
    
    # Load model with LoRA
    try:
        pipe = load_flux_with_lora(
            model_name=args.model_name,
            lora_path=str(lora_path),
            lora_scale=args.lora_scale,
            device=args.device,
            dtype=dtype,
        )
    except Exception as e:
        logger.error(f"Error loading model: {e}")
        logger.error("\nMake sure you've:")
        logger.error("1. Accepted the model license on Hugging Face")
        logger.error("2. Logged in: huggingface-cli login")
        logger.error("3. Have the correct LoRA file path")
        return
    
    # Generate images
    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    
    for i in range(args.num_images):
        logger.info(f"\nGenerating image {i+1}/{args.num_images}")
        
        # Adjust seed if generating multiple images
        seed = args.seed
        if seed is not None and args.num_images > 1:
            seed = seed + i
        
        # Generate
        try:
            image = generate_image(
                pipe=pipe,
                prompt=args.prompt,
                negative_prompt=args.negative_prompt,
                num_inference_steps=args.num_inference_steps,
                guidance_scale=args.guidance_scale,
                width=args.width,
                height=args.height,
                seed=seed,
            )
            
            # Save
            if args.num_images > 1:
                save_path = output_path.with_stem(f"{output_path.stem}_{i+1:03d}")
            else:
                save_path = output_path
            
            image.save(save_path)
            logger.info(f"✓ Saved to: {save_path}")
            
        except Exception as e:
            logger.error(f"Error generating image: {e}")
            continue
    
    logger.info("\n✓ Generation complete!")
    logger.info(f"Generated {args.num_images} image(s)")
    logger.info(f"\nPrompt: {args.prompt}")
    logger.info(f"LoRA scale: {args.lora_scale}")


if __name__ == "__main__":
    main()
