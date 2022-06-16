package lockfile

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

var errCSVRecordNotEnoughFields = errors.New("not enough fields (missing at least ecosystem and package name)")
var errCSVRecordMissingEcosystemField = errors.New("field 1 is empty (must be the name of an ecosystem)")
var errCSVRecordMissingPackageField = errors.New("field 2 is empty (must be the name of a package)")

func fromCSVRecord(lines []string) (PackageDetails, error) {
	if len(lines) < 2 {
		return PackageDetails{}, errCSVRecordNotEnoughFields
	}

	if lines[0] == "" {
		return PackageDetails{}, errCSVRecordMissingEcosystemField
	}

	if lines[1] == "" {
		return PackageDetails{}, errCSVRecordMissingPackageField
	}

	return PackageDetails{
		Name:      lines[1],
		Version:   lines[2],
		Ecosystem: Ecosystem(lines[0]),
	}, nil
}

func fromCSV(reader io.Reader) ([]PackageDetails, error) {
	var packages []PackageDetails

	i := 0
	r := csv.NewReader(reader)

	for {
		i++
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return packages, fmt.Errorf("%w", err)
		}

		details, err := fromCSVRecord(record)
		if err != nil {
			return packages, fmt.Errorf("row %d: %w", i, err)
		}

		packages = append(packages, details)
	}

	sort.Slice(packages, func(i, j int) bool {
		if packages[i].Name == packages[j].Name {
			return packages[i].Version < packages[j].Version
		}

		return packages[i].Name < packages[j].Name
	})

	return packages, nil
}

func FromCSVRows(filePath string, parseAs string, rows []string) (Lockfile, error) {
	packages, err := fromCSV(strings.NewReader(strings.Join(rows, "\n")))

	return Lockfile{
		FilePath: filePath,
		ParsedAs: parseAs,
		Packages: packages,
	}, err
}

func FromCSVFile(pathToCSV string, parseAs string) (Lockfile, error) {
	reader, err := os.Open(pathToCSV)

	if err != nil {
		return Lockfile{}, fmt.Errorf("could not read %s: %w", pathToCSV, err)
	}

	packages, err := fromCSV(reader)

	if err != nil {
		err = fmt.Errorf("%s: %w", pathToCSV, err)
	}

	return Lockfile{
		FilePath: pathToCSV,
		ParsedAs: parseAs,
		Packages: packages,
	}, err
}
