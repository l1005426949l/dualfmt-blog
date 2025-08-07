package parser

import (
	"bytes"
	"testing"
)

func TestParse(t *testing.T) {
	input := []byte("## Title\n\nThis is a paragraph.\n\n<script>alert('XSS');</script>\n\n* List item 1\n* List item 2")
	expected := `<h2>Title</h2>

<p>This is a paragraph.</p>



<ul>
<li>List item 1</li>
<li>List item 2</li>
</ul>
`

	output := Parse(input)

	if !bytes.Equal([]byte(expected), output) {
		t.Errorf("Parse() = %q, want %q", output, expected)
	}
}
