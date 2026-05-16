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
	Images  storage.Storage
	Log     utils.Logger
	Cfg     *config.Config
	Storage storage.Storage
}

func (s *ImageService) PushImage(
	ctx context.Context,
	req *pb.PushImageRequest,
) (*pb.PushImageResponse, error) {
	imageName := req.GetImageName()
	imageData := req.GetImageData()
	if imageName == "" {
		err := status.Errorf(codes.InvalidArgument, "Image name is empty")
		s.Log.Error(err.Error())
		return nil, err
	}
	err := s.Storage.SaveImage(imageName, imageData)
	if err != nil {
		s.Log.Error(err.Error())
		return nil, err
	}

	s.Log.Info(fmt.Sprintf("Image %s pushed", imageName))
	return &pb.PushImageResponse{Message: "Image pushed successfully"}, nil
}

func (s *ImageService) PullImage(
	ctx context.Context,
	req *pb.PullImageRequest,
) (*pb.PullImageResponse, error) {
	imageName := req.GetImageName()
	imageData, err := s.Storage.LoadImage(imageName)
	if err != nil {
		s.Log.Error(err.Error())
		return nil, err
	}

	s.Log.Info(fmt.Sprintf("Image %s pulled", imageName))
	return &pb.PullImageResponse{ImageData: []byte(imageData)}, nil
}

func (s *ImageService) ListImages(
	ctx context.Context,
	req *pb.ListImagesRequest,
) (*pb.ListImagesResponse, error) {
	images, err := s.Storage.ListImages()
	if err != nil {
		s.Log.Error(err.Error())
		return nil, err
	}
	var pbImages []*pb.ImageInfo
	for _, img := range images {
		pbImages = append(pbImages, &pb.ImageInfo{
			ImageName: img.ImageName,
			Tags:      img.Tags,
		})
	}

	s.Log.Info("Images listed")
	return &pb.ListImagesResponse{Images: pbImages}, nil
}

func (s *ImageService) DeleteImage(
	ctx context.Context,
	req *pb.DeleteImageRequest,
) (*pb.DeleteImageResponse, error) {
	imageName := req.GetImageName()
	err := s.Storage.DeleteImage(imageName)
	if err != nil {
		s.Log.Error(err.Error())
		return nil, err
	}
	s.Log.Info(fmt.Sprintf("Image %s deleted", imageName))
	return &pb.DeleteImageResponse{Message: "Image deleted successfully"}, nil
}

func (s *ImageService) DeleteImageTag(
	ctx context.Context,
	req *pb.DeleteImageTagRequest,
) (*pb.DeleteImageTagResponse, error) {
	imageName := req.GetImageName()
	tag := req.GetTag()
	err := s.Storage.DeleteImageTag(imageName, tag)
	if err != nil {
		s.Log.Error(err.Error())
		return nil, err
	}

	s.Log.Info(fmt.Sprintf("Tag %s deleted from image %s", tag, imageName))
	return &pb.DeleteImageTagResponse{Message: "Tag deleted successfully"}, nil
}
