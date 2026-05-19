package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dkr290/peridot-app/grpc-docker-registry/internal/storage"
	pb "github.com/dkr290/peridot-app/grpc-docker-registry/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ==================== Main Client ====================

type RegistryClient struct {
	grpcClient pb.ImageRegistryServiceClient
	client     *http.Client // Add http.Client wrapper
}

func NewRegistryClient(grpcAddr string) *RegistryClient {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Sprintf("failed to connect to gRPC server: %v", err))
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &RegistryClient{
		grpcClient: pb.NewImageRegistryServiceClient(conn),
		client:     &http.Client{Transport: transport}, // Wrap transport in client
	}
}

// ==================== Download from Docker Hub ====================

// DownloadImage downloads a Docker image from Docker Hub
func (c *RegistryClient) DownloadImage(imageRef string) error {
	fmt.Printf("📥 Downloading image: %s\n", imageRef)

	// Parse image reference
	repo, tag := parseImageRef(imageRef)
	if tag == "" {
		tag = "latest"
	}

	// Step 1: Get the manifest
	manifestURL := fmt.Sprintf("https://registry-1.docker.io/v2/%s/manifests/%s", repo, tag)
	manifest, manifestDigest, err := c.getManifest(manifestURL)
	if err != nil {
		return fmt.Errorf("failed to get manifest: %v", err)
	}

	fmt.Printf("   Manifest downloaded: %s\n", manifestDigest)

	// Step 2: Parse manifest to get layers and config
	type ManifestLayer struct {
		Digest string `json:"digest"`
		Size   int64  `json:"size"`
	}

	type ManifestConfig struct {
		Digest string `json:"digest"`
		Size   int64  `json:"size"`
	}

	type Manifest struct {
		SchemaVersion int             `json:"schemaVersion"`
		MediaType     string          `json:"mediaType"`
		Config        ManifestConfig  `json:"config"`
		Layers        []ManifestLayer `json:"layers"`
	}

	var m Manifest
	if err := json.Unmarshal(manifest, &m); err != nil {
		return fmt.Errorf("failed to parse manifest: %v", err)
	}

	fmt.Printf("   Layers: %d, Config: %s\n", len(m.Layers), m.Config.Digest)

	// Step 3: Download and push config blob
	fmt.Println("   Downloading config...")
	config, err := c.downloadBlob(
		"https://registry-1.docker.io/v2/" + repo + "/blobs/" + m.Config.Digest,
	)
	if err != nil {
		return fmt.Errorf("failed to download config: %v", err)
	}

	// Push config to registry
	_, err = c.grpcClient.PushBlob(context.Background(), &pb.PushBlobRequest{
		Data:   config,
		Digest: m.Config.Digest,
	})
	if err != nil {
		return fmt.Errorf("failed to push config: %v", err)
	}
	fmt.Printf("   Config pushed: %s\n", m.Config.Digest)

	// Step 4: Download and push each layer
	var layerDigests []string
	for i, layer := range m.Layers {
		fmt.Printf("   Downloading layer %d/%d...\n", i+1, len(m.Layers))
		layerData, err := c.downloadBlob(
			"https://registry-1.docker.io/v2/" + repo + "/blobs/" + layer.Digest,
		)
		if err != nil {
			return fmt.Errorf("failed to download layer %s: %v", layer.Digest, err)
		}

		// Push layer to registry
		_, err = c.grpcClient.PushBlob(context.Background(), &pb.PushBlobRequest{
			Data:   layerData,
			Digest: layer.Digest,
		})
		if err != nil {
			return fmt.Errorf("failed to push layer %s: %v", layer.Digest, err)
		}
		layerDigests = append(layerDigests, layer.Digest)
		fmt.Printf("   Layer %d pushed: %s\n", i+1, layer.Digest)
	}

	// Step 5: Push the image (manifest + layers)
	fmt.Println("   Pushing image manifest...")
	_, err = c.grpcClient.PushImage(context.Background(), &pb.PushImageRequest{
		Repository:   repo,
		Tag:          tag,
		Manifest:     manifest,
		Config:       config,
		LayerDigests: layerDigests,
	})
	if err != nil {
		return fmt.Errorf("failed to push image: %v", err)
	}

	fmt.Printf("   ✅ Image pushed successfully: %s:%s\n", repo, tag)
	return nil
}

// ==================== Helper Methods ====================

func (c *RegistryClient) downloadBlob(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)
}

func (c *RegistryClient) getManifest(url string) ([]byte, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	// Add authorization header if needed
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP %d for manifest", resp.StatusCode)
	}

	// Read manifest
	manifest, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// Get digest from response header
	digest := resp.Header.Get("Docker-Content-Digest")
	if digest == "" {
		digest = storage.ComputeDigest(manifest)
	}

	return manifest, digest, nil
}

