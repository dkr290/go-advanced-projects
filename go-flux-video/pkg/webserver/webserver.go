// Package webserver provides a simple web UI for viewing and downloading generated images
package webserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(s.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", s.OutputDir, err)
	}

	// Serve static files (images)
	fs := http.FileServer(http.Dir(s.OutputDir))
	http.Handle("/images/", http.StripPrefix("/images/", fs))

	// API endpoints
	http.HandleFunc("/api/images", s.handleListImages)
	http.HandleFunc("/api/download/", s.handleDownload)
	http.HandleFunc("/api/delete/", s.handleDelete)
	http.HandleFunc("/api/upload", s.handleUpload)
	http.HandleFunc("/api/upload-images", s.handleUploadToImagesDir)

	// Main gallery page
	http.HandleFunc("/", s.handleGallery)

	addr := fmt.Sprintf(":%d", s.Port)
	fmt.Printf("\n🌐 Web UI started at http://localhost%s\n", addr)
	fmt.Printf("📁 Serving images from: %s\n", s.OutputDir)
	fmt.Println("📤 Upload images to ./images/ directory for img2img processing")
	fmt.Println("Press Ctrl+C to stop the server")

	return http.ListenAndServe(addr, nil)
}

// loadTemplates loads the HTML template from the templates directory
func (s *Server) loadTemplates() (*template.Template, error) {
	// Parse the single template file
	tmpl, err := template.ParseFiles("templates/gallery.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

// handleGallery serves the main gallery HTML page
func (s *Server) handleGallery(w http.ResponseWriter, r *http.Request) {
	tmpl, err := s.loadTemplates()
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to load templates: %v", err),
			http.StatusInternalServerError,
		)
		return
	}

	data := map[string]interface{}{
		"Title":     "FLUX Image Gallery",
		"OutputDir": s.OutputDir,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to render template: %v", err),
			http.StatusInternalServerError,
		)
	}
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

// handleUpload handles image uploads to the output directory
func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form with 10MB max memory
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from form data
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No image file provided: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowedExts := []string{".png", ".jpg", ".jpeg", ".webp", ".gif", ".bmp"}
	validExt := false
	for _, allowed := range allowedExts {
		if ext == allowed {
			validExt = true
			break
		}
	}
	if !validExt {
		http.Error(
			w,
			"Invalid file type. Allowed: PNG, JPG, JPEG, WebP, GIF, BMP",
			http.StatusBadRequest,
		)
		return
	}

	// Create unique filename
	timestamp := time.Now().Format("20060102_150405")
	uniqueFilename := fmt.Sprintf("upload_%s_%s", timestamp, handler.Filename)
	destPath := filepath.Join(s.OutputDir, uniqueFilename)

	// Create destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// Copy the uploaded file to destination
	_, err = io.Copy(destFile, file)
	if err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"status":   "success",
		"filename": uniqueFilename,
		"path":     "/images/" + uniqueFilename,
		"size":     handler.Size,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleUploadToImagesDir handles image uploads specifically to the ./images/ directory
func (s *Server) handleUploadToImagesDir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form with 10MB max memory
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from form data
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No image file provided: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowedExts := []string{".png", ".jpg", ".jpeg", ".webp", ".gif", ".bmp"}
	validExt := false
	for _, allowed := range allowedExts {
		if ext == allowed {
			validExt = true
			break
		}
	}
	if !validExt {
		http.Error(
			w,
			"Invalid file type. Allowed: PNG, JPG, JPEG, WebP, GIF, BMP",
			http.StatusBadRequest,
		)
		return
	}

	// Create images directory if it doesn't exist
	imagesDir := filepath.Join(".", "images")
	if err := os.MkdirAll(imagesDir, 0o755); err != nil {
		http.Error(
			w,
			"Failed to create images directory: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	// Use original filename (or create unique if needed)
	destFilename := handler.Filename
	destPath := filepath.Join(imagesDir, destFilename)

	// Check if file already exists, add timestamp if it does
	if _, err := os.Stat(destPath); err == nil {
		timestamp := time.Now().Format("20060102_150405")
		destFilename = fmt.Sprintf("%s_%s", timestamp, handler.Filename)
		destPath = filepath.Join(imagesDir, destFilename)
	}

	// Create destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// Copy the uploaded file to destination
	_, err = io.Copy(destFile, file)
	if err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"status":        "success",
		"filename":      destFilename,
		"path":          destPath,
		"relative_path": "images/" + destFilename,
		"size":          handler.Size,
		"message":       "Image uploaded to ./images/ directory for img2img processing",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
