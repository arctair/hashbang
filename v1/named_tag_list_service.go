package v1

// NamedTagListService ...
type NamedTagListService interface {
	Create(namedTagList NamedTagList) (*NamedTagList, error)
}

type namedTagListService struct {
	namedTagListRepository NamedTagListRepository
	uuidGenerator          UUIDGenerator
}

func (s *namedTagListService) Create(namedTagList NamedTagList) (*NamedTagList, error) {
	namedTagList.ID = s.uuidGenerator.Generate()
	return &namedTagList, s.namedTagListRepository.Create(namedTagList)
}

// NewNamedTagListService ...
func NewNamedTagListService(
	namedTagListRepository NamedTagListRepository,
	uuidGenerator UUIDGenerator,
) NamedTagListService {
	return &namedTagListService{
		namedTagListRepository,
		uuidGenerator,
	}
}
