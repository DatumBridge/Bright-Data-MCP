package brightdata

import "testing"

func TestNormalizeBatchQueriesAcceptsStrings(t *testing.T) {
	got, err := NormalizeBatchQueries([]interface{}{
		"Vietnam GDP growth 2024",
		map[string]interface{}{"query": "site:worldbank.org Vietnam GDP", "engine": "bing"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Query != "Vietnam GDP growth 2024" || got[0].Engine != "" {
		t.Fatalf("queries %+v", got)
	}
	if got[1].Engine != "bing" {
		t.Fatalf("engine %q", got[1].Engine)
	}
}

func TestNormalizeBatchQueriesRejectsEmpty(t *testing.T) {
	if _, err := NormalizeBatchQueries([]interface{}{map[string]interface{}{"engine": "google"}}); err == nil {
		t.Fatal("expected query is required")
	}
}
