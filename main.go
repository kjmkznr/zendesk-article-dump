package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func cleanupHTML(content string) string {
	// Remove script tags and their content
	content = regexp.MustCompile(`(?s)<script.*?</script>`).ReplaceAllString(content, "")

	// Convert images to markdown format
	content = regexp.MustCompile(`<img[^>]*src="([^"]*)"[^>]*alt="([^"]*)"[^>]*(?:width="([^"]*)")?[^>]*(?:height="([^"]*)")?[^>]*>`).ReplaceAllStringFunc(content, func(img string) string {
		re := regexp.MustCompile(`src="([^"]*)"`)
		srcMatch := re.FindStringSubmatch(img)
		src := srcMatch[1]

		re = regexp.MustCompile(`alt="([^"]*)"`)
		altMatch := re.FindStringSubmatch(img)
		alt := ""
		if len(altMatch) > 1 {
			alt = altMatch[1]
		}

		return fmt.Sprintf("![%s](%s)", alt, src)
	})

	// Convert headers (h1-h6) with any attributes
	for i := 6; i >= 1; i-- {
		pattern := fmt.Sprintf(`(?s)<h%d[^>]*>(.*?)</h%d>`, i, i)
		replacement := fmt.Sprintf("\n%s $1\n", strings.Repeat("#", i))
		content = regexp.MustCompile(pattern).ReplaceAllString(content, replacement)
	}

	// Convert pre tags to code blocks with any attributes
	content = regexp.MustCompile(`(?s)<pre[^>]*>(.*?)</pre>`).ReplaceAllString(content, "\n```\n$1\n```\n")

	// Convert links with any attributes
	content = regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*>(.*?)</a>`).ReplaceAllString(content, "[$2]($1)")

	// Convert lists (handle nested lists)
	content = regexp.MustCompile(`(?s)<[uo]l[^>]*>.*?</[uo]l>`).ReplaceAllStringFunc(content, func(list string) string {
		// Convert list items to markdown format
		list = regexp.MustCompile(`<li[^>]*>(.*?)</li>`).ReplaceAllString(list, "- $1\n")
		// Remove the ul/ol tags
		list = regexp.MustCompile(`</?[uo]l[^>]*>`).ReplaceAllString(list, "\n")
		return list
	})

	// Convert br tags to newlines (handle self-closing tags)
	content = regexp.MustCompile(`<br[^>]*>`).ReplaceAllString(content, "\n")

	// Convert paragraphs and other block elements with any attributes
	content = regexp.MustCompile(`<(?:p|div|span)[^>]*>(.*?)</(?:p|div|span)>`).ReplaceAllString(content, "$1\n\n")

	// Remove any remaining HTML tags and their attributes
	content = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(content, "")

	// Fix multiple newlines (3 or more become 2)
	content = regexp.MustCompile(`\n{3,}`).ReplaceAllString(content, "\n\n")

	// Fix multiple spaces
	content = regexp.MustCompile(`[ \t]+`).ReplaceAllString(content, " ")

	// Trim spaces and normalize newlines
	return strings.TrimSpace(content)
}

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
		subdomain     string
		outputDir     string
		combineOutput bool
	)
	flag.StringVar(&subdomain, "subdomain", "", "Zendesk subdomain (required)")
	flag.StringVar(&outputDir, "output", "articles", "Output directory for markdown files")
	flag.BoolVar(&combineOutput, "combine", false, "Combine all articles into a single file")
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

	// Initialize articles slice for combined output
	var allArticles []Article

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
		if combineOutput {
			// For combined output, collect all articles first
			allArticles = append(allArticles, articles...)
		} else {
			// For individual files, save each article separately
			for _, article := range articles {
				if err := saveArticleAsMarkdown(article, outputDir); err != nil {
					fmt.Printf("Error saving article %d: %v\n", article.ID, err)
					continue
				}
			}
		}

		nextPage = next
	}

	// If combining output, save all articles to a single file
	if combineOutput {
		combinedFile := filepath.Join(outputDir, "articles.md")
		if err := saveArticlesCombined(allArticles, combinedFile); err != nil {
			fmt.Printf("Error saving combined articles: %v\n", err)
			os.Exit(1)
		}
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
	// Create filename using only article ID
	filename := fmt.Sprintf("%d.md", article.ID)
	outputPath := filepath.Join(outputDir, filename)

	// Create markdown content
	content := fmt.Sprintf("# %s\n\n", article.Title)
	content += fmt.Sprintf("- ID: %d\n", article.ID)
	content += fmt.Sprintf("- URL: %s\n", article.HTMLUrl)
	content += fmt.Sprintf("- Created: %s\n", article.CreatedAt)
	content += fmt.Sprintf("- Updated: %s\n", article.UpdatedAt)
	content += fmt.Sprintf("- Locale: %s\n\n", article.Locale)
	content += fmt.Sprintf("---\n\n%s\n", cleanupHTML(article.Body))

	// Write to file
	return os.WriteFile(outputPath, []byte(content), 0644)
}

func saveArticlesCombined(articles []Article, outputPath string) error {
	if len(articles) == 0 {
		return fmt.Errorf("no articles to save")
	}

	var content string
	for i, article := range articles {
		// Add separator between articles
		if i > 0 {
			content += "\n\n" + strings.Repeat("-", 50) + "\n\n"
		}

		// Add article content in the same format as individual files
		content += fmt.Sprintf("# %s\n\n", article.Title)
		content += fmt.Sprintf("- ID: %d\n", article.ID)
		content += fmt.Sprintf("- URL: %s\n", article.HTMLUrl)
		content += fmt.Sprintf("- Created: %s\n", article.CreatedAt)
		content += fmt.Sprintf("- Updated: %s\n", article.UpdatedAt)
		content += fmt.Sprintf("- Locale: %s\n\n", article.Locale)
		content += fmt.Sprintf("---\n\n%s", cleanupHTML(article.Body))
	}

	// Ensure content ends with a newline
	content += "\n"

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	// Write to file
	return os.WriteFile(outputPath, []byte(content), 0644)
}
