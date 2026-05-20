package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DockerHub struct {
	GenericRegistry
}

func NewDockerHub(client *http.Client) *GenericRegistry {
	r := NewGenericRegistry("registry-1.docker.io", client)
	r.normalizeRepo = func(repo string) string {
		if !strings.Contains(repo, "/") {
			return "library/" + repo
		}
		return repo
	}
	r.tokenFetcher = func(repo string) (string, error) {
		url := fmt.Sprintf(
			"https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull",
			repo,
		)
		resp, err := client.Get(url)
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

	return r
}

func (r *DockerHub) NormalizeRepo(repo string) string {
	if !strings.Contains(repo, "/") {
		return "library/" + repo
	}
	return repo
}

func (r *DockerHub) GetToken(repo string) (string, error) {
	url := fmt.Sprintf(
		"https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull",
		repo,
	)
	resp, err := r.client.Get(url)
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
