package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/pantheon-bolt/bifrost/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func apiIDKey(id string) string {
	return fmt.Sprintf("api:%s", id)
}

func (r *RedisRepo) Insert(ctx context.Context, api model.Api) error {
	data, err := json.Marshal(api)
	if err != nil {
		return fmt.Errorf("[>] failed to encode api: %w", err)
	}

	key := apiIDKey(api.ApiID)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("[>] failed to set: %w", err)
	}

	if err := txn.SAdd(ctx, "apis", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("[>] failed to add to apis set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("[>] failed to exec: %w", err)
	}

	return nil
}

var ErrNotExist = errors.New("[>] api does not exist")

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Api, error) {
	key := apiIDKey(strconv.FormatUint(id, 10))

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Api{}, ErrNotExist
	} else if err != nil {
		return model.Api{}, fmt.Errorf("[>] find api: %w", err)
	}

	var api model.Api
	err = json.Unmarshal([]byte(value), &api)
	if err != nil {
		return model.Api{}, fmt.Errorf("[>] failed to decode api json: %w", err)
	}

	return api, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := apiIDKey(strconv.FormatUint(id, 10))

	txn := r.Client.TxPipeline()

	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("[>] delete api: %w", err)
	}

	if err := txn.SRem(ctx, "apis", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("[>] failed to remove from apis set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("[>] failed to exec: %w", err)
	}

	return nil
}

func (r *RedisRepo) Update(ctx context.Context, api model.Api) error {
	data, err := json.Marshal(api)
	if err != nil {
		return fmt.Errorf("[>] failed to encode api: %w", err)
	}

	key := apiIDKey(api.ApiID)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("[>] update api: %w", err)
	}

	return nil
}

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Apis   []model.Api
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "apis", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("[>] failed to get api ids: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{
			Apis: []model.Api{},
		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("[>] failed to get apis: %w", err)
	}

	apis := make([]model.Api, len(xs))

	for i, x := range xs {
		x := x.(string)
		var api model.Api

		err := json.Unmarshal([]byte(x), &api)
		if err != nil {
			return FindResult{}, fmt.Errorf("[>] failed to decode api json: %w", err)
		}

		apis[i] = api
	}

	return FindResult{
		Apis:   apis,
		Cursor: cursor,
	}, nil
}
