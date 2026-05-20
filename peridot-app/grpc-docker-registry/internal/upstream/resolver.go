package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ParseRef splits an image ref into (host, repo, tag/digest).
// e.g. "quay.io/prometheus/prometheus:latest" -> ("quay.io", "prometheus/prometheus", "latest")
// e.g. "alpine:latest"                        -> ("", "alpine", "latest")
// ParseRef splits "[[host/]repo[:tag]]" into host, repo, reference.
func ParseRef(ref string) (host, repo, reference string) {
	reference = "latest"

	// Split off tag (last colon not inside a host:port)
	if idx := strings.LastIndex(ref, ":"); idx != -1 && !strings.Contains(ref[idx:], "/") {
		reference = ref[idx+1:]
		ref = ref[:idx]
	}

	segments := strings.SplitN(ref, "/", 2)
	if len(segments) == 2 && (strings.Contains(segments[0], ".") || strings.Contains(segments[0], ":")) {
		host = segments[0]
		repo = segments[1]
	} else {
		host = ""
		repo = ref
	}
	return
}

func NewFromRef(ref string, client *http.Client) (UpstreamRegistry, string, string, error) {
	host, repo, reference := ParseRef(ref)

	var reg UpstreamRegistry
	switch host {
	case "", "docker.io", "registry-1.docker.io":
		reg = NewDockerHub(client)
	case "quay.io":
		reg = NewQuayRegistry(client)
	default:
		reg = NewGenericRegistry(host, client)
	}

	repo = reg.NormalizeRepo(repo)
	return reg, repo, reference, nil
}

func ResolvePlatformDigest(manifestList []byte, osName, arch string) (string, error) {
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
		if m.Platform.OS == osName && m.Platform.Architecture == arch {
			return m.Digest, nil
		}
	}
	return "", fmt.Errorf("no manifest found for %s/%s", osName, arch)
}


