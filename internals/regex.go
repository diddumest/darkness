package internals

import "regexp"

var (
	// LinkRegexp is the regexp for matching links
	LinkRegexp = regexp.MustCompile(`\[\[([^][]+)\]\[([^][]+)\]\]`)
	// URLRegexp is yoinked from https://ihateregex.io/expr/url/
	URLRegexp = regexp.MustCompile(`(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*))`)
	// BoldText is the regexp for matching bold text
	BoldText = regexp.MustCompile(`(?mU)(^|[ ()_%<>])\*(\S|\S\S|\S.+\S)\*($|[ (),.!?;&_%<>])`)
	// ItalicText is the regexp for matching italic text
	ItalicText = regexp.MustCompile(`(?mU)(^|[ ()_%<>])/(\S|\S\S|\S.+\S)/($|[ (),.!?;&_%<>])`)
	// BoldItalicText is the regexp for matching bold-italic text from the left
	BoldItalicTextBegin = regexp.MustCompile(`(?mU)(^|[ ()_%<>])\*/`)
	// BoldItalicTextEnd is the regexp for matching bold-italic text from the right
	BoldItalicTextEnd = regexp.MustCompile(`(?mU)/\*($|[ (),.!?;&_%<>])`)
	// VerbatimText is the regexp for matching verbatim text
	VerbatimText = regexp.MustCompile(`(?mU)(^|[ ()_%<>])=(\S|\S\S|\S.+\S)=($|[ (),.!?;&_%<>])`)
	// KeyboardRegexp is the regexp for matching keyboard text
	KeyboardRegexp = regexp.MustCompile(`kbd:\[([^][]+)\]`)
	// MathRegexp is the regexp for matching math text
	MathRegexp = regexp.MustCompile(`(?mU)\$(.+)\$`)
	// ImageRegexp is the regexp for matching images (png, gif, jpg, jpeg, svg, webp)
	ImageExtRegexp = regexp.MustCompile(`\.(png|gif|jpg|jpeg|svg|webp)$`)
	// AudioRegexp is the regexp for matching audio (mp3, flac, midi)
	AudioFileExtRegexp = regexp.MustCompile(`\.(mp3|flac|midi)$`)

	// FootnoteRegexp is the regexp for matching footnotes
	FootnoteRegexp = regexp.MustCompile(`\[fn::([^][]+)\]`)
	// FootnoteReferenceRegexp is the regexp for matching footnotes references
	FootnotePostProcessingRegexp = regexp.MustCompile(`!(\d+)!`)
)
