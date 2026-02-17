[doc('Show available commands')]
@default:
    just --list

[doc('Run linter')]
@lint:
    echo "Lint..."
    command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is not installed or in PATH"; exit 1; }
    GOARCH=wasm GOOS=wasip1 golangci-lint run
