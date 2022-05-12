package lockfile_test

import (
	"errors"
	"io/ioutil"
	"osv-detector/pkg/lockfile"
	"reflect"
	"strings"
	"testing"
)

func expectNumberOfParsersCalled(t *testing.T, numberOfParsersCalled int) {
	t.Helper()

	directories, err := ioutil.ReadDir(".")

	if err != nil {
		t.Fatalf("unable to read current directory: ")
	}

	count := 0

	for _, directory := range directories {
		if strings.HasPrefix(directory.Name(), "parse-") &&
			!strings.HasSuffix(directory.Name(), "_test.go") {
			count++
		}
	}

	if numberOfParsersCalled != count {
		t.Errorf(
			"Expected %d parsers to have been called, but had %d",
			count,
			numberOfParsersCalled,
		)
	}
}

func TestFindParser(t *testing.T) {
	t.Parallel()

	lockfiles := []string{
		"cargo.lock",
		"package-lock.json",
		"yarn.lock",
		"pnpm-lock.yaml",
		"composer.lock",
		"Gemfile.lock",
		"go.mod",
		"pom.xml",
		"requirements.txt",
	}

	for _, file := range lockfiles {
		parser, parsedAs := lockfile.FindParser("/path/to/my/"+file, "")

		if parser == nil {
			t.Errorf("Expected a parser to be found for %s but did not", file)
		}

		if file != parsedAs {
			t.Errorf("Expected parsedAs to be %s but got %s instead", file, parsedAs)
		}
	}
}

func TestFindParser_ExplicitParseAs(t *testing.T) {
	t.Parallel()

	parser, parsedAs := lockfile.FindParser("/path/to/my/package-lock.json", "composer.lock")

	if parser == nil {
		t.Errorf("Expected a parser to be found for package-lock.json (overridden as composer.json) but did not")
	}

	if parsedAs != "composer.lock" {
		t.Errorf("Expected parsedAs to be composer.lock but got %s instead", parsedAs)
	}
}

func TestParse_FindsExpectedParsers(t *testing.T) {
	t.Parallel()

	lockfiles := []string{
		"cargo.lock",
		"package-lock.json",
		"yarn.lock",
		"pnpm-lock.yaml",
		"composer.lock",
		"Gemfile.lock",
		"go.mod",
		"pom.xml",
		"requirements.txt",
	}

	count := 0

	for _, file := range lockfiles {
		_, err := lockfile.Parse("/path/to/my/"+file, "")

		if errors.Is(err, lockfile.ErrParserNotFound) {
			t.Errorf("No parser was found for %s", file)
		}

		count++
	}

	expectNumberOfParsersCalled(t, count)
}

func TestParse_ParserNotFound(t *testing.T) {
	t.Parallel()

	_, err := lockfile.Parse("/path/to/my/", "")

	if err == nil {
		t.Errorf("Expected to get an error but did not")
	}

	if !errors.Is(err, lockfile.ErrParserNotFound) {
		t.Errorf("Did not get the expected ErrParserNotFound error - got %v instead", err)
	}
}

func TestListParsers(t *testing.T) {
	t.Parallel()

	parsers := lockfile.ListParsers()

	if first := parsers[0]; first != "cargo.lock" {
		t.Errorf("Expected first element to be cargo.lock, but got %s", first)
	}

	if last := parsers[len(parsers)-1]; last != "yarn.lock" {
		t.Errorf("Expected last element to be requirements.txt, but got %s", last)
	}
}

