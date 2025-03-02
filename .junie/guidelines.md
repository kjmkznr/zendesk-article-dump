# Zendesk Article Dump - Project Guidelines

## Project Overview
Zendesk Article Dump is a Go application designed to efficiently download and convert publicly available Zendesk Help Center articles into Markdown format. This tool streamlines the process of content migration and backup from Zendesk Help Center.

## Project Structure
```
zendesk-article-dump/
├── main.go           # Main application entry point
├── articles/         # Default output directory for downloaded articles
├── .junie/          # Project guidelines and documentation
└── README.md        # Project documentation
```

## Development Guidelines

### Code Standards
- Use Go 1.24 or later
- Follow Go best practices and idioms
- Maintain clear error handling as demonstrated in the existing codebase
- Document new features and changes

### Feature Implementation Guidelines
When implementing new features:
1. Maintain backward compatibility with existing command-line parameters
2. Follow the established error handling patterns
3. Consider pagination for large data sets
4. Preserve article metadata in the output
5. Handle special characters in filenames appropriately

### Testing
- Add tests for new functionality
- Ensure existing tests pass before submitting changes
- Test with various Zendesk Help Center configurations

### Documentation
- Update README.md for significant changes
- Document new command-line parameters
- Include examples for new features
- Keep code comments current

## Contribution Guidelines
1. Fork the repository
2. Create a feature branch
3. Follow the development guidelines
4. Submit a pull request with a clear description of changes
5. Ensure all tests pass

## Build and Deployment
1. Build using `go build`
2. Test the binary with various parameters
3. Verify output formatting matches specifications
4. Check error handling for edge cases

## Support
For issues and feature requests:
- Use GitHub issues
- Provide clear reproduction steps
- Include relevant error messages
- Specify the Go version and environment details