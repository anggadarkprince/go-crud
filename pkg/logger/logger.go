package logger

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anggadarkprince/crud-employee-go/configs"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *slog.Logger

// Initialize sets up the global logger
func Initialize() error {
	var handler slog.Handler

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// Setup log rotation using lumberjack
	logFile := &lumberjack.Logger{
		Filename: filepath.Join("logs", "app.log"),
		MaxSize: 100, // megabytes
		MaxBackups: 3,   // number of backups
		MaxAge: 30,  // days
		Compress: true,
	}

	if configs.Get().App.Environment == "production" {
		// Production: JSON format to both file and stdout
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			AddSource: true,
		})
	} else {
		multiWriter := io.MultiWriter(logFile)
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelDebug,
			AddSource: true,
		})
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)

	return nil
}

func LogError(message string, err error, r *http.Request) {
    stackTrace := fmt.Sprintf("%+v", err)
    
    Log.Error(message,
        slog.String("method", r.Method),
        slog.String("path", r.URL.Path),
        slog.String("error", err.Error()),
        slog.String("stack_trace", stackTrace),
    )
    
    // Also print to console
	if configs.Get().App.Environment == "development" {
    	fmt.Printf("\n=== Error ===\n%s %s\n%+v\n========\n", r.Method, r.URL.Path, err)
	}
}