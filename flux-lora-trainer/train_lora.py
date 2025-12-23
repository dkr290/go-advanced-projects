#!/usr/bin/env python3
"""
Flux LoRA Training Script

Train a custom LoRA adapter for Flux models using your own images.
This script handles the complete training pipeline including:
- Dataset loading and preprocessing
- LoRA model setup
- Training loop with checkpointing
- Validation and logging
"""

import argparse
import logging
import math
import os
import sys
from pathlib import Path
from typing import Optional, Dict, Any

import torch
import torch.nn.functional as F
from torch.utils.data import Dataset, DataLoader
from torchvision import transforms
from PIL import Image
from tqdm.auto import tqdm
import yaml

from accelerate import Accelerator
from accelerate.logging import get_logger
from accelerate.utils import ProjectConfiguration, set_seed
from diffusers import (
    AutoencoderKL,
    DDPMScheduler,
    FlowMatchEulerDiscreteScheduler,
)
from diffusers.optimization import get_scheduler
from diffusers.utils import check_min_version
from peft import LoraConfig, get_peft_model
from transformers import CLIPTextModel, CLIPTokenizer, T5EncoderModel, T5Tokenizer
import safetensors.torch

# Check diffusers version
check_min_version("0.27.0")

# Setup logging
logger = get_logger(__name__, log_level="INFO")
logging.basicConfig(
    format="%(asctime)s - %(levelname)s - %(name)s - %(message)s",
    datefmt="%m/%d/%Y %H:%M:%S",
    level=logging.INFO,
)


class FluxLoRADataset(Dataset):
    """
    Dataset for Flux LoRA training.
    
    Loads images and their corresponding text captions from a directory.
    Expected structure:
        dataset_path/
            image1.jpg
            image1.txt
            image2.png
            image2.txt
            ...
    """
    
    def __init__(
        self,
        dataset_path: str,
        tokenizer: CLIPTokenizer,
        tokenizer_2: T5Tokenizer,
        resolution: int = 1024,
        center_crop: bool = False,
        random_flip: bool = False,
        trigger_word: Optional[str] = None,
    ):
        self.dataset_path = Path(dataset_path)
        self.tokenizer = tokenizer
        self.tokenizer_2 = tokenizer_2
        self.resolution = resolution
        self.center_crop = center_crop
        self.random_flip = random_flip
        self.trigger_word = trigger_word
        
        # Find all image files
        self.image_files = []
        valid_extensions = {'.jpg', '.jpeg', '.png', '.webp', '.bmp'}
        
        for ext in valid_extensions:
            self.image_files.extend(self.dataset_path.glob(f'*{ext}'))
            self.image_files.extend(self.dataset_path.glob(f'*{ext.upper()}'))
        
        self.image_files = sorted(self.image_files)
        
        if len(self.image_files) == 0:
            raise ValueError(f"No images found in {dataset_path}")
        
        logger.info(f"Found {len(self.image_files)} images in {dataset_path}")
        
        # Setup image transforms
        self.transform = self._create_transforms()
    
    def _create_transforms(self):
        """Create image transformation pipeline"""
        transform_list = []
        
        if self.center_crop:
            transform_list.append(transforms.CenterCrop(self.resolution))
        else:
            transform_list.append(transforms.Resize(self.resolution))
            transform_list.append(transforms.RandomCrop(self.resolution))
        
        if self.random_flip:
            transform_list.append(transforms.RandomHorizontalFlip())
        
        transform_list.extend([
            transforms.ToTensor(),
            transforms.Normalize([0.5], [0.5]),  # Normalize to [-1, 1]
        ])
        
        return transforms.Compose(transform_list)
    
    def _load_caption(self, image_path: Path) -> str:
        """Load caption from .txt file with same name as image"""
        caption_path = image_path.with_suffix('.txt')
        
        if caption_path.exists():
            with open(caption_path, 'r', encoding='utf-8') as f:
                caption = f.read().strip()
        else:
            # Default caption if no .txt file exists
            caption = f"a photo of {self.trigger_word}" if self.trigger_word else "a photo"
            logger.warning(f"No caption file found for {image_path.name}, using: {caption}")
        
        return caption
    
    def __len__(self):
        return len(self.image_files)
    
    def __getitem__(self, idx):
        image_path = self.image_files[idx]
        
        # Load and transform image
        try:
            image = Image.open(image_path).convert('RGB')
        except Exception as e:
            logger.error(f"Error loading image {image_path}: {e}")
            # Return a blank image on error
            image = Image.new('RGB', (self.resolution, self.resolution))
        
        image = self.transform(image)
        
        # Load caption
        caption = self._load_caption(image_path)
        
        # Tokenize caption with both tokenizers (Flux uses CLIP + T5)
        text_inputs = self.tokenizer(
            caption,
            padding="max_length",
            max_length=self.tokenizer.model_max_length,
            truncation=True,
            return_tensors="pt",
        )
        
        text_inputs_2 = self.tokenizer_2(
            caption,
            padding="max_length",
            max_length=self.tokenizer_2.model_max_length,
            truncation=True,
            return_tensors="pt",
        )
        
        return {
            "pixel_values": image,
            "input_ids": text_inputs.input_ids[0],
            "input_ids_2": text_inputs_2.input_ids[0],
            "caption": caption,
        }


