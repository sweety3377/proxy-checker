package proxy_service

// Implementation of repository
type Repository interface {
	StartChecker(proxiesList []string)
}

type ProxyService struct {
	repository Repository
}

func New(repository Repository) *ProxyService {
	return &ProxyService{repository: repository}
}

func (c *ProxyService) StartChecker(proxiesList []string) {
	c.repository.StartChecker(proxiesList)
}
