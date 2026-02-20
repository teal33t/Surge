# New Features Added

## 1. Skip TLS Verification

### Description
Added a setting to skip TLS certificate verification for downloads from servers with certificate issues (e.g., self-signed certificates, expired certificates, hostname mismatches).

### Settings Location
- **Category**: Network
- **Setting Key**: `skip_tls_verification`
- **Label**: "Skip TLS Verification"
- **Type**: Boolean (checkbox)
- **Default**: `false`
- **Description**: "Skip TLS certificate verification (insecure, use only for trusted sources with certificate issues)."

### Implementation Details
- Added `SkipTLSVerification` field to `NetworkSettings` struct
- Propagated through `RuntimeConfig` to download engines
- Applied to:
  - Concurrent downloader HTTP client
  - Single downloader HTTP client
  - Probe server HTTP client
  - Mirror probe HTTP client

### Files Modified
- `internal/config/settings.go` - Added setting definition
- `internal/engine/types/config.go` - Added to RuntimeConfig
- `internal/engine/concurrent/downloader.go` - Applied TLS config
- `internal/engine/single/downloader.go` - Applied TLS config
- `internal/engine/probe.go` - Applied TLS config and updated signature
- `internal/download/manager.go` - Updated ProbeServer calls

### Security Warning
This feature should only be used for trusted sources. Skipping TLS verification makes connections vulnerable to man-in-the-middle attacks.

---

## 2. Preserve URL Path Structure

### Description
Added a setting to preserve the URL path structure when saving downloaded files. When enabled, files are saved in subdirectories that mirror the URL structure.

### Example
```
URL: https://aaaa.onion/a/b/d.zip
Without feature: download_dir/d.zip
With feature:    download_dir/aaaa.onion/a/b/d.zip
```

### Settings Location
- **Category**: General
- **Setting Key**: `preserve_url_path`
- **Label**: "Preserve URL Path"
- **Type**: Boolean (checkbox)
- **Default**: `false`
- **Description**: "Preserve the URL path structure when saving files (e.g., example.com/a/b/file.zip → download_dir/example.com/a/b/file.zip)."

### Implementation Details
- Added `PreserveURLPath` field to `GeneralSettings` struct
- Propagated through `RuntimeConfig`
- Created utility function `ExtractURLPath()` to parse URL and extract host + path
- Modified download manager to create subdirectories based on URL structure
- Automatically creates necessary directories
- Falls back to simple path if URL parsing or directory creation fails

### Files Created
- `internal/utils/urlpath.go` - URL path extraction utility
- `internal/utils/urlpath_test.go` - Comprehensive tests

### Files Modified
- `internal/config/settings.go` - Added setting definition
- `internal/engine/types/config.go` - Added to RuntimeConfig
- `internal/download/manager.go` - Implemented path preservation logic

### Benefits
- Organizes downloads by source domain
- Maintains original URL structure for easier file management
- Useful for mirroring websites or organizing downloads from multiple sources
- Prevents filename conflicts from different sources

---

## Testing

Both features have been implemented with:
- Proper error handling and fallback mechanisms
- Integration with existing settings system
- Backward compatibility (both default to `false`)
- No breaking changes to existing functionality

## Usage

Users can enable these features through the Settings UI:
1. Open Settings (press 's' in TUI)
2. Navigate to the appropriate tab (Network or General)
3. Toggle the checkbox for the desired feature
4. Settings are saved automatically
