#!/usr/bin/env python3
"""
Wan2.1 Video Generation Backend
Python backend for handling LTX-Video model inference with GPU acceleration
"""

import os
import sys
import logging
import torch
from pathlib import Path
from flask import Flask, request, jsonify, send_file
from werkzeug.utils import secure_filename
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
app.config['MAX_CONTENT_LENGTH'] = 100 * 1024 * 1024  # 100MB max file size
app.config['UPLOAD_FOLDER'] = './uploads'
app.config['OUTPUT_FOLDER'] = './outputs'

# Ensure directories exist
os.makedirs(app.config['UPLOAD_FOLDER'], exist_ok=True)
os.makedirs(app.config['OUTPUT_FOLDER'], exist_ok=True)

# Model configuration
MODEL_ID = os.getenv('HUGGINGFACE_MODEL_ID', 'Lightricks/LTX-Video')
CACHE_DIR = os.getenv('MODEL_CACHE_DIR', './models')
ENABLE_GPU = os.getenv('ENABLE_GPU', 'true').lower() == 'true'
GPU_DEVICE_ID = int(os.getenv('GPU_DEVICE_ID', '0'))

# Global model instance
model = None
pipeline = None


class VideoGenerationModel:
    """Wrapper for video generation model"""
    
    def __init__(self, model_id, cache_dir, use_gpu=True, device_id=0):
        self.model_id = model_id
        self.cache_dir = cache_dir
        self.device = self._setup_device(use_gpu, device_id)
        self.pipeline = None
        self.low_memory_mode = os.getenv('LOW_MEMORY_MODE', 'false').lower() == 'true'
        
    def _setup_device(self, use_gpu, device_id):
        """Setup compute device (GPU/CPU)"""
        if use_gpu and torch.cuda.is_available():
            device = f"cuda:{device_id}"
            logger.info(f"Using GPU: {torch.cuda.get_device_name(device_id)}")
            logger.info(f"GPU Memory: {torch.cuda.get_device_properties(device_id).total_memory / 1024**3:.2f} GB")
        else:
            device = "cpu"
            logger.info("Using CPU (GPU not available or disabled)")
        return device
    
    def load_model(self):
        """Load the video generation model"""
        try:
            logger.info(f"Loading model: {self.model_id}")
            
            # Try to import diffusers
            try:
                from diffusers import DiffusionPipeline
                from diffusers.utils import export_to_video
                
                # Load the pipeline
                self.pipeline = DiffusionPipeline.from_pretrained(
                    self.model_id,
                    cache_dir=self.cache_dir,
                    torch_dtype=torch.float16 if self.device.startswith('cuda') else torch.float32,
                )
                
                # Move to device
                self.pipeline = self.pipeline.to(self.device)
                
                # Enable optimizations
                if self.device.startswith('cuda'):
                    try:
                        self.pipeline.enable_xformers_memory_efficient_attention()
                        logger.info("Enabled xformers memory efficient attention")
                    except Exception as e:
                        logger.warning(f"Could not enable xformers: {e}")
                    
                    # Low memory optimizations
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
                        
                        try:
                            # Enable gradient checkpointing if available
                            if hasattr(self.pipeline.unet, 'enable_gradient_checkpointing'):
                                self.pipeline.unet.enable_gradient_checkpointing()
                                logger.info("Enabled gradient checkpointing")
                        except Exception as e:
                            logger.warning(f"Could not enable gradient checkpointing: {e}")
                
                logger.info("Model loaded successfully")
                return True
                
            except ImportError:
                logger.error("diffusers library not installed. Install with: pip install diffusers transformers accelerate")
                return False
                
        except Exception as e:
            logger.error(f"Failed to load model: {e}")
            return False
    
    def generate_text_to_video(self, prompt, negative_prompt="", num_frames=64, 
                              fps=24, width=512, height=512, seed=-1,
                              guidance_scale=7.5, num_inference_steps=50):
        """Generate video from text prompt"""
        try:
            if self.pipeline is None:
                raise Exception("Model not loaded")
            
            # Set random seed if provided
            generator = None
            if seed > 0:
                generator = torch.Generator(device=self.device).manual_seed(seed)
            
            logger.info(f"Generating video: prompt='{prompt[:50]}...', frames={num_frames}, fps={fps}")
            
            # Generate video frames
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
            output_filename = f"text2video_{timestamp}_{job_id}.mp4"
            output_path = os.path.join(app.config['OUTPUT_FOLDER'], output_filename)
            
            # Export video
            from diffusers.utils import export_to_video
            export_to_video(output.frames[0], output_path, fps=fps)
            
            logger.info(f"Video generated: {output_path}")
            
            return {
                'output_path': output_path,
                'frames': num_frames,
                'fps': fps,
                'duration': num_frames / fps
            }
            
        except Exception as e:
            logger.error(f"Generation failed: {e}")
            raise
    
    def generate_image_to_video(self, image_path, prompt="", negative_prompt="",
                               num_frames=64, fps=24, width=512, height=512,
                               seed=-1, guidance_scale=7.5, num_inference_steps=50):
        """Generate video from image and prompt"""
        try:
            from PIL import Image
            
            if self.pipeline is None:
                raise Exception("Model not loaded")
            
            # Load image
            image = Image.open(image_path).convert('RGB')
            image = image.resize((width, height))
            
            # Set random seed if provided
            generator = None
            if seed > 0:
                generator = torch.Generator(device=self.device).manual_seed(seed)
            
            logger.info(f"Generating video from image: {image_path}")
            
            # Generate video frames
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
            output_filename = f"img2video_{timestamp}_{job_id}.mp4"
            output_path = os.path.join(app.config['OUTPUT_FOLDER'], output_filename)
            
            # Export video
            from diffusers.utils import export_to_video
            export_to_video(output.frames[0], output_path, fps=fps)
            
            logger.info(f"Video generated: {output_path}")
            
            return {
                'output_path': output_path,
                'frames': num_frames,
                'fps': fps,
                'duration': num_frames / fps
            }
            
        except Exception as e:
            logger.error(f"Generation failed: {e}")
            raise
    
    def generate_video_to_video(self, video_path, prompt="", negative_prompt="",
                               fps=24, strength=0.8, seed=-1,
                               guidance_scale=7.5, num_inference_steps=50):
        """Generate video from another video and prompt"""
        try:
            if self.pipeline is None:
                raise Exception("Model not loaded")
            
            logger.info(f"Generating video from video: {video_path}")
            
            # This is a placeholder - actual implementation depends on the model's capabilities
            # For now, return a mock response
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            job_id = str(uuid.uuid4())[:8]
            output_filename = f"vid2video_{timestamp}_{job_id}.mp4"
            output_path = os.path.join(app.config['OUTPUT_FOLDER'], output_filename)
            
            logger.warning("Video-to-video not fully implemented yet")
            
            return {
                'output_path': output_path,
                'frames': 64,
                'fps': fps,
                'duration': 64 / fps
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
        'model_loaded': model is not None and model.pipeline is not None,
        'device': str(model.device) if model else 'not initialized',
        **gpu_info
    })


