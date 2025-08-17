# Portfolio Helm Chart

A Helm chart for the Portfolio application with API, Frontend, Database, and Monitoring.

## Installation

### Add the Helm repository

```bash
helm repo add bruno-site https://brunovlucena.github.io/bruno-site
helm repo update
```

### Install the chart

```bash
helm install portfolio bruno-site/portfolio
```

### Install with custom values

```bash
helm install portfolio bruno-site/portfolio -f values.yaml
```

## Configuration

See the [values.yaml](values.yaml) file for all available configuration options.

## Development

### Local development

```bash
# Install dependencies
helm dependency build

# Lint the chart
helm lint .

# Dry run installation
helm install portfolio . --dry-run --debug

# Package the chart
helm package .
```

### Releasing

To release a new version:

1. Update the version in `Chart.yaml`
2. Create and push a new tag:
   ```bash
   git tag v0.1.1
   git push origin v0.1.1
   ```

The GitHub Actions workflow will automatically build and publish the chart to GitHub Pages.

## Repository Structure

```
.
├── Chart.yaml          # Chart metadata
├── values.yaml         # Default configuration values
├── templates/          # Kubernetes manifests
└── .helmignore         # Files to ignore when packaging
```

## Support

For issues and questions, please open an issue on the [GitHub repository](https://github.com/brunovlucena/bruno-site).