func TestLockfile_ToString(t *testing.T) {
	t.Parallel()

	expected := strings.Join([]string{
		"  crates.io: addr2line@0.15.2",
		"  npm: @typescript-eslint/types@5.13.0",
		"  crates.io: wasi@0.10.2+wasi-snapshot-preview1",
		"  Packagist: sentry/sdk@2.0.4",
	}, "\n")

	lockf := lockfile.Lockfile{
		Packages: []lockfile.PackageDetails{
			{
				Name:      "addr2line",
				Version:   "0.15.2",
				Ecosystem: lockfile.CargoEcosystem,
			},
			{
				Name:      "@typescript-eslint/types",
				Version:   "5.13.0",
				Ecosystem: lockfile.PnpmEcosystem,
			},
			{
				Name:      "wasi",
				Version:   "0.10.2+wasi-snapshot-preview1",
				Ecosystem: lockfile.CargoEcosystem,
			},
			{
				Name:      "sentry/sdk",
				Version:   "2.0.4",
				Ecosystem: lockfile.ComposerEcosystem,
			},
		},
	}

	if actual := lockf.ToString(); expected != actual {
		t.Errorf("\nExpected:\n%s\nActual:\n%s", expected, actual)
	}
}

func TestFromCSVRows(t *testing.T) {
	t.Parallel()

	type args struct {
		filePath string
		parseAs  string
		rows     []string
	}
	tests := []struct {
		name    string
		args    args
		want    lockfile.Lockfile
		wantErr bool
	}{
		{
			name: "",
			args: args{
				filePath: "-",
				parseAs:  "-",
				rows:     nil,
			},
			want: lockfile.Lockfile{
				FilePath: "-",
				ParsedAs: "-",
				Packages: nil,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				filePath: "-",
				parseAs:  "csv-row",
				rows: []string{
					"crates.io,addr2line,0.15.2",
					"npm,@typescript-eslint/types,5.13.0",
					"crates.io,wasi,0.10.2+wasi-snapshot-preview1",
					"Packagist,sentry/sdk,2.0.4",
				},
			},
			want: lockfile.Lockfile{
				FilePath: "-",
				ParsedAs: "csv-row",
				Packages: []lockfile.PackageDetails{
					{
						Name:      "@typescript-eslint/types",
						Version:   "5.13.0",
						Ecosystem: lockfile.PnpmEcosystem,
					},
					{
						Name:      "addr2line",
						Version:   "0.15.2",
						Ecosystem: lockfile.CargoEcosystem,
					},
					{
						Name:      "sentry/sdk",
						Version:   "2.0.4",
						Ecosystem: lockfile.ComposerEcosystem,
					},
					{
						Name:      "wasi",
						Version:   "0.10.2+wasi-snapshot-preview1",
						Ecosystem: lockfile.CargoEcosystem,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				filePath: "-",
				parseAs:  "-",
				rows: []string{
					"NuGet,Yarp.ReverseProxy,",
					"npm,@typescript-eslint/types,5.13.0",
				},
			},
			want: lockfile.Lockfile{
				FilePath: "-",
				ParsedAs: "-",
				Packages: []lockfile.PackageDetails{
					{
						Name:      "@typescript-eslint/types",
						Version:   "5.13.0",
						Ecosystem: lockfile.PnpmEcosystem,
					},
					{
						Name:      "Yarp.ReverseProxy",
						Version:   "",
						Ecosystem: "NuGet",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := lockfile.FromCSVRows(tt.args.filePath, tt.args.parseAs, tt.args.rows)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromCSVRows() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromCSVRows() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromCSVRows_Errors(t *testing.T) {
	t.Parallel()

	type args struct {
		filePath string
		parseAs  string
		rows     []string
	}
	tests := []struct {
		name       string
		args       args
		wantErrMsg string
	}{
		{
			name: "",
			args: args{
				filePath: "",
				parseAs:  "",
				rows:     []string{"one,,"},
			},
			wantErrMsg: "row 1: field 2 is empty (must be the name of a package)",
		},
		{
			name: "",
			args: args{
				filePath: "",
				parseAs:  "",
				rows: []string{
					"crates.io,addr2line,",
					",,",
				},
			},
			wantErrMsg: "row 2: field 1 is empty (must be the name of an ecosystem)",
		},
		{
			name: "",
			args: args{
				filePath: "",
				parseAs:  "",
				rows: []string{
					"crates.io,addr2line,",
					",,,",
				},
			},
			wantErrMsg: "record on line 2: wrong number of fields",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := lockfile.FromCSVRows(tt.args.filePath, tt.args.parseAs, tt.args.rows)

			if err == nil {
				t.Errorf("FromCSVRows() did not error")

				return
			}

			if !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("FromCSVRows() error = \"%v\", wanted \"%s\"", err, tt.wantErrMsg)
			}
		})
	}
}

func TestFromCSVFile(t *testing.T) {
	t.Parallel()

	type args struct {
		pathToCSV string
		parseAs   string
	}
	tests := []struct {
		name    string
		args    args
		want    lockfile.Lockfile
		wantErr bool
	}{
		{
			name: "",
			args: args{
				pathToCSV: "fixtures/csv/does-not-exist",
				parseAs:   "csv-file",
			},
			want:    lockfile.Lockfile{},
			wantErr: true,
		},
		{
			name: "",
			args: args{
				pathToCSV: "fixtures/csv/empty.csv",
				parseAs:   "csv-file",
			},
			want: lockfile.Lockfile{
				FilePath: "fixtures/csv/empty.csv",
				ParsedAs: "csv-file",
				Packages: nil,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				pathToCSV: "fixtures/csv/multiple-rows.csv",
				parseAs:   "csv-file",
			},
			want: lockfile.Lockfile{
				FilePath: "fixtures/csv/multiple-rows.csv",
				ParsedAs: "csv-file",
				Packages: []lockfile.PackageDetails{
					{
						Name:      "@typescript-eslint/types",
						Version:   "4.9.0",
						Ecosystem: lockfile.PnpmEcosystem,
					},
					{
						Name:      "@typescript-eslint/types",
						Version:   "5.13.0",
						Ecosystem: lockfile.PnpmEcosystem,
					},
					{
						Name:      "addr2line",
						Version:   "0.15.2",
						Ecosystem: lockfile.CargoEcosystem,
					},
					{
						Name:      "sentry/sdk",
						Version:   "2.0.4",
						Ecosystem: lockfile.ComposerEcosystem,
					},
					{
						Name:      "wasi",
						Version:   "0.10.2+wasi-snapshot-preview1",
						Ecosystem: lockfile.CargoEcosystem,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				pathToCSV: "fixtures/csv/with-extra-columns.csv",
				parseAs:   "csv-file",
			},
			want: lockfile.Lockfile{
				FilePath: "fixtures/csv/with-extra-columns.csv",
				ParsedAs: "csv-file",
				Packages: []lockfile.PackageDetails{
					{
						Name:      "@typescript-eslint/types",
						Version:   "5.13.0",
						Ecosystem: lockfile.PnpmEcosystem,
					},
					{
						Name:      "addr2line",
						Version:   "0.15.2",
						Ecosystem: lockfile.CargoEcosystem,
					},
					{
						Name:      "sentry/sdk",
						Version:   "2.0.4",
						Ecosystem: lockfile.ComposerEcosystem,
					},
					{
						Name:      "wasi",
						Version:   "0.10.2+wasi-snapshot-preview1",
						Ecosystem: lockfile.CargoEcosystem,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				pathToCSV: "fixtures/csv/one-row.csv",
				parseAs:   "-",
			},
			want: lockfile.Lockfile{
				FilePath: "fixtures/csv/one-row.csv",
				ParsedAs: "-",
				Packages: []lockfile.PackageDetails{
					{
						Name:      "@typescript-eslint/types",
						Version:   "5.13.0",
						Ecosystem: lockfile.PnpmEcosystem,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				pathToCSV: "fixtures/csv/two-rows.csv",
				parseAs:   "-",
			},
			want: lockfile.Lockfile{
				FilePath: "fixtures/csv/two-rows.csv",
				ParsedAs: "-",
				Packages: []lockfile.PackageDetails{
					{
						Name:      "@typescript-eslint/types",
						Version:   "5.13.0",
						Ecosystem: lockfile.PnpmEcosystem,
					},
					{
						Name:      "Yarp.ReverseProxy",
						Version:   "",
						Ecosystem: "NuGet",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				pathToCSV: "fixtures/csv/with-headers.csv",
				parseAs:   "-",
			},
			want: lockfile.Lockfile{
				FilePath: "fixtures/csv/with-headers.csv",
				ParsedAs: "-",
				Packages: []lockfile.PackageDetails{
					{
						Name:      "@typescript-eslint/types",
						Version:   "5.13.0",
						Ecosystem: lockfile.PnpmEcosystem,
					},
					{
						Name:      "Package",
						Version:   "Version",
						Ecosystem: "Ecosystem",
					},
					{
						Name:      "sentry/sdk",
						Version:   "2.0.4",
						Ecosystem: lockfile.ComposerEcosystem,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := lockfile.FromCSVFile(tt.args.pathToCSV, tt.args.parseAs)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromCSVFile() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromCSVFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackages_Ecosystems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ps   lockfile.Packages
		want []lockfile.Ecosystem
	}{
		{name: "", ps: lockfile.Packages{}, want: []lockfile.Ecosystem{}},
		{
			name: "",
			ps: lockfile.Packages{
				{
					Name:      "addr2line",
					Version:   "0.15.2",
					Ecosystem: lockfile.CargoEcosystem,
				},
			},
			want: []lockfile.Ecosystem{
				lockfile.CargoEcosystem,
			},
		},
		{
			name: "",
			ps: lockfile.Packages{
				{
					Name:      "addr2line",
					Version:   "0.15.2",
					Ecosystem: lockfile.CargoEcosystem,
				},
				{
					Name:      "wasi",
					Version:   "0.10.2+wasi-snapshot-preview1",
					Ecosystem: lockfile.CargoEcosystem,
				},
			},
			want: []lockfile.Ecosystem{
				lockfile.CargoEcosystem,
			},
		},
		{
			name: "",
			ps: lockfile.Packages{
				{
					Name:      "addr2line",
					Version:   "0.15.2",
					Ecosystem: lockfile.CargoEcosystem,
				},
				{
					Name:      "@typescript-eslint/types",
					Version:   "5.13.0",
					Ecosystem: lockfile.PnpmEcosystem,
				},
				{
					Name:      "wasi",
					Version:   "0.10.2+wasi-snapshot-preview1",
					Ecosystem: lockfile.CargoEcosystem,
				},
				{
					Name:      "sentry/sdk",
					Version:   "2.0.4",
					Ecosystem: lockfile.ComposerEcosystem,
				},
			},
			want: []lockfile.Ecosystem{
				lockfile.ComposerEcosystem,
				lockfile.CargoEcosystem,
				lockfile.PnpmEcosystem,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.ps.Ecosystems(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ecosystems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromCSVFile_Errors(t *testing.T) {
	t.Parallel()

	type args struct {
		pathToCSV string
		parseAs   string
	}
	tests := []struct {
		name       string
		args       args
		wantErrMsg string
	}{
		{
			name: "",
			args: args{
				pathToCSV: "fixtures/csv/does-not-exist",
				parseAs:   "csv-file",
			},
			wantErrMsg: "could not read fixtures/csv/does-not-exist",
		},
		{
			name: "",
			args: args{
				pathToCSV: "fixtures/csv/not-a-csv.xml",
				parseAs:   "csv-file",
			},
			wantErrMsg: "fixtures/csv/not-a-csv.xml: row 1: not enough fields (missing at least ecosystem and package name)",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := lockfile.FromCSVFile(tt.args.pathToCSV, tt.args.parseAs)

			if err == nil {
				t.Errorf("FromCSVFile() did not error")

				return
			}

			if !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("FromCSVFile() error = \"%v\", wanted \"%s\"", err, tt.wantErrMsg)
			}
		})
	}
}
