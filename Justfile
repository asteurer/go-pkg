[doc('Show available commands')]
@default:
    just --list

[doc('Run linter')]
@lint:
    echo "Lint..."
    command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is not installed or in PATH"; exit 1; }
    GOARCH=wasm GOOS=wasip1 golangci-lint run

[doc('Validates that a componentize-go.toml file is NOT present in the root of the repo, as it causes issues with some of the defaults in componentize-go')]
@ensure_no_componentize_go_toml:
    test ! -e componentize-go.toml || { echo "componentize-go.toml must not exist"; exit 1; }