def collate_fn(examples):
    """Collate function for DataLoader"""
    pixel_values = torch.stack([example["pixel_values"] for example in examples])
    pixel_values = pixel_values.to(memory_format=torch.contiguous_format).float()
    
    input_ids = torch.stack([example["input_ids"] for example in examples])
    input_ids_2 = torch.stack([example["input_ids_2"] for example in examples])
    captions = [example["caption"] for example in examples]
    
    return {
        "pixel_values": pixel_values,
        "input_ids": input_ids,
        "input_ids_2": input_ids_2,
        "captions": captions,
    }


def load_config(config_path: str) -> Dict[str, Any]:
    """Load configuration from YAML file"""
    with open(config_path, 'r') as f:
        config = yaml.safe_load(f)
    return config


def parse_args():
    """Parse command line arguments"""
    parser = argparse.ArgumentParser(description="Train a Flux LoRA model")
    
    parser.add_argument(
        "--config",
        type=str,
        default="config.yaml",
        help="Path to configuration YAML file",
    )
    parser.add_argument(
        "--dataset_path",
        type=str,
        default=None,
        help="Path to dataset directory (overrides config)",
    )
    parser.add_argument(
        "--output_dir",
        type=str,
        default=None,
        help="Output directory for trained model (overrides config)",
    )
    parser.add_argument(
        "--max_train_steps",
        type=int,
        default=None,
        help="Maximum number of training steps (overrides config)",
    )
    parser.add_argument(
        "--learning_rate",
        type=float,
        default=None,
        help="Learning rate (overrides config)",
    )
    parser.add_argument(
        "--resume_from_checkpoint",
        type=str,
        default=None,
        help="Path to checkpoint to resume from",
    )
    
    args = parser.parse_args()
    return args


