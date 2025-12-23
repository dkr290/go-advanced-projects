// Package webserver provides a simple web UI for viewing and downloading generated images
package webserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ImageInfo represents metadata about a generated image
type ImageInfo struct {
	Filename    string    `json:"filename"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"mod_time"`
	Thumbnail   string    `json:"thumbnail"`
	DownloadURL string    `json:"download_url"`
}

// Server represents the web server
type Server struct {
	OutputDir string
	Port      int
}

// NewServer creates a new web server instance
func NewServer(outputDir string, port int) *Server {
	return &Server{
		OutputDir: outputDir,
		Port:      port,
	}
}

// Start starts the web server
func (s *Server) Start() error {
	// Serve static files (images)
	fs := http.FileServer(http.Dir(s.OutputDir))
	http.Handle("/images/", http.StripPrefix("/images/", fs))

	// API endpoints
	http.HandleFunc("/api/images", s.handleListImages)
	http.HandleFunc("/api/download/", s.handleDownload)
	http.HandleFunc("/api/delete/", s.handleDelete)

	// Main gallery page
	http.HandleFunc("/", s.handleGallery)

	addr := fmt.Sprintf(":%d", s.Port)
	fmt.Printf("\n🌐 Web UI started at http://localhost%s\n", addr)
	fmt.Printf("📁 Serving images from: %s\n", s.OutputDir)
	fmt.Println("Press Ctrl+C to stop the server")

	return http.ListenAndServe(addr, nil)
}

// handleGallery serves the main gallery HTML page
func (s *Server) handleGallery(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("gallery").Parse(galleryHTML))
	data := map[string]interface{}{
		"Title":     "FLUX Image Gallery",
		"OutputDir": s.OutputDir,
	}
	tmpl.Execute(w, data)
}

// handleListImages returns JSON list of all images
func (s *Server) handleListImages(w http.ResponseWriter, r *http.Request) {
	images, err := s.getImages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

// handleDownload handles image downloads
func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/api/download/")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	filepath := filepath.Join(s.OutputDir, filename)
	
	// Security check - prevent directory traversal
	if !strings.HasPrefix(filepath, s.OutputDir) {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, filepath)
}

