package board

import (
	"fmt"
	"rory-pearson/pkg/log"
	"rory-pearson/plugins"
	"sort"
	"time"
)

// Singleton instance of the Board
var instance *Board

// Config struct to initialize the Board with necessary dependencies
type Config struct {
	Log log.Log // Logger instance for the Board
}

// Board represents a collection of posts and includes logging
type Board struct {
	Log   log.Log     // Logger instance for the Board
	Posts []BoardPost // Slice to store posts
}

// PaginatedPosts holds the posts along with pagination metadata
type PaginatedPosts struct {
	Posts      []BoardPost `json:"posts"`       // Slice of posts in the current page
	TotalPosts int         `json:"total_posts"` // Total number of posts
	Page       int         `json:"page"`        // Current page number
	PageSize   int         `json:"page_size"`   // Number of posts per page
}

// Constants for pagination
const DefaultPageSize = 10 // Default page size if not specified
const MaxPageSize = 100    // Maximum page size to prevent excessive data retrieval

// Initialize sets up the Board instance with a given configuration
// It logs the initialization process
func Initialize(cfg Config) error {
	var board = &Board{
		Log: cfg.Log,
	}

	// Assign board to the global instance
	instance = board

	// Log successful initialization
	instance.Log.Info().Msg("Board initialized")

	plugins.GetInstance().Commands.RegisterCommand(plugins.Command{
		ID:          "create_post",
		Name:        "Create Posts",
		Description: "Create a new post with a title, content, and author",
		ArgTypes:    []string{"string", "string", "string"}, // Specify argument types
		Function: func(args ...any) error {
			title := args[0].(string)
			content := args[1].(string)
			author := args[2].(string)

			err := CreatePost(CreateBoardPost{
				Title:   title,
				Content: content,
				Author:  author,
			})

			if err != nil {
				return err
			}

			return nil
		},
	})

	return nil
}

// GetPosts retrieves paginated posts based on the specified page and pageSize
// It ensures that posts are returned in reverse chronological order based on CreatedAt
func GetPosts(page, pageSize int) (*PaginatedPosts, error) {
	// Check if the Board has been initialized
	if instance == nil {
		return nil, ErrorInstanceNotInitialized
	}

	// Ensure pageSize is within allowed limits
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

	// Sort the posts by CreatedAt in descending order
	sortedPosts := make([]BoardPost, len(instance.Posts))
	copy(sortedPosts, instance.Posts)
	sort.Slice(sortedPosts, func(i, j int) bool {
		return sortedPosts[i].CreatedAt > sortedPosts[j].CreatedAt
	})

	// Calculate pagination indices
	startIndex := (page - 1) * pageSize
	if startIndex >= len(sortedPosts) {
		// If the start index is beyond the available posts, return empty result
		return &PaginatedPosts{
			Posts:      []BoardPost{},
			TotalPosts: len(sortedPosts),
			Page:       page,
			PageSize:   pageSize,
		}, nil
	}

	// Ensure the end index doesn't exceed the available posts
	endIndex := startIndex + pageSize
	if endIndex > len(sortedPosts) {
		endIndex = len(sortedPosts)
	}

	// Return the paginated result
	return &PaginatedPosts{
		Posts:      sortedPosts[startIndex:endIndex],
		TotalPosts: len(sortedPosts),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// CreatePost adds a new post to the Board
// It validates that Title, Content, and Author are non-empty before adding the post
func CreatePost(post CreateBoardPost) error {
	// Check if the Board has been initialized
	if instance == nil {
		return ErrorInstanceNotInitialized
	}

	// Validate that required fields are present
	if post.Title == "" {
		return fmt.Errorf("title is required")
	}
	if post.Content == "" {
		return fmt.Errorf("content is required")
	}
	if post.Author == "" {
		return fmt.Errorf("author is required")
	}

	// Assign a new ID and creation timestamp for the post
	id := len(instance.Posts) + 1
	createdAt := time.Now().Unix()

	// Create a new BoardPost
	boardPost := BoardPost{
		ID:        id,
		Title:     post.Title,
		Content:   post.Content,
		Author:    post.Author,
		CreatedAt: createdAt,
	}

	// Append the post to the Board's list of posts
	instance.Posts = append(instance.Posts, boardPost)

	// Log the creation of the post
	instance.Log.Info().Msgf("Created post: %s", boardPost.Title)

	return nil
}

// GenerateRandomPosts generates a given number of posts for testing or seeding
// It creates generic posts with placeholder content and adds them to the Board
func GenerateRandomPosts(amount int) error {
	// Check if the Board has been initialized
	if instance == nil {
		return ErrorInstanceNotInitialized
	}

	// Loop to generate and append random posts
	for i := 0; i < amount; i++ {
		id := len(instance.Posts) + 1
		createdAt := time.Now().Unix()

		// Create a new BoardPost with placeholder content
		boardPost := BoardPost{
			ID:        id,
			Title:     fmt.Sprintf("Post %d", id),
			Content:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			Author:    "Anonymous",
			CreatedAt: createdAt,
		}

		// Append the post to the Board's list of posts
		instance.Posts = append(instance.Posts, boardPost)
	}

	// Log the generation of random posts
	instance.Log.Info().Msgf("Generated %d random posts", amount)

	return nil
}
