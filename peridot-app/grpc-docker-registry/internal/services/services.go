package services

import (
	"context"
	"fmt"

	"github.com/dkr290/peridot-app/grpc-docker-registry/config"
	"github.com/dkr290/peridot-app/grpc-docker-registry/internal/storage"
	pb "github.com/dkr290/peridot-app/grpc-docker-registry/proto/gen"
	"github.com/dkr290/peridot-app/grpc-docker-registry/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ImageService struct {
	pb.UnimplementedImageRegistryServiceServer
	Log           utils.Logger
	Cfg           *config.Config
	BlobStore     storage.BlobStore
	ManifestStore storage.ManifestStore
}

// ==================== Blob Operations ====================

func (s *ImageService) PushBlob(
	ctx context.Context,
	req *pb.PushBlobRequest,
) (*pb.PushBlobResponse, error) {
	digest := req.GetDigest()
	if digest == "" {
		digest = storage.ComputeDigest(req.GetData())
	}

	if err := s.BlobStore.StoreBlob(digest, req.GetData()); err != nil {
		s.Log.Error(fmt.Sprintf("Failed to store blob: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to store blob: %v", err)
	}

	s.Log.Info(fmt.Sprintf("Blob stored: %s (%d bytes)", digest, len(req.GetData())))
	return &pb.PushBlobResponse{Digest: digest, Size: int64(len(req.GetData()))}, nil
}

func (s *ImageService) PullBlob(
	ctx context.Context,
	req *pb.PullBlobRequest,
) (*pb.PullBlobResponse, error) {
	digest := req.GetDigest()
	data, err := s.BlobStore.GetBlob(digest)
	if err != nil {
		s.Log.Error(fmt.Sprintf("Blob not found: %s", digest))
		return nil, status.Errorf(codes.NotFound, "Blob not found: %s", digest)
	}

	return &pb.PullBlobResponse{Data: data}, nil
}

func (s *ImageService) DeleteBlob(
	ctx context.Context,
	req *pb.DeleteBlobRequest,
) (*pb.DeleteBlobResponse, error) {
	digest := req.GetDigest()
	if err := s.BlobStore.DeleteBlob(digest); err != nil {
		s.Log.Error(fmt.Sprintf("Failed to delete blob: %v", err))
		return nil, status.Errorf(codes.NotFound, "Blob not found: %s", digest)
	}

	s.Log.Info(fmt.Sprintf("Blob deleted: %s", digest))
	return &pb.DeleteBlobResponse{Message: "Blob deleted successfully"}, nil
}

func (s *ImageService) ListBlobs(
	ctx context.Context,
	req *pb.ListBlobsRequest,
) (*pb.ListBlobsResponse, error) {
	prefix := req.GetDigestPrefix()
	digests, err := s.BlobStore.ListBlobs(prefix)
	if err != nil {
		s.Log.Error(fmt.Sprintf("Failed to list blobs: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to list blobs: %v", err)
	}

	return &pb.ListBlobsResponse{Digests: digests}, nil
}

// ==================== Image Operations ====================

func (s *ImageService) PushImage(
	ctx context.Context,
	req *pb.PushImageRequest,
) (*pb.PushImageResponse, error) {
	repository := req.GetRepository()
	tag := req.GetTag()

	if repository == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Repository is required")
	}
	if tag == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Tag is required")
	}

	// Validate layers exist
	for _, layerDigest := range req.GetLayerDigests() {
		if !s.BlobStore.Exists(layerDigest) {
			return nil, status.Errorf(codes.NotFound, "Layer not found: %s", layerDigest)
		}
	}

	// Store config
	configDigest := storage.ComputeDigest(req.GetConfig())
	if err := s.BlobStore.StoreBlob(configDigest, req.GetConfig()); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to store config: %v", err)
	}

	// Store manifest
	manifestDigest := storage.ComputeDigest(req.GetManifest())
	if err := s.ManifestStore.StoreManifest(
		repository,
		tag,
		manifestDigest,
		req.GetManifest(),
	); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to store manifest: %v", err)
	}

	s.Log.Info(
		fmt.Sprintf(
			"Image pushed: %s:%s (manifest: %s, layers: %d)",
			repository,
			tag,
			manifestDigest,
			len(req.GetLayerDigests()),
		),
	)
	return &pb.PushImageResponse{
		Reference:      fmt.Sprintf("%s:%s", repository, tag),
		ManifestDigest: manifestDigest,
	}, nil
}

