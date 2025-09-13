package services

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/dkr290/go-advanced-projects/api-docker-builder/internal/models"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type DockerService struct {
	client *client.Client
	builds map[string]*models.BuildStatus
	clog   *logrus.Logger
}

func NewDockerService(clog *logrus.Logger) (*DockerService, error) {
	var cli *client.Client
	var err error
	if err := ensureDockerDaemon(clog); err != nil {
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
		clog:   clog,
	}, nil
}

func ensureDockerDaemon(clog *logrus.Logger) error {
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
