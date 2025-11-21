#!/usr/bin/env python3
"""
Multi-Model Video Generation Backend
Supports multiple video generation models dynamically
"""

import os
import sys
import logging
import torch
from pathlib import Path
from flask import Flask, request, jsonify
from datetime import datetime
import uuid

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Initialize Flask app
app = Flask(__name__)
app.config['MAX_CONTENT_LENGTH'] = 100 * 1024 * 1024
app.config['UPLOAD_FOLDER'] = './uploads'
app.config['OUTPUT_FOLDER'] = './outputs'

os.makedirs(app.config['UPLOAD_FOLDER'], exist_ok=True)
os.makedirs(app.config['OUTPUT_FOLDER'], exist_ok=True)

# Model configuration
MODEL_CONFIGS = {
    'ltx-video': {
        'model_id': 'Lightricks/LTX-Video',
        'min_vram': 12,
        'type': 'diffusers',
        'description': 'LTX-Video (Wan2.1) - High quality, needs 12GB+ VRAM'
    },
    'modelscope': {
        'model_id': 'damo-vilab/text-to-video-ms-1.7b',
        'min_vram': 4,
        'type': 'diffusers',
        'description': 'ModelScope - Works on 4GB GPU, good quality'
    },
    'zeroscope': {
        'model_id': 'cerspense/zeroscope_v2_576w',
        'min_vram': 4,
        'type': 'diffusers',
        'description': 'ZeroScope v2 - 576x320, works on 4GB GPU'
    },
    'svd': {
        'model_id': 'stabilityai/stable-video-diffusion-img2vid',
        'min_vram': 12,
        'type': 'diffusers',
        'description': 'Stable Video Diffusion - Image to video'
    },
    'svd-xt': {
        'model_id': 'stabilityai/stable-video-diffusion-img2vid-xt',
        'min_vram': 12,
        'type': 'diffusers',
        'description': 'Stable Video Diffusion XT - Extended frames'
    }
}

# Current model from environment
CURRENT_MODEL = os.getenv('VIDEO_MODEL', 'ltx-video')
CACHE_DIR = os.getenv('MODEL_CACHE_DIR', './models')
ENABLE_GPU = os.getenv('ENABLE_GPU', 'true').lower() == 'true'
GPU_DEVICE_ID = int(os.getenv('GPU_DEVICE_ID', '0'))
LOW_MEMORY_MODE = os.getenv('LOW_MEMORY_MODE', 'false').lower() == 'true'

# Global model instance
current_pipeline = None
current_model_name = None


