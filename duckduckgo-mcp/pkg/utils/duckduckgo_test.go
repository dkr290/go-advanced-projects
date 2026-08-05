package utils

import "testing"

func TestCleanDuckDuckGoRedirect(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "protocol relative redirect with uddg",
			in:   "//duckduckgo.com/l/?uddg=https%3A%2F%2Fgolang.org%2F&rut=abc",
			want: "https://golang.org/",
		},
		{
			name: "plain url passthrough",
			in:   "https://example.com/page",
			want: "https://example.com/page",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := cleanDuckDuckGoRedirect(tc.in)
			if got != tc.want {
				t.Fatalf("cleanDuckDuckGoRedirect(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
