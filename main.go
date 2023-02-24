package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anchore/stereoscope"
	"github.com/anchore/stereoscope/pkg/file"
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/fatih/color"
	"github.com/g-rath/osv-detector/internal"
	"github.com/g-rath/osv-detector/internal/configer"
	"github.com/g-rath/osv-detector/internal/reporter"
	"github.com/g-rath/osv-detector/pkg/database"
	"github.com/g-rath/osv-detector/pkg/lockfile"
	"golang.org/x/exp/slices"
)

// these come from goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func makeAPIDBConfig() database.Config {
	return database.Config{
		Name: "osv.dev v1 API",
		Type: "api",
		URL:  "https://api.osv.dev/v1",
	}
}

func makeEcosystemDBConfig(ecosystem internal.Ecosystem) database.Config {
	return database.Config{
		Name: string(ecosystem),
		Type: "zip",
		URL:  fmt.Sprintf("https://osv-vulnerabilities.storage.googleapis.com/%s/all.zip", ecosystem),
	}
}

type OSVDatabases []database.DB

func (dbs OSVDatabases) transposePkgResults(
	pkg internal.PackageDetails,
	ignores []string,
	packageIndex int,
	allVulns [][]database.Vulnerabilities,
) reporter.PackageDetailsWithVulnerabilities {
	vulnerabilities := make(database.Vulnerabilities, 0)
	ignored := make(database.Vulnerabilities, 0)

	for _, vulns1 := range allVulns {
		vulns := vulns1[packageIndex]

		for _, vulnerability := range vulns {
			// skip vulnerabilities that were already included from a previous database
			if vulnerabilities.Includes(vulnerability) || ignored.Includes(vulnerability) {
				continue
			}

			if slices.Contains(ignores, vulnerability.ID) {
				ignored = append(ignored, vulnerability)
			} else {
				vulnerabilities = append(vulnerabilities, vulnerability)
			}
		}
	}

	return reporter.PackageDetailsWithVulnerabilities{
		PackageDetails:  pkg,
		Vulnerabilities: vulnerabilities,
		Ignored:         ignored,
	}
}

func (dbs OSVDatabases) check(r *reporter.Reporter, lockf lockfile.Lockfile, ignores []string) reporter.Report {
	report := reporter.Report{
		Lockfile: lockf,
		Packages: make([]reporter.PackageDetailsWithVulnerabilities, 0, len(lockf.Packages)),
	}

	vulns := make([][]database.Vulnerabilities, 0, len(dbs))

	for _, db := range dbs {
		results, err := db.Check(lockf.Packages)

		if err != nil {
			r.PrintError(color.RedString(fmt.Sprintf(
				"  an api error occurred while trying to check the packages listed in %s: %v\n",
				lockf.FilePath,
				err,
			)))

			continue
		}

		vulns = append(vulns, results)
	}

	for i, pkg := range lockf.Packages {
		report.Packages = append(
			report.Packages,
			dbs.transposePkgResults(pkg, ignores, i, vulns),
		)
	}

	return report
}

// returns the OSV databases to use for the given database configs,
// assuming they have already been loaded
func (dbs OSVDatabases) forConfigs(dbConfigs []database.Config) OSVDatabases {
	specificDBs := make(OSVDatabases, 0)

	for _, db := range dbs {
		for _, dbConfig := range dbConfigs {
			if dbConfig.Identifier() == db.Identifier() {
				specificDBs = append(specificDBs, db)
			}
		}
	}

	return specificDBs
}

func uniqueDBConfigs(configs []*configer.Config) []database.Config {
	var dbConfigs []database.Config

	for _, config := range configs {
		for _, dbConfig := range config.Databases {
			alreadyExists := false

			for _, dbc := range dbConfigs {
				if alreadyExists {
					continue
				}

				if dbc.Identifier() == dbConfig.Identifier() {
					alreadyExists = true
				}
			}

			if alreadyExists {
				continue
			}

			dbConfigs = append(dbConfigs, dbConfig)
		}
	}

	return dbConfigs
}

func describeDB(db database.DB) string {
	switch tt := db.(type) {
	case *database.APIDB:
		return fmt.Sprintf(
			"using batches of %s",
			color.YellowString("%d", tt.BatchSize),
		)
	case *database.ZipDB:
		count := tt.VulnerabilitiesCount

		return fmt.Sprintf(
			"%s %s, including withdrawn - last updated %s",
			color.YellowString("%d", count),
			reporter.Form(count, "vulnerability", "vulnerabilities"),
			tt.UpdatedAt,
		)
	case *database.DirDB:
		count := tt.VulnerabilitiesCount

		return fmt.Sprintf(
			"%s %s, including withdrawn",
			color.YellowString("%d", count),
			reporter.Form(count, "vulnerability", "vulnerabilities"),
		)
	}

	return ""
}

