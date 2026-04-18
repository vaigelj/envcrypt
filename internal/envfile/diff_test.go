package envfile

import (
	"testing"
)

func TestCompareAdded(t *testing.T) {
	base := map[string]string{"A": "1"}
	updated := map[string]string{"A": "1", "B": "2"}
	changes := Compare(base, updated)
	if len(changes) != 1 || changes[0].Type != ChangeAdded || changes[0].Key != "B" {
		t.Fatalf("expected one added change for B, got %+v", changes)
	}
}

func TestCompareRemoved(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	updated := map[string]string{"A": "1"}
	changes := Compare(base, updated)
	if len(changes) != 1 || changes[0].Type != ChangeRemoved || changes[0].Key != "B" {
		t.Fatalf("expected one removed change for B, got %+v", changes)
	}
}

func TestCompareUpdated(t *testing.T) {
	base := map[string]string{"A": "old"}
	updated := map[string]string{"A": "new"}
	changes := Compare(base, updated)
	if len(changes) != 1 || changes[0].Type != ChangeUpdated {
		t.Fatalf("expected one updated change, got %+v", changes)
	}
	if changes[0].OldVal != "old" || changes[0].NewVal != "new" {
		t.Errorf("unexpected values: %+v", changes[0])
	}
}

func TestCompareNoChanges(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	changes := Compare(base, base)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %+v", changes)
	}
}

func TestSummary(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded},
		{Type: ChangeAdded},
		{Type: ChangeUpdated},
		{Type: ChangeRemoved},
	}
	a, u, r := Summary(changes)
	if a != 2 || u != 1 || r != 1 {
		t.Errorf("unexpected summary: added=%d updated=%d removed=%d", a, u, r)
	}
}
