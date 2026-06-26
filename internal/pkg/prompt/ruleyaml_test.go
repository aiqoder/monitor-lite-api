package prompt

import "testing"

func TestDefaultConfigFromDefaultGroup(t *testing.T) {
	cfg := DefaultConfig()
	if len(cfg.Groups) == 0 {
		t.Fatal("expected groups from defaultgroup")
	}
}

func TestConfigFromDefaultGroup(t *testing.T) {
	cfg, err := ConfigFromDefaultGroup()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Groups) == 0 {
		t.Fatal("expected groups")
	}
	if cfg.Groups[0].Name != "央视频道" {
		t.Fatalf("unexpected first group: %q", cfg.Groups[0].Name)
	}
}

func TestParseGroupsJSON(t *testing.T) {
	raw := `{"groups":[{"name":"央视频道","channels":["CCTV-1综合","CCTV-2财经"]}]}`
	cfg, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Groups) != 1 || len(cfg.Groups[0].Channels) != 2 {
		t.Fatalf("unexpected config: %+v", cfg.Groups)
	}
}
