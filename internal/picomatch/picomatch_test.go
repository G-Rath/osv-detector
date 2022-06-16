package picomatch_test

import (
	"osv-detector/internal/picomatch"
	"testing"
)

func TestCompiledMatchers_Matches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		patterns []string
		path     string
		want     bool
	}{
		{name: "", patterns: []string{}, path: "path/to/my/file", want: false},
		{name: "", patterns: []string{"path/to/my/file"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"path/to/my/file"}, path: "path/to/my/file/", want: false},
		{name: "", patterns: []string{"path/to/my/file"}, path: "path/to/my", want: false},
		{name: "", patterns: []string{"path/to/my/file"}, path: "path/to/my/", want: false},
		{name: "", patterns: []string{"*"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"**"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"!*"}, path: "path/to/my/file", want: false},
		{name: "", patterns: []string{"!**"}, path: "path/to/my/file", want: false},
		{name: "", patterns: []string{"path/to/my/*"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"*/to/my/file"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"*/to/my/*"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"path/to/*/file"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"path/**/file"}, path: "path/to/my/file", want: true},
		{
			name: "",
			patterns: []string{
				"!**/*",
				"path/to/my/file",
			},
			path: "path/to/my/file",
			want: true,
		},
		{
			name: "",
			patterns: []string{
				"!**/*",
				"path/to/my/file",
				"!path/to/my/file",
			},
			path: "path/to/my/file",
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cms := picomatch.FromPatterns(tt.patterns)

			if got := cms.Matches(tt.path); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func expectNoDuplicates(t *testing.T, slice []string) {
	t.Helper()

	set := make(map[string]bool)

	for _, s := range slice {
		if _, ok := set[s]; ok {
			t.Errorf("Value \"%s\" occurred more than once in slice", s)
		}

		set[s] = true
	}
}

func TestCompiledMatchers_Matches_Advanced(t *testing.T) {
	t.Parallel()

	fsTree := []string{
		".github/workflows/create_staging_branch.yaml",
		// "advisories/",
		"advisories/github-reviewed/2021/01/GHSA-29v9-2fpx-j5g9/GHSA-29v9-2fpx-j5g9.json",
		"advisories/github-reviewed/2021/01/GHSA-2ccx-2gf3-8xvv/GHSA-2ccx-2gf3-8xvv.json",
		"advisories/github-reviewed/2021/04/GHSA-22cm-3qf2-2wc7/GHSA-22cm-3qf2-2wc7.json",
		"advisories/github-reviewed/2021/04/GHSA-22wc-c9wj-6q2v/GHSA-22wc-c9wj-6q2v.json",
		"advisories/github-reviewed/2022/02/GHSA-227w-wv4j-67h4/GHSA-227w-wv4j-67h4.json",
		"advisories/github-reviewed/2022/02/GHSA-23hm-7w47-xw72/GHSA-23hm-7w47-xw72.json",
		"advisories/github-reviewed/2022/04/GHSA-29f8-q7mf-7cqj/GHSA-29f8-q7mf-7cqj.json",
		"advisories/github-reviewed/2022/04/GHSA-2cfc-865j-gm4w/GHSA-2cfc-865j-gm4w.json",
		"advisories/unreviewed/2021/04/GHSA-m5pg-8h68-j225/GHSA-m5pg-8h68-j225.json",
		"advisories/unreviewed/2021/05/GHSA-2jx2-76rc-2v7v/GHSA-2jx2-76rc-2v7v.json",
		"advisories/unreviewed/2021/05/GHSA-2r5r-x58v-cx3w/GHSA-2r5r-x58v-cx3w.json",
		"advisories/unreviewed/2022/01/GHSA-2222-76gx-28mm/GHSA-2222-76gx-28mm.json",
		"advisories/unreviewed/2022/01/GHSA-226r-cf9r-2r9j/GHSA-226r-cf9r-2r9j.json",
		"advisories/unreviewed/2022/02/GHSA-224h-mqw6-642p/GHSA-224h-mqw6-642p.json",
		"advisories/unreviewed/2022/03/GHSA-2369-w664-2vw7/GHSA-2369-w664-2vw7.json",
		"advisories/unreviewed/2022/04/GHSA-222r-4v3h-874c/GHSA-222r-4v3h-874c.json",
		"advisories/unreviewed/2022/05/GHSA-2224-j5w9-6w4m/GHSA-2224-j5w9-6w4m.json",
		"README.md",
		"tsconfig.json",
	}

	/*
	  "!**": len(matches) == 0
	   "**": len(matches) == len(fsTree)
	   "*.md": len(matches) == 1 ("README.md")
	   "*.json": len(matches) == 1 ("tsconfig.json")
	   "**.md": len(matches) == 1 ("README.md")
	   "**.json": len(matches) == <a bunch>
	   "advisories/**.json": len(matches) == <a bunch> - 1 (not "tsconfig.json")
	   "advisories/github-reviewed/**.json": len(matches) == <a bunch> - 1 (not "tsconfig.json")


	*/

	type testCase struct {
		patterns
	}

	tests := []struct {
		name     string
		patterns []string
		path     string
		want     bool
	}{
		{name: "", patterns: []string{}, path: "path/to/my/file", want: false},
		{name: "", patterns: []string{"path/to/my/file"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"path/to/my/file"}, path: "path/to/my/file/", want: false},
		{name: "", patterns: []string{"path/to/my/file"}, path: "path/to/my", want: false},
		{name: "", patterns: []string{"path/to/my/file"}, path: "path/to/my/", want: false},
		{name: "", patterns: []string{"*"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"**"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"!*"}, path: "path/to/my/file", want: false},
		{name: "", patterns: []string{"!**"}, path: "path/to/my/file", want: false},
		{name: "", patterns: []string{"path/to/my/*"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"*/to/my/file"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"*/to/my/*"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"path/to/*/file"}, path: "path/to/my/file", want: true},
		{name: "", patterns: []string{"path/**/file"}, path: "path/to/my/file", want: true},
		{
			name: "",
			patterns: []string{
				"!**/*",
				"path/to/my/file",
			},
			path: "path/to/my/file",
			want: true,
		},
		{
			name: "",
			patterns: []string{
				"!**/*",
				"path/to/my/file",
				"!path/to/my/file",
			},
			path: "path/to/my/file",
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cms := picomatch.FromPatterns(tt.patterns)

			if got := cms.Matches(tt.path); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
