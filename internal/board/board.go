package board

import (
	"fmt"
	"rory-pearson/pkg/log"
	"sort"
	"time"
)

var instance *Board

type Config struct {
	Log log.Log
}

type Board struct {
	Log   log.Log
	Posts []BoardPost
}

type PaginatedPosts struct {
	Posts      []BoardPost `json:"posts"`
	TotalPosts int         `json:"total_posts"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}

const DefaultPageSize = 10 // Set a reasonable default for page size
const MaxPageSize = 100    // Limit the maximum page size to prevent abuse

func Initialize(cfg Config) error {
	var board = &Board{
		Log: cfg.Log,
	}

	instance = board

	// GenerateRandomPosts(100)

	instance.Log.Info().Msg("Board initialized")

	return nil
}

func GetPosts(page, pageSize int) (*PaginatedPosts, error) {
	if instance == nil {
		return nil, ErrorInstanceNotInitialized
	}

	// Ensure pageSize is within limits
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	// Ensure valid pagination values
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}

	// Sort posts by CreatedAt in descending order (latest first)
	sortedPosts := make([]BoardPost, len(instance.Posts))
	copy(sortedPosts, instance.Posts)

	sort.Slice(sortedPosts, func(i, j int) bool {
		return sortedPosts[i].CreatedAt > sortedPosts[j].CreatedAt
	})

	// Calculate pagination indices
	startIndex := (page - 1) * pageSize
	if startIndex >= len(sortedPosts) {
		return &PaginatedPosts{
			Posts:      []BoardPost{},
			TotalPosts: len(sortedPosts),
			Page:       page,
			PageSize:   pageSize,
		}, nil
	}

	endIndex := startIndex + pageSize
	if endIndex > len(sortedPosts) {
		endIndex = len(sortedPosts)
	}

	return &PaginatedPosts{
		Posts:      sortedPosts[startIndex:endIndex],
		TotalPosts: len(sortedPosts),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func CreatePost(post CreateBoardPost) error {
	if instance == nil {
		return ErrorInstanceNotInitialized
	}

	if post.Title == "" {
		return fmt.Errorf("title is required")
	}
	if post.Content == "" {
		return fmt.Errorf("content is required")
	}
	if post.Author == "" {
		return fmt.Errorf("author is required")
	}

	id := len(instance.Posts) + 1
	createdAt := time.Now().Unix()

	boardPost := BoardPost{
		ID:        id,
		Title:     post.Title,
		Content:   post.Content,
		Author:    post.Author,
		CreatedAt: createdAt,
	}

	instance.Posts = append(instance.Posts, boardPost)

	instance.Log.Info().Msgf("Created post: %s", boardPost.Title)

	return nil
}

func GenerateRandomPosts(amount int) error {
	if instance == nil {
		return ErrorInstanceNotInitialized
	}

	for i := 0; i < amount; i++ {
		id := len(instance.Posts) + 1
		createdAt := time.Now().Unix()

		boardPost := BoardPost{
			ID:        id,
			Title:     fmt.Sprintf("Post %d", id),
			Content:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			Author:    "Anonymous",
			CreatedAt: createdAt,
		}

		instance.Posts = append(instance.Posts, boardPost)
	}

	instance.Log.Info().Msgf("Generated %d random posts", amount)

	return nil
}
