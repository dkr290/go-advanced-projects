import os
import subprocess

import torch

print("=== RunPod Diagnostic ===")
print(f"RUNPOD_POD_ID: {os.environ.get('RUNPOD_POD_ID', 'Not set')}")
print(f"RUNPOD_API_KEY: {'Set' if os.environ.get('RUNPOD_API_KEY') else 'Not set'}")

# Check GPU
print(f"\nCUDA available: {torch.cuda.is_available()}")
if torch.cuda.is_available():
    print(f"Device count: {torch.cuda.device_count()}")
    for i in range(torch.cuda.device_count()):
        print(f"  Device {i}: {torch.cuda.get_device_name(i)}")

    # Try to create context
    try:
        x = torch.tensor([1.0]).cuda()
        print(f"✓ Can create CUDA tensor: {x.device}")
    except Exception as e:
        print(f"✗ CUDA tensor creation failed: {e}")

# Check nvidia-smi
print("\n=== nvidia-smi ===")
try:
    result = subprocess.run(["nvidia-smi"], capture_output=True, text=True)
    print(result.stdout[:500])  # First 500 chars
except Exception as e:
    print(f"nvidia-smi failed: {e}")

# Check device files
print("\n=== Device Files ===")
for f in ["/dev/nvidia0", "/dev/nvidiactl", "/dev/nvidia-uvm"]:
    if os.path.exists(f):
        print(f"✓ {f} exists")
    else:
        print(f"✗ {f} missing")
