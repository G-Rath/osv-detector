package database

func (db *OSVDatabase) Vulnerabilities(includeWithdrawn bool) []OSV {
	if includeWithdrawn {
		return db.vulnerabilities
	}

	var vulnerabilities []OSV

	for _, vulnerability := range db.vulnerabilities {
		if vulnerability.Withdrawn == nil {
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	return vulnerabilities
}

func (db *OSVDatabase) VulnerabilitiesAffectingPackage(ecosystem Ecosystem, pkg string, version string) []OSV {
	var vulnerabilities []OSV

	for _, vulnerability := range db.Vulnerabilities(false) {
		if vulnerability.IsAffected(ecosystem, pkg, version) {
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	return vulnerabilities
}
