package memory

import (
	"context"
	"github.com/JMURv/service-discovery/internal/repo"
	md "github.com/JMURv/service-discovery/pkg/model"
	"sync"
)

type Repository struct {
	sync.RWMutex
	services []md.Service
	rrIndex  map[string]int
}

func New() *Repository {
	return &Repository{
		services: make([]md.Service, 0, 10),
		rrIndex:  make(map[string]int),
	}
}

func (r *Repository) Close() error {
	r.Lock()
	defer r.Unlock()

	r.services = nil
	r.rrIndex = nil
	return nil
}

func (r *Repository) Register(_ context.Context, name, addr string) error {
	r.Lock()
	defer r.Unlock()

	for _, registered := range r.services {
		if registered.Address == addr {
			return repo.ErrAlreadyExists
		}
	}

	r.services = append(r.services, md.Service{Name: name, Address: addr, IsActive: true})
	return nil
}

func (r *Repository) Deregister(_ context.Context, name, addr string) error {
	r.Lock()
	defer r.Unlock()

	for i, v := range r.services {
		if v.Name == name && v.Address == addr {
			r.services = append(r.services[:i], r.services[i+1:]...)
			delete(r.rrIndex, name)
			return nil
		}
	}

	return repo.ErrNotFound
}

func (r *Repository) FindServiceByName(_ context.Context, name string) (string, error) {
	r.RLock()
	defer r.RUnlock()

	var availableServices []md.Service
	for _, svc := range r.services {
		if svc.Name == name && svc.IsActive {
			availableServices = append(availableServices, svc)
		}
	}

	if len(availableServices) == 0 {
		return "", repo.ErrNotFound
	}

	currentIndex := r.rrIndex[name]
	selectedSvc := availableServices[currentIndex]

	r.rrIndex[name] = (currentIndex + 1) % len(availableServices)
	return selectedSvc.Address, nil
}

func (r *Repository) ListServices(_ context.Context) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	namesMap := make(map[string]struct{})
	for _, svc := range r.services {
		namesMap[svc.Name] = struct{}{}
	}

	names := make([]string, 0, len(namesMap))
	for name := range namesMap {
		names = append(names, name)
	}

	return names, nil
}

func (r *Repository) ListAddrs(_ context.Context, name string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	var addrs []string
	for _, svc := range r.services {
		if svc.Name == name && svc.IsActive {
			addrs = append(addrs, svc.Address)
		}
	}

	if len(addrs) == 0 {
		return []string{}, repo.ErrNotFound
	}

	return addrs, nil
}

func (r *Repository) DeactivateSvc(_ context.Context, name, addr string) error {
	r.Lock()
	defer r.Unlock()

	for i, svc := range r.services {
		if svc.Name == name && svc.Address == addr {
			r.services[i].IsActive = false
			return nil
		}
	}

	return repo.ErrNotFound
}

func (r *Repository) ActivateSvc(_ context.Context, name, addr string) error {
	r.Lock()
	defer r.Unlock()

	for i, svc := range r.services {
		if svc.Name == name && svc.Address == addr {
			r.services[i].IsActive = true
			return nil
		}
	}

	return repo.ErrNotFound
}
