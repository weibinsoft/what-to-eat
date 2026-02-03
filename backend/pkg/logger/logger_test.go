package logger

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "default config",
			cfg: &Config{
				Level:      "info",
				Format:     "console",
				OutputPath: "stdout",
			},
			wantErr: false,
		},
		{
			name: "json format",
			cfg: &Config{
				Level:      "debug",
				Format:     "json",
				OutputPath: "stdout",
			},
			wantErr: false,
		},
		{
			name: "stderr output",
			cfg: &Config{
				Level:      "warn",
				Format:     "console",
				OutputPath: "stderr",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset for each test
			Log = nil

			err := initLogger(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && Log == nil {
				t.Error("Log should not be nil after successful init")
			}
		})
	}
}

func TestInitDefault(t *testing.T) {
	Log = nil
	InitDefault()

	if Log == nil {
		t.Error("Log should not be nil after InitDefault()")
	}
}

func TestLogFunctions(t *testing.T) {
	// Reset and init
	Log = nil
	InitDefault()

	// These should not panic
	t.Run("Info", func(t *testing.T) {
		Info("test info message")
	})

	t.Run("Debug", func(t *testing.T) {
		Debug("test debug message")
	})

	t.Run("Warn", func(t *testing.T) {
		Warn("test warn message")
	})

	t.Run("Error", func(t *testing.T) {
		Error("test error message")
	})
}

func TestWith(t *testing.T) {
	Log = nil
	InitDefault()

	logger := With()
	if logger == nil {
		t.Error("With() should return a logger")
	}
}

func TestSync(t *testing.T) {
	Log = nil
	InitDefault()

	err := Sync()
	if err != nil {
		// Sync may return error for stdout/stderr, which is ok
		t.Logf("Sync() returned error (expected for stdout): %v", err)
	}
}

func TestInit_FileOutput(t *testing.T) {
	// Create a temp file
	tmpFile, err := os.CreateTemp("", "logger_test_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	Log = nil
	err = initLogger(&Config{
		Level:      "info",
		Format:     "json",
		OutputPath: tmpFile.Name(),
	})

	if err != nil {
		t.Errorf("Init() with file output failed: %v", err)
	}

	// Write a log entry
	Info("test file output")
	Sync()

	// Check file has content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Errorf("Failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Log file should have content")
	}
}

func TestLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"debug", "debug"},
		{"info", "info"},
		{"warn", "warn"},
		{"error", "error"},
		{"invalid", "invalid"}, // should default to info
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Log = nil
			err := initLogger(&Config{
				Level:      tt.level,
				Format:     "console",
				OutputPath: "stdout",
			})

			if err != nil {
				t.Errorf("Init() with level %s failed: %v", tt.level, err)
			}
		})
	}
}
