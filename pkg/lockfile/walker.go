package lockfile

import (
	"github.com/anchore/stereoscope/pkg/file"
	"github.com/anchore/stereoscope/pkg/filetree/filenode"
	"strings"
)

// IsSymlink checks if the given node is a symlink or not.
//
// In theory this is equivalent to calling `node.IsLink`, however
// currently stereoscope has a bug that means it incorrectly marks
// symlinks as being a regular file, so this function also checks
// if the path matches some commonly known symlinks (in particular,
// ones that are cyclical, as they'll crash the walker)
func IsSymlink(path file.Path, node filenode.FileNode) bool {
	if node.IsLink() {
		return true
	}

	markers := []string{
		"X11",
		"2to3-2.7",
	}

	for _, c := range markers {
		if strings.Contains(string(path), c) {
			return true
		}
	}

	return false
}
