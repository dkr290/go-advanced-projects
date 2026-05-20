package upstream

type UpstreamRegistry interface {
	GetToken(repo string) (string, error)
	GetManifest(repo, reference string) ([]byte, string, string, error)
	GetBlob(repo, digest string) ([]byte, error)
	NormalizeRepo(repo string) string
}

