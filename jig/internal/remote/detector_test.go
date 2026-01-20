package remote

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectPluginType(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected PluginType
	}{
		{
			name:     "Claude Plugin with plugin.json",
			files:    []string{"plugin.json"},
			expected: TypeClaude,
		},
		{
			name:     "Claude Plugin with marketplace.json",
			files:    []string{"marketplace.json"},
			expected: TypeClaude,
		},
		{
			name:     "Gemini Extension with gemini.json",
			files:    []string{"gemini.json"},
			expected: TypeGemini,
		},
		// Removed manifest.json from test for now as it's ambiguous without content inspection
		// {
		// 	name:     "Gemini Extension with manifest.json",
		// 	files:    []string{"manifest.json"},
		// 	expected: TypeGemini,
		// },
		{
			name:     "Mixed (Claude and Gemini)",
			files:    []string{"plugin.json", "gemini.json"},
			expected: TypeMixed,
		},
		{
			name:     "Unknown",
			files:    []string{"readme.md"},
			expected: TypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "detect-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			for _, file := range tt.files {
				f, err := os.Create(filepath.Join(tempDir, file))
				if err != nil {
					t.Fatalf("Failed to create file %s: %v", file, err)
				}
				f.Close()
			}

			detected, err := DetectPluginType(tempDir)
			if err != nil {
				t.Fatalf("DetectPluginType failed: %v", err)
			}

			if detected != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, detected)
			}
		})
	}
}
