package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wan2-video-server/pkg/config"
	"github.com/wan2-video-server/pkg/logger"
	"github.com/wan2-video-server/pkg/model"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the Wan2.1 model from Hugging Face",
	Long:  `Download and cache the Wan2.1 model files from Hugging Face Hub.`,
	Run: func(cmd *cobra.Command, args []string) {
		downloadModel(cmd)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().String("model-id", "Lightricks/LTX-Video", "Hugging Face model ID")
	downloadCmd.Flags().String("cache-dir", "./models", "Model cache directory")
}

func downloadModel(cmd *cobra.Command) {
	log := logger.NewLogger()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}

	modelID, err := cmd.Flags().GetString("model-id")
	if err != nil {
		log.Fatalf("Failed to get model-id flag: %v", err)
	}
	cacheDir, err := cmd.Flags().GetString("cache-dir")
	if err != nil {
		log.Fatalf("Failed to get cache-dir flag: %v", err)
	}

	if modelID != "" {
		cfg.Model.HuggingFaceModelID = modelID
	}
	if cacheDir != "" {
		cfg.Model.CacheDir = cacheDir
	}

	log.Infof("Downloading model: %s", cfg.Model.HuggingFaceModelID)
	log.Infof("Cache directory: %s", cfg.Model.CacheDir)

	downloader := model.NewHuggingFaceDownloader(cfg)
	if err := downloader.Download(); err != nil {
		log.Fatalf("Failed to download model: %v", err)
		return
	}

	fmt.Println("Model downloaded successfully!")
}
