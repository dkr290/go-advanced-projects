package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	// Step 0: Get auth token
	token, err := c.getDockerToken(repo)
	if err != nil {
		return fmt.Errorf("failed to get auth token: %v", err)
	}

	// Step 1: Get the manifest
	manifestURL := fmt.Sprintf("https://registry-1.docker.io/v2/%s/manifests/%s", repo, tag)
	manifest, manifestDigest, mediaType, err := c.getManifest(manifestURL, token)
	if err != nil {
		return fmt.Errorf("failed to get manifest: %v", err)
	}
	// If manifest list, resolve to linux/amd64 platform manifest
	if mediaType == "application/vnd.docker.distribution.manifest.list.v2+json" ||
		mediaType == "application/vnd.oci.image.index.v1+json" {
		digest, err := resolvePlatformDigest(manifest, "linux", "amd64")
		if err != nil {
			return fmt.Errorf("failed to resolve platform manifest: %v", err)
		}
		fmt.Printf("   Resolved linux/amd64 manifest: %s\n", digest)
		platformURL := fmt.Sprintf("https://registry-1.docker.io/v2/%s/manifests/%s", repo, digest)
		manifest, manifestDigest, _, err = c.getManifest(platformURL, token)
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
	config, err := c.downloadBlobWithToken(
		"https://registry-1.docker.io/v2/"+repo+"/blobs/"+m.Config.Digest, token,
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
		layerData, err := c.downloadBlobWithToken(
			"https://registry-1.docker.io/v2/"+repo+"/blobs/"+layer.Digest, token,
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
func (c *RegistryClient) downloadBlobWithToken(url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}

func (c *RegistryClient) getManifest(url string, token string) ([]byte, string, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", "", err
	}

	// Add authorization header if needed
	req.Header.Set("Accept", strings.Join([]string{
		"application/vnd.docker.distribution.manifest.v2+json",
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.oci.image.index.v1+json",
		"application/vnd.oci.image.manifest.v1+json",
	}, ","))

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", "", fmt.Errorf("HTTP %d for manifest", resp.StatusCode)
	}

	// Read manifest
	manifest, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", "", err
	}

	// Get digest from response header
	digest := resp.Header.Get("Docker-Content-Digest")
	if digest == "" {
		digest = storage.ComputeDigest(manifest)
	}
	mediaType := resp.Header.Get("Content-Type")

	return manifest, digest, mediaType, nil
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

func (c *RegistryClient) getDockerToken(repo string) (string, error) {
	url := fmt.Sprintf(
		"https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull",
		repo,
	)
	resp, err := c.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Token == "" {
		return "", fmt.Errorf("empty token received")
	}
	return result.Token, nil
}

func resolvePlatformDigest(manifestList []byte, os, arch string) (string, error) {
	var ml struct {
		Manifests []struct {
			Digest   string `json:"digest"`
			Platform struct {
				OS           string `json:"os"`
				Architecture string `json:"architecture"`
			} `json:"platform"`
		} `json:"manifests"`
	}
	if err := json.Unmarshal(manifestList, &ml); err != nil {
		return "", err
	}
	for _, m := range ml.Manifests {
		if m.Platform.OS == os && m.Platform.Architecture == arch {
			return m.Digest, nil
		}
	}
	return "", fmt.Errorf("no manifest found for %s/%s", os, arch)
}
