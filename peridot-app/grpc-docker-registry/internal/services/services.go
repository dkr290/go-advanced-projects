package services

import (
	"context"
	"fmt"

	"github.com/dkr290/peridot-app/grpc-docker-registry/config"
	pb "github.com/dkr290/peridot-app/grpc-docker-registry/proto/gen"
	"github.com/dkr290/peridot-app/grpc-docker-registry/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ImageService struct {
	pb.UnimplementedImageRegistryServiceServer
	Images map[string][]string
	Log    utils.Logger
	Cfg    *config.Config
}

func (s *ImageService) PushImage(ctx context.Context,	req *pb.PushImageRequest) (*pb.PushImageResponse, error) {

	imageName := req.GetImageName()
	imageData := req.GetImageData()
  if imageName == "" {
		err := status.Errorf(codes.InvalidArgument, "Image name is empty")
		s.Log.Error(err.Error())
		return nil, err
	}
	s.Images[imageName] = append(s.Images[imageName], string(imageData))
	s.Log.Info(fmt.Sprintf("Image %s pushed", imageName))
	return &pb.PushImageResponse{Message: "Image pushed successfully"}, nil
}

func (s *ImageService) PullImage(ctx context.Context, req *pb.PullImageRequest) (*pb.PullImageResponse, error) {
	imageName := req.GetImageName()
	imageData := s.Images[imageName][len(s.Images[imageName])-1]
	s.Log.Info(fmt.Sprintf("Image %s pulled", imageName))
	return &pb.PullImageResponse{ImageData: []byte(imageData)}, nil
}


func (s *ImageService) ListImages(ctx context.Context, req *pb.ListImagesRequest) (*pb.ListImagesResponse, error) {
	var imagesResponse pb.ListImagesResponse
	for imageName, tags := range s.Images {
		imagesResponse.Images = append(imagesResponse.Images, &pb.ImageInfo{
			ImageName: imageName,
			Tags:      tags,
		})
	}
	s.Log.Info("Images listed")
	return &imagesResponse, nil
}


func (s *ImageService) DeleteImage(ctx context.Context, req *pb.DeleteImageRequest) (*pb.DeleteImageResponse, error) {
	imageName := req.GetImageName()
	delete(s.Images, imageName)
	s.Log.Info(fmt.Sprintf("Image %s deleted", imageName))
	return &pb.DeleteImageResponse{Message: "Image deleted successfully"}, nil
}


func (s *ImageService) DeleteImageTag(ctx context.Context, req *pb.DeleteImageTagRequest) (*pb.DeleteImageTagResponse, error) {
	imageName := req.GetImageName()
	tag := req.GetTag()
	tags := s.Images[imageName]
	for i, t := range tags {
		if t == tag {
			s.Images[imageName] = append(tags[:i], tags[i+1:]...)
			break
		}
	}
	s.Log.Info(fmt.Sprintf("Tag %s deleted from image %s", tag, imageName))
	return &pb.DeleteImageTagResponse{Message: "Tag deleted successfully"}, nil
}
