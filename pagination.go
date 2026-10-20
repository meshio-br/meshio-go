package meshio

import "context"

// Iter percorre todos os itens de uma listagem paginada, buscando as páginas sob demanda.
//
//	it := client.Businesses.ListAll(ctx, nil)
//	for it.Next() {
//		b := it.Current()
//		_ = b
//	}
//	if err := it.Err(); err != nil { ... }
type Iter[T any] struct {
	ctx     context.Context
	fetch   func(ctx context.Context, cursor string) ([]T, *string, error)
	buf     []T
	idx     int
	cur     T
	cursor  string
	started bool
	done    bool
	err     error
}

func newIter[T any](ctx context.Context, startCursor string, fetch func(ctx context.Context, cursor string) ([]T, *string, error)) *Iter[T] {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Iter[T]{ctx: ctx, fetch: fetch, cursor: startCursor}
}

// Next avança para o próximo item; devolve false ao terminar ou em erro (veja Err).
func (it *Iter[T]) Next() bool {
	for {
		if it.err != nil {
			return false
		}
		if it.idx < len(it.buf) {
			it.cur = it.buf[it.idx]
			it.idx++
			return true
		}
		if it.done {
			return false
		}
		if it.started && it.cursor == "" {
			it.done = true
			return false
		}
		it.started = true
		items, next, err := it.fetch(it.ctx, it.cursor)
		if err != nil {
			it.err = err
			return false
		}
		it.buf, it.idx = items, 0
		if next == nil || *next == "" {
			it.cursor = ""
			it.done = true
		} else {
			it.cursor = *next
		}
	}
}

// Current devolve o item atual (válido após Next devolver true).
func (it *Iter[T]) Current() T { return it.cur }

// Err devolve o erro que interrompeu a iteração, se houver.
func (it *Iter[T]) Err() error { return it.err }

// Collect lê todos os itens restantes (use com cuidado em listagens grandes).
func (it *Iter[T]) Collect() ([]T, error) {
	var out []T
	for it.Next() {
		out = append(out, it.Current())
	}
	return out, it.Err()
}
