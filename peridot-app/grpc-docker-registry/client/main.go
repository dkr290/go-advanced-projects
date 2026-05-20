package main

import "fmt"

// ==================== Main ====================

func main() {
	client := NewRegistryClient("localhost:50051")

	// Example usage:
	// 1. Download an image from Docker Hub
	 err := client.DownloadImage("ubuntu:latest")
	if err != nil {
		fmt.Printf("❌ Failed to download image: %v\n", err)
		return
	}


	// 2. List images
	 	images, err := client.ListImages("*")
	if err != nil {
		fmt.Printf("❌ Failed to list images: %v\n", err)
		return
	}
	fmt.Printf("Images found: %d\n", len(images))
	for _, img := range images {
		fmt.Printf("  Image: %s:%s\n", img.GetRepository(), img.GetTag())
	}

	// 3. List tags
	 tags, _ := client.ListTags("library/alpine")
	 fmt.Printf("Tags: %v\n", tags)


	// 5. Delete image
	// client.DeleteImage("alpine", "latest")

	// 6. List blobs
	 blobs, _ := client.ListBlobs("sha256:")
	 fmt.Printf("Blobs: %d\n", len(blobs))
}