func parseImageRef(ref string) (string, string) {
	parts := strings.Split(ref, ":")
	tag := ""
	if len(parts) == 2 {
		tag = parts[1]
	}
	return parts[0], tag
}

// ==================== Image Operations ====================

func (c *RegistryClient) PushImage(
	repo, tag string,
	manifest, config []byte,
	layerDigests []string,
) error {
	resp, err := c.grpcClient.PushImage(context.Background(), &pb.PushImageRequest{
		Repository:   repo,
		Tag:          tag,
		Manifest:     manifest,
		Config:       config,
		LayerDigests: layerDigests,
	})
	if err != nil {
		return fmt.Errorf("push image failed: %v", err)
	}
	fmt.Printf("✅ Image pushed: %s (manifest: %s)\n", resp.GetReference(), resp.GetManifestDigest())
	return nil
}

func (c *RegistryClient) PullImage(repo, tag string) ([]byte, []byte, []string, error) {
	resp, err := c.grpcClient.PullImage(context.Background(), &pb.PullImageRequest{
		Repository: repo,
		Tag:        tag,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("pull image failed: %v", err)
	}
	return resp.GetManifest(), resp.GetConfig(), resp.GetLayerDigests(), nil
}

func (c *RegistryClient) ListImages(repo string) ([]*pb.ImageInfo, error) {
	resp, err := c.grpcClient.ListImages(context.Background(), &pb.ListImagesRequest{
		Repository: repo,
	})
	if err != nil {
		return nil, fmt.Errorf("list images failed: %v", err)
	}
	return resp.GetImages(), nil
}

func (c *RegistryClient) ListTags(repo string) ([]string, error) {
	resp, err := c.grpcClient.ListTags(context.Background(), &pb.ListTagsRequest{
		Repository: repo,
	})
	if err != nil {
		return nil, fmt.Errorf("list tags failed: %v", err)
	}
	return resp.GetTags(), nil
}

func (c *RegistryClient) DeleteImage(repo, tag string) error {
	resp, err := c.grpcClient.DeleteImage(context.Background(), &pb.DeleteImageRequest{
		Repository: repo,
		Tag:        tag,
	})
	if err != nil {
		return fmt.Errorf("delete image failed: %v", err)
	}
	fmt.Printf("✅ Image deleted: %s\n", resp.GetMessage())
	return nil
}

func (c *RegistryClient) DeleteBlob(digest string) error {
	resp, err := c.grpcClient.DeleteBlob(context.Background(), &pb.DeleteBlobRequest{
		Digest: digest,
	})
	if err != nil {
		return fmt.Errorf("delete blob failed: %v", err)
	}
	fmt.Printf("✅ Blob deleted: %s\n", resp.GetMessage())
	return nil
}

func (c *RegistryClient) ListBlobs(prefix string) ([]string, error) {
	resp, err := c.grpcClient.ListBlobs(context.Background(), &pb.ListBlobsRequest{
		DigestPrefix: prefix,
	})
	if err != nil {
		return nil, fmt.Errorf("list blobs failed: %v", err)
	}
	return resp.GetDigests(), nil
}

// ==================== Save Image Locally ====================

func (c *RegistryClient) SaveImageLocally(repo, tag, outputDir string) error {
	manifest, config, layerDigests, err := c.PullImage(repo, tag)
	if err != nil {
		return err
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Save manifest
	manifestPath := filepath.Join(outputDir, "manifest.json")
	if err := os.WriteFile(manifestPath, manifest, 0644); err != nil {
		return err
	}
	fmt.Printf("   Saved manifest to: %s\n", manifestPath)

	// Save config
	configPath := filepath.Join(outputDir, "config.json")
	if err := os.WriteFile(configPath, config, 0644); err != nil {
		return err
	}
	fmt.Printf("   Saved config to: %s\n", configPath)

	// Save layers
	layersDir := filepath.Join(outputDir, "layers")
	if err := os.MkdirAll(layersDir, 0755); err != nil {
		return err
	}

	for _, digest := range layerDigests {
		// Get layer data
		layerResp, err := c.grpcClient.PullBlob(context.Background(), &pb.PullBlobRequest{
			Digest: digest,
		})
		if err != nil {
			return fmt.Errorf("failed to pull layer %s: %v", digest, err)
		}

		// Save layer
		layerPath := filepath.Join(layersDir, digest)
		if err := os.WriteFile(layerPath, layerResp.GetData(), 0644); err != nil {
			return err
		}
		fmt.Printf("   Saved layer: %s\n", layerPath)
	}

	fmt.Printf("   ✅ Image saved locally to: %s\n", outputDir)
	return nil
}
