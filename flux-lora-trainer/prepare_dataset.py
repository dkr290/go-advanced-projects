#!/usr/bin/env python3
"""
Dataset preparation utility for Flux LoRA training.

This script helps prepare and validate your dataset before training.
"""

import argparse
import logging
from pathlib import Path
from typing import List, Tuple
import json

from PIL import Image
from tqdm import tqdm

# Setup logging
logging.basicConfig(
    format="%(asctime)s - %(levelname)s - %(message)s",
    level=logging.INFO,
)
logger = logging.getLogger(__name__)


def find_images(dataset_path: Path) -> List[Path]:
    """Find all valid images in dataset"""
    valid_extensions = {'.jpg', '.jpeg', '.png', '.webp', '.bmp'}
    image_files = []
    
    for ext in valid_extensions:
        image_files.extend(dataset_path.glob(f'*{ext}'))
        image_files.extend(dataset_path.glob(f'*{ext.upper()}'))
    
    return sorted(image_files)


def check_image(image_path: Path) -> Tuple[bool, str, Tuple[int, int]]:
    """
    Check if image is valid.
    
    Returns:
        (is_valid, error_message, dimensions)
    """
    try:
        img = Image.open(image_path)
        width, height = img.size
        
        # Check if image can be loaded
        img.verify()
        
        # Reload after verify
        img = Image.open(image_path)
        
        # Check minimum size
        if width < 512 or height < 512:
            return False, f"Image too small: {width}x{height} (minimum 512x512)", (width, height)
        
        # Check if corrupted
        img.load()
        
        return True, "", (width, height)
        
    except Exception as e:
        return False, str(e), (0, 0)


def check_caption(image_path: Path, trigger_word: str = None) -> Tuple[bool, str, str]:
    """
    Check if caption exists and is valid.
    
    Returns:
        (has_caption, warning_message, caption_text)
    """
    caption_path = image_path.with_suffix('.txt')
    
    if not caption_path.exists():
        return False, "No caption file", ""
    
    try:
        with open(caption_path, 'r', encoding='utf-8') as f:
            caption = f.read().strip()
        
        if not caption:
            return False, "Empty caption", ""
        
        # Check caption length
        if len(caption) < 5:
            return True, "Caption very short", caption
        
        if len(caption) > 200:
            return True, "Caption very long (>200 chars)", caption
        
        # Check for trigger word
        if trigger_word and trigger_word not in caption:
            return True, f"Trigger word '{trigger_word}' not in caption", caption
        
        return True, "", caption
        
    except Exception as e:
        return False, f"Error reading caption: {e}", ""


