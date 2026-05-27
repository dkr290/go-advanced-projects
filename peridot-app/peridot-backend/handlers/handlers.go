package handlers

import (
	"context"
	"log"

	"github.com/dkr290/peridot-app/peridot-backend/patcher"
	"github.com/dkr290/peridot-app/peridot-backend/registry"
	"github.com/dkr290/peridot-app/peridot-backend/scanner"
	"github.com/humaio/huma/v4"
)

type ImageHandler struct {
	scanner  *scanner.Scanner
	patcher  *patcher.Patcher
	registry *registry.Client
}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		scanner:  scanner.NewScanner(),
		patcher:  patcher.NewPatcher(),
		registry: registry.NewClient("localhost:50051"),
	}
}

// AddImageHandler registers the image-related endpoints
func (h *ImageHandler) Register(api huma.API) {
	// POST /images - Add image from Docker Hub with scanning and patching
	api.Post("/images", h.AddImage)

	// GET /images - List all images
	api.Get("/images", h.ListImages)

	// GET /images/{repo}/{tag}/scan - Scan specific image
	api.Get("/images/{repo}/{tag}/scan", h.ScanImage)

	// POST /images/{repo}/{tag}/patch - Patch specific image
	api.Post("/images/{repo}/{tag}/patch", h.PatchImage)
}

func (h *ImageHandler) AddImage(
	ctx context.Context,
	input *AddImageInput,
) (*AddImageOutput, error) {
	log.Printf("Adding image: %s", input.Body.ImageRef)

	// Step 1: Pull image from Docker Hub
	log.Println("Pulling image from Docker Hub...")
	// TODO: Implement pull logic using existing client

	// Step 2: Scan for vulnerabilities
	log.Println("Scanning for vulnerabilities...")
	scanResults, err := h.scanner.ScanImage(ctx, input.Body.ImageRef)
	if err != nil {
		return nil, huma.Error500InternalServerError("Scan failed: " + err.Error())
	}

	// Step 3: Patch vulnerabilities if found
	log.Println("Patching vulnerabilities...")
	patchedImage, err := h.patcher.PatchImage(ctx, input.Body.ImageRef, scanResults)
	if err != nil {
		return nil, huma.Error500InternalServerError("Patching failed: " + err.Error())
	}

	// Step 4: Push patched image to gRPC registry
	log.Println("Pushing patched image to registry...")
	err = h.registry.PushImage(ctx, patchedImage)
	if err != nil {
		return nil, huma.Error500InternalServerError("Registry push failed: " + err.Error())
	}

	return &AddImageOutput{
		Message:         "Image added, scanned, patched, and stored successfully",
		Vulnerabilities: scanResults.Vulnerabilities,
	}, nil
}

func (h *ImageHandler) ListImages(
	ctx context.Context,
	input *ListImagesInput,
) (*ListImagesOutput, error) {
	images, err := h.registry.ListImages(ctx, input.Repository)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to list images: " + err.Error())
	}
	return &ListImagesOutput{Images: images}, nil
}

func (h *ImageHandler) ScanImage(
	ctx context.Context,
	input *ScanImageInput,
) (*ScanImageOutput, error) {
	results, err := h.scanner.ScanImage(ctx, input.Repo+":"+input.Tag)
	if err != nil {
		return nil, huma.Error500InternalServerError("Scan failed: " + err.Error())
	}
	return &ScanImageOutput{Results: results}, nil
}

func (h *ImageHandler) PatchImage(
	ctx context.Context,
	input *PatchImageInput,
) (*PatchImageOutput, error) {
	// Scan first
	scanResults, err := h.scanner.ScanImage(ctx, input.Repo+":"+input.Tag)
	if err != nil {
		return nil, huma.Error500InternalServerError("Scan failed: " + err.Error())
	}

	// Patch
	patchedImage, err := h.patcher.PatchImage(ctx, input.Repo+":"+input.Tag, scanResults)
	if err != nil {
		return nil, huma.Error500InternalServerError("Patching failed: " + err.Error())
	}

	// Push to registry
	err = h.registry.PushImage(ctx, patchedImage)
	if err != nil {
		return nil, huma.Error500InternalServerError("Registry push failed: " + err.Error())
	}

	return &PatchImageOutput{
		Message: "Image patched and stored successfully",
	}, nil
}