@app.route('/api/generate/text-to-video', methods=['POST'])
def text_to_video():
    """Generate video from text prompt"""
    try:
        if model is None or model.pipeline is None:
            return jsonify({'error': 'Model not loaded'}), 503
        
        data = request.get_json()
        
        result = model.generate_text_to_video(
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
    """Generate video from image and prompt"""
    try:
        if model is None or model.pipeline is None:
            return jsonify({'error': 'Model not loaded'}), 503
        
        # Get uploaded image
        if 'image' not in request.files:
            return jsonify({'error': 'No image file provided'}), 400
        
        image_file = request.files['image']
        if image_file.filename == '':
            return jsonify({'error': 'No image file selected'}), 400
        
        # Save uploaded image
        filename = secure_filename(image_file.filename)
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        image_path = os.path.join(app.config['UPLOAD_FOLDER'], f"{timestamp}_{filename}")
        image_file.save(image_path)
        
        # Get parameters
        result = model.generate_image_to_video(
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


@app.route('/api/generate/video-to-video', methods=['POST'])
def video_to_video():
    """Generate video from another video and prompt"""
    try:
        if model is None or model.pipeline is None:
            return jsonify({'error': 'Model not loaded'}), 503
        
        # Get uploaded video
        if 'video' not in request.files:
            return jsonify({'error': 'No video file provided'}), 400
        
        video_file = request.files['video']
        if video_file.filename == '':
            return jsonify({'error': 'No video file selected'}), 400
        
        # Save uploaded video
        filename = secure_filename(video_file.filename)
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        video_path = os.path.join(app.config['UPLOAD_FOLDER'], f"{timestamp}_{filename}")
        video_file.save(video_path)
        
        # Get parameters
        result = model.generate_video_to_video(
            video_path=video_path,
            prompt=request.form.get('prompt', ''),
            negative_prompt=request.form.get('negative_prompt', ''),
            fps=int(request.form.get('fps', 24)),
            strength=float(request.form.get('strength', 0.8)),
            seed=int(request.form.get('seed', -1)),
            guidance_scale=float(request.form.get('guidance_scale', 7.5)),
            num_inference_steps=int(request.form.get('num_inference_steps', 50)),
        )
        
        return jsonify(result)
        
    except Exception as e:
        logger.error(f"Error in video-to-video: {e}")
        return jsonify({'error': str(e)}), 500


def main():
    """Main entry point"""
    global model
    
    logger.info("=" * 60)
    logger.info("Wan2.1 Video Generation Backend Starting")
    logger.info("=" * 60)
    
    # Initialize model
    model = VideoGenerationModel(
        model_id=MODEL_ID,
        cache_dir=CACHE_DIR,
        use_gpu=ENABLE_GPU,
        device_id=GPU_DEVICE_ID
    )
    
    # Load model
    if not model.load_model():
        logger.error("Failed to load model. Starting server anyway for health checks.")
    
    # Start Flask server
    host = os.getenv('SERVER_HOST', '0.0.0.0')
    port = int(os.getenv('PYTHON_BACKEND_PORT', '5000'))
    
    logger.info(f"Starting server on {host}:{port}")
    app.run(host=host, port=port, debug=False, threaded=True)


if __name__ == '__main__':
    main()
