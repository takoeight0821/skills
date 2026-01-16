package remote

import "context"

// Fetcher defines the interface for fetching remote skills/plugins.
type Fetcher interface {
	// Fetch downloads the remote content to a local destination.
	Fetch(ctx context.Context, url string, dest string) error
}
