package semver

import (
	"regexp"
	"strconv"
)

func Parse(line string) Version {
	var components []int

	numberReg := regexp.MustCompile(`\d`)

	current := ""
	foundBuild := false

	for _, c := range line {
		if foundBuild {
			current += string(c)

			continue
		}

		// this is part of a component version
		if numberReg.MatchString(string(c)) {
			current += string(c)

			continue
		}

		// at this point, we:
		//   1. might be parsing a component (as foundBuild != true)
		//   2. we're not looking at a part of a component (as c != number)
		//
		// so c must be either:
		//   1. a component terminator (.), or
		//   2. the start of the build string

		// this is a component terminator
		if c == '.' {
			if current != "" {
				v, _ := strconv.Atoi(current)

				components = append(components, v)
				current = ""
			}

			continue
		}

		// anything else is part of the build string
		foundBuild = true

		if current != "" {
			v, _ := strconv.Atoi(current)

			components = append(components, v)
		}

		current = string(c)
	}

	// if we looped over everything without finding a build string,
	// then what we were current parsing is actually a component
	if !foundBuild {
		v, _ := strconv.Atoi(current)

		components = append(components, v)
		current = ""
	}

	return Version{
		Components: components,
		Build:      current,
	}
}
