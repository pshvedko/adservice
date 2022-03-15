package service

type Storage interface {
	Search([]string, []string, int, int, []string) (interface{}, error)
	Store(float32, string, string, []string) (interface{}, error)
}

type Service struct {
	Storage
}

func New(storage Storage) *Service {
	return &Service{
		Storage: storage,
	}
}

func (s *Service) Read(ids, fields []string, limit, offset int, sorts []string) (interface{}, error) {
	return s.Storage.Search(ids, fields, limit, offset, sorts)
}

func (s *Service) Create(price float32, subject string, description string, photos []string) (interface{}, error) {
	return s.Storage.Store(price, subject, description, photos)
}