// handleDelete handles image deletion
func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/api/delete/")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	filepath := filepath.Join(s.OutputDir, filename)
	
	// Security check
	if !strings.HasPrefix(filepath, s.OutputDir) {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	if err := os.Remove(filepath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// getImages scans the output directory for images
func (s *Server) getImages() ([]ImageInfo, error) {
	var images []ImageInfo

	err := filepath.WalkDir(s.OutputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Only process image files
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(s.OutputDir, path)
		if err != nil {
			return err
		}

		images = append(images, ImageInfo{
			Filename:    filepath.Base(path),
			Path:        "/images/" + relPath,
			Size:        info.Size(),
			ModTime:     info.ModTime(),
			Thumbnail:   "/images/" + relPath,
			DownloadURL: "/api/download/" + relPath,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by modification time (newest first)
	sort.Slice(images, func(i, j int) bool {
		return images[i].ModTime.After(images[j].ModTime)
	})

	return images, nil
}

// galleryHTML is the embedded HTML template
const galleryHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
        }

        header {
            background: white;
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            margin-bottom: 30px;
        }

        h1 {
            color: #667eea;
            font-size: 2.5em;
            margin-bottom: 10px;
        }

        .stats {
            display: flex;
            gap: 20px;
            margin-top: 15px;
            flex-wrap: wrap;
        }

        .stat-item {
            background: #f7f7f7;
            padding: 10px 20px;
            border-radius: 8px;
            font-size: 0.9em;
            color: #666;
        }

        .stat-item strong {
            color: #667eea;
        }

        .controls {
            background: white;
            padding: 20px;
            border-radius: 15px;
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
            margin-bottom: 20px;
            display: flex;
            gap: 15px;
            flex-wrap: wrap;
            align-items: center;
        }

        .search-box {
            flex: 1;
            min-width: 250px;
            padding: 12px 20px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 1em;
            transition: border-color 0.3s;
        }

        .search-box:focus {
            outline: none;
            border-color: #667eea;
        }

        .btn {
            padding: 12px 24px;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-size: 1em;
            transition: all 0.3s;
            font-weight: 500;
        }

        .btn-primary {
            background: #667eea;
            color: white;
        }

        .btn-primary:hover {
            background: #5568d3;
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }

        .gallery {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
            gap: 25px;
            margin-top: 20px;
        }

        .image-card {
            background: white;
            border-radius: 15px;
            overflow: hidden;
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
            transition: transform 0.3s, box-shadow 0.3s;
            cursor: pointer;
        }

        .image-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 15px 30px rgba(0,0,0,0.2);
        }

        .image-wrapper {
            position: relative;
            width: 100%;
            padding-bottom: 100%; /* 1:1 Aspect Ratio */
            overflow: hidden;
            background: #f0f0f0;
        }

        .image-wrapper img {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            object-fit: cover;
        }

        .image-info {
            padding: 15px;
        }

        .image-name {
            font-weight: 600;
            color: #333;
            margin-bottom: 8px;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }

        .image-meta {
            font-size: 0.85em;
            color: #666;
            margin-bottom: 12px;
        }

        .image-actions {
            display: flex;
            gap: 10px;
        }

        .btn-small {
            flex: 1;
            padding: 8px 16px;
            font-size: 0.9em;
        }

        .btn-danger {
            background: #e74c3c;
            color: white;
        }

        .btn-danger:hover {
            background: #c0392b;
        }

        .loading {
            text-align: center;
            padding: 60px;
            color: white;
            font-size: 1.2em;
        }

        .empty-state {
            background: white;
            padding: 60px;
            border-radius: 15px;
            text-align: center;
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }

        .empty-state h2 {
            color: #667eea;
            margin-bottom: 10px;
        }

        .empty-state p {
            color: #666;
        }

        /* Modal styles */
        .modal {
            display: none;
            position: fixed;
            z-index: 1000;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0,0,0,0.9);
            align-items: center;
            justify-content: center;
        }

        .modal.active {
            display: flex;
        }

        .modal-content {
            max-width: 90%;
            max-height: 90%;
            object-fit: contain;
        }

        .close-modal {
            position: absolute;
            top: 30px;
            right: 40px;
            color: white;
            font-size: 40px;
            font-weight: bold;
            cursor: pointer;
            z-index: 1001;
        }

        .close-modal:hover {
            color: #ccc;
        }

        @media (max-width: 768px) {
            .gallery {
                grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
                gap: 15px;
            }

            h1 {
                font-size: 1.8em;
            }

            .controls {
                flex-direction: column;
            }

            .search-box {
                width: 100%;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🎨 {{.Title}}</h1>
            <div class="stats">
                <div class="stat-item">
                    <strong id="total-images">0</strong> images
                </div>
                <div class="stat-item">
                    Output: <strong>{{.OutputDir}}</strong>
                </div>
            </div>
        </header>

        <div class="controls">
            <input type="text" id="search" class="search-box" placeholder="🔍 Search images...">
            <button class="btn btn-primary" onclick="refreshGallery()">🔄 Refresh</button>
            <button class="btn btn-primary" onclick="downloadAll()">⬇️ Download All</button>
        </div>

        <div id="gallery" class="gallery">
            <div class="loading">Loading images...</div>
        </div>
    </div>

    <!-- Modal for full-size image view -->
    <div id="imageModal" class="modal" onclick="closeModal()">
        <span class="close-modal">&times;</span>
        <img class="modal-content" id="modalImage">
    </div>

    <script>
        let allImages = [];

        async function loadImages() {
            try {
                const response = await fetch('/api/images');
                allImages = await response.json();
                displayImages(allImages);
                document.getElementById('total-images').textContent = allImages.length;
            } catch (error) {
                console.error('Error loading images:', error);
                document.getElementById('gallery').innerHTML = 
                    '<div class="empty-state"><h2>Error</h2><p>Failed to load images</p></div>';
            }
        }

        function displayImages(images) {
            const gallery = document.getElementById('gallery');
            
            if (images.length === 0) {
                gallery.innerHTML = 
                    '<div class="empty-state"><h2>No images yet</h2><p>Generated images will appear here</p></div>';
                return;
            }

            gallery.innerHTML = images.map(img => {
                const size = (img.size / 1024 / 1024).toFixed(2);
                const date = new Date(img.mod_time).toLocaleString();
                
                return ` + "`" + `
                    <div class="image-card">
                        <div class="image-wrapper" onclick="openModal('${img.path}')">
                            <img src="${img.path}" alt="${img.filename}" loading="lazy">
                        </div>
                        <div class="image-info">
                            <div class="image-name" title="${img.filename}">${img.filename}</div>
                            <div class="image-meta">
                                ${size} MB • ${date}
                            </div>
                            <div class="image-actions">
                                <a href="${img.download_url}" class="btn btn-primary btn-small" download>
                                    ⬇️ Download
                                </a>
                                <button class="btn btn-danger btn-small" onclick="deleteImage('${img.filename}', event)">
                                    🗑️ Delete
                                </button>
                            </div>
                        </div>
                    </div>
                ` + "`" + `;
            }).join('');
        }

        function openModal(imagePath) {
            const modal = document.getElementById('imageModal');
            const modalImg = document.getElementById('modalImage');
            modal.classList.add('active');
            modalImg.src = imagePath;
        }

        function closeModal() {
            document.getElementById('imageModal').classList.remove('active');
        }

        async function deleteImage(filename, event) {
            event.stopPropagation();
            
            if (!confirm(` + "`" + `Delete ${filename}?` + "`" + `)) {
                return;
            }

            try {
                const response = await fetch(` + "`" + `/api/delete/${filename}` + "`" + `, {
                    method: 'DELETE'
                });

                if (response.ok) {
                    refreshGallery();
                } else {
                    alert('Failed to delete image');
                }
            } catch (error) {
                console.error('Error deleting image:', error);
                alert('Error deleting image');
            }
        }

        function refreshGallery() {
            loadImages();
        }

        function downloadAll() {
            allImages.forEach(img => {
                const a = document.createElement('a');
                a.href = img.download_url;
                a.download = img.filename;
                a.click();
            });
        }

        // Search functionality
        document.getElementById('search').addEventListener('input', (e) => {
            const searchTerm = e.target.value.toLowerCase();
            const filtered = allImages.filter(img => 
                img.filename.toLowerCase().includes(searchTerm)
            );
            displayImages(filtered);
        });

        // Close modal with Escape key
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                closeModal();
            }
        });

        // Load images on page load
        loadImages();

        // Auto-refresh every 5 seconds
        setInterval(loadImages, 5000);
    </script>
</body>
</html>
`
