# Zendesk Article Dump

A Go application that dumps publicly available Zendesk Help Center articles to Markdown format.

## Features

- Fetches all public articles from your Zendesk Help Center
- Converts articles to Markdown format
- Supports pagination for large article collections
- Preserves article metadata (ID, URL, creation date, etc.)
- Handles special characters in filenames

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

### Examples

Basic usage:
```bash
./zendesk-article-dump -subdomain=example
```

Specify custom output directory:
```bash
./zendesk-article-dump -subdomain=example -output=docs
```

The application will:
1. Create the specified output directory (default: `articles`) if it doesn't exist
2. Download all public articles from your Zendesk Help Center
3. Convert each article to Markdown format
4. Save articles as individual files named `{id}-{title}.md` in the output directory

## Output Format

Each article is saved as a Markdown file with the following structure:

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

## Error Handling

The application includes error handling for:
- Missing subdomain configuration
- API connection issues
- File system operations
- Invalid responses

If an error occurs while saving a specific article, the application will log the error and continue with the remaining articles.

## License

MIT License
