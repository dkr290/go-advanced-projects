# Proposed Changes for Go Program - Stable Diffusion Support

## 📋 Overview

I've created two new Python scripts for Stable Diffusion models:
- `scripts/sd_python_generate.py` - Text-to-image for SD 1.5/2.1/SDXL/SD3
- `scripts/sd_python_img2img.py` - Image-to-image for SD 1.5/2.1/SDXL/SD3

These scripts support ALL Stable Diffusion models with advanced optimizations.

---

## 🎯 What the New Python Scripts Support

### ✅ Supported Models
- **SD 1.5** (512x512 optimal)
- **SD 2.1** (768x768 optimal)
- **SDXL** (1024x1024 optimal)
- **SD 3** (1024x1024 optimal)
- **Community models** from Civitai (via safetensors)

### ✅ Optimizations
- **FP16** - 50% less VRAM, auto-enabled
- **CPU Offload** - Low VRAM mode (8GB)
- **Sequential Offload** - Ultra low VRAM mode (4GB)
- **Torch Compile** - 2x faster inference
- **VAE Tiling** - Extra memory savings for SDXL
- **xformers** - Memory efficient attention
- **LoRA Support** - All SD/SDXL LoRAs
- **Safetensors Loading** - Direct Civitai model support

### ✅ Auto-Detection
- Automatically detects model type (SD 1.5/2.1/SDXL/SD3)
- Sets optimal tokenizer settings per model
- Recommends correct resolution

---

## 🔧 Proposed Changes to Go Code

### Option 1: **Minimal Changes** (Recommended - Don't Break Existing)

Add a new flag `--use-sd` to switch between FLUX and SD scripts:

#### In `pkg/config/config.go`:

```go
type CmdConf struct {
    // ... existing fields
    UseSD             bool    // Use Stable Diffusion scripts instead of FLUX
    SafetensorsPath   string  // Path to safetensors file (Civitai models)
    SequentialOffload bool    // Ultra low VRAM mode
    CompileModel      bool    // Torch compile for speed
    DisableSafety     bool    // Disable safety checker
}
```

#### Add flags in `GetFlags()`:

```go
flag.BoolVar(
    &c.UseSD,
    "use-sd",
    false,
    "Use Stable Diffusion models instead of FLUX (SD 1.5/2.1/SDXL/SD3)",
)

flag.StringVar(
    &c.SafetensorsPath,
    "safetensors",
    "",
    "Path to safetensors model file (for Civitai models)",
)

flag.BoolVar(
    &c.SequentialOffload,
    "sequential-offload",
    false,
    "Enable sequential CPU offload for ultra low VRAM (4GB)",
)

flag.BoolVar(
    &c.CompileModel,
    "compile",
    false,
    "Compile model for 2x faster inference (requires PyTorch 2.0+)",
)

flag.BoolVar(
    &c.DisableSafety,
    "disable-safety",
    false,
    "Disable safety checker (allows NSFW content)",
)
```

#### In `pkg/generate/generate.go`:

Add a new function (or modify existing):

