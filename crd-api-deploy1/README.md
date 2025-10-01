# CRD API Deploy

A REST API for deploying SimpleAPI Custom Resource Definitions to Kubernetes clusters using the Huma framework and Go templates.

## Features

- **REST API**: Built with Huma framework for automatic OpenAPI documentation
- **Template-based CRD Generation**: Uses `.templ` files with Go template syntax (`{{.Value}}`)
- **Dynamic Parameter Substitution**: All CRD fields can be customized via API inputs
- **Kubernetes Integration**: Direct integration with Kubernetes API using client-go
- **Automatic Defaults**: Sensible default values for optional fields

## Quick Start

1. **Clone and setup:**
   ```bash
   git clone <repository-url>
   cd crd-api-deploy
   go mod tidy
   ```

2. **Test the template system:**
   ```bash
   make test-template
   ```

3. **Run the API server:**
   ```bash
   make run
   ```

4. **Access the documentation:**
   ```
   http://localhost:8080/docs
   ```

## Template System

The project uses Go templates (`.templ` files) for dynamic CRD generation:

### Template File: `templates/simpleapi.templ`

```yaml
apiVersion: apps.api.test/v1alpha1
kind: {{.Kind}}
metadata:
  labels:
{{- range .Labels}}
    {{.Key}}: {{.Value}}
{{- end}}
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  image: "{{.Image}}"
  version: "{{.Version}}"
  port: {{.Port}}
  replicas: {{.Replicas}}
  # ... more fields with {{.Value}} substitution
```

### Template Engine Usage

```go
engine, err := template.NewEngine()
crdYAML, err := engine.GenerateCRD(request)
```

The template engine automatically:
- Loads the `.templ` file from the filesystem
- Substitutes all `{{.Value}}` placeholders with API request data
- Handles conditionals and loops for complex fields
- Generates valid Kubernetes YAML

## API Endpoints

### Create SimpleAPI CRD
```http
POST /api/v1/simpleapis
Content-Type: application/json

{
  "kind": "Simpleapi",
  "name": "my-api",
  "namespace": "default",
  "image": "nginx",
  "version": "latest",
  "labels": [
    {"key": "app", "value": "my-app"}
  ]
}
```

### Get SimpleAPI Resource
```http
GET /api/v1/simpleapis/{name}?namespace=default
```

### List SimpleAPI Resources
```http
GET /api/v1/simpleapis?namespace=default
```

## Dynamic Parameters

All these fields are configurable via API:

### Required Parameters
- `kind`: CRD kind (e.g., "Simpleapi")
- `name`: Resource name
- `namespace`: Target namespace
- `image`: Container image
- `version`: Image version/tag

### Optional Parameters (with defaults)
- `labels[]`: Custom labels
- `port`: Container port (default: 8000)
- `replicas`: Number of replicas (default: 1)
- `resources`: CPU/memory limits and requests
- `affinity`: Pod affinity rules
- `tolerations`: Pod tolerations
- And many more...

## Examples

### Minimal API Call
```bash
curl -X POST http://localhost:8080/api/v1/simpleapis \
  -H "Content-Type: application/json" \
  -d '{
    "kind": "Simpleapi",
    "name": "minimal-api",
    "namespace": "default",
    "image": "nginx",
    "version": "latest"
  }'
```

### Full Configuration
See `examples/create-simpleapi.json` for a complete example with all parameters.

## Project Structure

```
├── templates/
│   └── simpleapi.templ          # CRD template with {{.Value}} syntax
├── internal/
│   ├── models/
│   │   └── api.go              # API request/response models
│   ├── template/
│   │   └── engine.go           # Template engine for .templ files
│   ├── k8s/
│   │   └── client.go           # Kubernetes client wrapper
│   └── service/
│       └── simpleapi.go        # Business logic service
├── cmd/
│   ├── api/
│   │   └── main.go             # Application entry point
│   └── router/
│       └── router.go           # HTTP routes and handlers
├── examples/                   # Example API payloads
└── scripts/                    # Utility scripts
```

## Template Customization

To modify the CRD template:

1. Edit `templates/simpleapi.templ`
2. Use Go template syntax: `{{.FieldName}}`
3. Add conditionals: `{{- if .Field}}...{{- end}}`
4. Add loops: `{{- range .Array}}...{{- end}}`
5. Restart the API server

Example template addition:
```yaml
# Add new field to template
{{- if .CustomField}}
  customField: {{.CustomField}}
{{- end}}
```

Then add the field to API model:
```go
type CreateSimpleAPIRequest struct {
    // ... existing fields
    CustomField string `json:"customField,omitempty"`
}
```

## Authentication

The API uses Kubernetes RBAC for authorization. Ensure the service account has permissions to:
- Create/update/get/list SimpleAPI CRDs
- Create namespaces (if needed)

## Development

```bash
# Build
make build

# Run locally  
make run

# Test template system
make test-template

# Clean up
make clean
```

## Docker

```bash
# Build image
make docker-build

# Run container
docker run -p 8080:8080 -v ~/.kube:/root/.kube:ro crd-api-deploy:latest
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes (including template modifications)
4. Test with `make test-template`
5. Submit a pull request

The template system makes it easy to extend the CRD structure without changing Go code - just modify the `.templ` file!