func loadDatabases(
	r *reporter.Reporter,
	dbConfigs []database.Config,
	listPackages bool,
	offline bool,
	batchSize int,
) (OSVDatabases, bool) {
	dbs := make(OSVDatabases, 0, len(dbConfigs))

	// an easy dirty little optimisation: we don't need any databases
	// if we're going to be listing packages, so return the empty slice
	if listPackages {
		return dbs, false
	}

	errored := false

	r.PrintText("Loaded the following OSV databases:\n")

	for _, dbConfig := range dbConfigs {
		r.PrintText(fmt.Sprintf("  %s", dbConfig.Name))

		db, err := database.Load(dbConfig, offline, batchSize)

		if err != nil {
			r.PrintDatabaseLoadErr(err)
			errored = true

			continue
		}

		desc := describeDB(db)

		if desc != "" {
			desc = fmt.Sprintf(" (%s)", desc)
		}

		r.PrintText(fmt.Sprintf("%s\n", desc))

		dbs = append(dbs, db)
	}

	r.PrintText("\n")

	return dbs, errored
}

const parseAsCsvFile = "csv-file"
const parseAsCsvRow = "csv-row"

func findLockfiles(r *reporter.Reporter, pathToLockOrDirectory string, parseAs string) ([]string, bool) {
	lockfiles := make([]string, 0, 1)
	file, err := os.Open(pathToLockOrDirectory)

	if err == nil {
		info, err := file.Stat()

		if err == nil {
			if info.IsDir() {
				if info.Name() == "node_modules" {
					lockfiles = append(lockfiles, pathToLockOrDirectory)
				} else {
					dirs, err := file.ReadDir(-1)

					if err == nil {
						for _, dir := range dirs {
							if dir.IsDir() {
								continue
							}

							if parseAs != parseAsCsvFile {
								if p, _ := lockfile.FindParser(dir.Name(), parseAs); p == nil {
									continue
								}
							}

							lockfiles = append(lockfiles, filepath.Join(pathToLockOrDirectory, dir.Name()))
						}
					}
				}
			} else {
				lockfiles = append(lockfiles, pathToLockOrDirectory)
			}
		}
	}

	if err != nil {
		r.PrintError(fmt.Sprintf("Error reading %s: %v\n", pathToLockOrDirectory, err))
	}

	sort.Slice(lockfiles, func(i, j int) bool {
		return lockfiles[i] < lockfiles[j]
	})

	return lockfiles, err != nil
}

// findAllLockfiles checks the provided paths to determine an absolute list of lockfiles
// to be checked, expanding any directories into a list of supported lockfiles
func findAllLockfilesOnHost(r *reporter.Reporter, pathsToCheck []string, parseAsGlobal string) ([]string, bool) {
	var paths []string

	if parseAsGlobal == parseAsCsvRow {
		return []string{parseAsCsvRow + ":-"}, false
	}

	errored := false

	for _, pathToLockOrDirectory := range pathsToCheck {
		parseAs, pathToLockOrDirectory := parseLockfilePathWithParseAs(pathToLockOrDirectory)

		if parseAs == "" {
			parseAs = parseAsGlobal
		}

		lps, erred := findLockfiles(r, pathToLockOrDirectory, parseAs)

		if erred {
			errored = true
		}

		for _, p := range lps {
			paths = append(paths, parseAs+":"+filepath.Clean(p))
		}
	}

	return paths, errored
}

func imageID(img *image.Image) string {
	for _, tag := range img.Metadata.Tags {
		return tag.Name()
	}

	return img.Metadata.ID
}

func findAllLockfiles(r *reporter.Reporter, pathsToCheck []string, parseAs string, img *image.Image) ([]string, bool) {
	if img == nil {
		return findAllLockfilesOnHost(r, pathsToCheck, parseAs)
	}

	r.PrintText(fmt.Sprintf("reading image %s\n", imageID(img)))

	paths := make([]string, 0, 2)

	if img.SquashedTree().HasPath("/lib/apk/db/installed") {
		paths = append(paths, "/lib/apk/db/installed")
	}

	if img.SquashedTree().HasPath("/var/lib/dpkg/status") {
		paths = append(paths, "/var/lib/dpkg/status")
	}

	// todo: we want to search the container for things like node_modules, lockfiles, etc

	return paths, false
}