```go
// GenerateWithPythonSD calls Python for Stable Diffusion generation
func GenerateWithPythonSD(
    cmdConf config.Config,
    promptConf config.PromptConfig,
    modelPath, loraDir string,
) error {
    scriptPath := filepath.Join("scripts", "sd_python_generate.py")

    // Check if script exists
    if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
        return fmt.Errorf("python SD script not found at %s", scriptPath)
    }

    if err := os.MkdirAll(cmdConf.OutputDir, 0o755); err != nil {
        return fmt.Errorf("failed to create output dir: %w", err)
    }

    // Use HuggingFace model ID directly (no GGUF for SD)
    hfModel := cmdConf.HfModelID

    // Find LoRA file if loraDir is provided
    var loraFilePath string
    if loraDir != "" {
        entries, err := os.ReadDir(loraDir)
        if err == nil {
            for _, e := range entries {
                if strings.HasSuffix(e.Name(), ".safetensors") {
                    loraFilePath = filepath.Join(loraDir, e.Name())
                    break
                }
            }
        }
    }

    fmt.Printf("\nStarting Stable Diffusion generation for %d images...\n", len(promptConf.Prompts))

    var promptsData []PromptData
    for i, p := range promptConf.Prompts {
        prompt := fmt.Sprintf("%s %s", p, promptConf.StyleSuffix)
        filename := utils.SanitizeFilenameForImage(p, i+1)
        promptsData = append(promptsData, PromptData{
            Prompt:   prompt,
            Filename: filename,
            Seed:     cmdConf.Seed + i,
        })
    }

    promptsDataJSON, err := json.Marshal(promptsData)
    if err != nil {
        return fmt.Errorf("failed to marshal prompts data to JSON: %w", err)
    }

    args := []string{
        scriptPath,
        "--model", hfModel,
        "--negative-prompt", promptConf.NegativePrompt,
        "--width", strconv.Itoa(cmdConf.Resolution[0]),
        "--height", strconv.Itoa(cmdConf.Resolution[1]),
        "--steps", strconv.Itoa(cmdConf.Steps),
        "--guidance-scale", fmt.Sprintf("%.2f", cmdConf.GuidanceScale),
        "--output-dir", cmdConf.OutputDir,
        "--prompts-data", string(promptsDataJSON),
    }

    // Add safetensors path if specified
    if cmdConf.SafetensorsPath != "" {
        args = append(args, "--safetensors", cmdConf.SafetensorsPath)
    }

    // Add LoRA if found
    if loraFilePath != "" {
        args = append(args, "--lora-file", loraFilePath)
    }

    // Add low VRAM flag if enabled
    if cmdConf.LowVRAM {
        args = append(args, "--low-vram")
    }

    // Add sequential offload if enabled
    if cmdConf.SequentialOffload {
        args = append(args, "--sequential-offload")
    }

    // Add compile flag if enabled
    if cmdConf.CompileModel {
        args = append(args, "--compile")
    }

    // Add disable safety flag if enabled
    if cmdConf.DisableSafety {
        args = append(args, "--disable-safety")
    }

    start := time.Now()

    cmd := exec.Command("python3", args...)

    stdoutPipe, err := cmd.StdoutPipe()
    if err != nil {
        return fmt.Errorf("failed to create stdout pipe: %w", err)
    }

    cmd.Stderr = os.Stderr

    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start python: %w", err)
    }

    var stdoutBuffer bytes.Buffer
    stdoutMulti := io.MultiWriter(&stdoutBuffer, os.Stdout)
    if _, err := io.Copy(stdoutMulti, stdoutPipe); err != nil {
        return fmt.Errorf("failed to read python output: %w", err)
    }

    if err := cmd.Wait(); err != nil {
        return fmt.Errorf("python execution failed: %v", err)
    }

    var overallResult PythonOverallResult
    if err := json.Unmarshal(stdoutBuffer.Bytes(), &overallResult); err != nil {
        return fmt.Errorf("failed to parse python output: %v", err)
    }

    if overallResult.OverallStatus != "success" {
        return fmt.Errorf("overall generation failed: %s", overallResult.Error)
    }

    for _, res := range overallResult.Generations {
        if res.Status != "success" {
            fmt.Printf("⚠ Generation for prompt index %d failed: %s\n", res.PromptIndex, res.Output)
        } else {
            fmt.Printf("    ✓ Saved to %s (Prompt %d) in %s\n", filepath.Base(res.Output), res.PromptIndex+1, time.Since(start))
        }
    }
    
    fmt.Printf(
        "\nTotal generation time for %d images: %s\n",
        len(promptConf.Prompts),
        time.Since(start),
    )

    return nil
}
```

#### Add similar function for img2img:

```go
// GenerateImg2ImgWithPythonSD calls Python for SD img2img generation
func GenerateImg2ImgWithPythonSD(
    cmdConf config.Config,
    promptConf config.PromptConfig,
    modelPath, loraDir string,
    inputImagesDir string,
) error {
    scriptPath := filepath.Join("scripts", "sd_python_img2img.py")
    
    // Same pattern as above, but:
    // 1. Use sd_python_img2img.py
    // 2. Add --strength parameter
    // 3. Include input_image in PromptData
    
    // ... (similar implementation to GenerateImg2ImgWithPython but for SD)
}
```

#### In `main.go`:

Modify the generation logic:

