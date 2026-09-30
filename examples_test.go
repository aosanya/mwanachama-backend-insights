package mwanachamainsights_test

import (
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
	"github.com/aosanya/mwanachama-backend-shared/spec"
)

func shippedSpecs(t *testing.T) []string {
	t.Helper()
	out := []string{"insights.wakala.json"}
	examples, err := filepath.Glob("spec/examples/*.json")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	return append(out, examples...)
}

func TestEveryExampleFitsTheTypes(t *testing.T) {
	for _, path := range shippedSpecs(t) {
		t.Run(path, func(t *testing.T) {
			s, err := mwanachamainsights.LoadSpec(path)
			if err != nil {
				t.Fatalf("LoadSpec(%s): %v", path, err)
			}
			db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			if err != nil {
				t.Fatalf("gorm.Open: %v", err)
			}
			if err := mwanachamainsights.Provision(db, s); err != nil {
				t.Fatalf("Provision(%s): %v", path, err)
			}
			if _, err := mwanachamainsights.NewInsightManager(db, s); err != nil {
				t.Fatalf("NewInsightManager(%s): %v", path, err)
			}
		})
	}
}

func TestTwoDomainsCoexist(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}

	tables := map[string]bool{}
	for _, path := range shippedSpecs(t) {
		s, err := mwanachamainsights.LoadSpec(path)
		if err != nil {
			t.Fatalf("LoadSpec(%s): %v", path, err)
		}
		if err := mwanachamainsights.Provision(db, s); err != nil {
			t.Fatalf("Provision(%s): %v", path, err)
		}
		for _, o := range s.Objects {
			name := s.TableFor(o)
			if tables[name] {
				t.Fatalf("%s: %s collides with a table another domain already declared", path, s.RawNameFor(o))
			}
			tables[name] = true
		}
	}
	if len(tables) < 4 {
		t.Fatalf("expected at least two domains of two objects each, got %d tables", len(tables))
	}
}

func TestRequiredFieldsHaveNoDefault(t *testing.T) {
	b, err := mwanachamainsights.Blueprint()
	if err != nil {
		t.Fatalf("blueprint: %v", err)
	}
	for _, o := range b.Objects {
		for _, f := range o.Fields {
			if f.Required && f.Default != "" {
				t.Errorf("%s.%s is both required and defaulted to %q — the default is exactly what would let an omitted value pass unnoticed",
					o.Role, f.Name, f.Default)
			}
		}
	}
}

func TestEveryDeclaredFieldReachesAColumn(t *testing.T) {
	s, err := mwanachamainsights.LoadSpec("insights.wakala.json")
	if err != nil {
		t.Fatalf("LoadSpec: %v", err)
	}
	for _, role := range []string{mwanachamainsights.RoleInsight, mwanachamainsights.RoleNote} {
		o, ok := s.ByRole(role)
		if !ok {
			t.Fatalf("the spec fills no object for the role %q", role)
		}
		if len(o.Fields) == 0 {
			t.Errorf("role %q arrived with no fields — a domain spec loaded through spec.Load rather than the blueprint does this", role)
		}
	}
}

func TestTableNamesAreDerivedNotLiteral(t *testing.T) {
	s, err := mwanachamainsights.LoadSpec("insights.wakala.json")
	if err != nil {
		t.Fatalf("LoadSpec: %v", err)
	}
	o, ok := s.ByRole(mwanachamainsights.RoleInsight)
	if !ok {
		t.Fatal("no object fills the insight role")
	}
	if got, want := s.RawNameFor(o), "insights_"+spec.DefaultMount+"_insights"; got != want {
		t.Fatalf("RawNameFor = %q, want %q", got, want)
	}
}

func TestSpecForCarriesTheDeclaredFields(t *testing.T) {
	s, err := mwanachamainsights.SpecFor("app")
	if err != nil {
		t.Fatalf("SpecFor: %v", err)
	}
	if s.Instance != "app" {
		t.Fatalf("Instance = %q, want %q", s.Instance, "app")
	}
	o, ok := s.ByRole(mwanachamainsights.RoleInsight)
	if !ok {
		t.Fatal("no object fills the insight role")
	}
	if len(o.Fields) == 0 {
		t.Fatal("the insight role arrived with no fields — SpecFor must go through the blueprint, not spec.Parse")
	}
	if o.Table != "insights" {
		t.Fatalf("Table = %q, want %q — the legacy adoption derives app_insights from this", o.Table, "insights")
	}
}
