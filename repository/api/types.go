package api

import (
	"errors"

	"github.com/pantheon-bolt/bifrost/model"
)

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Apis   []model.Api
	Cursor uint64
}

var ErrNotExist = errors.New("[>] api does not exist")
