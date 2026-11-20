package core

import (
	"os"
	"testing"
)

func TestParseSchemaExample(t *testing.T) {
	data, err := os.ReadFile("../examples/profile/schema.json")
	if err != nil {
		t.Fatalf("failed to read test data: %s", err)
	}
	s, err := ParseSchema(data)
	if err != nil {
		t.Fatalf("failed to parse schema: %s", err)
	}

	if len(s.Blocks) != 4 {
		t.Fatalf("blocks: got %d, want 4", len(s.Blocks))
	}
}

func TestParseSchemaRejects(t *testing.T) {
	cases := map[string]string{
		"unsupported format": `{"tessera":2,"name":"a","version":1}`,
		"missing name":       `{"tessera":1,"version":1}`,
		"version zero":       `{"tessera":1,"name":"a"}`,
		"unknown key": `{"tessera":1,"name":"a","version":1,
			"blocks":{"x":{"props":{"p":{"type":"string","maxLenght":5}}}}}`,
		"unknown block in allow": `{"tessera":1,"name":"a","version":1,
			"root":{"slots":{"content":{"allow":{"txt":{}}}}}}`,
		"unknown prop type": `{"tessera":1,"name":"a","version":1,
			"blocks":{"x":{"props":{"p":{"type":"color"}}}}}`,
		"url without protocols": `{"tessera":1,"name":"a","version":1,
			"blocks":{"x":{"props":{"p":{"type":"url"}}}}}`,
		"enum without options": `{"tessera":1,"name":"a","version":1,
			"blocks":{"x":{"props":{"p":{"type":"enum"}}}}}`,
		"object without fields": `{"tessera":1,"name":"a","version":1,
			"blocks":{"x":{"props":{"p":{"type":"object"}}}}}`,
		"list without of": `{"tessera":1,"name":"a","version":1,
			"blocks":{"x":{"props":{"p":{"type":"list"}}}}}`,
		"nested url without protocols": `{"tessera":1,"name":"a","version":1,
			"blocks":{"x":{"props":{"p":{"type":"list","of":
				{"type":"object","fields":{"h":{"type":"url"}}}}}}}}`,
	}

	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseSchema([]byte(data)); err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	}
}