```go
if cmdConf.ImageToImage {
    llogger.Logging.Infof("Starting Image to image mode")
    inputImagesDir := "./images"

    // Choose SD or FLUX script
    if cmdConf.UseSD {
        if err := generate.GenerateImg2ImgWithPythonSD(cmdConf, promptConf, modelPath, loraDir, inputImagesDir); err != nil {
            llogger.Logging.Errorf("SD image to image generation failed: %v", err)
            os.Exit(1)
        }
    } else {
        if err := generate.GenerateImg2ImgWithPython(cmdConf, promptConf, modelPath, loraDir, inputImagesDir); err != nil {
            llogger.Logging.Errorf("FLUX image to image generation failed: %v", err)
            os.Exit(1)
        }
    }
    fmt.Println("\n✅ Image-to-Image Generation Complete!")

} else {
    // Choose SD or FLUX script
    if cmdConf.UseSD {
        llogger.Logging.Infof("Starting Stable Diffusion model initialization")
        
        if err := generate.GenerateWithPythonSD(cmdConf, promptConf, modelPath, loraDir); err != nil {
            llogger.Logging.Errorf("SD generate images failed %v", err)
            os.Exit(1)
        }
    } else {
        llogger.Logging.Infof(
            "Starting %s model initialization",
            utils.GetFilenameFromURL(cmdConf.ModelURL),
        )

        if err := generate.GenerateWithPython(cmdConf, promptConf, modelPath, loraDir); err != nil {
            llogger.Logging.Errorf("FLUX generate images failed %v", err)
            os.Exit(1)
        }
    }

    fmt.Println("\n✅ Go image Generation Complete!")
}
```

---

## 📝 Usage Examples

### Text-to-Image with SDXL
```bash
go run main.go \
  -config prompts.json \
  --use-sd \
  --hf-model "stabilityai/stable-diffusion-xl-base-1.0" \
  --resolution 1024x1024 \
  --num_steps 40 \
  --guidence_scale 7.5
```

### Text-to-Image with SD 1.5 (Low VRAM)
```bash
go run main.go \
  -config prompts.json \
  --use-sd \
  --hf-model "runwayml/stable-diffusion-v1-5" \
  --resolution 512x512 \
  --low-vram
```

### Text-to-Image with Civitai Model (Safetensors)
```bash
# Download model from Civitai
wget https://civitai.com/api/download/models/xxx -O juggernaut.safetensors

# Use it
go run main.go \
  -config prompts.json \
  --use-sd \
  --hf-model "stabilityai/stable-diffusion-xl-base-1.0" \
  --safetensors ./juggernaut.safetensors \
  --resolution 1024x1024
```

### Image-to-Image with SDXL
```bash
go run main.go \
  -config prompts.json \
  --use-sd \
  --img2img \
  --hf-model "stabilityai/stable-diffusion-xl-base-1.0" \
  --strength 0.7 \
  --resolution 1024x1024
```

### Ultra Low VRAM (4GB) with Sequential Offload
```bash
go run main.go \
  -config prompts.json \
  --use-sd \
  --hf-model "stabilityai/stable-diffusion-xl-base-1.0" \
  --low-vram \
  --sequential-offload
```

### Fast Generation with Torch Compile
```bash
go run main.go \
  -config prompts.json \
  --use-sd \
  --hf-model "stabilityai/stable-diffusion-xl-base-1.0" \
  --compile
```

### FLUX (Original - No Changes)
```bash
# Still works exactly the same
go run main.go \
  -config prompts.json \
  --hf-model "black-forest-labs/FLUX.1-dev"
```

---

## 🎯 Key Differences: FLUX vs SD Scripts

| Feature | FLUX Scripts | SD Scripts |
|---------|-------------|------------|
| **GGUF Support** | ✅ Yes | ❌ No (use safetensors) |
| **FP16 Variant** | No | ✅ Auto-enabled |
| **Model Types** | FLUX only | SD 1.5/2.1/SDXL/SD3 |
| **Auto-Detection** | No | ✅ Yes |
| **Safetensors Loading** | No | ✅ Yes (Civitai) |
| **Sequential Offload** | No | ✅ Yes |
| **Torch Compile** | No | ✅ Yes |
| **VAE Tiling** | No | ✅ Yes (SDXL) |

---

