package factory

import "errors"

// ErrUnknownRestAPIEntityType signals that an unknown REST API entity type has been provided
var ErrUnknownRestAPIEntityType = errors.New("unknown REST API entity type")
