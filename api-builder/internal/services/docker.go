// Package services is the main docker build process
package services

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/dkr290/go-advanced-projects/api-builder/internal/logging"
	"github.com/dkr290/go-advanced-projects/api-builder/internal/models"
	"github.com/dkr290/go-advanced-projects/api-builder/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

type DockerService struct {
	client     *client.Client
	builds     map[string]*models.BuildStatus
	buildMutex sync.RWMutex
}

var clog = logging.Init(false)

func NewDockerService() (*DockerService, error) {
	var cli *client.Client
	var err error
	if err := ensureDockerDaemon(); err != nil {
		return nil, err
	}

	dockerHosts := []string{
		"unix:///var/run/docker.sock", // Sidecar approach
		"tcp://localhost:2375",        // DinD on same pod
		"tcp://docker:2375",           // If docker service exists
	}

	for _, host := range dockerHosts {
		clog.Infof("Trying Docker host: %s", host)
		opts := []client.Opt{
			client.WithHost(host),
			client.WithAPIVersionNegotiation(),
		}
		if strings.HasPrefix(host, "tcp://") {
			opts = append(opts, client.WithTimeout(10*time.Second))
		}
		cli, err = client.NewClientWithOpts(opts...)
		if err != nil {
			clog.Warnf("Failed to create client for %s: %v", host, err)
			continue

		}
		// Test connection with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = cli.Ping(ctx)
		cancel()
		if err != nil {
			clog.Warnf("Failed to ping Docker at %s: %v", host, err)
			cli.Close()
			continue
		} else {
			clog.Infof("Successfully connected to Docker at %s", host)
			break
		}
	}

	return &DockerService{
		client: cli,
		builds: make(map[string]*models.BuildStatus),
	}, nil
}

func ensureDockerDaemon() error {
	// We're in Docker, check if daemon is running
	if _, err := net.Dial("unix", "/var/run/docker.sock"); err != nil {
		clog.Info("Docker daemon not running, starting it...")

		// Start dockerd
		cmd := exec.Command("dockerd",
			"--host=unix:///var/run/docker.sock",
			"--host=tcp://0.0.0.0:2375",
			"--tls=false",
			"--storage-driver=vfs")

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start dockerd: %v", err)
		}

		// Wait for daemon to be ready
		for range 10 {
			if _, err := net.Dial("unix", "/var/run/docker.sock"); err == nil {
				clog.Info("✅ Docker daemon started successfully")
				return nil
			}
			time.Sleep(1 * time.Second)
		}
		return fmt.Errorf("docker daemon failed to start within 10 seconds")
	}
	return nil
}

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

	d.buildMutex.Lock()
	d.builds[buildID] = status
	d.buildMutex.Unlock()

	// Start build in goroutine
	go d.performBuild(ctx, buildID, req, imageName)

	return &models.BuildImageResponse{
		BuildID:   buildID,
		Status:    "pending",
		Message:   "Build started",
		ImageName: imageName,
		StartedAt: status.StartedAt,
	}, nil
}

func (d *DockerService) GetBuildStatus(buildID string) (*models.BuildStatus, error) {
	d.buildMutex.RLock()
	defer d.buildMutex.RUnlock()

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

	// Generate Dockerfile based on model version
	dockerfile := d.generateDockerfile(req)

	// Create build context
	buildContext, err := d.createBuildContext(dockerfile)
	if err != nil {
		d.updateBuildStatus(
			buildID,
			"failed",
			fmt.Sprintf("Failed to create build context: %v", err),
			nil,
		)
		return
	}

	// Build image
	buildOptions := types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: "Dockerfile",
		Remove:     true,
	}

	resp, err := d.client.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		d.updateBuildStatus(buildID, "failed", fmt.Sprintf("Failed to build image: %v", err), nil)
		return
	}
	defer resp.Body.Close()

	// Read build logs
	logs := []string{}
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			logLine := string(buf[:n])
			logs = append(logs, logLine)
			clog.Infof("Build %s: %s", buildID, strings.TrimSpace(logLine))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			d.updateBuildStatus(
				buildID,
				"failed",
				fmt.Sprintf("Error reading build logs: %v", err),
				logs,
			)
			return
		}
	}

	// Update final status
	completedAt := time.Now()
	d.buildMutex.Lock()
	if build, exists := d.builds[buildID]; exists {
		build.Status = "success"
		build.Message = "Image built successfully"
		build.CompletedAt = &completedAt
		build.Logs = logs
	}
	d.buildMutex.Unlock()

	clog.Infof("Build %s completed successfully for image %s", buildID, imageName)
}

func (d *DockerService) updateBuildStatus(buildID, status, message string, logs []string) {
	d.buildMutex.Lock()
	defer d.buildMutex.Unlock()

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

func (d *DockerService) generateDockerfile(req *models.BuildImageRequest) string {
	// Generate different Dockerfiles based on model version
	switch strings.ToLower(req.ModelVersion) {
	case "python-flask":
		return utils.GeneratePythonFlaskDockerfile(req)
	case "python-fastapi":
		return utils.GeneratePythonFastAPIDockerfile(req)
	case "nodejs":
		return utils.GenerateNodeJSDockerfile(req)
	default:
		// Default to Python Flask
		return utils.GeneratePythonFlaskDockerfile(req)
	}
}

func (d *DockerService) createBuildContext(dockerfile string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// Add Dockerfile
	header := &tar.Header{
		Name: "Dockerfile",
		Mode: 0o644,
		Size: int64(len(dockerfile)),
	}
	if err := tw.WriteHeader(header); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(dockerfile)); err != nil {
		return nil, err
	}

	// Add sample application files based on the dockerfile type
	if strings.Contains(dockerfile, "python") {
		// Add sample Python app
		appContent := `from flask import Flask, jsonify

app = Flask(__name__)

@app.route('/')
def hello():
    return jsonify({
        "message": "Hello from Docker built image!",
        "status": "success"
    })

@app.route('/health')
def health():
    return jsonify({"status": "healthy"})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
`
		appHeader := &tar.Header{
			Name: "app.py",
			Mode: 0o644,
			Size: int64(len(appContent)),
		}
		if err := tw.WriteHeader(appHeader); err != nil {
			return nil, err
		}
		if _, err := tw.Write([]byte(appContent)); err != nil {
			return nil, err
		}

		// Add requirements.txt
		reqContent := "Flask==2.3.3\n"
		reqHeader := &tar.Header{
			Name: "requirements.txt",
			Mode: 0o644,
			Size: int64(len(reqContent)),
		}
		if err := tw.WriteHeader(reqHeader); err != nil {
			return nil, err
		}
		if _, err := tw.Write([]byte(reqContent)); err != nil {
			return nil, err
		}
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
