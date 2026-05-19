package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// ==================== Interfaces ====================

type BlobStore interface {
	StoreBlob(digest string, data []byte) error
	GetBlob(digest string) ([]byte, error)
	Exists(digest string) bool
	DeleteBlob(digest string) error
	ListBlobs(prefix string) ([]string, error)
}

type ManifestStore interface {
	StoreManifest(repository, tag, digest string, manifest []byte) error
	GetManifest(repository, tag string) ([]byte, string, error)
	DeleteManifest(repository, tag string) error
	ListTags(repository string) ([]string, error)
	DeleteTag(repository, tag string) error
}

// ==================== File Blob Store ====================

type FileBlobStore struct {
	rootPath string
	mu       sync.RWMutex
}

func NewFileBlobStore(rootPath string) *FileBlobStore {
	return &FileBlobStore{rootPath: rootPath}
}

func (s *FileBlobStore) StoreBlob(digest string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.rootPath, "blobs", "sha256", digest)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *FileBlobStore) GetBlob(digest string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.rootPath, "blobs", "sha256", digest)
	return os.ReadFile(path)
}

func (s *FileBlobStore) Exists(digest string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.rootPath, "blobs", "sha256", digest)
	_, err := os.Stat(path)
	return err == nil
}

func (s *FileBlobStore) DeleteBlob(digest string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.rootPath, "blobs", "sha256", digest)
	return os.Remove(path)
}

func (s *FileBlobStore) ListBlobs(prefix string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dir := filepath.Join(s.rootPath, "blobs", "sha256")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var digests []string
	for _, e := range entries {
		name := e.Name()
		if len(prefix) == 0 || len(name) >= len(prefix) && name[:len(prefix)] == prefix {
			digests = append(digests, name)
		}
	}
	sort.Strings(digests)
	return digests, nil
}

// ==================== File Manifest Store ====================

type FileManifestStore struct {
	rootPath string
	mu       sync.RWMutex
}

func NewFileManifestStore(rootPath string) *FileManifestStore {
	return &FileManifestStore{rootPath: rootPath}
}

func (s *FileManifestStore) StoreManifest(repository, tag, digest string, manifest []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create repository directory
	repoDir := filepath.Join(s.rootPath, "manifests", repository)
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return err
	}

	// Store manifest file named by digest
	manifestPath := filepath.Join(repoDir, digest)
	if err := os.WriteFile(manifestPath, manifest, 0644); err != nil {
		return err
	}

	// Create/update tag file
	tagPath := filepath.Join(repoDir, "tags", tag)
	if err := os.MkdirAll(filepath.Dir(tagPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(tagPath, []byte(digest), 0644)
}

func (s *FileManifestStore) GetManifest(repository, tag string) ([]byte, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Read tag to get digest

	tagDir := filepath.Join(s.rootPath, "manifests", repository, "tags")
	if err := os.MkdirAll(tagDir, 0755); err != nil {
		return nil, "", err
	}

	tagPath := filepath.Join(tagDir, tag)
	digestBytes, err := os.ReadFile(tagPath)
	if err != nil {
		return nil, "", err
	}

	digest := string(digestBytes)
	manifestPath := filepath.Join(s.rootPath, "manifests", repository, digest)
	manifest, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, "", err
	}

	return manifest, digest, nil
}

func (s *FileManifestStore) DeleteManifest(repository, tag string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Read tag to get digest
	tagPath := filepath.Join(s.rootPath, "manifests", repository, "tags", tag)
	digestBytes, err := os.ReadFile(tagPath)
	if err != nil {
		return err
	}

	digest := string(digestBytes)
	// Delete manifest file
	manifestPath := filepath.Join(s.rootPath, "manifests", repository, digest)
	if err := os.Remove(manifestPath); err != nil {
		return err
	}

	// Delete tag file
	return os.Remove(tagPath)
}

func (s *FileManifestStore) ListTags(repository string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tagsDir := filepath.Join(s.rootPath, "manifests", repository, "tags")
// Ensure directory exists before reading
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(tagsDir)
	if err != nil {
		return nil, err
	}

	var tags []string
	for _, e := range entries {
		if !e.IsDir() {
			tags = append(tags, e.Name())
		}
	}
	sort.Strings(tags)
	return tags, nil
}

func (s *FileManifestStore) DeleteTag(repository, tag string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tagPath := filepath.Join(s.rootPath, "manifests", repository, "tags", tag)
	return os.Remove(tagPath)
}

// ==================== Helpers ====================

func ComputeDigest(data []byte) string {
	hash := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(hash[:])
}
