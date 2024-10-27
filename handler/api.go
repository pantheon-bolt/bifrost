package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/pantheon-bolt/bifrost/model"
	"github.com/pantheon-bolt/bifrost/repository/api"
)

type Repo interface {
	Insert(ctx context.Context, api Api) error
	FindByID(ctx context.Context, id string) (Api, error)
	DeleteByID(ctx context.Context, id string) error
	Update(ctx context.Context, api Api) error
	FindAll(ctx context.Context, page api.FindAllPage) (api.FindResult, error)
}

type Api struct {
	Repo *api.RedisRepo
}

func (a *Api) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[@] CREATE API")
	var body struct {
		Headers         []model.Header     `json:"headers"`
		QueryParams     []model.QueryParam `json:"query_params"`
		PathParams      []model.PathParam  `json:"path_params"`
		Target          string             `json:"target"`
		RootDomain      string             `json:"root_domain"`
		Domain          string             `json:"domain"`
		Protocol        string             `json:"protocol"`
		ProtocolVersion string             `json:"protocol_version"`
		Port            string             `json:"port"`
		Method          string             `json:"method"`
		Path            string             `json:"path"`
		Body            string             `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	api := model.Api{
		ApiID:           strconv.FormatUint(rand.Uint64(), 10),
		CreatedAt:       &now,
		UpdatedAt:       &now,
		Headers:         body.Headers,
		QueryParams:     body.QueryParams,
		PathParams:      body.PathParams,
		Target:          body.Target,
		RootDomain:      body.RootDomain,
		Domain:          body.Domain,
		Protocol:        body.Protocol,
		ProtocolVersion: body.ProtocolVersion,
		Port:            body.Port,
		Method:          body.Method,
		Path:            body.Path,
		Body:            body.Body,
	}

	err := a.Repo.Insert(r.Context(), api)
	if err != nil {
		fmt.Println("[>] failed to insert: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(api)
	if err != nil {
		fmt.Println("[>] failed to marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (a *Api) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[@] LIST API")

	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	const decimal = 10
	const bitSize = 64

	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := a.Repo.FindAll(r.Context(), api.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("[>] failed to find all: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Entries []model.Api `json:"entries"`
		Next    uint64      `json:"next,omitempty"`
	}

	response.Entries = res.Apis
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("[>] failed to marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (a *Api) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[@] GET BY ID API")

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	apiID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	o, err := a.Repo.FindByID(r.Context(), apiID)

	if errors.Is(err, api.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("[>] failed to find by id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(o); err != nil {
		fmt.Println("[>] failed to marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Api) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[@] UPDATE BY ID API")
	var body struct {
		Body string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	apiID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	theApi, err := a.Repo.FindByID(r.Context(), apiID)
	if errors.Is(err, api.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("[>] failed to find by id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const case_a = "a"
	const case_b = "b"
	now := time.Now().UTC()

	switch body.Body {
	case case_a:
		if theApi.UpdatedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theApi.UpdatedAt = &now
	case case_b:
		if theApi.UpdatedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theApi.UpdatedAt = &now
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.Repo.Update(r.Context(), theApi)
	if err != nil {
		fmt.Println("[>] failed to update: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(theApi); err != nil {
		fmt.Println("[>] failed to marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Api) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[@] DELETE API")

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	apiID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.Repo.DeleteByID(r.Context(), apiID)
	if errors.Is(err, api.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("[>] failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