func parseLockfile(pathToLock string, args []string, img *image.Image) (lockfile.Lockfile, error) {
	parseAs, pathToLock := parseLockfilePathWithParseAs(pathToLock)
	if parseAs == parseAsCsvRow {
		l, err := lockfile.FromCSVRows(pathToLock, parseAs, args)

		if err != nil {
			err = fmt.Errorf("%w", err)
		}

		return l, err
	}
	if parseAs == parseAsCsvFile {
		l, err := lockfile.FromCSVFile(pathToLock, parseAs)

		if err != nil {
			err = fmt.Errorf("%w", err)
		}

		return l, err
	}

	if img != nil {
		if pathToLock == "/lib/apk/db/installed" {
			r, err := img.OpenPathFromSquash(file.Path(pathToLock))
			if err != nil {
				err = fmt.Errorf("%w", err)

				return lockfile.Lockfile{}, err
			}

			l, err := lockfile.FromApkInstalled(r, pathToLock)

			if err != nil {
				err = fmt.Errorf("%w", err)
			}

			return l, err
		}

		if pathToLock == "/var/lib/dpkg/status" {
			r, err := img.OpenPathFromSquash(file.Path(pathToLock))
			if err != nil {
				err = fmt.Errorf("%w", err)

				return lockfile.Lockfile{}, err
			}

			l, err := lockfile.FromDpkgStatus(r, pathToLock)

			if err != nil {
				err = fmt.Errorf("%w", err)
			}

			return l, err
		}
	}

	// adding the separator to the current lock path is an easy way to handle "node_modules"
	// is passed by itself when checking that the full directory name is "node_modules",
	// since it doesn't matter for this check if we end up with two separators at the start
	if strings.HasSuffix(string(filepath.Separator)+pathToLock, string(filepath.Separator)+"node_modules") {
		l, err := lockfile.WalkNodeModules(pathToLock)

		if err != nil {
			err = fmt.Errorf("%w", err)
		}

		return l, err
	}

	l, err := lockfile.Parse(pathToLock, parseAs)

	if err != nil {
		err = fmt.Errorf("%w", err)
	}

	return l, err
}

type stringsFlag []string

func (s *stringsFlag) String() string {
	return fmt.Sprint(*s)
}

func (s *stringsFlag) Set(value string) error {
	*s = append(*s, value)

	return nil
}

type lockfileAndConfigOrErr struct {
	lockf  lockfile.Lockfile
	config *configer.Config
	err    error
}

type lockfileAndConfigOrErrs []lockfileAndConfigOrErr

func (files lockfileAndConfigOrErrs) getConfigs() []*configer.Config {
	configs := make([]*configer.Config, 0, len(files))

	for _, file := range files {
		if file.err != nil {
			continue
		}

		configs = append(configs, file.config)
	}

	return configs
}

func (files lockfileAndConfigOrErrs) adjustExtraDatabases(
	removeConfigDatabases bool,
	addDefaultAPIDatabase bool,
	addEcosystemDatabases bool,
) {
	for _, file := range files {
		if file.err != nil {
			continue
		}
		var extraDBConfigs []database.Config

		if removeConfigDatabases {
			file.config.Databases = []database.Config{}
		}

		if addDefaultAPIDatabase {
			extraDBConfigs = append(extraDBConfigs, makeAPIDBConfig())
		}

		if addEcosystemDatabases {
			ecosystems := collectEcosystems([]lockfileAndConfigOrErr{file})

			for _, ecosystem := range ecosystems {
				extraDBConfigs = append(extraDBConfigs, makeEcosystemDBConfig(ecosystem))
			}
		}

		// a bit of a hack to let us reuse this method...
		file.config.Databases = uniqueDBConfigs([]*configer.Config{
			file.config,
			{Databases: extraDBConfigs},
		})
	}
}

func parseLockfilePathWithParseAs(lockfilePathWithParseAs string) (string, string) {
	if !strings.Contains(lockfilePathWithParseAs, ":") {
		return "", lockfilePathWithParseAs
	}

	parseAs, path, _ := strings.Cut(lockfilePathWithParseAs, ":")

	return parseAs, path
}

