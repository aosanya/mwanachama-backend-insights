package mwanachamainsights_test

import (
	"sort"
	"testing"

	mwanachamainsights "github.com/aosanya/mwanachama-backend-insights"
	"github.com/aosanya/mwanachama-backend-insights/models"
	"github.com/aosanya/mwanachama-backend-shared/spec"
)

func TestVocabularyMatchesTheBlueprint(t *testing.T) {
	b, err := mwanachamainsights.Blueprint()
	if err != nil {
		t.Fatalf("blueprint: %v", err)
	}

	inGo := []string{models.SourceAuto, models.SourceUser}
	sort.Strings(inGo)

	for _, role := range []string{mwanachamainsights.RoleInsight, mwanachamainsights.RoleNote} {
		o, ok := b.Object(role)
		if !ok {
			t.Fatalf("the blueprint declares no %q", role)
		}
		var declared []string
		for _, f := range o.Fields {
			if f.Name == "source" {
				if f.Type != spec.TypeEnum {
					t.Fatalf("%s.source is %q, not an enum — the constants in models promise a closed set", role, f.Type)
				}
				declared = append(declared, f.Values...)
			}
		}
		sort.Strings(declared)

		if len(declared) != len(inGo) {
			t.Fatalf("%s.source declares %v, models compares against %v", role, declared, inGo)
		}
		for i := range inGo {
			if declared[i] != inGo[i] {
				t.Errorf("%s.source: blueprint has %q where models has %q", role, declared[i], inGo[i])
			}
		}
	}
}

func TestSourceDefaultsToAutoInTheBlueprintNotInGo(t *testing.T) {
	b, err := mwanachamainsights.Blueprint()
	if err != nil {
		t.Fatalf("blueprint: %v", err)
	}
	for _, role := range []string{mwanachamainsights.RoleInsight, mwanachamainsights.RoleNote} {
		o, _ := b.Object(role)
		for _, f := range o.Fields {
			if f.Name != "source" {
				continue
			}
			if f.Default != models.SourceAuto {
				t.Errorf("%s.source defaults to %q in the blueprint, but models.SourceAuto is %q — a create that omits it would land on the wrong value",
					role, f.Default, models.SourceAuto)
			}
		}
	}
}
