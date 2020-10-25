package v1

// PostRepository ...
type PostRepository interface {
	FindAll() []Post
}

type postRepository struct {
}

func (r *postRepository) FindAll() []Post {
	return []Post{
		{
			ImageUri: "https://images.unsplash.com/photo-1603316851229-26637b4bd1b8?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1400&q=80",
			Tags: []string{
				"#windy",
				"#tdd",
			},
		},
	}
}

// NewPostRepository ...
func NewPostRepository() PostRepository {
	return &postRepository{}
}
