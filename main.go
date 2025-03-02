package main

import (
	"encoding/json"
	"flag"
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
	// Parse command line arguments
	var (
		subdomain string
		outputDir string
	)
	flag.StringVar(&subdomain, "subdomain", "", "Zendesk subdomain (required)")
	flag.StringVar(&outputDir, "output", "articles", "Output directory for markdown files")
	flag.Parse()

	if subdomain == "" {
		fmt.Println("Error: Required argument is not provided")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create output directory
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
		articles, next, err := fetchArticles(client, nextPage)
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

func fetchArticles(client *http.Client, url string) ([]Article, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error making request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil && err == nil {
			err = fmt.Errorf("error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	var response ArticlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, "", fmt.Errorf("error decoding response: %v", err)
	}

	return response.Articles, response.Next, err
}

func saveArticleAsMarkdown(article Article, outputDir string) error {
	// Create filename from article title
	filename := fmt.Sprintf("%d-%s.md", article.ID, sanitizeFilename(article.Title))
	outputPath := filepath.Join(outputDir, filename)

	// Create markdown content
	content := fmt.Sprintf("# %s\n\n", article.Title)
	content += fmt.Sprintf("- ID: %d\n", article.ID)
	content += fmt.Sprintf("- URL: %s\n", article.HTMLUrl)
	content += fmt.Sprintf("- Created: %s\n", article.CreatedAt)
	content += fmt.Sprintf("- Updated: %s\n", article.UpdatedAt)
	content += fmt.Sprintf("- Locale: %s\n\n", article.Locale)
	content += fmt.Sprintf("---\n\n%s\n", article.Body)

	// Write to file
	return os.WriteFile(outputPath, []byte(content), 0644)
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
