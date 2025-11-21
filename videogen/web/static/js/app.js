// Wan2.1 Video Generator - Frontend JavaScript

// File Upload Handling
document.addEventListener('DOMContentLoaded', function() {
    // File upload preview
    const fileInputs = document.querySelectorAll('input[type="file"]');
    
    fileInputs.forEach(input => {
        input.addEventListener('change', function(e) {
            const file = e.target.files[0];
            if (file) {
                const label = this.closest('.file-upload-wrapper')?.querySelector('.file-upload-label');
                if (label) {
                    label.classList.add('has-file');
                    label.querySelector('.file-name')?.remove();
                    const fileName = document.createElement('div');
                    fileName.className = 'file-name text-primary mt-2';
                    fileName.innerHTML = `<i class="bi bi-file-earmark-check"></i> ${file.name}`;
                    label.appendChild(fileName);
                }
                
                // Preview for images
                if (file.type.startsWith('image/')) {
                    const preview = document.getElementById('imagePreview');
                    if (preview) {
                        const reader = new FileReader();
                        reader.onload = function(e) {
                            preview.src = e.target.result;
                            preview.classList.remove('d-none');
                        };
                        reader.readAsDataURL(file);
                    }
                }
                
                // Preview for videos
                if (file.type.startsWith('video/')) {
                    const preview = document.getElementById('videoPreview');
                    if (preview) {
                        const reader = new FileReader();
                        reader.onload = function(e) {
                            preview.src = e.target.result;
                            preview.classList.remove('d-none');
                        };
                        reader.readAsDataURL(file);
                    }
                }
            }
        });
    });
});

// Copy to clipboard
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(function() {
        showToast('Copied to clipboard!', 'success');
    }, function(err) {
        showToast('Failed to copy', 'danger');
    });
}

// Toast notifications
function showToast(message, type = 'info') {
    const toastContainer = document.getElementById('toastContainer');
    if (!toastContainer) {
        const container = document.createElement('div');
        container.id = 'toastContainer';
        container.className = 'toast-container position-fixed bottom-0 end-0 p-3';
        document.body.appendChild(container);
    }
    
    const toastHTML = `
        <div class="toast align-items-center text-white bg-${type} border-0" role="alert">
            <div class="d-flex">
                <div class="toast-body">${message}</div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        </div>
    `;
    
    const toastElement = document.createElement('div');
    toastElement.innerHTML = toastHTML;
    const toast = toastElement.firstElementChild;
    
    document.getElementById('toastContainer').appendChild(toast);
    
    const bsToast = new bootstrap.Toast(toast, { delay: 3000 });
    bsToast.show();
    
    toast.addEventListener('hidden.bs.toast', function() {
        toast.remove();
    });
}

// HTMX event listeners
document.body.addEventListener('htmx:beforeRequest', function(evt) {
    console.log('Request starting:', evt.detail.path);
});

document.body.addEventListener('htmx:afterRequest', function(evt) {
    console.log('Request completed:', evt.detail.path);
    
    // Show notification based on response
    if (evt.detail.successful) {
        const response = evt.detail.xhr.response;
        if (response.includes('alert-success')) {
            showToast('Video generated successfully!', 'success');
        }
    } else {
        showToast('Request failed. Please try again.', 'danger');
    }
});

document.body.addEventListener('htmx:responseError', function(evt) {
    console.error('Request error:', evt.detail);
    showToast('Connection error. Check if the API server is running.', 'danger');
});

// Form validation
function validateForm(formId) {
    const form = document.getElementById(formId);
    if (!form) return true;
    
    const inputs = form.querySelectorAll('[required]');
    let isValid = true;
    
    inputs.forEach(input => {
        if (!input.value.trim()) {
            input.classList.add('is-invalid');
            isValid = false;
        } else {
            input.classList.remove('is-invalid');
        }
    });
    
    return isValid;
}

// Dark/Light mode toggle (if implemented)
function toggleTheme() {
    const html = document.documentElement;
    const currentTheme = html.getAttribute('data-bs-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    html.setAttribute('data-bs-theme', newTheme);
    localStorage.setItem('theme', newTheme);
}

// Load saved theme
const savedTheme = localStorage.getItem('theme');
if (savedTheme) {
    document.documentElement.setAttribute('data-bs-theme', savedTheme);
}

// Utility: Format file size
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

// Utility: Format duration
function formatDuration(seconds) {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
}

// Example prompt insertion
function useExample(prompt) {
    const promptInput = document.getElementById('prompt');
    if (promptInput) {
        promptInput.value = prompt;
        promptInput.focus();
    }
}

// Keyboard shortcuts
document.addEventListener('keydown', function(e) {
    // Ctrl/Cmd + Enter to submit form
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        const activeForm = document.activeElement.closest('form');
        if (activeForm) {
            e.preventDefault();
            activeForm.requestSubmit();
        }
    }
});

// Auto-save form data to localStorage
function autoSaveForm(formId) {
    const form = document.getElementById(formId);
    if (!form) return;
    
    const inputs = form.querySelectorAll('input, textarea, select');
    
    // Load saved data
    inputs.forEach(input => {
        const savedValue = localStorage.getItem(`${formId}_${input.name}`);
        if (savedValue && !input.value) {
            input.value = savedValue;
        }
    });
    
    // Save on change
    inputs.forEach(input => {
        input.addEventListener('change', function() {
            localStorage.setItem(`${formId}_${this.name}`, this.value);
        });
    });
}

// Initialize auto-save for forms
document.addEventListener('DOMContentLoaded', function() {
    autoSaveForm('text-to-video-form');
    autoSaveForm('image-to-video-form');
    autoSaveForm('video-to-video-form');
});

console.log('🎬 Wan2.1 Video Generator - Web UI Loaded');
