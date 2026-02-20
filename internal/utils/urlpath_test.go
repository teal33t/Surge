package utils

import (
	"path/filepath"
	"testing"
)

func TestExtractURLPath(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple URL with path",
			url:      "https://example.com/a/b/file.zip",
			expected: filepath.Join("example.com", "a", "b"),
			wantErr:  false,
		},
		{
			name:     "Onion URL with path",
			url:      "https://aaaa.onion/a/b/d.zip",
			expected: filepath.Join("aaaa.onion", "a", "b"),
			wantErr:  false,
		},
		{
			name:     "URL with no subdirectories",
			url:      "https://example.com/file.zip",
			expected: "example.com",
			wantErr:  false,
		},
		{
			name:     "URL with deep path",
			url:      "https://cdn.example.com/downloads/2024/01/archive.tar.gz",
			expected: filepath.Join("cdn.example.com", "downloads", "2024", "01"),
			wantErr:  false,
		},
		{
			name:     "URL with port",
			url:      "https://example.com:8080/path/to/file.bin",
			expected: filepath.Join("example.com:8080", "path", "to"),
			wantErr:  false,
		},
		{
			name:     "URL with query parameters",
			url:      "https://example.com/download/file.zip?token=abc123",
			expected: filepath.Join("example.com", "download"),
			wantErr:  false,
		},
		{
			name:     "HTTP URL",
			url:      "http://mirror.org/files/data.tar",
			expected: filepath.Join("mirror.org", "files"),
			wantErr:  false,
		},
		{
			name:     "Invalid URL",
			url:      "://invalid",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "URL with trailing slash",
			url:      "https://example.com/folder/",
			expected: "example.com",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractURLPath(tt.url)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ExtractURLPath() expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("ExtractURLPath() unexpected error: %v", err)
				return
			}
			
			if result != tt.expected {
				t.Errorf("ExtractURLPath() = %q, want %q", result, tt.expected)
			}
		})
	}
}
