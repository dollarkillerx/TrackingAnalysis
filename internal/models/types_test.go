package models

import (
	"encoding/json"
	"testing"
)

func TestJSONMap_Value(t *testing.T) {
	m := JSONMap{"key": "value", "count": float64(42)}
	v, err := m.Value()
	if err != nil {
		t.Fatalf("Value(): %v", err)
	}
	bytes, ok := v.([]byte)
	if !ok {
		t.Fatalf("Value() returned %T, want []byte", v)
	}
	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(bytes, &parsed); err != nil {
		t.Fatalf("Value() produced invalid JSON: %v", err)
	}
	if parsed["key"] != "value" {
		t.Errorf("parsed[key] = %v, want %q", parsed["key"], "value")
	}
}

func TestJSONMap_Value_Nil(t *testing.T) {
	var m JSONMap // nil
	v, err := m.Value()
	if err != nil {
		t.Fatalf("Value(): %v", err)
	}
	if v != nil {
		t.Errorf("Value() = %v, want nil", v)
	}
}

func TestJSONMap_Scan(t *testing.T) {
	raw := []byte(`{"name":"test","active":true}`)
	var m JSONMap
	if err := m.Scan(raw); err != nil {
		t.Fatalf("Scan(): %v", err)
	}
	if m["name"] != "test" {
		t.Errorf("m[name] = %v, want %q", m["name"], "test")
	}
	if m["active"] != true {
		t.Errorf("m[active] = %v, want true", m["active"])
	}
}

func TestJSONMap_Scan_Nil(t *testing.T) {
	m := JSONMap{"existing": "data"}
	if err := m.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if m != nil {
		t.Errorf("Scan(nil) should set map to nil, got %v", m)
	}
}

func TestJSONMap_Scan_InvalidType(t *testing.T) {
	var m JSONMap
	err := m.Scan(12345) // not []byte
	if err == nil {
		t.Fatal("expected error for non-byte input")
	}
}
