#!/usr/bin/env python3
"""
Automatic Image Captioning Script

Automatically generate captions for training images using vision-language models.
Supports BLIP, BLIP-2, and GIT models.
"""

import argparse
import logging
from pathlib import Path
from typing import Optional, List

import torch
from PIL import Image
from tqdm.auto import tqdm
from transformers import (
    BlipProcessor,
    BlipForConditionalGeneration,
    AutoProcessor,
    AutoModelForVision2Seq,
)

# Setup logging
logging.basicConfig(
    format="%(asctime)s - %(levelname)s - %(message)s",
    level=logging.INFO,
)
logger = logging.getLogger(__name__)


class ImageCaptioner:
    """
    Automatic image captioning using vision-language models.
    """
    
    SUPPORTED_MODELS = {
        'blip-base': 'Salesforce/blip-image-captioning-base',
        'blip-large': 'Salesforce/blip-image-captioning-large',
        'blip2': 'Salesforce/blip2-opt-2.7b',
        'git-large': 'microsoft/git-large-coco',
    }
    
    def __init__(
        self,
        model_name: str = 'blip-base',
        device: Optional[str] = None,
    ):
        """
        Initialize the captioner.
        
        Args:
            model_name: Name of the model to use (blip-base, blip-large, blip2, git-large)
            device: Device to run on (cuda/cpu), auto-detected if None
        """
        if model_name not in self.SUPPORTED_MODELS:
            raise ValueError(
                f"Model {model_name} not supported. "
                f"Choose from: {list(self.SUPPORTED_MODELS.keys())}"
            )
        
        self.model_name = model_name
        self.model_path = self.SUPPORTED_MODELS[model_name]
        
        # Auto-detect device
        if device is None:
            self.device = "cuda" if torch.cuda.is_available() else "cpu"
        else:
            self.device = device
        
        logger.info(f"Loading model: {self.model_path}")
        logger.info(f"Using device: {self.device}")
        
        # Load model and processor
        if 'blip' in model_name.lower():
            self.processor = BlipProcessor.from_pretrained(self.model_path)
            self.model = BlipForConditionalGeneration.from_pretrained(
                self.model_path,
                torch_dtype=torch.float16 if self.device == "cuda" else torch.float32,
            ).to(self.device)
        else:
            self.processor = AutoProcessor.from_pretrained(self.model_path)
            self.model = AutoModelForVision2Seq.from_pretrained(
                self.model_path,
                torch_dtype=torch.float16 if self.device == "cuda" else torch.float32,
            ).to(self.device)
        
        self.model.eval()
        logger.info("Model loaded successfully")
    
    def caption_image(
        self,
        image_path: Path,
        max_length: int = 75,
        num_beams: int = 5,
    ) -> str:
        """
        Generate a caption for a single image.
        
        Args:
            image_path: Path to the image
            max_length: Maximum caption length
            num_beams: Number of beams for beam search
            
        Returns:
            Generated caption as string
        """
        try:
            image = Image.open(image_path).convert('RGB')
        except Exception as e:
            logger.error(f"Error loading image {image_path}: {e}")
            return ""
        
        # Process image
        inputs = self.processor(image, return_tensors="pt").to(self.device)
        
        # Generate caption
        with torch.no_grad():
            output_ids = self.model.generate(
                **inputs,
                max_length=max_length,
                num_beams=num_beams,
            )
        
        caption = self.processor.decode(output_ids[0], skip_special_tokens=True)
        return caption.strip()
    
    def caption_batch(
        self,
        image_paths: List[Path],
        batch_size: int = 4,
        max_length: int = 75,
        num_beams: int = 5,
    ) -> List[str]:
        """
        Generate captions for multiple images.
        
        Args:
            image_paths: List of image paths
            batch_size: Batch size for processing
            max_length: Maximum caption length
            num_beams: Number of beams for beam search
            
        Returns:
            List of generated captions
        """
        captions = []
        
        # Process in batches
        for i in tqdm(range(0, len(image_paths), batch_size), desc="Generating captions"):
            batch_paths = image_paths[i:i + batch_size]
            batch_images = []
            
            for path in batch_paths:
                try:
                    image = Image.open(path).convert('RGB')
                    batch_images.append(image)
                except Exception as e:
                    logger.error(f"Error loading image {path}: {e}")
                    batch_images.append(Image.new('RGB', (224, 224)))
            
            # Process batch
            inputs = self.processor(batch_images, return_tensors="pt", padding=True).to(self.device)
            
            # Generate captions
            with torch.no_grad():
                output_ids = self.model.generate(
                    **inputs,
                    max_length=max_length,
                    num_beams=num_beams,
                )
            
            # Decode captions
            batch_captions = self.processor.batch_decode(output_ids, skip_special_tokens=True)
            captions.extend([caption.strip() for caption in batch_captions])
        
        return captions


def find_images(dataset_path: Path) -> List[Path]:
    """Find all images in dataset directory"""
    valid_extensions = {'.jpg', '.jpeg', '.png', '.webp', '.bmp'}
    image_files = []
    
    for ext in valid_extensions:
        image_files.extend(dataset_path.glob(f'*{ext}'))
        image_files.extend(dataset_path.glob(f'*{ext.upper()}'))
    
    return sorted(image_files)


