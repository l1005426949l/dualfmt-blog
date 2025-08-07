package service

// Article is the Go struct corresponding to the GraphQL Article type.
type Article struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Version string `json:"version"`
	Format  string `json:"format"`
}
