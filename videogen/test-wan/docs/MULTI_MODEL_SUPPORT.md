# Multi-Model Support Guide

## Current Status

**Currently supports:** LTX-Video (Lightricks/LTX-Video) only

**Can be adapted for:** ANY video generation model from Hugging Face or custom models

---

## ✅ Compatible Video Generation Models

### Text-to-Video Models

| Model | Provider | VRAM | Status |
|-------|----------|------|--------|
| **LTX-Video** | Lightricks | 12GB+ | ✅ Currently supported |
| **ModelScope** | Alibaba DAMO | 4GB+ | ✅ Easy to add |
| **ZeroScope v2** | Cerspense | 4GB+ | ✅ Easy to add |
| **AnimateDiff** | Various | 8GB+ | ⚠️ Needs adapter |
| **Text2Video-Zero** | Picsart | 8GB+ | ⚠️ Needs adapter |
| **CogVideo** | Tsinghua | 16GB+ | ⚠️ Needs adapter |
| **VideoCrafter** | Various | 12GB+ | ⚠️ Needs adapter |
| **Show-1** | ShowLab | 16GB+ | ⚠️ Needs adapter |
| **Runway Gen-2** | Runway | API | ⚠️ API integration |
| **Stable Video Diffusion** | Stability AI | 12GB+ | ✅ Easy to add |

### Image-to-Video Models

| Model | VRAM | Status |
|-------|------|--------|
| **Stable Video Diffusion** | 12GB+ | ✅ Easy to add |
| **I2VGen-XL** | 16GB+ | ⚠️ Needs adapter |
| **DynamiCrafter** | 12GB+ | ⚠️ Needs adapter |
| **SEINE** | 12GB+ | ⚠️ Needs adapter |

### Custom Models

| Type | Compatibility |
|------|--------------|
| Custom Diffusers pipelines | ✅ Yes |
| Custom PyTorch models | ✅ Yes (wrapper needed) |
| ONNX models | ⚠️ Possible (needs adapter) |
| TensorFlow models | ⚠️ Possible (needs converter) |

---

## 🏗️ Architecture Overview

The project has a **modular architecture** that separates:

```
┌─────────────────────────────────────────────┐
│  Go Server (HTTP API)                       │
│  - Request handling                         │
│  - File uploads                             │
│  - Job management                           │
│  └─ Model-agnostic                          │
└──────────────┬──────────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────────┐
│  Python Backend (Model Inference)           │
│  - Model loading ← Needs modification       │
│  - Video generation ← Needs modification    │
│  - GPU management                           │
└─────────────────────────────────────────────┘
```

**What needs changing:**
- ✅ Go Server: Nothing! (already generic)
- ⚠️ Python Backend: Model loading code
- ⚠️ Configuration: Model parameters

---

## 🔧 How to Add a New Model

### Example: Adding ModelScope (4GB compatible!)

I'll create the files below to show you how...

---

## Models Comparison

### Currently Implemented: LTX-Video
```
VRAM:     12GB minimum
Quality:  Excellent
Speed:    Medium-slow
Setup:    ./setup.sh
```

### Easy to Add: ModelScope
```
VRAM:     4GB minimum!
Quality:  Good
Speed:    Fast
Setup:    Just change config + small code update
```

### Easy to Add: Stable Video Diffusion
```
VRAM:     12GB minimum
Quality:  Excellent (images → video)
Speed:    Medium
Setup:    Small code update
```

### Harder: AnimateDiff
```
VRAM:     8GB minimum
Quality:  Excellent
Speed:    Medium
Setup:    Needs custom pipeline
```

---

## Which Models Should We Support?

I can create multi-model support for:

### Option A: Popular Models Pack
- ✅ LTX-Video (current)
- ✅ ModelScope (4GB compatible)
- ✅ Stable Video Diffusion
- ✅ ZeroScope v2

### Option B: Low-Memory Models Pack
- ✅ ModelScope (4GB)
- ✅ ZeroScope v2 (4GB)
- ✅ AnimateDiff (8GB)

### Option C: Enterprise Pack
- ✅ All above
- ✅ Custom model loader
- ✅ API integrations (Runway, etc.)
- ✅ Model switching via API

---

## Let me create the multi-model support now...
