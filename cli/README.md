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

```plaintext
.
├── cmd/                    # CLI command entry points (for Cobra or similar)
├── go.mod                  # Go module file
├── internal/               # Internal packages (not importable by others)
│   ├── client/             # API/HTTP or external system interaction logic
│   ├── collector/          # Logic for gathering or aggregating data
│   ├── config/             # Configuration loading (e.g., Viper-based or custom)
│   └── formatter/          # Output formatting (JSON, table, CSV, etc.)
├── main.go                 # Entry point (invokes cmd/ logic)
├── pkg/                    # Public utility packages (can be reused externally)
│   └── utils/
│       └── io.go           # General-purpose IO helpers
├── README.md               # Documentation
└── tree.txt                # Tree snapshot (probably for reference)
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
