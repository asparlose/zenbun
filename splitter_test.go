package zenbun

import "testing"

func TestSplitter(t *testing.T) {
	tests := map[string][]string{
		"Hello, World!!":   []string{"Hello", "World"},
		"!!Hello, World!!": []string{"Hello", "World"},
		"こんにちは、世界！！":       []string{"こんにちは", "世界"},
		"！！こんにちは、世界！！":     []string{"こんにちは", "世界"},
	}
	for i, o := range tests {
		s := newSplitter(&i)
		for _, v := range o {
			if s.Next() != v {
				t.Fatalf("splitter failed: want \"%v\", got \"%v\"", v, s.Next())
			}
		}
		k := s.Next()
		if k != "" {
			t.Fatalf("splitter failed: want \"\", got \"%v\"", k)
		}
	}
}
