package v1

// NamedTagListService ...
type NamedTagListService interface {
	Create(namedTagList NamedTagList) (*NamedTagList, error)
}

type namedTagListService struct {
	namedTagListRepository NamedTagListRepository
}

func (s *namedTagListService) Create(namedTagList NamedTagList) (*NamedTagList, error) {
	return &namedTagList, s.namedTagListRepository.Create(namedTagList)
}

// NewNamedTagListService ...
func NewNamedTagListService(namedTagListRepository NamedTagListRepository) NamedTagListService {
	return &namedTagListService{
		namedTagListRepository,
	}
}