func readAllLockfiles(
	r *reporter.Reporter,
	pathsToLocksWithParseAs []string,
	img *image.Image,
	args []string,
	checkForLocalConfig bool,
	config *configer.Config,
) lockfileAndConfigOrErrs {
	lockfiles := make([]lockfileAndConfigOrErr, 0, len(pathsToLocksWithParseAs))

	// todo: for now lets just not do local configs
	if img != nil {
		checkForLocalConfig = false
	}

	for _, pathToLockWithParseAs := range pathsToLocksWithParseAs {
		_, pathToLock := parseLockfilePathWithParseAs(pathToLockWithParseAs)
		if checkForLocalConfig {
			base := filepath.Dir(pathToLock)
			con, err := configer.Find(r, base)

			if err != nil {
				// treat config errors as the same as if we failed to load the lockfile
				// as continuing without the desired config could cause different results
				// e.g. if the config has ignores or custom databases
				lockfiles = append(lockfiles, lockfileAndConfigOrErr{lockfile.Lockfile{}, config, err})

				continue
			}

			config = &con
		} else if config.FilePath != "" {
			// if there's a global config, then copy it - otherwise all lockfiles
			// will hold a reference to the same config, which can result in configs
			// for ecosystem-specific databases being used unnecessarily for lockfiles
			// that don't have any packages that are part of that ecosystem
			config = &configer.Config{
				FilePath:  config.FilePath,
				Ignore:    config.Ignore,
				Databases: config.Databases,
			}
		}

		lockf, err := parseLockfile(pathToLockWithParseAs, args, img)
		lockfiles = append(lockfiles, lockfileAndConfigOrErr{lockf, config, err})
	}

	return lockfiles
}

func collectEcosystems(files []lockfileAndConfigOrErr) []internal.Ecosystem {
	var ecosystems []internal.Ecosystem

	for _, result := range files {
		if result.err != nil {
			continue
		}

		for _, ecosystem := range result.lockf.Packages.Ecosystems() {
			alreadyExists := false

			for _, eco := range ecosystems {
				if alreadyExists {
					continue
				}

				if eco == ecosystem {
					alreadyExists = true
				}
			}

			if alreadyExists {
				continue
			}

			ecosystems = append(ecosystems, ecosystem)
		}
	}

	return ecosystems
}

