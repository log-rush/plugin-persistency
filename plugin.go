package pluginPersistency

import "errors"

// TODO: implement actual plugin interface
type xxx struct{}

func NewPersistencyPlugin() (xxx, error) {
	return xxx{}, errors.New("not implemented")
}
