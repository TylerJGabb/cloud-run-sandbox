package main

import (
	"testing"
)

func TestKeyValuesToJsonPayload(t *testing.T) {
	keyValue := []any{"key1", "value1", "key2", "value2"}
	result := KeyValuesToJsonPayload(keyValue...)
	// Output:
	// map[key1:value1 key2:value2]
	if len(result) != 2 {
		t.Errorf("Expected 2 items, but got %v", len(result))
	}
	if result["key1"] != "value1" {
		t.Errorf("Expected value1, but got %v", result["key1"])
	}
	if result["key2"] != "value2" {
		t.Errorf("Expected value2, but got %v", result["key2"])
	}
}

func TestKeyValuesToJsonPayloadWithEmpty(t *testing.T) {
	keyValue := []any{}
	result := KeyValuesToJsonPayload(keyValue...)
	// Output:
	// map[]
	if len(result) != 0 {
		t.Errorf("Expected 0 items, but got %v", len(result))
	}
}

func TestKeyValuesToJsonPayloadWithOneItem(t *testing.T) {
	keyValue := []any{"key1"}
	result := KeyValuesToJsonPayload(keyValue...)
	// Output:
	// map[key1:value1]
	if len(result) != 0 {
		t.Errorf("Expected 1 items, but got %v", len(result))
	}
}
