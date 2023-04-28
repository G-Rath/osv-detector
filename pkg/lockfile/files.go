package lockfile

import (
	"fmt"
	"github.com/anchore/stereoscope/pkg/file"
	"github.com/anchore/stereoscope/pkg/image"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Lockfile struct {
	FilePath string   `json:"filePath"`
	ParsedAs string   `json:"parsedAs"`
	Packages Packages `json:"packages"`
}

func (l Lockfile) String() string {
	lines := make([]string, 0, len(l.Packages))

	for _, details := range l.Packages {
		ecosystem := details.Ecosystem

		if ecosystem == "" {
			ecosystem = "<unknown>"
		}

		ln := fmt.Sprintf("  %s: %s", ecosystem, details.Name)

		if details.Version != "" {
			ln += "@" + details.Version
		}

		if details.Commit != "" {
			ln += " (" + details.Commit + ")"
		}

		lines = append(lines, ln)
	}

	return strings.Join(lines, "\n")
}

// LFS Local File System
// VFS Virtual File System
// IFS Image File System

// ParsableFile is an abstraction for a file that has been opened for parsing
// to create a Lockfile, and that knows how to open other ParsableFiles
// relative to itself.
type ParsableFile interface {
	io.ReadCloser

	OpenRelative(string) (ParsableFile, error)

	Path() string
}

// A LocalFile represents a file that exists on the local filesystem.
type LocalFile struct {
	io.ReadCloser

	path string
}

func (f LocalFile) OpenRelative(path string) (ParsableFile, error) {
	return OpenLocalFile(filepath.Join(filepath.Dir(f.path), path))
}

func (f LocalFile) Path() string { return f.path }

var _ ParsableFile = LocalFile{}

func OpenLocalFile(path string) (LocalFile, error) {
	r, err := os.Open(path)

	if err != nil {
		return LocalFile{}, err
	}

	return LocalFile{r, path}, nil
}

// An ImageFile represents a file that exists in an image.
type ImageFile struct {
	io.ReadCloser

	path string
	img  image.Image
}

func (f ImageFile) OpenRelative(path string) (ParsableFile, error) {
	return OpenImageFile(f.img, filepath.Join(filepath.Dir(f.path), path))
}

func (f ImageFile) Path() string { return f.path }

var _ ParsableFile = ImageFile{}

func OpenImageFile(img image.Image, path string) (ImageFile, error) {
	r, err := img.OpenPathFromSquash(file.Path(path))

	if err != nil {
		return ImageFile{}, err
	}

	return ImageFile{r, path, img}, nil
}
