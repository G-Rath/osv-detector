package semver

import (
	"regexp"
	"strconv"
)

// ([\d]+)\.([\d]+)\.

type char = int32

func isComponentTerminator(c char) bool {
	// regexp.MustCompile(`[.-]`).

	return c == '.' || c == '-'
}

func terminatesComponent(c char) bool {
	return isComponentTerminator(c)
}

// func (v *Version) updateComponent

func Parse(line string) Version {
	var components []int
	// components := []int{}
	// currentCom := 0

	numberReg := regexp.MustCompile(`\d`)

	// build := ""
	current := ""
	foundBuild := false

	for _, c := range line {
		// fmt.Printf("%s (%d)\n", string(c))

		if foundBuild {
			current += string(c)

			continue
		}

		// this is part of a component version
		if numberReg.MatchString(string(c)) {
			current += string(c)

			continue
		}

		// this is a component terminator
		if c == '.' {
			v, _ := strconv.Atoi(current)

			components = append(components, v)
			current = ""

			continue
		}

		// anything else is part of the build string
		foundBuild = true
		current += string(c)
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
