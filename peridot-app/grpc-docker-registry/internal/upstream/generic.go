package upstream

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dkr290/peridot-app/grpc-docker-registry/internal/storage"
)

type GenericRegistry struct {
	Host          string
	client        *http.Client
	normalizeRepo func(string) string
	tokenFetcher  func(repo string) (string, error)
}

func NewGenericRegistry(host string, client *http.Client) *GenericRegistry {
	r := &GenericRegistry{Host: host, client: client}
	r.normalizeRepo = func(repo string) string { return repo }
	r.tokenFetcher = r.defaultGetToken
	return r
}

func (r *GenericRegistry) NormalizeRepo(repo string) string {
	return r.normalizeRepo(repo)
}
func (r *GenericRegistry) defaultGetToken(repo string) (string, error) {
	url := fmt.Sprintf(
		"https://%s/v2/auth?service=%s&scope=repository:%s:pull",
		r.Host, r.Host, repo,
	)
	resp, err := r.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Token       string `json:"token"`
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Token != "" {
		return result.Token, nil
	}
	if result.AccessToken != "" {
		return result.AccessToken, nil
	}
	return "", fmt.Errorf("empty token received from %s", r.Host)
}
func (r *GenericRegistry) GetToken(repo string) (string, error) {
	return r.tokenFetcher(repo)
}



func (r *GenericRegistry) GetManifest(
	repo string,
	reference string,
) ([]byte, string, string, error) {
	token, err := r.tokenFetcher(repo)
	if err != nil {
		return nil, "", "", fmt.Errorf("auth failed: %v", err)
	}
	url := fmt.Sprintf("https://%s/v2/%s/manifests/%s", r.Host, repo, reference)

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

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := r.client.Do(req)
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

func (r *GenericRegistry) GetBlob(repo, digest string) ([]byte, error) {
	token, err := r.tokenFetcher(repo)
	if err != nil {
		return nil, fmt.Errorf("auth failed: %v", err)
	}

	url := fmt.Sprintf("https://%s/v2/%s/blobs/%s", r.Host, repo, digest)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for blob %s", resp.StatusCode, digest)
	}
	return io.ReadAll(resp.Body)
}
