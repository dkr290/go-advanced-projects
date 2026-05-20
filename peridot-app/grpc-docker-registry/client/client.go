package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dkr290/peridot-app/grpc-docker-registry/internal/upstream"
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
	reg, repo, reference, err := upstream.NewFromRef(imageRef, c.client)
	if err != nil {
		return err
	}

	// Step 1: Get the manifest
	manifest, manifestDigest, mediaType, err := reg.GetManifest(repo, reference)
	if err != nil {
		return fmt.Errorf("failed to get manifest: %v", err)
	}
	// If manifest list, resolve to linux/amd64 platform manifest
	if mediaType == "application/vnd.docker.distribution.manifest.list.v2+json" ||
		mediaType == "application/vnd.oci.image.index.v1+json" {
		digest, err := upstream.ResolvePlatformDigest(manifest, "linux", "amd64")
		if err != nil {
			return fmt.Errorf("failed to resolve platform manifest: %v", err)
		}
		fmt.Printf("   Resolved linux/amd64 manifest: %s\n", digest)
		manifest, manifestDigest, _, err = reg.GetManifest(repo, digest)
		if err != nil {
			return fmt.Errorf("failed to get platform manifest: %v", err)
		}
	}

	fmt.Printf("   Manifest downloaded: %s\n", manifestDigest)

	// Step 2: Parse manifest to get layers and config
	type ManifestLayer struct {
		Digest string `json:"digest"`
		Size   int64  `json:"size"`
	}

	type ManifestConfig struct {
		Digest string `json:"digest"`
	}

	type Manifest struct {
		Config ManifestConfig  `json:"config"`
		Layers []ManifestLayer `json:"layers"`
	}

	var m Manifest
	if err := json.Unmarshal(manifest, &m); err != nil {
		return fmt.Errorf("failed to parse manifest: %v", err)
	}

	fmt.Printf("   Layers: %d, Config: %s\n", len(m.Layers), m.Config.Digest)
	fmt.Println("   Downloading config...")
	config, err := reg.GetBlob(repo, m.Config.Digest)
	if err != nil {
		return fmt.Errorf("failed to download config: %v", err)
	}
	if _, err = c.grpcClient.PushBlob(context.Background(), &pb.PushBlobRequest{
		Data:   config,
		Digest: m.Config.Digest,
	}); err != nil {
		return fmt.Errorf("failed to push config: %v", err)
	}
	fmt.Printf("   Config pushed: %s\n", m.Config.Digest)

	// Step 4: Download and push each layer
	var layerDigests []string
	for i, layer := range m.Layers {
		fmt.Printf("   Downloading layer %d/%d...\n", i+1, len(m.Layers))
		layerData, err := reg.GetBlob(repo, layer.Digest)
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
		Tag:          reference,
		Config:       config,
		LayerDigests: layerDigests,
	})
	if err != nil {
		return fmt.Errorf("failed to push image: %v", err)
	}

	fmt.Printf("   ✅ Image pushed successfully: %s:%s\n", repo, reference)
	return nil
}

// ==================== Helper Methods ====================

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
