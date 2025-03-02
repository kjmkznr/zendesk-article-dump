package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Article struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	HTMLUrl   string `json:"html_url"`
	Locale    string `json:"locale"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
	Position  int    `json:"position"`
	SectionID int64  `json:"section_id"`
}

type ArticlesResponse struct {
	Articles []Article `json:"articles"`
	Next     string    `json:"next_page"`
	Previous string    `json:"previous_page"`
}

func main() {
	// Get environment variables
	subdomain := os.Getenv("ZENDESK_SUBDOMAIN")
	email := os.Getenv("ZENDESK_EMAIL")
	token := os.Getenv("ZENDESK_API_TOKEN")

	if subdomain == "" || email == "" || token == "" {
		fmt.Println("Error: Required environment variables are not set")
		fmt.Println("Please set ZENDESK_SUBDOMAIN, ZENDESK_EMAIL, and ZENDESK_API_TOKEN")
		os.Exit(1)
	}

	// Create output directory
	outputDir := "articles"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize HTTP client
	client := &http.Client{}

	// Fetch articles
	baseURL := fmt.Sprintf("https://%s.zendesk.com/api/v2/help_center/articles.json", subdomain)
	nextPage := baseURL

	for nextPage != "" {
		articles, next, err := fetchArticles(client, nextPage, email, token)
		if err != nil {
			fmt.Printf("Error fetching articles: %v\n", err)
			os.Exit(1)
		}

		// Process articles
		for _, article := range articles {
			if err := saveArticleAsMarkdown(article, outputDir); err != nil {
				fmt.Printf("Error saving article %d: %v\n", article.ID, err)
				continue
			}
		}

		nextPage = next
	}

	fmt.Println("Article dump completed successfully!")
}

func fetchArticles(client *http.Client, url, email, token string) ([]Article, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("error creating request: %v", err)
	}

	req.SetBasicAuth(email+"/token", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	var response ArticlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, "", fmt.Errorf("error decoding response: %v", err)
	}

	return response.Articles, response.Next, nil
}

func saveArticleAsMarkdown(article Article, outputDir string) error {
	// Create filename from article title
	filename := fmt.Sprintf("%d-%s.md", article.ID, sanitizeFilename(article.Title))
	filepath := filepath.Join(outputDir, filename)

	// Create markdown content
	content := fmt.Sprintf("# %s\n\n", article.Title)
	content += fmt.Sprintf("- ID: %d\n", article.ID)
	content += fmt.Sprintf("- URL: %s\n", article.HTMLUrl)
	content += fmt.Sprintf("- Created: %s\n", article.CreatedAt)
	content += fmt.Sprintf("- Updated: %s\n", article.UpdatedAt)
	content += fmt.Sprintf("- Locale: %s\n\n", article.Locale)
	content += fmt.Sprintf("---\n\n%s\n", article.Body)

	// Write to file
	return os.WriteFile(filepath, []byte(content), 0644)
}

func sanitizeFilename(filename string) string {
	// Replace invalid characters with underscore
	filename = strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) {
			return '_'
		}
		return r
	}, filename)

	// Trim spaces and limit length
	filename = strings.TrimSpace(filename)
	if len(filename) > 100 {
		filename = filename[:100]
	}

	return filename
}
