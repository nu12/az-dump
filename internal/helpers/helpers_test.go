package helpers

import "testing"

func TestGetResourceGroupNameFromFileName(t *testing.T) {
	tests := []struct {
		fileName string
		expected string
	}{
		{"mygroup.json", "mygroup"},
		{"another-resource-group.json", "another-resource-group"},
		{"file.with.dots.json", "file.with.dots"},
		{"noextension", "noextension"},
		{"endswithjson", "endswithjson"},
		{"", ""},
	}

	for _, tt := range tests {
		result := GetResourceGroupNameFromFileName(tt.fileName)
		if result != tt.expected {
			t.Errorf("GetResourceGroupNameFromFileName(%q) = %q; want %q", tt.fileName, result, tt.expected)
		}
	}
}
func TestComaListContains(t *testing.T) {
	tests := []struct {
		commaList string
		item      string
		expected  bool
	}{
		{"a,b,c", "a", true},
		{"a,b,c", "b", true},
		{"a,b,c", "c", true},
		{"a,b,c", "d", false},
		{"", "a", false},
		{"a", "a", true},
		{"a", "b", false},
		{",a,", "a", true},
		{"a,,b", "", true},
		{"a,b,c", "", false},
		{"a,b,c", "A", false}, // case-sensitive
	}

	for _, tt := range tests {
		result := ComaListContains(tt.commaList, tt.item)
		if result != tt.expected {
			t.Errorf("ComaListContains(%q, %q) = %v; want %v", tt.commaList, tt.item, result, tt.expected)
		}
	}
}
