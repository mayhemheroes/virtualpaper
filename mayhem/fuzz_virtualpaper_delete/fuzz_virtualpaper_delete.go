package fuzz_virtualpaper_delete

import (
	"os"

	fuzz "github.com/AdaLogics/go-fuzz-headers"

	"tryffel.net/go/virtualpaper/config"
	"tryffel.net/go/virtualpaper/services/process"
)

// process.DeleteDocument -> storage.PreviewPath/DocumentPath dereference the global
// config.C (a *Config that is nil until the app loads its config), so without setup
// every non-trivial docId would nil-panic before exercising any path logic. Initialize
// a minimal config pointing at a throwaway temp dir so the harness actually fuzzes the
// id->path construction + os.Remove handling instead of crashing on a nil deref.
func init() {
	dir, err := os.MkdirTemp("", "vp-fuzz-delete")
	if err != nil {
		dir = os.TempDir()
	}
	config.C = &config.Config{}
	config.C.Processing.DocumentsDir = dir
	config.C.Processing.PreviewsDir = dir
}

func mayhemit(bytes []byte) int {

	fuzzConsumer := fuzz.NewConsumer(bytes)

	fuzzId, err := fuzzConsumer.GetString()
	if err != nil {
		return 0
	}

	process.DeleteDocument(fuzzId)
	return 0
}

func Fuzz(data []byte) int {
	_ = mayhemit(data)
	return 0
}
