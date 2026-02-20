package utils

import (
	"net/url"
	"path/filepath"
	"strings"
)

// ExtractURLPath extracts the full path from a URL including the host
// Example: https://example.com/a/b/file.zip -> example.com/a/b
func ExtractURLPath(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Get the host (domain)
	host := parsed.Host
	
	// Get the path without the filename
	urlPath := parsed.Path
	
	// Remove leading slash
	urlPath = strings.TrimPrefix(urlPath, "/")
	
	// Get directory part (without filename)
	dir := filepath.Dir(urlPath)
	
	// If dir is ".", it means no subdirectories
	if dir == "." {
		return host, nil
	}
	
	// Combine host with directory path
	// Use filepath.Join to handle path separators correctly
	fullPath := filepath.Join(host, dir)
	
	return fullPath, nil
}
