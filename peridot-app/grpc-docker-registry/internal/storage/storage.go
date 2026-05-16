package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Storage interface {
	SaveImage(imageName string, imageData []byte) error
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

func (fs *FileStorage) SaveImage(imageName string, imageData []byte) error {
	filePath := filepath.Join(fs.BasePath, imageName)
	return os.WriteFile(filePath, imageData, 0o644)
}

func (fs *FileStorage) LoadImage(imageName string) ([]byte, error) {
	filePath := filepath.Join(fs.BasePath, imageName)
	return os.ReadFile(filePath)
}

func (fs *FileStorage) ListImages() ([]ImageInfo, error) {
	files, err := os.ReadDir(fs.BasePath)
	if err != nil {
		return nil, err
	}

	var images []ImageInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		imageData, err := os.ReadFile(filepath.Join(fs.BasePath, file.Name()))
		if err != nil {
			return nil, err
		}
		images = append(images, ImageInfo{
			ImageName: file.Name(),
			Tags:      []string{string(imageData)}, // Simplified for now
		})
	}
	return images, nil
}

func (fs *FileStorage) DeleteImage(imageName string) error {
	filePath := filepath.Join(fs.BasePath, imageName)
	return os.Remove(filePath)
}

func (fs *FileStorage) DeleteImageTag(imageName, tag string) error {
	filePath := filepath.Join(fs.BasePath, imageName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var tags []string
	err = json.Unmarshal([]byte(data), &tags)
	if err != nil {
		return err
	}
	newTags := []string{}
	for _, t := range tags {
		if t != tag {
			newTags = append(newTags, t)
		}
	}
	newData, err := json.Marshal(newTags)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, newData, 0o644)
}
