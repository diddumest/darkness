package makima

// Woof is the interface that wraps the basic methods of a makima parser.
type Woof interface {
	Read() (Woof, error)
	// Parse parses the input internally.
	Parse() Woof
	// Store the parsed objects in a simple json blob
	Cache() Woof
	// Export exports the result internally.
	Export() Woof
	// Write flushes the exported data.
	Write() error
}
