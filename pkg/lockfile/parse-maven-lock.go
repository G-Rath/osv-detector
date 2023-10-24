package lockfile

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/g-rath/osv-detector/internal/cachedregexp"
)

type MavenLockDependency struct {
	XMLName    xml.Name `xml:"dependency"`
	GroupID    string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Version    string   `xml:"version"`
}

func (mld MavenLockDependency) parseResolvedVersion(version string) string {
	versionRequirementReg := cachedregexp.MustCompile(`[[(]?(.*?)(?:,|[)\]]|$)`)

	results := versionRequirementReg.FindStringSubmatch(version)

	if results == nil || results[1] == "" {
		return "0"
	}

	return results[1]
}

func (mld MavenLockDependency) resolveVersionValue(lockfile MavenLockFile) string {
	interpolationReg := cachedregexp.MustCompile(`\${(.+)}`)

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
		"Failed to resolve version of %s: property \"%s\" could not be found\n",
		mld.GroupID+":"+mld.ArtifactID,
		results[1],
	)

	return "0"
}

func (mld MavenLockDependency) ResolveVersion(lockfile MavenLockFile) string {
	version := mld.resolveVersionValue(lockfile)

	return mld.parseResolvedVersion(version)
}

type MavenLockFile struct {
	XMLName             xml.Name              `xml:"project"`
	ModelVersion        string                `xml:"modelVersion"`
	Properties          MavenLockProperties   `xml:"properties"`
	Dependencies        []MavenLockDependency `xml:"dependencies>dependency"`
	ManagedDependencies []MavenLockDependency `xml:"dependencyManagement>dependencies>dependency"`
}

const MavenEcosystem Ecosystem = "Maven"

type MavenLockProperties struct {
	m map[string]string
}

func (p *MavenLockProperties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	p.m = map[string]string{}

	for {
		t, err := d.Token()

		if err != nil {
			return fmt.Errorf("%w", err)
		}

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

func ParseMavenLockFile(pathToLockfile string) ([]PackageDetails, error) {
	return parseFile(pathToLockfile, ParseMavenLock)
}

func ParseMavenLock(r io.Reader) ([]PackageDetails, error) {
	var parsedLockfile *MavenLockFile

	err := xml.NewDecoder(r).Decode(&parsedLockfile)

	if err != nil {
		return []PackageDetails{}, fmt.Errorf("could not parse: %w", err)
	}

	details := map[string]PackageDetails{}

	for _, lockPackage := range parsedLockfile.Dependencies {
		finalName := lockPackage.GroupID + ":" + lockPackage.ArtifactID

		details[finalName] = PackageDetails{
			Name:      finalName,
			Version:   lockPackage.ResolveVersion(*parsedLockfile),
			Ecosystem: MavenEcosystem,
			CompareAs: MavenEcosystem,
		}
	}

	// managed dependencies take precedent over standard dependencies
	for _, lockPackage := range parsedLockfile.ManagedDependencies {
		finalName := lockPackage.GroupID + ":" + lockPackage.ArtifactID

		details[finalName] = PackageDetails{
			Name:      finalName,
			Version:   lockPackage.ResolveVersion(*parsedLockfile),
			Ecosystem: MavenEcosystem,
			CompareAs: MavenEcosystem,
		}
	}

	return pkgDetailsMapToSlice(details), nil
}
