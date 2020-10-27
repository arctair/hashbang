package v1

// Post ...
type Post struct {
	Id       string   `json:"id"`
	ImageUri string   `json:"imageUri"`
	Tags     []string `json:"tags"`
}
