package storage

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Storage interface {
	SaveImage(imageName string, imageData []byte, tag string) error
	LoadImage(imageName string) ([]byte, error)
	ListImages() ([]ImageInfo, error)
	DeleteImage(imageName string) error
	DeleteImageTag(imageName, tag string) error
}

type ImageInfo struct {
	ImageName string
	Tags      []string
}

type FileStorage struct {
	BasePath string
}

func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{BasePath: basePath}
}

func (fs *FileStorage) SaveImage(imageName string, imageData []byte, tag string) error {
	manifestPath := filepath.Join(fs.BasePath, imageName+".manifest")
	if err := os.WriteFile(manifestPath, imageData, 0o644); err != nil {
		return err
	}
	tagPath := filepath.Join(fs.BasePath, imageName+".tags")
	tags, err := fs.readTags(tagPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if !contains(tags, tag) {
		tags = append(tags, tag)
	}
	return os.WriteFile(tagPath, []byte(strings.Join(tags, "\n")), 0o644)
}

func (fs *FileStorage) LoadImage(imageName string) ([]byte, error) {
	filePath := filepath.Join(fs.BasePath, imageName)
	return os.ReadFile(filePath)
}

func (fs *FileStorage) ListImages() ([]ImageInfo, error) {
	files, err := os.ReadDir(fs.BasePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []ImageInfo{}, nil
		}

		return nil, err
	}

	var images []ImageInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// Only list files that are image manifests (not tag files)
		if filepath.Ext(file.Name()) != ".manifest" {
			continue
		}
		images = append(images, ImageInfo{
			ImageName: strings.TrimSuffix(file.Name(), ".manifest"),
			Tags:      []string{}, // Tags tracked separately
		})
	}
	return images, nil
}

func (fs *FileStorage) DeleteImageTag(imageName, tag string) error {
	tagPath := filepath.Join(fs.BasePath, imageName+".tags")
	tags, err := fs.readTags(tagPath)
	if err != nil {
		return err
	}
	newTags := []string{}
	for _, t := range tags {
		if t != tag {
			newTags = append(newTags, t)
		}
	}
	if len(newTags) == 0 {
		return os.Remove(tagPath)
	}
	return os.WriteFile(tagPath, []byte(strings.Join(newTags, "\n")), 0o644)
}

func (fs *FileStorage) DeleteImage(imageName string) error {
	manifestPath := filepath.Join(fs.BasePath, imageName+".manifest")
	tagPath := filepath.Join(fs.BasePath, imageName+".tags")

	if err := os.Remove(manifestPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Remove(tagPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (fs *FileStorage) readTags(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var tags []string
	for _, line := range lines {
		if line != "" {
			tags = append(tags, line)
		}
	}
	return tags, nil
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
