package defaultgroup

import (
	"testing"
)

func TestAll(t *testing.T) {
	groups, err := All()
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 42 {
		t.Fatalf("expected 42 groups, got %d", len(groups))
	}
	if groups[0].GroupName != "央视频道" {
		t.Fatalf("expected first group 央视频道, got %q", groups[0].GroupName)
	}
	if len(groups[0].TvNames) == 0 {
		t.Fatal("expected tv names in first group")
	}
}

func TestGroupMap(t *testing.T) {
	m, err := GroupMap()
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 42 {
		t.Fatalf("expected 42 groups, got %d", len(m))
	}
	if _, ok := m["央视频道"]; !ok {
		t.Fatal("missing 央视频道")
	}
}
