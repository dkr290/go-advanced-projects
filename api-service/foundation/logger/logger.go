package logger

func New(service string, outputPaths ...string) (*zap.SugaredLogger, error) {
  config := zap.NewProductionConfig()
  config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
  config.DisableStacktrace = true
}
