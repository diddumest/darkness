package yunyun

import (
	"path/filepath"
	"time"

	"github.com/thecsw/gana"
)

// Page is a struct for holding the page contents.
type Page struct {
	// To prevent unkeyed literars.
	_ struct{}
	// Accoutrement are additional options enabled on a page.
	Accoutrement *Accoutrement

	// Location is the Location of the page.
	Location RelativePathDir
	// Title is the title of the page.
	Title string
	// Author is the author of the page.
	Author string
	// Date is the date of the page.
	Date string
	// File is the original filename of the page (optional).
	File RelativePathFile
	// Contents is the contents of the page.
	Contents Contents
	// Scripts is the scripts of the page.
	Scripts []string
	// Stylesheets is the list of css of the page.
	Stylesheets []string
	// HtmlHead is the list of extra HTML declaration to add in the head.
	HtmlHead []string
	// Footnotes is the footnotes of the page.
	Footnotes []string
	// DateHoloscene tells us whether the first paragraph
	// on the page is given as holoscene date stamp.
	DateHoloscene bool
	// Tracked changes determined by the delta in page cache
	Changes []Change
}

// Represents explicit and implicit changes to single pages
type Change struct {
	// Optional provided reasoning set via orgmode statements/suggestion things
	Reason string
	// Date of compilation or static change information
	Date time.Time
}

func (c Change) Compare(cc Change) int {
	return c.Date.Compare(cc.Date)
}

// MetaTag is a struct for holding the meta tag.
type MetaTag struct {
	// To prevent unkeyed literars.
	_ struct{}
	// Name is the name of the meta tag.
	Name string
	// Content is the content of the meta tag.
	Content string
	// Propery is the property of the meta tag.
	Property string
}

// Link is a struct for holding the link tag.
type Link struct {
	// To prevent unkeyed literars.
	_ struct{}
	// Rel is the rel of the link tag.
	Rel string
	// Type is the type of the link tag.
	Type string
	// Href is the href of the link tag.
	Href string
}

// RelativePathDir is used in `Page` for `Location` to make sure
// that we are passing correct understanding of that it should
// only have the relative path to the workspace -- this should
// be a directory with NO filename and NO base.
type RelativePathDir string

// RelativePathFile is similar to `RelativePath` but also
// includes the filename in the end as base.
type RelativePathFile string

// RelativePath is a type constraint for `RelativePath`-style paths.
type RelativePath interface {
	RelativePathDir | RelativePathFile
}

// FullPathDir is the result of joining emilia root with `RelativePath`.
type FullPathDir string

// FullPathFile is the result of joining emilia root with `RelativePathFile`.
type FullPathFile string

// FullPath is a type constraint for `FullPath`-style types.
type FullPath interface {
	FullPathDir | FullPathFile
}

// AnyPath is a full generalization of relative and full paths.
type AnyPath interface {
	FullPath | RelativePath
}

// RelativePathTrim returns the directory of the relative file.
func RelativePathTrim(filename RelativePathFile) RelativePathDir {
	return RelativePathDir(filepath.Dir(string(filename)))
}

// JoinRelativePaths joins relative paths.
func JoinRelativePaths(dir RelativePathDir, file RelativePathFile) RelativePathFile {
	return JoinPaths(RelativePathFile(dir), file)
}

// JoinPaths joins relative paths.
func JoinPaths[T AnyPath](what ...T) T {
	return T(filepath.Join(AnyPathsToStrings(what)...))
}

// AnyPathsToStrings converts an array of `AnyPath` to `string`.
func AnyPathsToStrings[T AnyPath](what []T) []string {
	return gana.Map(func(t T) string { return string(t) }, what)
}
