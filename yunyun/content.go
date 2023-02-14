package yunyun

import (
	"github.com/thecsw/gana"
)

// Content is a piece of content of a page.
type Content struct {
	// Type is the type of content.
	Type TypeContent

	// Options tells us about the options enabled on the type.
	Options Bits
	// HeadingLevel is the heading level of the content (1 being the title, starts at 2).
	HeadingLevel int
	// HeadingLevelAdjusted is adjusted `HeadingLevel` by `emilia`.
	HeadingLevelAdjusted int
	// HeadingChild tells us if the current heading is a child of some previous heading.
	HeadingChild bool
	// HeadingFirst tells us if the current heading is the first heading on the page.
	HeadingFirst bool
	// HeadingLast tells us if the current heading is the last heading on the page.
	HeadingLast bool
	// Heading is the heading text.
	Heading string
	// Paragraph is the paragraph text.
	Paragraph string
	// List is the list of items, unordered.
	List []string
	// ListNumbered is the list of items, numbered.
	ListNumbered []string
	// Link is the link text.
	Link string
	// LinkTitle is the link title.
	LinkTitle string
	// LinkDescription is the optional description of the link.
	LinkDescription string
	// SourceCode is the source code.
	SourceCode string
	// SourceCodeLanguage is the language of the source code.
	SourceCodeLang string
	// RawHTML is the raw HTML.
	RawHTML string
	// AttentionTitle is the attention text title (IMPORTANT, WARNING, etc.).
	AttentionTitle string
	// AttentionText is the attention text.
	AttentionText string
	// Table is the table of items.
	Table [][]string
	// TableHeaders tell us whether the table has headers
	// (use the first row as headers instead of data).
	TableHeaders bool
	// Caption is the current caption.
	Caption string
	// Summary is the current summary, used by details to denote
	// the title of the summary block.
	Summary string
	// GalleryPath stores the path declared by the gallery directive,
	// it can either be the relative path (to the page that it was
	// declared on) or some other http/absolute link.
	GalleryPath RelativePathDir
	// GalleryImagesPerRow stores the number of default images per row,
	// therefore what flex class to use -- defaults to 3.
	GalleryImagesPerRow uint
}

// Contents is a type of contents
type Contents []*Content

// Galleries returns all contents that are galleries AND proper list types.
func (c Contents) Galleries() Contents {
	return gana.Filter(func(v *Content) bool { return v.IsGallery() && v.IsList() }, c)
}

// Headings only returns headings of contents, useful for bulding tables
// of contents and alike.
func (c Contents) Headings() Contents {
	return gana.Filter(func(v *Content) bool { return v.IsHeading() }, c)
}

// SourceCodeBlocks returns all source code blocks from contents.
func (c Contents) SourceCodeBlocks() Contents {
	return gana.Filter(func(v *Content) bool { return v.IsSourceCode() }, c)
}

// IsHeading tells us if the content is a heading.
func (c Content) IsHeading() bool { return c.Type == TypeHeading }

// IsParagraph tells us if the content is a paragraph.
func (c Content) IsParagraph() bool { return c.Type == TypeParagraph }

// IsList tells us if the content is a list.
func (c Content) IsList() bool { return c.Type == TypeList }

// IsListNumbered tells us if the content is a numbered list.
func (c Content) IsListNumbered() bool { return c.Type == TypeListNumbered }

// IsLink tells us if the content is a link.
func (c Content) IsLink() bool { return c.Type == TypeLink }

// IsSourceCode tells us if the content is a source code block.
func (c Content) IsSourceCode() bool { return c.Type == TypeSourceCode }

// IsRawHTML tells us if the content is a raw HTML block.
func (c Content) IsRawHTML() bool { return c.Type == TypeRawHTML }

// IsRawHTMLUnsafe tells us if the html block is raw and unsafe.
func (c Content) IsRawHTMLUnsafe() bool { return HasFlag(&c.Options, InRawHtmlFlagUnsafe) }

// IsHorizontalLine tells us if the content is a horizontal line.
func (c Content) IsHorizontalLine() bool { return c.Type == TypeHorizontalLine }

// IsAttentionBlock tells us if the content is an attention text block.
func (c Content) IsAttentionBlock() bool { return c.Type == TypeAttentionText }

// IsTable tells us if the content block is a table.
func (c Content) IsTable() bool { return c.Type == TypeTable }

// IsCentered returns true if the content should be centered, false otherwise.
func (c Content) IsCentered() bool { return HasFlag(&c.Options, InCenterFlag) }

// IsDetails returns true if content is a part of a details block, false otherwise.
func (c Content) IsDetails() bool { return HasFlag(&c.Options, InDetailsFlag) }

// IsGallery returns true if content is a gallery, false otherwise.
func (c Content) IsGallery() bool { return HasFlag(&c.Options, InGalleryFlag) }

// IsDropCap returns true if the first letter should be a drop cap, false otherwise.
func (c Content) IsDropCap() bool { return HasFlag(&c.Options, InDropCapFlag) }

// IsQuote returns true if the content is a quotation, false otherwise.
func (c Content) IsQuote() bool { return HasFlag(&c.Options, InQuoteFlag) }
