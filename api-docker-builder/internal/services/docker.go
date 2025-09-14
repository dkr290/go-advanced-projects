// Package services is the main docker build process
package services

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dkr290/go-advanced-projects/api-docker-builder/internal/models"
	"github.com/docker/docker/api/types"
	"github.com/google/uuid"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func (d *DockerService) Close() error {
	return d.client.Close()
}

func (d *DockerService) BuildImage(
	ctx context.Context,
	req *models.BuildImageRequest,
) (*models.BuildImageResponse, error) {
	buildID := uuid.New().String()
	imageName := fmt.Sprintf("%s:%s", req.Name, req.Tag)

	// Initialize build status
	status := &models.BuildStatus{
		BuildID:   buildID,
		Status:    "pending",
		Message:   "Build initiated",
		ImageName: imageName,
		StartedAt: time.Now(),
		Logs:      []string{},
	}

	d.builds[buildID] = status

	// Start build in goroutine
	d.clog.Infof("Starting the build with buildID: %s for image: %s", buildID, imageName)
	d.performBuild(ctx, buildID, req, imageName)

	return &models.BuildImageResponse{
		BuildID:   buildID,
		Status:    "pending",
		Message:   "Build started",
		ImageName: imageName,
		StartedAt: status.StartedAt,
	}, nil
}

func (d *DockerService) GetBuildStatus(buildID string) (*models.BuildStatus, error) {
	status, exists := d.builds[buildID]
	if !exists {
		return nil, fmt.Errorf("build not found")
	}

	return status, nil
}

func (d *DockerService) performBuild(
	ctx context.Context,
	buildID string,
	req *models.BuildImageRequest,
	imageName string,
) {
	d.updateBuildStatus(buildID, "building", "Creating Dockerfile and building image", nil)
	targetDir := "./data"
	// clone repo to some temp dir
	err := d.cloneRepository(
		req.RepoURL,
		targetDir,
		strings.TrimSpace(req.RepoUsername),
		strings.TrimSpace(req.RepoPassword),
		req.UserAuth,
	)
	if err != nil {
		d.updateBuildStatus(
			buildID,
			"failed",
			fmt.Sprintf("failed to clone repository: %v", err),
			nil,
		)
	}
	defer func() {
		_ = os.RemoveAll(targetDir)
	}()
	// Find Dockerfile path (relative to repoDir)
	dockerfilePath, err := findDockerfilePathFromDir(targetDir)
	if err != nil {
		d.updateBuildStatus(buildID, "failed", fmt.Sprintf("no Dockerfile found: %v", err), nil)
		return
	}

	// Build image and the Dockerfile should be in the archive
	buildOptions := types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: dockerfilePath,
		Remove:     true,
	}
	// Create build context
	buildContext, err := d.createBuildContext(dockerfilePath)
	if err != nil {
		d.updateBuildStatus(
			buildID,
			"failed",
			fmt.Sprintf("Failed to create build context: %v", err),
			nil,
		)
		return
	}
	resp, err := d.client.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		d.updateBuildStatus(buildID, "failed", fmt.Sprintf("Failed to build image: %v", err), nil)
		return
	}
	defer resp.Body.Close()

	// Read build logs
	logs, err := d.readBuildLogs(buildID, resp.Body)
	if err != nil {
		d.updateBuildStatus(
			buildID,
			"failed",
			fmt.Sprintf("error reading build logs: %v", err),
			logs,
		)
		return
	}

	// Update final status
	completedAt := time.Now()
	if build, exists := d.builds[buildID]; exists {
		build.Status = "success"
		build.Message = "Image built successfully"
		build.CompletedAt = &completedAt
		build.Logs = logs
	}

	d.clog.Infof("Build %s completed successfully for image %s", buildID, imageName)
}

func (d *DockerService) updateBuildStatus(buildID, status, message string, logs []string) {
	if build, exists := d.builds[buildID]; exists {
		build.Status = status
		build.Message = message
		if logs != nil {
			build.Logs = logs
		}
		if status == "failed" || status == "success" {
			completedAt := time.Now()
			build.CompletedAt = &completedAt
		}
	}
}

func (d *DockerService) cloneRepository(
	repoURL, targetDir, user, pass string, useAuth bool,
) error {
	// Prepare clone options
	cloneOptions := &git.CloneOptions{
		URL: repoURL,
	}
	// Add authentication if flag is true
	if useAuth {
		cloneOptions.Auth = &http.BasicAuth{
			Username: user,
			Password: pass,
		}
	}
	d.clog.Infof("Cloning the repository into  %s", targetDir)
	_, err := git.PlainClone(targetDir, false, cloneOptions)
	if err != nil {
		d.clog.Errorf("Error cloning repository: %v\n", err)
		os.Exit(1)
	}
	d.clog.Infof("The repository %s cloned sucessfully", repoURL)

	return nil
}

func findDockerfilePathFromDir(targetDir string) (dockerFilePath string, err error) {
	// Check if directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return "", fmt.Errorf("directory %s does not exist", targetDir)
	}

	// Walk through the directory recursively
	err = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories (like .git)
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
			return filepath.SkipDir
		}

		// Check for Dockerfile (case-insensitive)
		if !info.IsDir() {
			filename := strings.ToLower(info.Name())
			if filename == "dockerfile" || strings.HasPrefix(filename, "dockerfile.") {
				dockerFilePath = path
			}
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error walking directory: %v", err)
	}

	return dockerFilePath, err
}

func (d *DockerService) createBuildContext(
	dockerfilePath string,
) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFileReader, err := os.Open(dockerfilePath)
	if err != nil {
		return nil, err
	}
	readDockerFile, err := io.ReadAll(dockerFileReader)
	if err != nil {
		return nil, err
	}
	// Add Dockerfile
	header := &tar.Header{
		Name: filepath.Base(dockerfilePath),
		Mode: 0o644,
		Size: int64(len(readDockerFile)),
	}
	if err := tw.WriteHeader(header); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(filepath.Base(dockerfilePath))); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

func (d *DockerService) readBuildLogs(buildID string, r io.Reader) ([]string, error) {
	logs := []string{}
	buf := make([]byte, 4096)

	for {
		n, err := r.Read(buf)
		if n > 0 {
			logLine := string(buf[:n])
			logs = append(logs, logLine)
			d.clog.Infof("Build %s: %s", buildID, strings.TrimSpace(logLine))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return logs, err
		}
	}
	return logs, nil
}
