package backup

import "context"

// Creator интерфейс создателя резервной копии.
type Creator interface {
	// Create создать резервную копию.
	Create(ctx context.Context) error
}