func run(args []string, stdout, stderr io.Writer) int {
	var globalIgnores stringsFlag
	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	offline := cli.Bool("offline", false, "Perform checks using only the cached databases on disk")
	parseAs := cli.String("parse-as", "", "Name of a supported lockfile to parse the input files as")
	configPath := cli.String("config", "", "Path to a config file to use for all lockfiles")
	noConfig := cli.Bool("no-config", false, "Disable loading of any config files")
	noConfigIgnores := cli.Bool("no-config-ignores", false, "Don't respect any OSVs listed as ignored in configs")
	noConfigDatabases := cli.Bool("no-config-databases", false, "Don't load any extra databases listed in configs")
	printVersion := cli.Bool("version", false, "Print version information")
	listEcosystems := cli.Bool("list-ecosystems", false, "List all of the known ecosystems that are supported by the detector")
	listPackages := cli.Bool("list-packages", false, "List the packages that are parsed from the input files")
	outputAsJSON := cli.Bool("json", false, "Output the results in JSON format")
	useDatabases := cli.Bool("use-dbs", true, "Use the databases from osv.dev to check for known vulnerabilities")
	useAPI := cli.Bool("use-api", false, "Use the osv.dev API to check for known vulnerabilities")
	batchSize := cli.Int("batch-size", 1000, "The number of packages to include in each batch when using the api database")
	imagePath := cli.String("image", "", "Path to an Docker save tar, docker image, or OCI directory/archive to check")

	cli.Var(&globalIgnores, "ignore", `ID of an OSV to ignore when determining exit codes.
This flag can be passed multiple times to ignore different vulnerabilities`)

	// cli is set for ExitOnError so this will never return an error
	_ = cli.Parse(args)

	r := reporter.New(stdout, stderr, *outputAsJSON)
	if *outputAsJSON {
		defer r.PrintJSONResults()
	}

	if *printVersion {
		r.PrintText(fmt.Sprintf("osv-detector %s (%s, commit %s)\n", version, date, commit))

		return 0
	}

	if *listEcosystems {
		r.PrintKnownEcosystems()

		return 0
	}

	// ensure that if the global parseAs is set, it is one of the supported values
	if *parseAs != "" && *parseAs != parseAsCsvFile && *parseAs != parseAsCsvRow {
		if parser, parsedAs := lockfile.FindParser("", *parseAs); parser == nil {
			r.PrintError(fmt.Sprintf("Don't know how to parse files as \"%s\" - supported values are:\n", parsedAs))

			for _, s := range lockfile.ListParsers() {
				r.PrintError(fmt.Sprintf("  %s\n", s))
			}

			r.PrintError(fmt.Sprintf("  %s\n", parseAsCsvFile))
			r.PrintError(fmt.Sprintf("  %s\n", parseAsCsvRow))

			return 127
		}
	}

	var img *image.Image

	if *imagePath != "" {
		// context for network requests
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var err error

		r.PrintText(fmt.Sprintf("opening image %s\n", *imagePath))
		img, err = stereoscope.GetImage(ctx, *imagePath)

		if err != nil {
			return 1
		}

		// note: we are writing out temp files which should be cleaned up after you're done with the image object
		defer img.Cleanup()
	}

	pathsToLocksWithParseAs, errored := findAllLockfiles(r, cli.Args(), *parseAs, img)

	if len(pathsToLocksWithParseAs) == 0 {
		r.PrintError(
			"You must provide at least one path to either a lockfile or a directory containing at least one lockfile (see --help for usage and flags)\n",
		)

		// being provided with at least one path and not hitting an error on any of those
		// paths means everything was valid, we just didn't find any parsable lockfiles
		// in any of the directories
		if len(cli.Args()) > 0 && !errored {
			// so we want to use a specific exit code to represent this state
			return 128
		}

		return 127
	}

	exitCode := 0

	var config configer.Config
	loadLocalConfig := !*noConfig

	// if we're listing packages, then we don't need to do _any_ config loading
	if *listPackages {
		loadLocalConfig = false
	} else if loadLocalConfig && *configPath != "" {
		con, err := configer.Load(r, *configPath)

		if err != nil {
			r.PrintError(fmt.Sprintf("Error, %s\n", err))

			return 127
		}

		config = con
		loadLocalConfig = false
	}

	files := readAllLockfiles(r, pathsToLocksWithParseAs, img, cli.Args(), loadLocalConfig, &config)

	files.adjustExtraDatabases(*noConfigDatabases, *useAPI, *useDatabases)

	dbs, errored := loadDatabases(
		r,
		uniqueDBConfigs(files.getConfigs()),
		*listPackages,
		*offline,
		*batchSize,
	)

	if errored {
		exitCode = 127
	}

	for i, result := range files {
		if i >= 1 {
			r.PrintText("\n")
		}

		if result.err != nil {
			r.PrintError(fmt.Sprintf("Error, %s\n", result.err))
			exitCode = 127

			continue
		}

		config := result.config
		lockf := result.lockf

		if *noConfigIgnores {
			config.Ignore = []string{}
		}

		r.PrintText(fmt.Sprintf(
			"%s: found %s %s\n",
			color.MagentaString("%s", lockf.FilePath),
			color.YellowString("%d", len(lockf.Packages)),
			reporter.Form(len(lockf.Packages), "package", "packages"),
		))

		if *listPackages {
			r.PrintResult(lockf)

			continue
		}

		ignores := make(
			[]string,
			0,
			// len cannot return negative numbers, but the types can't reflect that
			uint64(len(globalIgnores))+uint64(len(config.Ignore)),
		)

		// an empty FilePath means we didn't load a config
		if config.FilePath != "" {
			var ignoresStr string

			if *noConfigIgnores {
				ignoresStr = "skipping any ignores"
			} else {
				ignores = append(ignores, config.Ignore...)
				ignoresStr = color.YellowString("%d %s",
					len(config.Ignore),
					reporter.Form(len(config.Ignore), "ignore", "ignores"),
				)
			}

			r.PrintText(fmt.Sprintf(
				"  Using config at %s (%s)\n",
				color.MagentaString(config.FilePath),
				ignoresStr,
			))
		}

		ignores = append(ignores, globalIgnores...)

		dbs := dbs.forConfigs(config.Databases)
		for _, db := range dbs {
			desc := describeDB(db)

			if desc != "" {
				desc = fmt.Sprintf(" (%s)", desc)
			}

			r.PrintText(fmt.Sprintf(
				"  Using db %s%s\n",
				color.HiCyanString(db.Name()),
				desc,
			))
		}
		r.PrintText("\n")

		report := dbs.check(r, lockf, ignores)

		r.PrintResult(report)

		if report.HasKnownVulnerabilities() && exitCode == 0 {
			exitCode = 1
		}
	}

	return exitCode
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
