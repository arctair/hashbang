package v1

// NamedTagListService ...
type NamedTagListService interface {
	Create(bucket string, namedTagList NamedTagList) (*NamedTagList, error)
	CreateOld(namedTagList NamedTagList) (*NamedTagList, error)
}

type namedTagListService struct {
	namedTagListRepository NamedTagListRepository
	uuidGenerator          UUIDGenerator
}

func (s *namedTagListService) Create(bucket string, namedTagList NamedTagList) (*NamedTagList, error) {
	return nil, nil
}

func (s *namedTagListService) CreateOld(namedTagList NamedTagList) (*NamedTagList, error) {
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