class MultiModelVideoGenerator:
    """Wrapper for multiple video generation models"""
    
    def __init__(self, model_name='ltx-video', cache_dir='./models', use_gpu=True, device_id=0):
        self.model_name = model_name
        self.model_config = MODEL_CONFIGS.get(model_name)
        
        if not self.model_config:
            raise ValueError(f"Unknown model: {model_name}. Available: {list(MODEL_CONFIGS.keys())}")
        
        self.model_id = self.model_config['model_id']
        self.cache_dir = cache_dir
        self.device = self._setup_device(use_gpu, device_id)
        self.pipeline = None
        self.low_memory_mode = LOW_MEMORY_MODE
        
    def _setup_device(self, use_gpu, device_id):
        """Setup compute device (GPU/CPU)"""
        if use_gpu and torch.cuda.is_available():
            device = f"cuda:{device_id}"
            logger.info(f"Using GPU: {torch.cuda.get_device_name(device_id)}")
            vram_gb = torch.cuda.get_device_properties(device_id).total_memory / 1024**3
            logger.info(f"GPU Memory: {vram_gb:.2f} GB")
            
            # Check if enough VRAM
            min_vram = self.model_config.get('min_vram', 12)
            if vram_gb < min_vram:
                logger.warning(f"Low VRAM! Model needs {min_vram}GB, you have {vram_gb:.1f}GB")
        else:
            device = "cpu"
            logger.info("Using CPU (GPU not available or disabled)")
        return device
    
    def load_model(self):
        """Load the video generation model"""
        try:
            logger.info(f"Loading model: {self.model_name} ({self.model_id})")
            
            from diffusers import DiffusionPipeline
            from diffusers.utils import export_to_video
            
            # Model-specific loading
            if self.model_name == 'svd' or self.model_name == 'svd-xt':
                # Stable Video Diffusion needs special pipeline
                from diffusers import StableVideoDiffusionPipeline
                self.pipeline = StableVideoDiffusionPipeline.from_pretrained(
                    self.model_id,
                    cache_dir=self.cache_dir,
                    torch_dtype=torch.float16 if self.device.startswith('cuda') else torch.float32,
                )
            else:
                # Standard diffusers pipeline
                self.pipeline = DiffusionPipeline.from_pretrained(
                    self.model_id,
                    cache_dir=self.cache_dir,
                    torch_dtype=torch.float16 if self.device.startswith('cuda') else torch.float32,
                )
            
            # Move to device
            self.pipeline = self.pipeline.to(self.device)
            
            # Enable optimizations
            if self.device.startswith('cuda'):
                self._enable_optimizations()
            
            logger.info(f"Model {self.model_name} loaded successfully")
            return True
            
        except Exception as e:
            logger.error(f"Failed to load model: {e}")
            return False
    
    def _enable_optimizations(self):
        """Enable GPU optimizations"""
        try:
            self.pipeline.enable_xformers_memory_efficient_attention()
            logger.info("Enabled xformers memory efficient attention")
        except Exception as e:
            logger.warning(f"Could not enable xformers: {e}")
        
        if self.low_memory_mode:
            try:
                self.pipeline.enable_attention_slicing(1)
                logger.info("Enabled attention slicing for low memory")
            except Exception as e:
                logger.warning(f"Could not enable attention slicing: {e}")
            
            try:
                self.pipeline.enable_vae_slicing()
                logger.info("Enabled VAE slicing for low memory")
            except Exception as e:
                logger.warning(f"Could not enable VAE slicing: {e}")
    
    def generate_text_to_video(self, prompt, negative_prompt="", num_frames=64,
                              fps=24, width=512, height=512, seed=-1,
                              guidance_scale=7.5, num_inference_steps=50):
        """Generate video from text prompt"""
        try:
            if self.pipeline is None:
                raise Exception("Model not loaded")
            
            generator = None
            if seed > 0:
                generator = torch.Generator(device=self.device).manual_seed(seed)
            
            logger.info(f"Generating with {self.model_name}: '{prompt[:50]}...'")
            
            # Model-specific generation
            if self.model_name == 'modelscope':
                # ModelScope has different parameter names
                output = self.pipeline(
                    prompt=prompt,
                    negative_prompt=negative_prompt,
                    num_frames=num_frames,
                    height=height,
                    width=width,
                    num_inference_steps=num_inference_steps,
                    guidance_scale=guidance_scale,
                    generator=generator,
                )
            elif self.model_name == 'zeroscope':
                # ZeroScope fixed resolution
                output = self.pipeline(
                    prompt=prompt,
                    negative_prompt=negative_prompt,
                    num_frames=num_frames,
                    num_inference_steps=num_inference_steps,
                    guidance_scale=guidance_scale,
                    generator=generator,
                )
            else:
                # LTX-Video and others
                output = self.pipeline(
                    prompt=prompt,
                    negative_prompt=negative_prompt,
                    num_frames=num_frames,
                    height=height,
                    width=width,
                    num_inference_steps=num_inference_steps,
                    guidance_scale=guidance_scale,
                    generator=generator,
                )
            
            # Save video
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            job_id = str(uuid.uuid4())[:8]
            output_filename = f"{self.model_name}_text2video_{timestamp}_{job_id}.mp4"
            output_path = os.path.join(app.config['OUTPUT_FOLDER'], output_filename)
            
            from diffusers.utils import export_to_video
            export_to_video(output.frames[0], output_path, fps=fps)
            
            logger.info(f"Video generated: {output_path}")
            
            return {
                'output_path': output_path,
                'frames': num_frames,
                'fps': fps,
                'duration': num_frames / fps,
                'model': self.model_name
            }
            
        except Exception as e:
            logger.error(f"Generation failed: {e}")
            raise
    
    def generate_image_to_video(self, image_path, prompt="", negative_prompt="",
                               num_frames=64, fps=24, width=512, height=512,
                               seed=-1, guidance_scale=7.5, num_inference_steps=50):
        """Generate video from image"""
        try:
            from PIL import Image
            
            if self.pipeline is None:
                raise Exception("Model not loaded")
            
            # Check if model supports image-to-video
            if self.model_name not in ['svd', 'svd-xt', 'ltx-video']:
                raise Exception(f"{self.model_name} doesn't support image-to-video. Use 'svd' or 'ltx-video'")
            
            image = Image.open(image_path).convert('RGB')
            image = image.resize((width, height))
            
            generator = None
            if seed > 0:
                generator = torch.Generator(device=self.device).manual_seed(seed)
            
            logger.info(f"Generating from image with {self.model_name}")
            
            if self.model_name in ['svd', 'svd-xt']:
                # Stable Video Diffusion
                output = self.pipeline(
                    image=image,
                    decode_chunk_size=8,
                    generator=generator,
                )
            else:
                # LTX-Video
                output = self.pipeline(
                    prompt=prompt,
                    image=image,
                    negative_prompt=negative_prompt,
                    num_frames=num_frames,
                    height=height,
                    width=width,
                    num_inference_steps=num_inference_steps,
                    guidance_scale=guidance_scale,
                    generator=generator,
                )
            
            # Save video
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            job_id = str(uuid.uuid4())[:8]
            output_filename = f"{self.model_name}_img2video_{timestamp}_{job_id}.mp4"
            output_path = os.path.join(app.config['OUTPUT_FOLDER'], output_filename)
            
            from diffusers.utils import export_to_video
            export_to_video(output.frames[0], output_path, fps=fps)
            
            logger.info(f"Video generated: {output_path}")
            
            return {
                'output_path': output_path,
                'frames': num_frames,
                'fps': fps,
                'duration': num_frames / fps,
                'model': self.model_name
            }
            
        except Exception as e:
            logger.error(f"Generation failed: {e}")
            raise


