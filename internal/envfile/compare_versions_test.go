package envfile

import (
	"strings"
	"testing"
)

func TestCompareVersionsAdded(t *testing.T) {
	from := Version{Name: "v1", Vars: map[string]string{"A": "1"}}
	to := Version{Name: "v2", Vars: map[string]string{"A": "1", "B": "2"}}
	d := CompareVersions(from, to)
	if len(d.Changes) != 1 || d.Changes[0].Op != "added" || d.Changes[0].Key != "B" {
		t.Fatalf("expected added B, got %+v", d.Changes)
	}
}

func TestCompareVersionsRemoved(t *testing.T) {
	from := Version{Name: "v1", Vars: map[string]string{"A": "1", "B": "2"}}
	to := Version{Name: "v2", Vars: map[string]string{"A": "1"}}
	d := CompareVersions(from, to)
	if len(d.Changes) != 1 || d.Changes[0].Op != "removed" || d.Changes[0].Key != "B" {
		t.Fatalf("expected removed B, got %+v", d.Changes)
	}
}

func TestCompareVersionsUpdated(t *testing.T) {
	from := Version{Name: "v1", Vars: map[string]string{"A": "old"}}
	to := Version{Name: "v2", Vars: map[string]string{"A": "new"}}
	d := CompareVersions(from, to)
	if len(d.Changes) != 1 || d.Changes[0].Op != "updated" {
		t.Fatalf("expected updated A, got %+v", d.Changes)
	}
	if d.Changes[0].OldVal != "old" || d.Changes[0].NewVal != "new" {
		t.Fatalf("unexpected values: %+v", d.Changes[0])
	}
}

func TestCompareVersionsNoChanges(t *testing.T) {
	from := Version{Name: "v1", Vars: map[string]string{"A": "1"}}
	to := Version{Name: "v2", Vars: map[string]string{"A": "1"}}
	d := CompareVersions(from, to)
	if len(d.Changes) != 0 {
		t.Fatalf("expected no changes, got %+v", d.Changes)
	}
}

func TestFormatVersionDiff(t *testing.T) {
	from := Version{Name: "v1", Vars: map[string]string{"A": "1"}}
	to := Version{Name: "v2", Vars: map[string]string{"A": "1", "B": "2"}}
	d := CompareVersions(from, to)
	out := FormatVersionDiff(d)
	if !strings.Contains(out, "+ B=2") {
		t.Fatalf("expected added B in output, got: %s", out)
	}
}

func TestFormatVersionDiffNoChanges(t *testing.T) {
	from := Version{Name: "v1", Vars: map[string]string{}}
	to := Version{Name: "v2", Vars: map[string]string{}}
	d := CompareVersions(from, to)
	out := FormatVersionDiff(d)
	if !strings.Contains(out, "No changes") {
		t.Fatalf("expected no-changes message, got: %s", out)
	}
}
