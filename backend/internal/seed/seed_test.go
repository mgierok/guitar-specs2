package seed

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	input := `{
  "brands": [{"id":"b1","name":"Fender"}],
  "specs": [],
  "spec_options": [],
  "guitars": [],
  "guitar_spec_values": [],
  "guitar_media": []
}`
	reader := strings.NewReader(input)

	data, err := loadFromReader(reader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data.Brands) != 1 {
		t.Fatalf("expected 1 brand, got %d", len(data.Brands))
	}
	if data.Brands[0].Name != "Fender" {
		t.Fatalf("unexpected brand name: %s", data.Brands[0].Name)
	}
}
