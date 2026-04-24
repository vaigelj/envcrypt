package envfile

import (
	"testing"
)

func TestFlattenDefaultSeparator(t *testing.T) {
	entries := []Entry{
		{Key: "db.host", Value: "localhost"},
		{Key: "db.port", Value: "5432"},
	}
	result := Flatten(entries, FlattenOptions{})
	if result[0].Key != "db_host" {
		t.Errorf("expected db_host, got %s", result[0].Key)
	}
	if result[1].Key != "db_port" {
		t.Errorf("expected db_port, got %s", result[1].Key)
	}
}

func TestFlattenUppercase(t *testing.T) {
	entries := []Entry{
		{Key: "app.name", Value: "envcrypt"},
	}
	result := Flatten(entries, FlattenOptions{Uppercase: true})
	if result[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %s", result[0].Key)
	}
}

func TestFlattenPrefix(t *testing.T) {
	entries := []Entry{
		{Key: "host", Value: "localhost"},
	}
	result := Flatten(entries, FlattenOptions{Prefix: "MY_"})
	if result[0].Key != "MY_host" {
		t.Errorf("expected MY_host, got %s", result[0].Key)
	}
}

func TestFlattenSlashSeparator(t *testing.T) {
	entries := []Entry{
		{Key: "aws/region", Value: "us-east-1"},
	}
	result := Flatten(entries, FlattenOptions{Uppercase: true})
	if result[0].Key != "AWS_REGION" {
		t.Errorf("expected AWS_REGION, got %s", result[0].Key)
	}
}

func TestFlattenMapBasic(t *testing.T) {
	m := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
	}
	result := FlattenMap(m, FlattenOptions{Uppercase: true})
	keys := make(map[string]string)
	for _, e := range result {
		keys[e.Key] = e.Value
	}
	if keys["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", keys["DB_HOST"])
	}
	if keys["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %s", keys["DB_PORT"])
	}
}

func TestUnflattenToMap(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	}
	result := UnflattenToMap(entries, "_")
	db, ok := result["DB"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected nested map under DB")
	}
	if db["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %v", db["HOST"])
	}
}

func TestFormatFlattened(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
	}
	out := FormatFlattened(entries)
	if out != "FOO=bar\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
