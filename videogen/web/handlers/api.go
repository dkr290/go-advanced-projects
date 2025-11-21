package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var apiBaseURL = getEnv("API_BASE_URL", "http://localhost:8080")

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GenerateTextToVideoHandler handles text-to-video generation
func GenerateTextToVideoHandler(c *gin.Context) {
	prompt := c.PostForm("prompt")
	negativePrompt := c.PostForm("negative_prompt")
	numFrames := c.PostForm("num_frames")
	fps := c.PostForm("fps")
	width := c.PostForm("width")
	height := c.PostForm("height")
	guidanceScale := c.PostForm("guidance_scale")
	numInferenceSteps := c.PostForm("num_inference_steps")
	seed := c.PostForm("seed")

	// Prepare request body
	requestBody := map[string]interface{}{
		"prompt":              prompt,
		"negative_prompt":     negativePrompt,
		"num_frames":          parseIntOrDefault(numFrames, 64),
		"fps":                 parseIntOrDefault(fps, 24),
		"width":               parseIntOrDefault(width, 512),
		"height":              parseIntOrDefault(height, 512),
		"guidance_scale":      parseFloatOrDefault(guidanceScale, 7.5),
		"num_inference_steps": parseIntOrDefault(numInferenceSteps, 50),
		"seed":                parseIntOrDefault(seed, -1),
	}

	jsonData, _ := json.Marshal(requestBody)

	// Call API
	resp, err := http.Post(
		apiBaseURL+"/api/v1/generate/text-to-video",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		c.HTML(http.StatusOK, "components/error.html", gin.H{
			"Error": "Failed to connect to API server: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != 202 {
		c.HTML(http.StatusOK, "components/error.html", gin.H{
			"Error": result["error"],
		})
		return
	}

	// Return job status component
	c.HTML(http.StatusOK, "components/job-status.html", gin.H{
		"JobID":   result["job_id"],
		"Status":  result["status"],
		"Message": result["message"],
	})
}

// GenerateImageToVideoHandler handles image-to-video generation
func GenerateImageToVideoHandler(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		c.HTML(http.StatusOK, "components/error.html", gin.H{
			"Error": "No image file uploaded",
		})
		return
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	fileContent, _ := file.Open()
	defer fileContent.Close()
	part, _ := writer.CreateFormFile("image", file.Filename)
	io.Copy(part, fileContent)

	// Add other fields
	writer.WriteField("prompt", c.PostForm("prompt"))
	writer.WriteField("negative_prompt", c.PostForm("negative_prompt"))
	writer.WriteField("num_frames", c.PostForm("num_frames"))
	writer.WriteField("fps", c.PostForm("fps"))
	writer.WriteField("width", c.PostForm("width"))
	writer.WriteField("height", c.PostForm("height"))
	writer.WriteField("guidance_scale", c.PostForm("guidance_scale"))
	writer.WriteField("num_inference_steps", c.PostForm("num_inference_steps"))
	writer.WriteField("seed", c.PostForm("seed"))

	writer.Close()

	// Call API
	req, _ := http.NewRequest("POST", apiBaseURL+"/api/v1/generate/image-to-video", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.HTML(http.StatusOK, "components/error.html", gin.H{
			"Error": "Failed to connect to API server: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != 202 {
		c.HTML(http.StatusOK, "components/error.html", gin.H{
			"Error": result["error"],
		})
		return
	}

	c.HTML(http.StatusOK, "components/job-status.html", gin.H{
		"JobID":   result["job_id"],
		"Status":  result["status"],
		"Message": result["message"],
	})
}

// GenerateVideoToVideoHandler handles video-to-video generation
func GenerateVideoToVideoHandler(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("video")
	if err != nil {
		c.HTML(http.StatusOK, "components/error.html", gin.H{
			"Error": "No video file uploaded",
		})
		return
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	fileContent, _ := file.Open()
	defer fileContent.Close()
	part, _ := writer.CreateFormFile("video", file.Filename)
	io.Copy(part, fileContent)

	// Add other fields
	writer.WriteField("prompt", c.PostForm("prompt"))
	writer.WriteField("negative_prompt", c.PostForm("negative_prompt"))
	writer.WriteField("fps", c.PostForm("fps"))
	writer.WriteField("strength", c.PostForm("strength"))
	writer.WriteField("guidance_scale", c.PostForm("guidance_scale"))
	writer.WriteField("num_inference_steps", c.PostForm("num_inference_steps"))
	writer.WriteField("seed", c.PostForm("seed"))

	writer.Close()

	// Call API
	req, _ := http.NewRequest("POST", apiBaseURL+"/api/v1/generate/video-to-video", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.HTML(http.StatusOK, "components/error.html", gin.H{
			"Error": "Failed to connect to API server: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != 202 {
		c.HTML(http.StatusOK, "components/error.html", gin.H{
			"Error": result["error"],
		})
		return
	}

	c.HTML(http.StatusOK, "components/job-status.html", gin.H{
		"JobID":   result["job_id"],
		"Status":  result["status"],
		"Message": result["message"],
	})
}

// GetJobStatusHandler checks job status
func GetJobStatusHandler(c *gin.Context) {
	jobID := c.Param("id")

	resp, err := http.Get(apiBaseURL + "/api/v1/job/" + jobID)
	if err != nil {
		c.HTML(http.StatusOK, "components/error.html", gin.H{
			"Error": "Failed to check job status",
		})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	status := result["status"].(string)

	if status == "completed" {
		// Return completed video component
		c.HTML(http.StatusOK, "components/video-result.html", gin.H{
			"JobID":      result["job_id"],
			"OutputPath": result["output_path"],
			"VideoURL":   apiBaseURL + "/" + result["output_path"].(string),
		})
	} else if status == "failed" {
		c.HTML(http.StatusOK, "components/error.html", gin.H{
			"Error": result["message"],
		})
	} else {
		// Still processing
		c.HTML(http.StatusOK, "components/job-status.html", gin.H{
			"JobID":   result["job_id"],
			"Status":  status,
			"Message": result["message"],
		})
	}
}

// GetModelsHandler gets available models
func GetModelsHandler(c *gin.Context) {
	resp, err := http.Get(apiBaseURL + "/api/v1/models")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	c.JSON(http.StatusOK, result)
}

// SwitchModelHandler switches the active model
func SwitchModelHandler(c *gin.Context) {
	modelName := c.PostForm("model_name")

	requestBody := map[string]string{
		"model_name": modelName,
	}
	jsonData, _ := json.Marshal(requestBody)

	resp, err := http.Post(
		apiBaseURL+"/api/v1/models/switch",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	c.JSON(resp.StatusCode, result)
}

// GetGalleryListHandler returns list of generated videos
func GetGalleryListHandler(c *gin.Context) {
	// This is a simple implementation
	// In production, you'd want to store this in a database
	outputsDir := "./outputs"
	
	var videos []map[string]interface{}
	
	filepath.Walk(outputsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && (filepath.Ext(path) == ".mp4" || filepath.Ext(path) == ".webm") {
			videos = append(videos, map[string]interface{}{
				"filename":  info.Name(),
				"path":      path,
				"size":      info.Size(),
				"created":   info.ModTime().Format(time.RFC3339),
			})
		}
		return nil
	})

	c.HTML(http.StatusOK, "components/gallery-grid.html", gin.H{
		"Videos": videos,
	})
}

// Helper functions
func parseIntOrDefault(s string, defaultValue int) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	if result == 0 {
		return defaultValue
	}
	return result
}

func parseFloatOrDefault(s string, defaultValue float64) float64 {
	var result float64
	fmt.Sscanf(s, "%f", &result)
	if result == 0 {
		return defaultValue
	}
	return result
}
