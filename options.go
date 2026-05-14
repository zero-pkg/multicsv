package multicsv

// ReaderOption changes the default reader behavior.
type ReaderOption func(*readerOptions)

type readerOptions struct {
	skipHeader          bool
	autoDetectDelimiter bool
}

func newReaderOptions(opts ...ReaderOption) readerOptions {
	options := readerOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	return options
}

// WithSkipHeader skips the first CSV record before reading data.
func WithSkipHeader() ReaderOption {
	return func(o *readerOptions) {
		o.skipHeader = true
	}
}

// WithAutoDetectDelimiter detects the CSV delimiter before reading data.
func WithAutoDetectDelimiter() ReaderOption {
	return func(o *readerOptions) {
		o.autoDetectDelimiter = true
	}
}
