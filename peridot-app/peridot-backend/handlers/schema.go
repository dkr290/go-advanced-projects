package handlers

import (
	"github.com/humaio/huma/v4"
)

// AddImageInput - Input for adding a new image
type AddImageInput struct {
	Body struct {
		ImageRef string `json:"image_ref" example:"ubuntu:22.04" doc:"Docker image reference"`
	}
}

type AddImageOutput struct {
	Body struct {
		Message         string   `json:"message"`
		Vulnerabilities []string `json:"vulnerabilities"`
	}
}

type ListImagesInput struct {
	Repository string `query:"repository" example:"*" doc:"Repository to filter by"`
}

type ListImagesOutput struct {
	Body struct {
		Images []ImageInfo `json:"images"`
	}
}

type ImageInfo struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Digest     string `json:"digest"`
}

type ScanImageInput struct {
	Repo string `param:"repo" example:"ubuntu"`
	Tag  string `param:"tag" example:"22.04"`
}

type ScanImageOutput struct {
	Body struct {
		Vulnerabilities []Vulnerability `json:"vulnerabilities"`
		Score           int             `json:"score"`
	}
}

type Vulnerability struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

type PatchImageInput struct {
	Repo string `param:"repo" example:"ubuntu"`
	Tag  string `param:"tag" example:"22.04"`
}

type PatchImageOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

