package brightdata_test

import (
	"testing"

	"github.com/datumbridge/bright-data-mcp/internal/brightdata"
)

func TestLoadDatasetCatalog(t *testing.T) {
	catalog, err := brightdata.LoadDatasetCatalog()
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog) < 50 {
		t.Fatalf("expected 50 datasets, got %d", len(catalog))
	}
}