def analyze_dataset(dataset_path: Path, trigger_word: str = None):
    """Analyze and validate dataset"""
    
    logger.info(f"Analyzing dataset: {dataset_path}")
    logger.info("=" * 70)
    
    # Find images
    image_files = find_images(dataset_path)
    
    if not image_files:
        logger.error("❌ No images found!")
        return
    
    logger.info(f"Found {len(image_files)} images\n")
    
    # Statistics
    stats = {
        'total_images': len(image_files),
        'valid_images': 0,
        'corrupted_images': 0,
        'small_images': 0,
        'with_captions': 0,
        'without_captions': 0,
        'missing_trigger_word': 0,
        'resolutions': {},
    }
    
    issues = []
    
    # Check each image
    logger.info("Checking images and captions...\n")
    
    for image_path in tqdm(image_files, desc="Validating"):
        # Check image
        is_valid, error_msg, (width, height) = check_image(image_path)
        
        if not is_valid:
            stats['corrupted_images'] += 1
            issues.append({
                'file': image_path.name,
                'type': 'image_error',
                'message': error_msg,
            })
            continue
        
        if width < 512 or height < 512:
            stats['small_images'] += 1
        
        stats['valid_images'] += 1
        
        # Track resolutions
        resolution = f"{width}x{height}"
        stats['resolutions'][resolution] = stats['resolutions'].get(resolution, 0) + 1
        
        # Check caption
        has_caption, warning, caption = check_caption(image_path, trigger_word)
        
        if has_caption:
            stats['with_captions'] += 1
            
            if warning:
                issues.append({
                    'file': image_path.name,
                    'type': 'caption_warning',
                    'message': warning,
                })
                
                if trigger_word and trigger_word not in caption:
                    stats['missing_trigger_word'] += 1
        else:
            stats['without_captions'] += 1
            issues.append({
                'file': image_path.name,
                'type': 'caption_missing',
                'message': warning,
            })
    
    # Print results
    logger.info("\n" + "=" * 70)
    logger.info("DATASET ANALYSIS RESULTS")
    logger.info("=" * 70)
    
    logger.info(f"\n📊 Image Statistics:")
    logger.info(f"  Total images: {stats['total_images']}")
    logger.info(f"  ✓ Valid images: {stats['valid_images']}")
    
    if stats['corrupted_images'] > 0:
        logger.warning(f"  ❌ Corrupted images: {stats['corrupted_images']}")
    
    if stats['small_images'] > 0:
        logger.warning(f"  ⚠️  Small images (<512px): {stats['small_images']}")
    
    logger.info(f"\n📝 Caption Statistics:")
    logger.info(f"  ✓ With captions: {stats['with_captions']}")
    
    if stats['without_captions'] > 0:
        logger.warning(f"  ❌ Without captions: {stats['without_captions']}")
    
    if trigger_word and stats['missing_trigger_word'] > 0:
        logger.warning(f"  ⚠️  Missing trigger word '{trigger_word}': {stats['missing_trigger_word']}")
    
    # Resolution distribution
    logger.info(f"\n🖼️  Resolution Distribution:")
    for resolution, count in sorted(stats['resolutions'].items(), key=lambda x: -x[1])[:10]:
        logger.info(f"  {resolution}: {count} images")
    
    # Recommendations
    logger.info(f"\n💡 Recommendations:")
    
    if stats['valid_images'] < 10:
        logger.warning("  ⚠️  Less than 10 images - consider adding more for better results")
    elif stats['valid_images'] > 100:
        logger.info("  ℹ️  Large dataset - training may take longer")
    else:
        logger.info("  ✓ Good number of images (10-100)")
    
    if stats['without_captions'] > 0:
        logger.warning(f"  ⚠️  {stats['without_captions']} images missing captions")
        logger.info("     Run: python caption_images.py --dataset_path " + str(dataset_path))
    
    if trigger_word and stats['missing_trigger_word'] > 0:
        logger.warning(f"  ⚠️  Some captions missing trigger word '{trigger_word}'")
        logger.info("     Consider regenerating captions with --trigger_word option")
    
    if stats['corrupted_images'] > 0:
        logger.warning(f"  ❌ {stats['corrupted_images']} corrupted images - remove or replace them")
    
    # Show issues
    if issues:
        logger.info(f"\n⚠️  Issues Found ({len(issues)}):")
        
        # Group by type
        issue_types = {}
        for issue in issues:
            issue_type = issue['type']
            if issue_type not in issue_types:
                issue_types[issue_type] = []
            issue_types[issue_type].append(issue)
        
        # Show first few of each type
        for issue_type, type_issues in issue_types.items():
            logger.info(f"\n  {issue_type.replace('_', ' ').title()}:")
            for issue in type_issues[:5]:
                logger.info(f"    - {issue['file']}: {issue['message']}")
            
            if len(type_issues) > 5:
                logger.info(f"    ... and {len(type_issues) - 5} more")
    
    logger.info("\n" + "=" * 70)
    
    # Save report
    report_path = dataset_path / 'dataset_report.json'
    with open(report_path, 'w') as f:
        json.dump({
            'statistics': stats,
            'issues': issues,
        }, f, indent=2)
    
    logger.info(f"📄 Full report saved to: {report_path}")
    
    # Overall status
    logger.info("\n" + "=" * 70)
    
    ready_for_training = (
        stats['valid_images'] >= 10 and
        stats['corrupted_images'] == 0 and
        stats['without_captions'] == 0
    )
    
    if ready_for_training:
        logger.info("✅ Dataset is ready for training!")
    else:
        logger.warning("⚠️  Dataset has issues - address them before training")
        logger.info("\nTo prepare dataset:")
        logger.info("1. Remove corrupted images")
        logger.info("2. Generate captions: python caption_images.py --dataset_path " + str(dataset_path))
        logger.info("3. Run this script again to verify")


def main():
    parser = argparse.ArgumentParser(
        description="Prepare and validate dataset for Flux LoRA training"
    )
    
    parser.add_argument(
        "--dataset_path",
        type=str,
        required=True,
        help="Path to dataset directory",
    )
    parser.add_argument(
        "--trigger_word",
        type=str,
        default=None,
        help="Trigger word to check for in captions",
    )
    
    args = parser.parse_args()
    
    dataset_path = Path(args.dataset_path)
    
    if not dataset_path.exists():
        logger.error(f"Dataset path does not exist: {dataset_path}")
        return
    
    if not dataset_path.is_dir():
        logger.error(f"Dataset path is not a directory: {dataset_path}")
        return
    
    analyze_dataset(dataset_path, args.trigger_word)


if __name__ == "__main__":
    main()
