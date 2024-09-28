package memory

import (
	"context"
	"github.com/JMURv/service-discovery/internal/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository(t *testing.T) {
	ctx := context.Background()
	r := New()

	t.Run("Register services", func(t *testing.T) {
		err := r.Register(ctx, "service1", "addr1")
		assert.NoError(t, err)

		err = r.Register(ctx, "service1", "addr1")
		assert.Equal(t, repo.ErrAlreadyExists, err)

		err = r.Register(ctx, "service1", "addr2")
		assert.NoError(t, err)
	})

	t.Run("Deregister services", func(t *testing.T) {
		err := r.Deregister(ctx, "non-existing-service", "non-existing-addr")
		assert.Equal(t, repo.ErrNotFound, err)

		err = r.Deregister(ctx, "service1", "addr1")
		assert.NoError(t, err)

		err = r.Deregister(ctx, "service1", "non-existing-addr")
		assert.Equal(t, repo.ErrNotFound, err)

		err = r.Deregister(ctx, "service1", "addr2")
		assert.NoError(t, err)

		_, err = r.FindServiceByName(ctx, "service1")
		assert.Equal(t, repo.ErrNotFound, err)
	})

	t.Run("Find service with round-robin", func(t *testing.T) {
		r.Register(ctx, "service2", "addr3")
		r.Register(ctx, "service2", "addr4")

		addr, err := r.FindServiceByName(ctx, "service2")
		assert.NoError(t, err)
		assert.Equal(t, "addr3", addr)

		addr, err = r.FindServiceByName(ctx, "service2")
		assert.NoError(t, err)
		assert.Equal(t, "addr4", addr)

		addr, err = r.FindServiceByName(ctx, "service2")
		assert.NoError(t, err)
		assert.Equal(t, "addr3", addr)
	})

	t.Run("List services", func(t *testing.T) {
		r.Register(ctx, "service3", "addr5")
		r.Register(ctx, "service4", "addr6")

		services, err := r.ListServices(ctx)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"service2", "service3", "service4"}, services)
	})

	t.Run("List addresses for a service", func(t *testing.T) {
		addrs, err := r.ListAddrs(ctx, "service2")
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"addr3", "addr4"}, addrs)

		addrs, err = r.ListAddrs(ctx, "non-existing-service")
		assert.Equal(t, repo.ErrNotFound, err)
		assert.Empty(t, addrs)
	})

	t.Run("Deactivate service", func(t *testing.T) {
		r.Register(ctx, "service5", "addr7")

		err := r.DeactivateSvc(ctx, "service5", "addr7")
		assert.NoError(t, err)

		addr, err := r.FindServiceByName(ctx, "service5")
		assert.Equal(t, repo.ErrNotFound, err)
		assert.Empty(t, addr)

		err = r.DeactivateSvc(ctx, "non-existing-service", "non-existing-addr")
		assert.Equal(t, repo.ErrNotFound, err)

		err = r.DeactivateSvc(ctx, "service5", "non-existing-addr")
		assert.Equal(t, repo.ErrNotFound, err)

		err = r.Register(ctx, "service5", "addr8")
		assert.NoError(t, err)

		err = r.DeactivateSvc(ctx, "service5", "addr8")
		assert.NoError(t, err)

		addr, err = r.FindServiceByName(ctx, "service5")
		assert.Equal(t, repo.ErrNotFound, err)
	})

	t.Run("Activate service", func(t *testing.T) {
		r.Register(ctx, "service6", "addr9")
		err := r.DeactivateSvc(ctx, "service6", "addr9")
		assert.NoError(t, err)

		addr, err := r.FindServiceByName(ctx, "service6")
		assert.Equal(t, repo.ErrNotFound, err)
		assert.Empty(t, addr)

		err = r.ActivateSvc(ctx, "service6", "addr9")
		assert.NoError(t, err)

		addr, err = r.FindServiceByName(ctx, "service6")
		assert.NoError(t, err)
		assert.Equal(t, "addr9", addr)

		err = r.ActivateSvc(ctx, "non-existing-service", "non-existing-addr")
		assert.Equal(t, repo.ErrNotFound, err)

		err = r.ActivateSvc(ctx, "service6", "non-existing-addr")
		assert.Equal(t, repo.ErrNotFound, err)
	})

	t.Run("Close", func(t *testing.T) {
		err := r.Close()
		assert.Nil(t, err)
	})

}
