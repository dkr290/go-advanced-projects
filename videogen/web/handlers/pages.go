package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IndexHandler renders the home page
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/index.html", gin.H{
		"Title": "Wan2.1 Video Generator",
		"Page":  "home",
	})
}

// TextToVideoHandler renders the text-to-video page
func TextToVideoHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/text-to-video.html", gin.H{
		"Title": "Text to Video - Wan2.1",
		"Page":  "text-to-video",
	})
}

// ImageToVideoHandler renders the image-to-video page
func ImageToVideoHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/image-to-video.html", gin.H{
		"Title": "Image to Video - Wan2.1",
		"Page":  "image-to-video",
	})
}

// VideoToVideoHandler renders the video-to-video page
func VideoToVideoHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/video-to-video.html", gin.H{
		"Title": "Video to Video - Wan2.1",
		"Page":  "video-to-video",
	})
}

// GalleryHandler renders the gallery page
func GalleryHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/gallery.html", gin.H{
		"Title": "Video Gallery - Wan2.1",
		"Page":  "gallery",
	})
}

// SettingsHandler renders the settings page
func SettingsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/settings.html", gin.H{
		"Title": "Settings - Wan2.1",
		"Page":  "settings",
	})
}