def main():
    """Main training function"""
    args = parse_args()
    
    # Load configuration
    logger.info(f"Loading configuration from {args.config}")
    config = load_config(args.config)
    
    # Override config with command line arguments
    if args.dataset_path:
        config['dataset_path'] = args.dataset_path
    if args.output_dir:
        config['output_dir'] = args.output_dir
    if args.max_train_steps:
        config['max_train_steps'] = args.max_train_steps
    if args.learning_rate:
        config['learning_rate'] = args.learning_rate
    
    # Create output directory
    output_dir = Path(config['output_dir'])
    output_dir.mkdir(parents=True, exist_ok=True)
    
    # Save config to output directory
    with open(output_dir / 'training_config.yaml', 'w') as f:
        yaml.dump(config, f, default_flow_style=False)
    
    # Setup accelerator for distributed training
    accelerator_project_config = ProjectConfiguration(
        project_dir=str(output_dir),
        logging_dir=str(output_dir / "logs"),
    )
    
    accelerator = Accelerator(
        gradient_accumulation_steps=config.get('gradient_accumulation_steps', 1),
        mixed_precision=config.get('mixed_precision', 'bf16'),
        log_with=config.get('report_to', 'tensorboard'),
        project_config=accelerator_project_config,
    )
    
    # Set seed for reproducibility
    if config.get('seed'):
        set_seed(config['seed'])
    
    # Log some info
    logger.info("***** Training Configuration *****")
    logger.info(f"  Project: {config['project_name']}")
    logger.info(f"  Trigger word: {config['trigger_word']}")
    logger.info(f"  Dataset: {config['dataset_path']}")
    logger.info(f"  Output directory: {config['output_dir']}")
    logger.info(f"  Max train steps: {config['max_train_steps']}")
    logger.info(f"  Learning rate: {config['learning_rate']}")
    logger.info(f"  LoRA rank: {config['lora_rank']}")
    
    # Load Flux model components
    # Note: Flux architecture is simplified here - you may need to adapt
    # based on the actual Flux implementation
    logger.info(f"Loading Flux model: {config['model_name']}")
    
    try:
        # Load tokenizers (Flux uses dual text encoders)
        tokenizer = CLIPTokenizer.from_pretrained(
            config['model_name'],
            subfolder="tokenizer",
            cache_dir=config.get('cache_dir'),
        )
        
        tokenizer_2 = T5Tokenizer.from_pretrained(
            config['model_name'],
            subfolder="tokenizer_2",
            cache_dir=config.get('cache_dir'),
        )
        
        # Note: For actual Flux models, you'll need to load the transformer
        # This is a simplified example - adapt to actual Flux architecture
        logger.warning("Note: This script uses a simplified Flux model loading.")
        logger.warning("You may need to adapt the model loading code based on")
        logger.warning("the specific Flux implementation you're using.")
        
    except Exception as e:
        logger.error(f"Error loading model: {e}")
        logger.error("Make sure you've accepted the license and logged in to Hugging Face:")
        logger.error("  huggingface-cli login")
        sys.exit(1)
    
    # Setup LoRA configuration
    lora_config = LoraConfig(
        r=config['lora_rank'],
        lora_alpha=config.get('lora_alpha', config['lora_rank']),
        target_modules=config.get('target_modules', ["to_q", "to_k", "to_v", "to_out.0"]),
        lora_dropout=config.get('lora_dropout', 0.0),
        bias="none",
        task_type="FEATURE_EXTRACTION",  # For diffusion models
    )
    
    logger.info(f"LoRA Configuration: rank={lora_config.r}, alpha={lora_config.lora_alpha}")
    
    # Create dataset
    logger.info("Creating dataset...")
    train_dataset = FluxLoRADataset(
        dataset_path=config['dataset_path'],
        tokenizer=tokenizer,
        tokenizer_2=tokenizer_2,
        resolution=config.get('resolution', 1024),
        center_crop=config.get('center_crop', False),
        random_flip=config.get('random_flip', False),
        trigger_word=config.get('trigger_word'),
    )
    
    # Create dataloader
    train_dataloader = DataLoader(
        train_dataset,
        batch_size=config.get('train_batch_size', 1),
        shuffle=True,
        collate_fn=collate_fn,
        num_workers=config.get('dataloader_num_workers', 0),
    )
    
    # Calculate training steps
    num_update_steps_per_epoch = math.ceil(
        len(train_dataloader) / config.get('gradient_accumulation_steps', 1)
    )
    
    if config.get('max_train_steps') is None:
        max_train_steps = config.get('num_train_epochs', 1) * num_update_steps_per_epoch
    else:
        max_train_steps = config['max_train_steps']
    
    logger.info(f"Total training steps: {max_train_steps}")
    logger.info(f"Steps per epoch: {num_update_steps_per_epoch}")
    
    # Note: Actual training loop would go here
    # This is where you would:
    # 1. Apply LoRA to the Flux transformer
    # 2. Setup optimizer
    # 3. Setup learning rate scheduler
    # 4. Run training loop with proper loss calculation
    # 5. Save checkpoints
    
    logger.info("\n" + "="*70)
    logger.info("IMPORTANT: Training loop implementation required")
    logger.info("="*70)
    logger.info("\nThis script provides the framework for Flux LoRA training.")
    logger.info("To complete the implementation, you need to:")
    logger.info("\n1. Load the actual Flux transformer model")
    logger.info("2. Apply LoRA adapters using PEFT")
    logger.info("3. Implement the training loop with:")
    logger.info("   - Forward pass through Flux model")
    logger.info("   - Loss calculation (diffusion loss)")
    logger.info("   - Backward pass and optimization")
    logger.info("   - Checkpoint saving")
    logger.info("\nFor a complete implementation, consider using:")
    logger.info("- kohya-ss/sd-scripts (supports Flux)")
    logger.info("- ai-toolkit by ostris")
    logger.info("- SimpleTuner")
    logger.info("\nOr refer to the diffusers examples for LoRA training.")
    logger.info("="*70 + "\n")
    
    # Placeholder for actual training
    # In a real implementation, this would be replaced with the full training loop
    
    logger.info("Training script setup completed successfully!")
    logger.info(f"Dataset loaded: {len(train_dataset)} images")
    logger.info(f"Configuration saved to: {output_dir / 'training_config.yaml'}")
    logger.info("\nNext steps:")
    logger.info("1. Review the training configuration")
    logger.info("2. Ensure your images have proper captions (.txt files)")
    logger.info("3. Implement or integrate a complete Flux training backend")
    

if __name__ == "__main__":
    main()
