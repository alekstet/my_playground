package code

import "github.com/alekstet/my_playground/config"

type CodeStore struct {
	Config *config.PlayConfig
	Coder  Coder
}

func NewCodeStore(cfg config.PlayConfig) (*CodeStore, error) {
	store := &CodeStore{
		Config: &cfg,
	}
	store.Coder = store

	return store, nil
}
