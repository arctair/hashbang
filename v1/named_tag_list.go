package v1

// NamedTagList ...
type NamedTagList struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
