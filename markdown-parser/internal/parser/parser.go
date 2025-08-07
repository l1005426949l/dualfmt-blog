package parser

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

// Parse converts a markdown byte slice to sanitized HTML.
func Parse(input []byte) []byte {
	// Use blackfriday to generate HTML from markdown.
	// CommonExtensions are a good set of default extensions.
	unsafeHTML := blackfriday.Run(input, blackfriday.WithExtensions(blackfriday.CommonExtensions))

	// Use bluemonday to sanitize the HTML.
	// The UGC policy is a good starting point for user-generated content.
	p := bluemonday.UGCPolicy()
	sanitizedHTML := p.SanitizeBytes(unsafeHTML)

	return sanitizedHTML
}
