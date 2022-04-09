package lockfile

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type MavenLockDependency struct {
	// GroupId    string `xml:"groupId"`
	XMLName xml.Name `xml:"dependency"`
	Name    string   `xml:"artifactId"`
	Version string   `xml:"version"`
}

func (mld MavenLockDependency) ResolveVersion(lockfile MavenLockFile) string {
	interpolationReg := regexp.MustCompile(`\${(.+)}`)

	results := interpolationReg.FindStringSubmatch(mld.Version)

	// no interpolation, so just return the version as-is
	if results == nil {
		return mld.Version
	}
	if val, ok := lockfile.Properties.m[results[1]]; ok {
		return val
	}

	fmt.Fprintf(
		os.Stderr,
		"Failed to resolve version of %s: property \"%s\" could not be found",
		mld.Name,
		results[1],
	)

	return "0"
}

type MavenLockFile struct {
	XMLName      xml.Name              `xml:"project"`
	ModelVersion string                `xml:"modelVersion"`
	Properties   MavenLockProperties   `xml:"properties"`
	Dependencies []MavenLockDependency `xml:"dependencies>dependency"`
}

const MavenEcosystem Ecosystem = "Maven"

type MavenLockProperties struct {
	m map[string]string
}

func (p *MavenLockProperties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	p.m = map[string]string{}

	for {
		t, _ := d.Token()

		switch tt := t.(type) {
		case xml.StartElement:
			var s string

			if err := d.DecodeElement(&s, &tt); err != nil {
				return fmt.Errorf("%w", err)
			}

			p.m[tt.Name.Local] = s

		case xml.EndElement:
			if tt.Name == start.Name {
				return nil
			}
		}
	}
}

func ParseMavenLock(pathToLockfile string) ([]PackageDetails, error) {
	var parsedLockfile *MavenLockFile

	lockfileContents, err := ioutil.ReadFile(pathToLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not read %s: %w", pathToLockfile, err)
	}

	err = xml.Unmarshal(lockfileContents, &parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse %s: %w", pathToLockfile, err)
	}

	packages := make([]PackageDetails, 0, len(parsedLockfile.Dependencies))

	for _, lockPackage := range parsedLockfile.Dependencies {
		packages = append(packages, PackageDetails{
			Name:      lockPackage.Name,
			Version:   lockPackage.ResolveVersion(*parsedLockfile),
			Ecosystem: MavenEcosystem,
		})
	}

	return packages, nil
}
