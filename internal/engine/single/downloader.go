package single

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"

	"github.com/surge-downloader/surge/internal/engine/types"
	"github.com/surge-downloader/surge/internal/utils"
)

// SingleDownloader handles single-threaded downloads for servers that don't support range requests.
// NOTE: Pause/resume is NOT supported because this downloader is only used when
// the server doesn't support Range headers. If interrupted, the download must restart.
type SingleDownloader struct {
	Client       *http.Client
	ProgressChan chan<- any           // Channel for events (start/complete/error)
	ID           string               // Download ID
	State        *types.ProgressState // Shared state for TUI polling
	Runtime      *types.RuntimeConfig
	Headers      map[string]string // Custom HTTP headers (cookies, auth, etc.)
}

// NewSingleDownloader creates a new single-threaded downloader with all required parameters
func NewSingleDownloader(id string, progressCh chan<- any, state *types.ProgressState, runtime *types.RuntimeConfig) *SingleDownloader {
	transport := &http.Transport{}
	
	// Configure proxy if runtime config is provided
	if runtime != nil && runtime.ProxyURL != "" {
		parsedURL, err := url.Parse(runtime.ProxyURL)
		if err != nil {
			utils.Debug("Single downloader: Invalid proxy URL %s: %v", runtime.ProxyURL, err)
			transport.Proxy = http.ProxyFromEnvironment
		} else if strings.HasPrefix(parsedURL.Scheme, "socks5") {
			utils.Debug("Single downloader: Using SOCKS5 proxy: %s", runtime.ProxyURL)
			dialer, dialErr := proxy.SOCKS5("tcp", parsedURL.Host, nil, proxy.Direct)
			if dialErr != nil {
				utils.Debug("Single downloader: Failed to create SOCKS5 dialer: %v", dialErr)
				transport.Proxy = http.ProxyFromEnvironment
			} else {
				transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				}
			}
		} else {
			transport.Proxy = http.ProxyURL(parsedURL)
		}
	} else {
		transport.Proxy = http.ProxyFromEnvironment
	}
	
	// Configure TLS if runtime config is provided
	if runtime != nil && runtime.SkipTLSVerification {
		utils.Debug("Single downloader: TLS verification disabled")
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		} else {
			transport.TLSClientConfig.InsecureSkipVerify = true
		}
	}
	
	client := &http.Client{
		Timeout:   0,
		Transport: transport,
	}
	
	return &SingleDownloader{
		Client:       client,
		ProgressChan: progressCh,
		ID:           id,
		State:        state,
		Runtime:      runtime,
	}
}

// Download downloads a file using a single connection.
// This is used for servers that don't support Range requests.
// If interrupted, the download cannot be resumed and must restart from the beginning.
func (d *SingleDownloader) Download(ctx context.Context, rawurl, destPath string, fileSize int64, filename string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawurl, nil)
	if err != nil {
		return err
	}

	for key, val := range d.Headers {
		req.Header.Set(key, val)
	}
	req.Header.Set("User-Agent", d.Runtime.GetUserAgent())

	resp, err := d.Client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			utils.Debug("Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Use .surge extension for incomplete file
	workingPath := destPath + types.IncompleteSuffix
	outFile, err := os.Create(workingPath)
	if err != nil {
		return err
	}

	// Track whether we completed successfully for cleanup
	success := false
	defer func() {
		_ = outFile.Close()
		if !success {
			_ = os.Remove(workingPath)
		}
	}()

	start := time.Now()

	// Copy response body to file with context cancellation support
	var written int64
	buf := make([]byte, d.Runtime.GetWorkerBufferSize())

	for {
		// Check for context cancellation (allows clean shutdown)
		select {
		case <-ctx.Done():
			// Can't resume - server doesn't support Range requests
			return ctx.Err()
		default:
		}

		nr, readErr := resp.Body.Read(buf)
		if nr > 0 {
			nw, writeErr := outFile.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				if d.State != nil {
					d.State.Downloaded.Store(written)
					d.State.VerifiedProgress.Store(written)
				}
			}
			if writeErr != nil {
				return fmt.Errorf("write error: %w", writeErr)
			}
			if nr != nw {
				return io.ErrShortWrite
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break // Done reading
			}
			return fmt.Errorf("read error: %w", readErr)
		}
	}

	if err := outFile.Sync(); err != nil {
		return fmt.Errorf("sync error: %w", err)
	}
	if err := outFile.Close(); err != nil {
		return fmt.Errorf("close error: %w", err)
	}

	// Rename .surge file to final destination
	if err := os.Rename(workingPath, destPath); err != nil {
		// Fallback: copy if rename fails (cross-device)
		if copyErr := copyFile(workingPath, destPath); copyErr != nil {
			return fmt.Errorf("failed to finalize file: %w", copyErr)
		}
		_ = os.Remove(workingPath)
	}

	success = true // Mark successful so defer doesn't clean up

	elapsed := time.Since(start)
	speed := float64(written) / elapsed.Seconds()
	utils.Debug("\nDownloaded %s in %s (%s/s)\n",
		destPath,
		elapsed.Round(time.Second),
		utils.ConvertBytesToHumanReadable(int64(speed)),
	)

	return nil
}

// copyFile copies a file from src to dst (fallback when rename fails)
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := in.Close(); err != nil {
			utils.Debug("Error closing input file: %v", err)
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			utils.Debug("Error closing output file: %v", err)
		}
	}()

	buf := make([]byte, 1024*1024)
	if _, err := io.CopyBuffer(out, in, buf); err != nil {
		return err
	}
	return out.Sync()
}
