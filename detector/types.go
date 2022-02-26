package detector

type PackageDetails struct {
	Name      string
	Version   string
	Ecosystem Ecosystem
}

type Ecosystem string
