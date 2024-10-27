package model

import (
	"time"

	"github.com/google/uuid"
)

type Api struct {
	ApiID           string       `json:"api_id"`
	CreatedAt       *time.Time   `json:"created_at"`
	UpdatedAt       *time.Time   `json:"updated_at"`
	Headers         []Header     `json:"headers"`
	QueryParams     []QueryParam `json:"query_params"`
	PathParams      []PathParam  `json:"path_params"`
	Target          string       `json:"target"`
	RootDomain      string       `json:"root_domain"`
	Domain          string       `json:"domain"`
	Protocol        string       `json:"protocol"`
	ProtocolVersion string       `json:"protocol_version"`
	Port            string       `json:"port"`
	Method          string       `json:"method"`
	Path            string       `json:"path"`
	Body            string       `json:"body"`
}

type Header struct {
	HeaderID  uuid.UUID  `json:"header_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Name      string     `json:"name"`
	Value     string     `json:"value"`
}

type QueryParam struct {
	QueryParamID uuid.UUID  `json:"query_param_id"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	Name         string     `json:"name"`
	Value        string     `json:"value"`
}

type PathParam struct {
	PathParamID uuid.UUID  `json:"path_param_id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Name        string     `json:"name"`
	Value       string     `json:"value"`
}