func (s *ImageService) PullImage(
	ctx context.Context,
	req *pb.PullImageRequest,
) (*pb.PullImageResponse, error) {
	repository := req.GetRepository()
	tag := req.GetTag()

	manifest, manifestDigest, err := s.ManifestStore.GetManifest(repository, tag)
	if err != nil {
		s.Log.Error(fmt.Sprintf("Image not found: %s:%s", repository, tag))
		return nil, status.Errorf(codes.NotFound, "Image not found: %s:%s", repository, tag)
	}

	// Parse layer digests from manifest (simplified - in real impl, parse manifest JSON)
	// For now, we'll store layer digests separately or extract from manifest
	// This is a placeholder - you'll need to implement manifest parsing

	s.Log.Info(fmt.Sprintf("Image pulled: %s:%s (manifest: %s)", repository, tag, manifestDigest))
	return &pb.PullImageResponse{
		Manifest:     manifest,
		Config:       []byte{},   // Would need to store config separately
		LayerDigests: []string{}, // Would need to extract from manifest
	}, nil
}

func (s *ImageService) ListImages(
	ctx context.Context,
	req *pb.ListImagesRequest,
) (*pb.ListImagesResponse, error) {
	repository := req.GetRepository()
	tags, err := s.ManifestStore.ListTags(repository)
	if err != nil {
		s.Log.Error(fmt.Sprintf("Failed to list tags: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to list tags: %v", err)
	}

	var images []*pb.ImageInfo
	for _, tag := range tags {
		manifest, digest, err := s.ManifestStore.GetManifest(repository, tag)
		if err != nil {
			continue
		}
		_ = manifest
		images = append(images, &pb.ImageInfo{
			Repository:     repository,
			Tag:            tag,
			ManifestDigest: digest,
		})
	}

	s.Log.Info(fmt.Sprintf("Images listed: %d", len(images)))
	return &pb.ListImagesResponse{Images: images}, nil
}

func (s *ImageService) DeleteImage(
	ctx context.Context,
	req *pb.DeleteImageRequest,
) (*pb.DeleteImageResponse, error) {
	repository := req.GetRepository()
	tag := req.GetTag()

	if err := s.ManifestStore.DeleteManifest(repository, tag); err != nil {
		s.Log.Error(fmt.Sprintf("Failed to delete image: %v", err))
		return nil, status.Errorf(codes.NotFound, "Image not found: %s:%s", repository, tag)
	}

	s.Log.Info(fmt.Sprintf("Image deleted: %s:%s", repository, tag))
	return &pb.DeleteImageResponse{Message: "Image deleted successfully"}, nil
}

func (s *ImageService) ListTags(
	ctx context.Context,
	req *pb.ListTagsRequest,
) (*pb.ListTagsResponse, error) {
	repository := req.GetRepository()
	tags, err := s.ManifestStore.ListTags(repository)
	if err != nil {
		s.Log.Error(fmt.Sprintf("Failed to list tags: %v", err))
		return nil, status.Errorf(codes.Internal, "Failed to list tags: %v", err)
	}

	s.Log.Info(fmt.Sprintf("Tags listed: %d", len(tags)))
	return &pb.ListTagsResponse{Tags: tags}, nil
}

func (s *ImageService) DeleteTag(
	ctx context.Context,
	req *pb.DeleteTagRequest,
) (*pb.DeleteTagResponse, error) {
	repository := req.GetRepository()
	tag := req.GetTag()

	if err := s.ManifestStore.DeleteTag(repository, tag); err != nil {
		s.Log.Error(fmt.Sprintf("Failed to delete tag: %v", err))
		return nil, status.Errorf(codes.NotFound, "Tag not found: %s:%s", repository, tag)
	}

	s.Log.Info(fmt.Sprintf("Tag deleted: %s:%s", repository, tag))
	return &pb.DeleteTagResponse{Message: "Tag deleted successfully"}, nil
}


