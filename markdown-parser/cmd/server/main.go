package main

import (
	"fmt"
	"github.com/example/markdown-parser/internal/parser"
)

func main() {
	markdown := []byte("## Hello World\n\nThis is a test of the markdown parser.\n\n<script>alert('xss')</script>")
	html := parser.Parse(markdown)
	fmt.Println(string(html))
}
