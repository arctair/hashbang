package v1

// PostRepository ...
type PostRepository interface {
	FindAll() []Post
	Create(post Post)
	DeleteAll()
}

type postRepository struct {
	posts []Post
}

func (r *postRepository) FindAll() []Post {
	return r.posts
}

func (r *postRepository) Create(post Post) {
	r.posts = append(r.posts, post)
}

func (r *postRepository) DeleteAll() {
	r.posts = []Post{}
}

// NewPostRepository ...
func NewPostRepository() PostRepository {
	return &postRepository{
		posts: []Post{},
	}
}
