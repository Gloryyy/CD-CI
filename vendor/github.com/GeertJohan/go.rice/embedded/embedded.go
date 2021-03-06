// Package embedded defines embedded data types that are shared between the go.rice package and generated code.
package embedded

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const (
	EmbedTypeGo   = 0
	EmbedTypeSyso = 1
)

// EmbeddedBox defines an embedded box
type EmbeddedBox struct {
	Name      string                   // box name
	Time      time.Time                // embed time
	EmbedType int                      // kind of embedding
	Files     map[string]*EmbeddedFile // ALL embedded files by full path
	Dirs      map[string]*EmbeddedDir  // ALL embedded dirs by full path
}

// Link creates the ChildDirs and ChildFiles links in all EmbeddedDir's
func (e *EmbeddedBox) Link() {
	for path, ed := range e.Dirs {
		fmt.Println(path)
		ed.ChildDirs = make([]*EmbeddedDir, 0)
		ed.ChildFiles = make([]*EmbeddedFile, 0)
	}
	for path, ed := range e.Dirs {
		parentDirpath, _ := filepath.Split(path)
		if strings.HasSuffix(parentDirpath, "/") {
			parentDirpath = parentDirpath[:len(parentDirpath)-1]
		}
		parentDir := e.Dirs[parentDirpath]
		if parentDir == nil {
			panic("parentDir `" + parentDirpath + "` is missing in embedded box")
		}
		parentDir.ChildDirs = append(parentDir.ChildDirs, ed)
	}
	for path, ef := range e.Files {
		dirpath, _ := filepath.Split(path)
		if strings.HasSuffix(dirpath, "/") {
			dirpath = dirpath[:len(dirpath)-1]
		}
		dir := e.Dirs[dirpath]
		if dir == nil {
			panic("dir `" + dirpath + "` is missing in embedded box")
		}
		dir.ChildFiles = append(dir.ChildFiles, ef)
	}
}

// EmbeddedDir is instanced in the code generated by the rice tool and contains all necicary information about an embedded file
type EmbeddedDir struct {
	Filename   string
	DirModTime time.Time
	ChildDirs  []*EmbeddedDir  // direct childs, as returned by virtualDir.Readdir()
	ChildFiles []*EmbeddedFile // direct childs, as returned by virtualDir.Readdir()
}

// EmbeddedFile is instanced in the code generated by the rice tool and contains all necicary information about an embedded file
type EmbeddedFile struct {
	Filename    string // filename
	FileModTime time.Time
	Content     string
}

// EmbeddedBoxes is a public register of embedded boxes
var EmbeddedBoxes = make(map[string]*EmbeddedBox)

// RegisterEmbeddedBox registers an EmbeddedBox
func RegisterEmbeddedBox(name string, box *EmbeddedBox) {
	if _, exists := EmbeddedBoxes[name]; exists {
		panic(fmt.Sprintf("EmbeddedBox with name `%s` exists already", name))
	}
	EmbeddedBoxes[name] = box
}
