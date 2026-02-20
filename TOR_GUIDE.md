# 🧅 Downloading from Tor Hidden Services (.onion)

Surge now supports downloading from Tor hidden services (`.onion` domains) via SOCKS5 proxy support!

## Prerequisites

### 1. Install Tor

**Ubuntu/Debian/WSL:**
```bash
sudo apt update
sudo apt install tor
sudo service tor start
```

**macOS:**
```bash
brew install tor
brew services start tor
```

**Arch Linux:**
```bash
sudo pacman -S tor
sudo systemctl start tor
sudo systemctl enable tor
```

**Windows:**
Download from https://www.torproject.org/download/

### 2. Verify Tor is Running

Tor runs on `127.0.0.1:9050` by default (SOCKS5 proxy)

```bash
# Check if Tor is listening
netstat -tuln | grep 9050
# or
ss -tuln | grep 9050
```

You should see output like:
```
tcp   0   0 127.0.0.1:9050   0.0.0.0:*   LISTEN
```

## Configuration

### Method 1: Via TUI (Recommended)

1. Start Surge:
   ```bash
   surge
   ```

2. Press `s` to open Settings

3. Navigate to **Network** tab (use `→` or Tab)

4. Configure Proxy:
   - Select **Proxy URL**
   - Press `Enter`
   - Type: `socks5://127.0.0.1:9050`
   - Press `Enter` to confirm

5. Enable TLS Skip (recommended for .onion sites):
   - Select **Skip TLS Verification**
   - Press `Enter` to toggle to `True`

6. Press `Esc` to save and close

### Method 2: Via Settings File

Edit `~/.surge/settings.json`:

```json
{
  "general": {
    ...
  },
  "network": {
    "proxy_url": "socks5://127.0.0.1:9050",
    "skip_tls_verification": true,
    ...
  },
  "performance": {
    ...
  }
}
```

### Method 3: Via Environment Variable

```bash
export SURGE_PROXY=socks5://127.0.0.1:9050
surge https://example.onion/file.zip
```

## Usage Examples

### Basic Download
```bash
surge https://example.onion/file.zip
```

### Download to Specific Directory
```bash
surge https://example.onion/file.zip -o ~/Downloads/
```

### Multiple Downloads
```bash
surge https://site1.onion/file1.zip https://site2.onion/file2.zip
```

### Batch Download from File
Create a file `onion-urls.txt`:
```
https://example1.onion/file1.zip
https://example2.onion/file2.zip
https://example3.onion/archive.tar.gz
```

Then:
```bash
surge --batch onion-urls.txt
```

## Troubleshooting

### Connection Refused
```
Error: dial tcp 127.0.0.1:9050: connect: connection refused
```

**Solution:** Tor is not running. Start it:
```bash
sudo service tor start  # Linux
brew services start tor # macOS
```

### Timeout Errors
```
Error: context deadline exceeded
```

**Solutions:**
1. Tor network might be slow - this is normal for .onion sites
2. The .onion site might be down
3. Try increasing timeout (future feature)

### Certificate Errors
```
Error: x509: certificate is not valid
```

**Solution:** Enable "Skip TLS Verification" in settings (see Configuration above)

### Slow Downloads
.onion downloads are typically slower than clearnet due to:
- Multiple hops through Tor network
- Bandwidth limitations of Tor relays
- Server-side bandwidth limits

**Tips:**
- Be patient - Tor is designed for anonymity, not speed
- Reduce `max_connections_per_host` to avoid overwhelming Tor
- Consider using fewer concurrent downloads

## Security Considerations

### ⚠️ Important Notes

1. **Anonymity**: Using Surge with Tor provides anonymity for your downloads, but:
   - Your ISP can see you're using Tor (but not what you're downloading)
   - Downloaded files are stored on your disk unencrypted
   - File metadata may contain identifying information

2. **TLS Skip**: When you enable "Skip TLS Verification":
   - You're vulnerable to man-in-the-middle attacks
   - Only use for trusted .onion sites
   - Tor provides encryption at the network layer, but TLS adds another layer

3. **Legal**: 
   - Tor is legal in most countries
   - Downloading illegal content is illegal regardless of the network used
   - Be aware of your local laws

## Recommended Settings for Tor

```json
{
  "network": {
    "max_connections_per_host": 8,
    "max_concurrent_downloads": 2,
    "proxy_url": "socks5://127.0.0.1:9050",
    "skip_tls_verification": true
  }
}
```

Lower connection counts help avoid overwhelming Tor relays and improve stability.

## Advanced: Using Tor Browser's SOCKS Port

If you have Tor Browser installed, it runs on port `9150`:

```bash
surge --proxy socks5://127.0.0.1:9150 https://example.onion/file.zip
```

Or in settings:
```json
{
  "network": {
    "proxy_url": "socks5://127.0.0.1:9150"
  }
}
```

## Testing Your Setup

Test with a known .onion site:

```bash
# Test connectivity (this is a test .onion address)
surge https://www.thehiddenwiki.org/
```

If you see the download starting, your Tor setup is working!

## Performance Tips

1. **Reduce Connections**: Lower `max_connections_per_host` to 4-8
2. **Fewer Concurrent Downloads**: Set `max_concurrent_downloads` to 1-2
3. **Patience**: .onion downloads are inherently slower
4. **Check Tor Status**: Ensure Tor has established circuits
   ```bash
   sudo systemctl status tor
   ```

## Support

If you encounter issues:
1. Check Tor is running: `netstat -tuln | grep 9050`
2. Test Tor connectivity: `curl --socks5 127.0.0.1:9050 https://check.torproject.org`
3. Check Surge logs in `~/.surge/logs/`
4. Enable debug mode (if available)

---

**Happy anonymous downloading! 🧅**
