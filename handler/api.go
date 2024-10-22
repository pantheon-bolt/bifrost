package handler

import (
	"fmt"
	"net/http"
)

type Api struct{}

func (a *Api) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CREATE API")
}

func (a *Api) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LIST API")
}

func (a *Api) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET BY ID API")
}

func (a *Api) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UPDATE BY ID API")
}

func (a *Api) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE API")
}
