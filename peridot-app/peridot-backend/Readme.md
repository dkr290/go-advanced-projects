**Trivy** | Vulnerability scanning | `apt install trivy` or download binary |
| **Syft** | SBOM generation | `go install github.com/anchore/syft/cmd/syft@latest` |
| **Anchore** | Image analysis | Docker container or CLI |
| **Distroless** | Minimal base images | Docker official images |

---

## Next Steps

1. **Install Trivy** for vulnerability scanning
2. **Complete the scanner module** - integrate Trivy CLI
3. **Implement patcher** - apply security patches to images
4. **Test the API** with `curl http://localhost:8080/images`

Would you like me to continue with the Trivy scanner implementation and patcher module?

