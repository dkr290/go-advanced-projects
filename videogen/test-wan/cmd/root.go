package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wan2-video-server/pkg/config"
	"github.com/wan2-video-server/pkg/logger"
	"github.com/wan2-video-server/pkg/server"
)

var (
	cfgFile string
	log     = logger.NewLogger()
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "wan2-video-server",
	Short: "Wan2.1 Video Generation Server",
	Long: `A high-performance video generation server for Wan2.1 model.
Supports text-to-video, image-to-video, and video-to-video generation
with GPU acceleration and Hugging Face integration.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .env)")
	rootCmd.PersistentFlags().String("host", "0.0.0.0", "Server host")
	rootCmd.PersistentFlags().Int("port", 8080, "Server port")
	rootCmd.PersistentFlags().Bool("gpu", true, "Enable GPU acceleration")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")

	viper.BindPFlag("server.host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("server.port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("gpu.enabled", rootCmd.PersistentFlags().Lookup("gpu"))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}
}

func runServer() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	// Set log level
	logger.SetLogLevel(cfg.Log.Level)

	log.Infof("Starting Wan2.1 Video Server...")
	log.Infof("Configuration loaded: GPU=%v, Backend=%s", cfg.GPU.Enabled, cfg.Model.Provider)

	// Create and start server
	srv := server.NewServer(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
