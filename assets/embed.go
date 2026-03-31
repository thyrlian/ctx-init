package assetdata

import (
	"embed"
	"fmt"
)

//go:embed adapters/*.md
var fs embed.FS

func ReadAdapter(name string) ([]byte, error) {
	path := "adapters/" + name
	b, err := fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read asset %q: %w", path, err)
	}
	return b, nil
}
