package board

import "errors"

var (
	ErrorInstanceNotInitialized = errors.New("board not initialized")
)

type BoardPost struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt int64  `json:"created_at"`
}

type CreateBoardPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}
