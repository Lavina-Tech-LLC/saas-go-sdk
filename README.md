# saas-go-sdk

Go SDK for Lavina SaaS platform.

## Install

```bash
go get github.com/Lavina-Tech-LLC/saas-go-sdk@latest
```

## Publish a new version

```bash
git tag v1.0.0
git push --tags
```

GitHub Actions will run tests and create a GitHub Release automatically.
The module is immediately available via `pkg.go.dev` after tagging.
