# Test Results CLI

A command-line interface for collecting, formatting, and sending test results to the Test Results API.

## 📆 Overview

This CLI tool is intended to be used by developers, CI pipelines, or automation scripts to:

* Parse and collect test result artifacts (e.g., `junit.xml`, `coverage.xml`)
* Optionally format results for console or JSON output
* Post results to a central API endpoint for aggregation and reporting

## 🚀 Getting Started

### Build the CLI

From the `cli/` directory:

```bash
go build -o testresults ./cmd/testresults
```

### Run the CLI

```bash
./testresults --help
```

## 🔧 Usage

> Example usage (actual subcommands and flags will depend on your implementation):

```bash
./testresults collect --input tests/junit.xml --format json
./testresults push --api-url https://api.example.com/results
```

## 🧪 Testing

To run unit tests:

```bash
go test ./...
```

## 📁 Project Structure

```
cli/
├── cmd/               # Main entry point
├── internal/          # CLI internals (client, collector, config, etc.)
├── pkg/               # Optional shared utilities
├── tests/             # Unit and integration tests
└── go.mod             # Go module definition
```

## ✅ Requirements

* Go 1.21+
* Access to the Test Results API

## 🔐 Environment Variables

You can optionally configure the CLI with environment variables:

| Variable    | Description                       |
| ----------- | --------------------------------- |
| `API_URL`   | Base URL of the Test Results API  |
| `API_TOKEN` | Optional token for authentication |

## 📄 License

MIT © 2025 \[Your Name or Organization]
