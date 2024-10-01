package model

type Post struct {
	baseModel
	Title   string
	Desc    string
	Content string
}
