package envfile

import (
	"testing"
)

func TestMergePreferBase(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	override := map[string]string{"B": "99", "C": "3"}
	out := Merge(base, override, PreferBase)
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %s", out["A"])
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2 (base wins), got %s", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("expected C=3 (new key), got %s", out["C"])
	}
}

func TestMergePreferOverride(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	override := map[string]string{"B": "99", "C": "3"}
	out := Merge(base, override, PreferOverride)
	if out["B"] != "99" {
		t.Errorf("expected B=99 (override wins), got %s", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("expected C=3, got %s", out["C"])
	}
}

func TestMergeEmptyBase(t *testing.T) {
	out := Merge(map[string]string{}, map[string]string{"X": "10"}, PreferBase)
	if out["X"] != "10" {
		t.Errorf("expected X=10, got %s", out["X"])
	}
}

func TestDiffChanged(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2", "C": "3"}
	b := map[string]string{"A": "1", "B": "99", "D": "4"}
	d := Diff(a, b)
	if _, ok := d["A"]; ok {
		t.Error("A is unchanged, should not appear in diff")
	}
	if d["B"] != "99" {
		t.Errorf("expected B=99 in diff, got %s", d["B"])
	}
	if d["C"] != "" {
		t.Errorf("expected C='' (removed), got %s", d["C"])
	}
	if d["D"] != "4" {
		t.Errorf("expected D=4 (added), got %s", d["D"])
	}
}

func TestDiffIdentical(t *testing.T) {
	a := map[string]string{"A": "1"}
	if len(Diff(a, a)) != 0 {
		t.Error("expected empty diff for identical maps")
	}
}
