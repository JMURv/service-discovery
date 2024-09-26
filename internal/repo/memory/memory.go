package memory

import (
	"context"
	"github.com/JMURv/service-discovery/internal/repo"
	"sync"
)

type Repository struct {
	sync.RWMutex
	services map[string][]string
	rrIndex  map[string]int
}

func New() *Repository {
	return &Repository{
		services: make(map[string][]string),
		rrIndex:  make(map[string]int),
	}
}

func (r *Repository) Register(_ context.Context, name, addr string) error {
	r.Lock()
	defer r.Unlock()

	for _, registered := range r.services[name] {
		if registered == addr {
			return repo.ErrAlreadyExists
		}
	}

	r.services[name] = append(r.services[name], addr)
	return nil
}

func (r *Repository) Deregister(_ context.Context, name, addr string) error {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.services[name]; !exists {
		return repo.ErrNotFound
	}

	found := false
	newAddrs := make([]string, 0, len(r.services[name]))
	for _, registered := range r.services[name] {
		if registered == addr {
			found = true
			continue
		}
		newAddrs = append(newAddrs, registered)
	}

	if !found {
		return repo.ErrNotFound
	}

	if len(newAddrs) > 0 {
		r.services[name] = newAddrs
	} else {
		delete(r.services, name)
		delete(r.rrIndex, name)
	}

	return nil
}

func (r *Repository) FindServiceByName(_ context.Context, name string) (string, error) {
	r.RLock()
	defer r.RUnlock()

	if service, exists := r.services[name]; exists {
		currentIndex := r.rrIndex[name]
		selectedAddr := service[currentIndex]

		r.rrIndex[name] = (currentIndex + 1) % len(service)
		return selectedAddr, nil
	}

	return "", repo.ErrNotFound
}

func (r *Repository) ListServices(_ context.Context) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	names := make([]string, 0, len(r.services))
	for name := range r.services {
		names = append(names, name)
	}

	return names, nil
}

func (r *Repository) ListAddrs(_ context.Context, name string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if addrs, exists := r.services[name]; exists {
		return addrs, nil
	}

	return []string{}, repo.ErrNotFound
}
