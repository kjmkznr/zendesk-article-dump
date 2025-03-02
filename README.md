# Zendesk Article Dump

A Go application that dumps Zendesk Help Center articles to Markdown format.

## Features

- Fetches all articles from your Zendesk Help Center
- Converts articles to Markdown format
- Supports pagination for large article collections
- Preserves article metadata (ID, URL, creation date, etc.)
- Handles special characters in filenames

## Prerequisites

- Go 1.24 or later
- Zendesk account with API access
- API token from Zendesk

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/kjmkznr/zendesk-article-dump.git
   cd zendesk-article-dump
   ```

2. Copy the example environment file and edit it with your Zendesk credentials:
   ```bash
   cp .env.example .env
   ```

3. Edit `.env` file with your Zendesk credentials:
   ```
   ZENDESK_SUBDOMAIN=your-subdomain
   ZENDESK_EMAIL=your-email@example.com
   ZENDESK_API_TOKEN=your-api-token
   ```

   - `ZENDESK_SUBDOMAIN`: Your Zendesk subdomain (if your Zendesk URL is `https://example.zendesk.com`, then your subdomain is `example`)
   - `ZENDESK_EMAIL`: Your Zendesk email address
   - `ZENDESK_API_TOKEN`: Your Zendesk API token

## Usage

1. Build the application:
   ```bash
   go build
   ```

2. Run the application:
   ```bash
   ./zendesk-article-dump
   ```

The application will:
1. Create an `articles` directory in the current folder
2. Download all articles from your Zendesk Help Center
3. Convert each article to Markdown format
4. Save articles as individual files named `{id}-{title}.md`

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
- Missing environment variables
- API connection issues
- File system operations
- Invalid responses

If an error occurs while saving a specific article, the application will log the error and continue with the remaining articles.

## License

MIT License