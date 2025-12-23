#!/usr/bin/env python3
"""
Flux LoRA Training - Interactive Quick Start

This script guides you through the entire process interactively.
"""

import os
import sys
from pathlib import Path
import subprocess


def print_header(text):
    """Print a formatted header"""
    print("\n" + "=" * 70)
    print(f"  {text}")
    print("=" * 70 + "\n")


def ask_yes_no(question, default=True):
    """Ask a yes/no question"""
    suffix = "[Y/n]" if default else "[y/N]"
    while True:
        response = input(f"{question} {suffix}: ").strip().lower()
        if response == "":
            return default
        if response in ["y", "yes"]:
            return True
        if response in ["n", "no"]:
            return False
        print("Please answer 'y' or 'n'")


def run_command(cmd, check=True):
    """Run a shell command"""
    print(f"Running: {cmd}")
    result = subprocess.run(cmd, shell=True, check=False)
    if check and result.returncode != 0:
        print(f"⚠️  Command failed with code {result.returncode}")
        return False
    return True


def main():
    print_header("Flux LoRA Training - Quick Start Wizard")
    
    print("This wizard will guide you through:")
    print("  1. Setting up your dataset")
    print("  2. Generating captions")
    print("  3. Validating your data")
    print("  4. Configuring training")
    print("  5. Starting training (or showing you how)")
    print()
    
    if not ask_yes_no("Continue?", True):
        print("Exiting.")
        return
    
    # Step 1: Dataset setup
    print_header("Step 1: Dataset Setup")
    
    print("First, let's set up your training images.")
    print()
    
    dataset_name = input("Enter a name for your subject (e.g., 'my_dog', 'my_face'): ").strip()
    if not dataset_name:
        dataset_name = "my_subject"
    
    dataset_path = Path("dataset") / dataset_name
    dataset_path.mkdir(parents=True, exist_ok=True)
    
    print(f"\n✓ Created dataset directory: {dataset_path}")
    print()
    print("Now, please copy 10-50 images of your subject to:")
    print(f"  {dataset_path.absolute()}")
    print()
    print("Tips:")
    print("  - Use high quality images (1024x1024 or larger)")
    print("  - Include variety (different angles, lighting, backgrounds)")
    print("  - Supported formats: .jpg, .png, .webp")
    print()
    
    if not ask_yes_no("Have you copied your images?", False):
        print()
        print("Please copy your images and run this script again.")
        print(f"  cp /path/to/your/images/* {dataset_path}/")
        return
    
    # Check if images exist
    image_files = list(dataset_path.glob("*.jpg")) + list(dataset_path.glob("*.png")) + list(dataset_path.glob("*.webp"))
    if len(image_files) == 0:
        print(f"\n❌ No images found in {dataset_path}")
        print("Please add images and try again.")
        return
    
    print(f"\n✓ Found {len(image_files)} images")
    
    # Step 2: Trigger word
    print_header("Step 2: Trigger Word")
    
    print("Choose a unique trigger word to identify your subject.")
    print("This word should be:")
    print("  - Unique (not a common word)")
    print("  - Easy to remember")
    print("  - UPPERCASE by convention")
    print()
    print("Examples: SKSDOG, MYFACE, ARTXYZ, TOK123")
    print()
    
    trigger_word = input("Enter your trigger word: ").strip().upper()
    if not trigger_word:
        trigger_word = "MYSUBJECT"
    
    print(f"\n✓ Trigger word: {trigger_word}")
    
    # Step 3: Captions
    print_header("Step 3: Generate Captions")
    
    print("Now we'll generate captions for your images.")
    print("Captions help the model understand what's in each image.")
    print()
    
    if ask_yes_no("Auto-generate captions using AI?", True):
        print("\nGenerating captions...")
        cmd = f"python caption_images.py --dataset_path {dataset_path} --trigger_word {trigger_word}"
        
        if not run_command(cmd, check=False):
            print("\n⚠️  Caption generation had issues.")
            print("You can manually create .txt files for each image.")
        else:
            print("\n✓ Captions generated!")
    else:
        print("\nYou'll need to manually create caption files.")
        print(f"For each image (e.g., image.jpg), create a text file (image.txt) with:")
        print(f"  {trigger_word} [description of the image]")
    
    # Step 4: Validation
    print_header("Step 4: Validate Dataset")
    
    print("Let's validate your dataset is ready for training.")
    print()
    
    if ask_yes_no("Run validation?", True):
        cmd = f"python prepare_dataset.py --dataset_path {dataset_path} --trigger_word {trigger_word}"
        run_command(cmd, check=False)
        print()
        
        if not ask_yes_no("Does the validation look good?", True):
            print("\nPlease fix any issues and run the wizard again.")
            return
    
    # Step 5: Training configuration
    print_header("Step 5: Training Configuration")
    
    print("Let's configure the training parameters.")
    print()
    
    # Ask about GPU
    print("What GPU do you have?")
    print("  1. 12GB VRAM (RTX 3060, RTX 4060)")
    print("  2. 16GB VRAM (RTX 4060 Ti)")
    print("  3. 24GB+ VRAM (RTX 4090, A5000)")
    print("  4. Not sure / Other")
    
    gpu_choice = input("Enter number [1-4]: ").strip()
    
    # Set training parameters based on GPU
    if gpu_choice == "3":
        batch_size = 2
        grad_accum = 2
        train_steps = 1000
    elif gpu_choice == "2":
        batch_size = 1
        grad_accum = 4
        train_steps = 1000
    else:
        batch_size = 1
        grad_accum = 4
        train_steps = 500  # Conservative for safety
    
    # Calculate based on dataset size
    if len(image_files) < 15:
        train_steps = 500
    elif len(image_files) > 30:
        train_steps = 1500
    
    print(f"\nRecommended settings:")
    print(f"  Batch size: {batch_size}")
    print(f"  Gradient accumulation: {grad_accum}")
    print(f"  Training steps: {train_steps}")
    print()
    
    # Create config
    project_name = dataset_name + "_lora"
    output_dir = f"outputs/{project_name}"
    
    config_content = f"""# Auto-generated configuration
project_name: "{project_name}"
trigger_word: "{trigger_word}"

# Model
model_name: "black-forest-labs/FLUX.1-dev"

# Dataset
dataset_path: "{dataset_path}"
output_dir: "{output_dir}"

# Training
max_train_steps: {train_steps}
learning_rate: 1.0e-4
train_batch_size: {batch_size}
gradient_accumulation_steps: {grad_accum}

# LoRA
lora_rank: 16
lora_alpha: 16

# Memory optimization
mixed_precision: "bf16"
gradient_checkpointing: true
use_8bit_adam: true

# Checkpoints
save_steps: 200
save_total_limit: 3
logging_steps: 50

# Image settings
resolution: 1024
center_crop: false
random_flip: false

# Seed
seed: 42
"""
    
    config_path = Path(f"config_{dataset_name}.yaml")
    with open(config_path, 'w') as f:
        f.write(config_content)
    
    print(f"✓ Configuration saved to: {config_path}")
    
    # Step 6: Training
    print_header("Step 6: Start Training")
    
    print("Everything is ready to start training!")
    print()
    print("⚠️  IMPORTANT: The provided train_lora.py is a framework.")
    print("For actual training, you have two options:")
    print()
    print("Option A (Recommended): Use ai-toolkit")
    print("  git clone https://github.com/ostris/ai-toolkit")
    print("  cd ai-toolkit")
    print("  # Copy your dataset and use their training script")
    print()
    print("Option B: Complete the training loop yourself")
    print("  See IMPLEMENTATION_GUIDE.md for details")
    print()
    
    if ask_yes_no("Would you like to see the training command?", True):
        print()
        print("To train with the framework (after completing it):")
        print(f"  python train_lora.py --config {config_path}")
        print()
        print("To test after training:")
        print(f"  python test_lora.py \\")
        print(f"    --lora_path {output_dir}/final.safetensors \\")
        print(f"    --prompt '{trigger_word} as an astronaut' \\")
        print(f"    --output test_output.png")
    
    # Summary
    print_header("Setup Complete!")
    
    print("Summary:")
    print(f"  ✓ Dataset: {len(image_files)} images in {dataset_path}")
    print(f"  ✓ Trigger word: {trigger_word}")
    print(f"  ✓ Configuration: {config_path}")
    print(f"  ✓ Output directory: {output_dir}")
    print()
    print("Next steps:")
    print("  1. Review the configuration file")
    print("  2. Choose a training backend (see IMPLEMENTATION_GUIDE.md)")
    print("  3. Start training!")
    print("  4. Test your LoRA")
    print()
    print("For detailed help, see:")
    print("  - README.md - Complete guide")
    print("  - IMPLEMENTATION_GUIDE.md - Training implementation")
    print("  - TROUBLESHOOTING.md - Common issues")
    print()
    print("Good luck with your training! 🚀")


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\nInterrupted by user. Exiting.")
        sys.exit(0)
    except Exception as e:
        print(f"\n❌ Error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
