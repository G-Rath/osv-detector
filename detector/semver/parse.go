package semver

import (
	"regexp"
	"strconv"
)

func Parse(line string) Version {
	var components []int

	numberReg := regexp.MustCompile(`\d`)

	currentCom := ""
	foundBuild := false
	emptyComponent := false

	for _, c := range line {
		if foundBuild {
			currentCom += string(c)

			continue
		}

		// this is part of a component version
		if numberReg.MatchString(string(c)) {
			currentCom += string(c)

			continue
		}

		// at this point, we:
		//   1. might be parsing a component (as foundBuild != true)
		//   2. we're not looking at a part of a component (as c != number)
		//
		// so c must be either:
		//   1. a component terminator (.), or
		//   2. the start of the build string
		//
		// either way, we will be terminating the current component being
		// parsed (if any), so let's do that first
		if currentCom != "" {
			v, _ := strconv.Atoi(currentCom)

			components = append(components, v)
			currentCom = ""

			emptyComponent = false
		}

		// a component terminator means there might be another component
		// afterwards, so don't start parsing the build string just yet
		if c == '.' {
			emptyComponent = true

			continue
		}

		// anything else is part of the build string
		foundBuild = true
		currentCom = string(c)
	}

	// if we looped over everything without finding a build string,
	// then what we were currently parsing is actually a component
	if !foundBuild && currentCom != "" {
		v, _ := strconv.Atoi(currentCom)

		components = append(components, v)
		currentCom = ""
		emptyComponent = false
	}

	// if we ended with an empty component section,
	// prefix the build string with a '.'
	if emptyComponent {
		currentCom = "." + currentCom
	}

	return Version{
		Components: components,
		Build:      currentCom,
	}
}