# API Routes

@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint"""
    gpu_info = {}
    if torch.cuda.is_available():
        gpu_info = {
            'gpu_available': True,
            'gpu_count': torch.cuda.device_count(),
            'gpu_name': torch.cuda.get_device_name(0),
            'gpu_memory_allocated': f"{torch.cuda.memory_allocated(0) / 1024**3:.2f} GB",
            'gpu_memory_reserved': f"{torch.cuda.memory_reserved(0) / 1024**3:.2f} GB",
        }
    else:
        gpu_info = {'gpu_available': False}
    
    return jsonify({
        'status': 'healthy',
        'model_loaded': current_pipeline is not None,
        'current_model': current_model_name,
        'device': str(current_pipeline.device) if current_pipeline else 'not initialized',
        'available_models': list(MODEL_CONFIGS.keys()),
        **gpu_info
    })


@app.route('/api/models', methods=['GET'])
def list_models():
    """List available models"""
    models_info = []
    for name, config in MODEL_CONFIGS.items():
        models_info.append({
            'name': name,
            'model_id': config['model_id'],
            'min_vram_gb': config['min_vram'],
            'description': config['description'],
            'loaded': name == current_model_name
        })
    
    return jsonify({
        'models': models_info,
        'current_model': current_model_name
    })


@app.route('/api/switch-model', methods=['POST'])
def switch_model():
    """Switch to a different model"""
    global current_pipeline, current_model_name
    
    data = request.get_json()
    new_model = data.get('model_name')
    
    if not new_model:
        return jsonify({'error': 'model_name required'}), 400
    
    if new_model not in MODEL_CONFIGS:
        return jsonify({
            'error': f'Unknown model: {new_model}',
            'available_models': list(MODEL_CONFIGS.keys())
        }), 400
    
    try:
        logger.info(f"Switching from {current_model_name} to {new_model}")
        
        # Unload current model
        if current_pipeline:
            del current_pipeline
            torch.cuda.empty_cache() if torch.cuda.is_available() else None
        
        # Load new model
        generator = MultiModelVideoGenerator(
            model_name=new_model,
            cache_dir=CACHE_DIR,
            use_gpu=ENABLE_GPU,
            device_id=GPU_DEVICE_ID
        )
        
        if generator.load_model():
            current_pipeline = generator
            current_model_name = new_model
            
            return jsonify({
                'message': f'Switched to {new_model}',
                'model': MODEL_CONFIGS[new_model]
            })
        else:
            return jsonify({'error': 'Failed to load model'}), 500
            
    except Exception as e:
        logger.error(f"Error switching model: {e}")
        return jsonify({'error': str(e)}), 500


@app.route('/api/generate/text-to-video', methods=['POST'])
def text_to_video():
    """Generate video from text prompt"""
    try:
        if current_pipeline is None:
            return jsonify({'error': 'No model loaded'}), 503
        
        data = request.get_json()
        
        result = current_pipeline.generate_text_to_video(
            prompt=data.get('prompt', ''),
            negative_prompt=data.get('negative_prompt', ''),
            num_frames=int(data.get('num_frames', 64)),
            fps=int(data.get('fps', 24)),
            width=int(data.get('width', 512)),
            height=int(data.get('height', 512)),
            seed=int(data.get('seed', -1)),
            guidance_scale=float(data.get('guidance_scale', 7.5)),
            num_inference_steps=int(data.get('num_inference_steps', 50)),
        )
        
        return jsonify(result)
        
    except Exception as e:
        logger.error(f"Error in text-to-video: {e}")
        return jsonify({'error': str(e)}), 500


@app.route('/api/generate/image-to-video', methods=['POST'])
def image_to_video():
    """Generate video from image"""
    try:
        if current_pipeline is None:
            return jsonify({'error': 'No model loaded'}), 503
        
        if 'image' not in request.files:
            return jsonify({'error': 'No image file provided'}), 400
        
        from werkzeug.utils import secure_filename
        
        image_file = request.files['image']
        filename = secure_filename(image_file.filename)
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        image_path = os.path.join(app.config['UPLOAD_FOLDER'], f"{timestamp}_{filename}")
        image_file.save(image_path)
        
        result = current_pipeline.generate_image_to_video(
            image_path=image_path,
            prompt=request.form.get('prompt', ''),
            negative_prompt=request.form.get('negative_prompt', ''),
            num_frames=int(request.form.get('num_frames', 64)),
            fps=int(request.form.get('fps', 24)),
            width=int(request.form.get('width', 512)),
            height=int(request.form.get('height', 512)),
            seed=int(request.form.get('seed', -1)),
            guidance_scale=float(request.form.get('guidance_scale', 7.5)),
            num_inference_steps=int(request.form.get('num_inference_steps', 50)),
        )
        
        return jsonify(result)
        
    except Exception as e:
        logger.error(f"Error in image-to-video: {e}")
        return jsonify({'error': str(e)}), 500


def main():
    """Main entry point"""
    global current_pipeline, current_model_name
    
    logger.info("=" * 60)
    logger.info("Multi-Model Video Generation Backend Starting")
    logger.info("=" * 60)
    
    # Initialize with default model
    try:
        logger.info(f"Loading default model: {CURRENT_MODEL}")
        generator = MultiModelVideoGenerator(
            model_name=CURRENT_MODEL,
            cache_dir=CACHE_DIR,
            use_gpu=ENABLE_GPU,
            device_id=GPU_DEVICE_ID
        )
        
        if generator.load_model():
            current_pipeline = generator
            current_model_name = CURRENT_MODEL
            logger.info(f"✅ Model {CURRENT_MODEL} loaded successfully")
        else:
            logger.error("Failed to load default model")
    except Exception as e:
        logger.error(f"Error loading model: {e}")
        logger.info("Starting server without model...")
    
    # Start Flask server
    host = os.getenv('SERVER_HOST', '0.0.0.0')
    port = int(os.getenv('PYTHON_BACKEND_PORT', '5000'))
    
    logger.info(f"Starting server on {host}:{port}")
    logger.info(f"Available models: {list(MODEL_CONFIGS.keys())}")
    app.run(host=host, port=port, debug=False, threaded=True)


if __name__ == '__main__':
    main()
