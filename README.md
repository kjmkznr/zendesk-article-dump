# Zendesk Article Dump

A Go application that dumps publicly available Zendesk Help Center articles to Markdown format.

## Features

- Fetches all public articles from your Zendesk Help Center
- Converts articles to Markdown format
- Supports pagination for large article collections
- Preserves article metadata (ID, URL, creation date, etc.)
- Handles special characters in filenames
- Option to combine all articles into a single markdown file

## Prerequisites

- Go 1.24 or later
- Public Zendesk Help Center URL

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/kjmkznr/zendesk-article-dump.git
   cd zendesk-article-dump
   ```

2. Build the application:
   ```bash
   go build
   ```

## Usage

Run the application with the required subdomain parameter:

```bash
./zendesk-article-dump -subdomain=your-subdomain
```

### Parameters

- `-subdomain` (required): Your Zendesk subdomain (if your Zendesk URL is `https://example.zendesk.com`, then your subdomain is `example`)
- `-output` (optional): Output directory for markdown files (default: "articles")
- `-combine` (optional): Boolean flag to combine all articles into a single file (default: false)

### Examples

Basic usage (separate files):
```bash
./zendesk-article-dump -subdomain=example
```

Specify custom output directory:
```bash
./zendesk-article-dump -subdomain=example -output=docs
```

Combine all articles into a single file:
```bash
./zendesk-article-dump -subdomain=example -combine
```

Note: When using -combine, articles will be combined into a file named 'combined.md' in the output directory

The application will:
1. Create the specified output directory (default: `articles`) if it doesn't exist
2. Download all public articles from your Zendesk Help Center
3. Convert each article to Markdown format
4. Either:
   - Save articles as individual files in the output directory (default behavior)
   - Combine all articles into a single file if -combine is specified

## Output Format

### Individual Files (Default)

When saving articles individually, each article is saved as a separate Markdown file with the following structure:

```markdown
# Article Title

- ID: 123456
- URL: https://example.zendesk.com/hc/en-us/articles/123456
- Created: 2024-03-02T10:00:00Z
- Updated: 2024-03-02T10:00:00Z
- Locale: en-us

---

Article content in Markdown format...
```

### Combined File

When using the `-combine` flag, all articles are combined into a single file named 'combined.md' in the output directory. Articles are separated by a horizontal line and follow the same format as individual files:

```markdown
# First Article Title

- ID: 123456
- URL: https://example.zendesk.com/hc/en-us/articles/123456
- Created: 2024-03-02T10:00:00Z
- Updated: 2024-03-02T10:00:00Z
- Locale: en-us

---

First article content...

--------------------------------------------------

# Second Article Title

- ID: 789012
- URL: https://example.zendesk.com/hc/en-us/articles/789012
- Created: 2024-03-02T11:00:00Z
- Updated: 2024-03-02T11:00:00Z
- Locale: en-us

---

Second article content...
```

## Error Handling

The application includes error handling for:
- Missing subdomain configuration
- API connection issues
- File system operations
- Invalid responses

If an error occurs while saving a specific article, the application will log the error and continue with the remaining articles.

## License

MIT License
