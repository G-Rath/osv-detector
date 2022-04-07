package internal

type PackageDetails struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Ecosystem Ecosystem `json:"ecosystem"`
}

type Ecosystem string
