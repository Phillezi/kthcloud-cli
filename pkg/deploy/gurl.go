//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run gurl.go <url> <output-file> <remove-path-1> [<remove-path-2> ...]")
		os.Exit(1)
	}

	url := os.Args[1]
	outFile := os.Args[2]
	pathsToRemove := os.Args[3:]

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		os.Exit(1)
	}

	// Check if we have cached headers
	cacheFile := outFile + ".cache"
	if _, err := os.Stat(cacheFile); err == nil {
		if data, err := os.ReadFile(cacheFile); err == nil {
			headers := map[string]string{}
			if err := json.Unmarshal(data, &headers); err == nil {
				if etag, ok := headers["etag"]; ok {
					req.Header.Set("If-None-Match", etag)
				}
				if lastMod, ok := headers["last-modified"]; ok {
					req.Header.Set("If-Modified-Since", lastMod)
				}
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to download: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		fmt.Printf("Spec not modified, using cached file: %s\n", outFile)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Bad status: %s\n", resp.Status)
		os.Exit(1)
	}

	contentType := resp.Header.Get("Content-Type")
	file, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if contentType == "application/json" || contentType == "application/openapi+json" {
		var spec map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&spec); err != nil {
			fmt.Printf("Failed to parse JSON: %v\n", err)
			os.Exit(1)
		}

		// Remove paths
		if paths, ok := spec["paths"].(map[string]interface{}); ok {
			for _, p := range pathsToRemove {
				delete(paths, p)
			}
		}

		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		if err := enc.Encode(spec); err != nil {
			fmt.Printf("Failed to write JSON: %v\n", err)
			os.Exit(1)
		}

		// Save caching headers
		headers := map[string]string{}
		if etag := resp.Header.Get("ETag"); etag != "" {
			headers["etag"] = etag
		}
		if lastMod := resp.Header.Get("Last-Modified"); lastMod != "" {
			headers["last-modified"] = lastMod
		}
		if data, err := json.Marshal(headers); err == nil {
			_ = os.WriteFile(cacheFile, data, 0644)
		}

		fmt.Printf("Downloaded JSON from %s, removed paths %v -> %s\n", url, pathsToRemove, outFile)
		return
	}

	// Otherwise, just copy the body
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Downloaded %s -> %s\n", url, outFile)
}
