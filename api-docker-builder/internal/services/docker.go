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
	go d.performBuild(context.Background(), buildID, req, imageName)

	return &models.BuildImageResponse{
		BuildID:   buildID,
		Status:    "pending",
		Message:   "Build started",
		ImageName: imageName,
		StartedAt: status.StartedAt,
	}, nil
}

func (d *DockerService) GetBuildStatus(buildID string) (*models.BuildStatus, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	status, exists := d.builds[buildID]
	if !exists {
		return nil, fmt.Errorf("build not found")
	}
	// Return a copy to avoid race conditions
	statusCopy := *status
	if status.CompletedAt != nil {
		completedAt := *status.CompletedAt
		statusCopy.CompletedAt = &completedAt
	}
	if status.Logs != nil {
		statusCopy.Logs = make([]string, len(status.Logs))
		copy(statusCopy.Logs, status.Logs)
	}

	return &statusCopy, nil
}

func (d *DockerService) performBuild(
	ctx context.Context,
	buildID string,
	req *models.BuildImageRequest,
	imageName string,
) {
	// Add timeout context for the build
	buildCtx, cancel := context.WithTimeout(ctx, 30*time.Minute) // 30 min timeout
	defer cancel()

	d.updateBuildStatus(buildID, "building", "Creating Dockerfile and building image", nil)
	// buildID to create unique directory - this prevents conflicts
	targetDir := fmt.Sprintf("./data/build-%s", buildID)
	if err := os.RemoveAll(targetDir); err != nil {
		d.clog.Warnf("Failed to clean directory %s: %v", targetDir, err)
	}
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		d.updateBuildStatus(
			buildID,
			"failed",
			fmt.Sprintf("failed to create build directory: %v", err),
			nil,
		)
		return
	}

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
		return
	}
	defer func() {
		err := os.RemoveAll(targetDir)
		d.clog.Warnf("Failed to cleanup directory %s: %v", targetDir, err)
	}()
	// Find Dockerfile path (relative to repoDir)
	dockerfilePath, err := findDockerfilePathFromDir(targetDir)
	fmt.Println(dockerfilePath)
	if err != nil {
		d.updateBuildStatus(buildID, "failed", fmt.Sprintf("no Dockerfile found: %v", err), nil)
		return
	}
	relativeDockerfilePath, err := filepath.Rel(targetDir, dockerfilePath)
	if err != nil {
		d.updateBuildStatus(
			buildID,
			"failed",
			fmt.Sprintf("failed to get relative path: %v", err),
			nil,
		)
		return
	}

	// Build image and the Dockerfile should be in the archive
	buildOptions := types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: relativeDockerfilePath,
		Remove:     true,
		NoCache:    true,
	}
	// Create build context
	buildContext, err := d.createBuildContext(targetDir)
	if err != nil {
		d.updateBuildStatus(
			buildID,
			"failed",
			fmt.Sprintf("Failed to create build context: %v", err),
			nil,
		)
		return
	}
	resp, err := d.client.ImageBuild(buildCtx, buildContext, buildOptions)
	if err != nil {
		if buildCtx.Err() == context.DeadlineExceeded {
			d.updateBuildStatus(buildID, "failed", "Build timeout exceeded", nil)
		} else {
			d.updateBuildStatus(buildID, "failed", fmt.Sprintf("Failed to build image: %v", err), nil)
		}
		d.clog.Errorf("Error building the image %v", err)
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
	d.mutex.Lock()
	defer d.mutex.Unlock()
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
	fmt.Println(repoURL)
	d.clog.Infof("Cloning the repository into  %s", targetDir)
	_, err := git.PlainClone(targetDir, false, cloneOptions)
	if err != nil {
		d.clog.Errorf("Error cloning repository: %v\n", err)
		return fmt.Errorf("failed to clone repository: %w", err)
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
	contextDir string,
) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// Walk through the entire context directory
	err := filepath.Walk(contextDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files/directories (like .git)
		if info.IsDir() {
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Get relative path from context directory
		relPath, err := filepath.Rel(contextDir, path)
		if err != nil {
			return err
		}

		// Read file content
		fileContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Create tar header
		header := &tar.Header{
			Name: filepath.ToSlash(relPath), // Use forward slashes for tar
			Mode: 0o644,
			Size: int64(len(fileContent)),
		}

		// Write header and content to tar
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header for %s: %w", relPath, err)
		}

		if _, err := tw.Write(fileContent); err != nil {
			return fmt.Errorf("failed to write file content for %s: %w", relPath, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create build context: %w", err)
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
