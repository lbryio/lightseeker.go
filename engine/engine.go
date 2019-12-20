package engine

type SearchResponse []struct {
	ClaimID   string
	ClaimName string
}

type SearchVersion struct {
	SemVer     string
	CommitHash string
}

type SearchEngine interface {
	Query(string) (SearchResponse, error)
	Version() (*SearchVersion, error)
	GetEndpoint() string
}
