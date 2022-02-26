package parsers

type EcosystemPackage struct {
	Name    string
	Version string
}

type EcosystemPackages struct {
	Packages  []EcosystemPackage
	Ecosystem string
}
