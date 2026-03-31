package core

import "testing"

func TestSanitizeBookTitleStripsKnownNoise(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "黎明之剑 (远瞳) (Z-Library)", want: "黎明之剑 (远瞳)"},
		{in: "末世_我的关键词比别人多一个_棉衣卫_Z_Library", want: "末世_我的关键词比别人多一个_棉衣卫"},
		{in: "鲁迅全集（全20册） [Z-Library]", want: "鲁迅全集（全20册）"},
		{in: "神墓（辰东）", want: "神墓（辰东）"},
	}

	for _, tc := range tests {
		if got := sanitizeBookTitle(tc.in); got != tc.want {
			t.Fatalf("sanitizeBookTitle(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
