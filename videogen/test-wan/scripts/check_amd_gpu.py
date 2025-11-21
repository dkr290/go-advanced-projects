#!/usr/bin/env python3
"""
AMD GPU Detection and Configuration Helper
"""

import sys
import subprocess
import os

def check_amd_gpu():
    """Check if AMD GPU is present"""
    try:
        result = subprocess.run(['lspci'], capture_output=True, text=True)
        if 'AMD' in result.stdout and 'VGA' in result.stdout:
            print("✅ AMD GPU detected")
            for line in result.stdout.split('\n'):
                if 'AMD' in line and 'VGA' in line:
                    print(f"   {line.strip()}")
            return True
        else:
            print("❌ No AMD GPU detected")
            return False
    except:
        print("⚠️  Could not detect GPU")
        return False

def check_rocm():
    """Check if ROCm is installed"""
    try:
        result = subprocess.run(['rocm-smi'], capture_output=True, text=True)
        if result.returncode == 0:
            print("✅ ROCm is installed")
            return True
        else:
            print("❌ ROCm is not installed")
            return False
    except FileNotFoundError:
        print("❌ ROCm is not installed")
        return False

def check_pytorch_rocm():
    """Check if PyTorch with ROCm support is installed"""
    try:
        import torch
        if torch.cuda.is_available():
            print("✅ PyTorch with GPU support detected")
            print(f"   Device: {torch.cuda.get_device_name(0)}")
            print(f"   VRAM: {torch.cuda.get_device_properties(0).total_memory / 1024**3:.1f}GB")
            return True
        else:
            print("⚠️  PyTorch installed but no GPU detected")
            return False
    except ImportError:
        print("❌ PyTorch is not installed")
        return False

def get_recommended_config():
    """Get recommended configuration"""
    has_gpu = check_amd_gpu()
    has_rocm = check_rocm()
    has_pytorch = check_pytorch_rocm()
    
    print("\n" + "="*60)
    print("RECOMMENDED CONFIGURATION")
    print("="*60)
    
    if has_gpu and has_rocm and has_pytorch:
        print("\n✅ AMD GPU Setup Complete!")
        print("\nAdd to .env:")
        print("  ENABLE_GPU=true")
        print("  GPU_BACKEND=rocm")
        
        # Check VRAM
        try:
            import torch
            vram_gb = torch.cuda.get_device_properties(0).total_memory / 1024**3
            
            if vram_gb < 8:
                print(f"\n⚠️  Low VRAM detected ({vram_gb:.1f}GB)")
                print("  Use CPU mode instead: ENABLE_GPU=false")
            elif vram_gb < 12:
                print(f"\n⚠️  Limited VRAM ({vram_gb:.1f}GB)")
                print("  Recommended settings:")
                print("    MAX_FRAMES=32")
                print("    DEFAULT_WIDTH=256")
                print("    DEFAULT_HEIGHT=256")
                print("    LOW_MEMORY_MODE=true")
            else:
                print(f"\n✅ Good VRAM ({vram_gb:.1f}GB)")
                print("  You can use standard settings")
        except:
            pass
            
    elif has_gpu and not has_rocm:
        print("\n⚠️  AMD GPU detected but ROCm not installed")
        print("\nOptions:")
        print("  1. Install ROCm: ./setup_rocm.sh")
        print("  2. Use CPU mode: ENABLE_GPU=false in .env")
        
    elif not has_gpu:
        print("\n⚠️  No AMD GPU detected")
        print("\nUse CPU mode:")
        print("  ENABLE_GPU=false")
        
    print("\n")

def main():
    print("AMD GPU Configuration Helper")
    print("="*60)
    print()
    
    get_recommended_config()

if __name__ == "__main__":
    main()
