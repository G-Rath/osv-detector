package reporter

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"strings"
)

func (r Report) FormatLineByLine() string {
	lines := make([]string, 0, len(r.Packages))

	for _, pkg := range r.Packages {
		if len(pkg.Vulnerabilities) == 0 {
			continue
		}

		lines = append(lines, fmt.Sprintf(
			"  %s %s",
			color.YellowString("%s@%s", pkg.Name, pkg.Version),
			color.RedString("is affected by the following vulnerabilities:"),
		))

		for _, vulnerability := range pkg.Vulnerabilities {
			lines = append(lines, fmt.Sprintf(
				"    %s %s",
				color.CyanString("%s:", vulnerability.ID),
				vulnerability.Describe(),
			))
		}
	}

	return strings.Join(lines, "\n")
}

func (r Report) FormatJSON() string {
	out, err := json.Marshal(r)

	if err != nil {
		panic("oh noes!")
	}

	return string(out)
}
