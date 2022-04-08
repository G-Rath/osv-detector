package lockfile

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type MavenLockDependency struct {
	// GroupId    string `xml:"groupId"`
	XMLName xml.Name `xml:"dependency"`
	Name    string   `xml:"artifactId"`
	Version string   `xml:"version"`
}

type MavenLockFile struct {
	XMLName      xml.Name              `xml:"project"`
	ModelVersion string                `xml:"modelVersion"`
	Dependencies []MavenLockDependency `xml:"dependencies>dependency"`
}

const MavenEcosystem Ecosystem = "Maven"

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
			Version:   lockPackage.Version,
			Ecosystem: MavenEcosystem,
		})
	}

	return packages, nil
}