def save_caption(
    image_path: Path,
    caption: str,
    prefix: str = "",
    suffix: str = "",
    trigger_word: Optional[str] = None,
):
    """
    Save caption to .txt file.
    
    Args:
        image_path: Path to the image
        caption: Generated caption
        prefix: Text to prepend to caption
        suffix: Text to append to caption
        trigger_word: Trigger word to insert at beginning
    """
    caption_path = image_path.with_suffix('.txt')
    
    # Build final caption
    parts = []
    if trigger_word:
        parts.append(trigger_word)
    if prefix:
        parts.append(prefix)
    parts.append(caption)
    if suffix:
        parts.append(suffix)
    
    final_caption = " ".join(parts)
    
    with open(caption_path, 'w', encoding='utf-8') as f:
        f.write(final_caption)


def main():
    parser = argparse.ArgumentParser(
        description="Automatically generate captions for training images"
    )
    
    parser.add_argument(
        "--dataset_path",
        type=str,
        required=True,
        help="Path to dataset directory containing images",
    )
    parser.add_argument(
        "--model",
        type=str,
        default="blip-base",
        choices=['blip-base', 'blip-large', 'blip2', 'git-large'],
        help="Captioning model to use",
    )
    parser.add_argument(
        "--batch_size",
        type=int,
        default=4,
        help="Batch size for processing (reduce if OOM)",
    )
    parser.add_argument(
        "--max_length",
        type=int,
        default=75,
        help="Maximum caption length",
    )
    parser.add_argument(
        "--num_beams",
        type=int,
        default=5,
        help="Number of beams for beam search",
    )
    parser.add_argument(
        "--trigger_word",
        type=str,
        default=None,
        help="Trigger word to add at beginning of captions",
    )
    parser.add_argument(
        "--prefix",
        type=str,
        default="",
        help="Prefix to add to all captions",
    )
    parser.add_argument(
        "--suffix",
        type=str,
        default="",
        help="Suffix to add to all captions",
    )
    parser.add_argument(
        "--overwrite",
        action="store_true",
        help="Overwrite existing caption files",
    )
    parser.add_argument(
        "--device",
        type=str,
        default=None,
        help="Device to use (cuda/cpu)",
    )
    
    args = parser.parse_args()
    
    dataset_path = Path(args.dataset_path)
    
    if not dataset_path.exists():
        logger.error(f"Dataset path does not exist: {dataset_path}")
        return
    
    # Find images
    logger.info(f"Searching for images in {dataset_path}")
    image_files = find_images(dataset_path)
    
    if len(image_files) == 0:
        logger.error("No images found!")
        return
    
    logger.info(f"Found {len(image_files)} images")
    
    # Filter out images that already have captions (if not overwriting)
    if not args.overwrite:
        images_to_caption = [
            img for img in image_files
            if not img.with_suffix('.txt').exists()
        ]
        
        if len(images_to_caption) < len(image_files):
            skipped = len(image_files) - len(images_to_caption)
            logger.info(f"Skipping {skipped} images with existing captions")
            logger.info("Use --overwrite to regenerate all captions")
    else:
        images_to_caption = image_files
    
    if len(images_to_caption) == 0:
        logger.info("All images already have captions!")
        return
    
    logger.info(f"Captioning {len(images_to_caption)} images")
    
    # Initialize captioner
    captioner = ImageCaptioner(model_name=args.model, device=args.device)
    
    # Generate captions
    if args.batch_size > 1 and len(images_to_caption) > 1:
        # Batch processing
        captions = captioner.caption_batch(
            images_to_caption,
            batch_size=args.batch_size,
            max_length=args.max_length,
            num_beams=args.num_beams,
        )
        
        # Save captions
        logger.info("Saving captions...")
        for image_path, caption in tqdm(
            zip(images_to_caption, captions),
            total=len(images_to_caption),
            desc="Saving"
        ):
            save_caption(
                image_path,
                caption,
                prefix=args.prefix,
                suffix=args.suffix,
                trigger_word=args.trigger_word,
            )
    else:
        # Single image processing
        for image_path in tqdm(images_to_caption, desc="Processing"):
            caption = captioner.caption_image(
                image_path,
                max_length=args.max_length,
                num_beams=args.num_beams,
            )
            
            save_caption(
                image_path,
                caption,
                prefix=args.prefix,
                suffix=args.suffix,
                trigger_word=args.trigger_word,
            )
    
    logger.info("✓ Captioning complete!")
    logger.info(f"Generated {len(images_to_caption)} captions")
    logger.info("\nExample captions:")
    
    # Show some examples
    for image_path in image_files[:3]:
        caption_path = image_path.with_suffix('.txt')
        if caption_path.exists():
            with open(caption_path, 'r', encoding='utf-8') as f:
                caption = f.read().strip()
            logger.info(f"  {image_path.name}: {caption}")


if __name__ == "__main__":
    main()