## 📊 Model Recommendations

### For SDXL (1024x1024)
```bash
--use-sd \
--hf-model "stabilityai/stable-diffusion-xl-base-1.0" \
--resolution 1024x1024 \
--num_steps 40 \
--guidence_scale 7.5
```

### For SD 1.5 (512x512, Fast)
```bash
--use-sd \
--hf-model "runwayml/stable-diffusion-v1-5" \
--resolution 512x512 \
--num_steps 25 \
--guidence_scale 7.5
```

### For SD 2.1 (768x768)
```bash
--use-sd \
--hf-model "stabilityai/stable-diffusion-2-1" \
--resolution 768x768 \
--num_steps 50 \
--guidence_scale 7.0
```

### For SD 3 (1024x1024)
```bash
--use-sd \
--hf-model "stabilityai/stable-diffusion-3-medium-diffusers" \
--resolution 1024x1024 \
--num_steps 28 \
--guidence_scale 7.0
```

---

## 🔧 Optional: Model Presets (Future Enhancement)

Instead of manually setting parameters, you could add presets:

```go
var SDModelPresets = map[string]ModelPreset{
    "sd15": {
        ModelID:       "runwayml/stable-diffusion-v1-5",
        Resolution:    []int{512, 512},
        Steps:         25,
        GuidanceScale: 7.5,
    },
    "sd21": {
        ModelID:       "stabilityai/stable-diffusion-2-1",
        Resolution:    []int{768, 768},
        Steps:         50,
        GuidanceScale: 7.0,
    },
    "sdxl": {
        ModelID:       "stabilityai/stable-diffusion-xl-base-1.0",
        Resolution:    []int{1024, 1024},
        Steps:         40,
        GuidanceScale: 7.5,
    },
}
```

Then use:
```bash
go run main.go -config config.json --model-preset sdxl
```

---

## ✅ Summary of Changes Needed

### New Files (Already Created)
- ✅ `scripts/sd_python_generate.py`
- ✅ `scripts/sd_python_img2img.py`

### Files to Modify
1. **`pkg/config/config.go`**
   - Add `UseSD` flag
   - Add `SafetensorsPath` field
   - Add `SequentialOffload` flag
   - Add `CompileModel` flag
   - Add `DisableSafety` flag

2. **`pkg/generate/generate.go`** (or new `pkg/generate/sd_generate.go`)
   - Add `GenerateWithPythonSD()` function
   - Add `GenerateImg2ImgWithPythonSD()` function

3. **`main.go`**
   - Add conditional logic to choose FLUX vs SD scripts

---

## 🎯 Migration Path

### Phase 1: Keep FLUX Working (Don't Break Existing)
- Add `--use-sd` flag
- Default to FLUX (existing behavior)
- SD only when flag is set

### Phase 2: Test SD Support
```bash
# Test SDXL
go run main.go -config config.json --use-sd --hf-model "stabilityai/stable-diffusion-xl-base-1.0"

# Test SD 1.5
go run main.go -config config.json --use-sd --hf-model "runwayml/stable-diffusion-v1-5" --resolution 512x512
```

### Phase 3: Add Safetensors Support
```bash
# Download Civitai model
go run main.go -config config.json --use-sd --safetensors ./model.safetensors
```

### Phase 4: Web UI Integration
```bash
# SD with web UI
go run main.go -config config.json --use-sd --hf-model "stabilityai/stable-diffusion-xl-base-1.0" --web
```

---

## 🚀 Benefits of This Approach

1. ✅ **Non-Breaking** - FLUX still works exactly as before
2. ✅ **Flexible** - Choose FLUX or SD with a flag
3. ✅ **Future-Proof** - Easy to add more models
4. ✅ **Optimized** - Best performance for each model type
5. ✅ **Community Models** - Direct Civitai support
6. ✅ **Low VRAM** - Works on 4GB cards with sequential offload

---

## 📚 Next Steps

1. **Add the new flags to `config.go`**
2. **Create `GenerateWithPythonSD()` function**
3. **Update `main.go` conditional logic**
4. **Test with SDXL first** (most popular)
5. **Test with Civitai models**
6. **Add model presets** (optional enhancement)

Would you like me to implement these Go changes, or do you prefer to do it yourself following this guide?
