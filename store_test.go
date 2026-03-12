package relationalpostgres

import (
	"strings"
	"testing"

	"github.com/BHRK-codelabs/corekit/configkit"
)

func TestOpenRequiresDatabaseURL(t *testing.T) {
	t.Parallel()

	if _, err := Open(""); err == nil {
		t.Fatal("expected database url error")
	}
}

func TestNewModuleRequiresDatabaseURL(t *testing.T) {
	t.Parallel()

	cfg := configkit.New()
	cfg.Database.URL = ""

	if _, err := NewModule(cfg); err == nil {
		t.Fatal("expected database url error")
	}
}

func TestNormalizeDatabaseURLAddsSimpleProtocol(t *testing.T) {
	t.Parallel()

	got := normalizeDatabaseURL("postgres://user:pass@localhost:5432/db?sslmode=require")
	if got == "postgres://user:pass@localhost:5432/db?sslmode=require" {
		t.Fatal("expected normalized url to add query exec mode")
	}
	if want := "default_query_exec_mode=simple_protocol"; !contains(got, want) {
		t.Fatalf("expected %q in %q", want, got)
	}
}

func contains(value, fragment string) bool {
	return strings.Contains(value, fragment)
}
