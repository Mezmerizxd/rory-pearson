package board

import (
	"testing"

	"rory-pearson/pkg/log"

	"github.com/stretchr/testify/assert"
)

// Helper function to initialize the Board with a logger
func getLogger() log.Log {
	return log.New(log.Config{
		ID:            "board_test",
		ConsoleOutput: false,
		FileOutput:    false,
		StoragePath:   "",
	})
}

func TestInitialize(t *testing.T) {
	// Mock configuration
	config := Config{
		Log: getLogger(),
	}

	// Initialize the Board
	err := Initialize(config)
	assert.NoError(t, err, "failed to initialize board")
	assert.NotNil(t, instance, "instance should not be nil")
}

func TestCreatePost(t *testing.T) {
	// Mock configuration
	config := Config{
		Log: getLogger(),
	}

	// Initialize the Board
	err := Initialize(config)
	assert.NoError(t, err, "failed to initialize board")

	// Create a valid post
	post := CreateBoardPost{
		Title:   "Test Title",
		Content: "Test Content",
		Author:  "Test Author",
	}
	err = CreatePost(post)
	assert.NoError(t, err, "failed to create post")

	// Check if the post was added
	assert.Equal(t, 1, len(instance.Posts), "there should be one post")
	assert.Equal(t, post.Title, instance.Posts[0].Title, "titles should match")
	assert.Equal(t, post.Content, instance.Posts[0].Content, "content should match")
	assert.Equal(t, post.Author, instance.Posts[0].Author, "author should match")
}

func TestCreatePostValidation(t *testing.T) {
	// Mock configuration
	config := Config{
		Log: getLogger(),
	}

	// Initialize the Board
	err := Initialize(config)
	assert.NoError(t, err, "failed to initialize board")

	// Create a post with missing fields
	post := CreateBoardPost{
		Title:   "",
		Content: "",
		Author:  "",
	}
	err = CreatePost(post)
	assert.Error(t, err, "post validation should fail")
}

func TestGetPosts(t *testing.T) {
	// Mock configuration
	config := Config{
		Log: getLogger(),
	}

	// Initialize the Board
	err := Initialize(config)
	assert.NoError(t, err, "failed to initialize board")

	// Generate random posts
	GenerateRandomPosts(25)

	// Fetch paginated posts
	posts, err := GetPosts(1, 10)
	assert.NoError(t, err, "failed to fetch posts")
	assert.Equal(t, 10, len(posts.Posts), "should return 10 posts")
	assert.Equal(t, 25, posts.TotalPosts, "total posts should be 25")

	// Check that the posts are sorted by CreatedAt in descending order
	for i := 0; i < len(posts.Posts)-1; i++ {
		assert.True(t, posts.Posts[i].CreatedAt >= posts.Posts[i+1].CreatedAt, "posts should be sorted by CreatedAt")
	}
}

func TestGenerateRandomPosts(t *testing.T) {
	// Mock configuration
	config := Config{
		Log: getLogger(),
	}

	// Initialize the Board
	err := Initialize(config)
	assert.NoError(t, err, "failed to initialize board")

	// Generate 5 random posts
	err = GenerateRandomPosts(5)
	assert.NoError(t, err, "failed to generate random posts")
	assert.Equal(t, 5, len(instance.Posts), "should have 5 posts")
}

// Test for empty GetPosts case
func TestGetPostsEmpty(t *testing.T) {
	// Mock configuration
	config := Config{
		Log: getLogger(),
	}

	// Initialize the Board
	err := Initialize(config)
	assert.NoError(t, err, "failed to initialize board")

	// Fetch paginated posts when no posts exist
	posts, err := GetPosts(1, 10)
	assert.NoError(t, err, "should not error on empty board")
	assert.Equal(t, 0, len(posts.Posts), "should return zero posts")
	assert.Equal(t, 0, posts.TotalPosts, "total posts should be zero")
}
