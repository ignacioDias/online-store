package router

import (
	"net/http"
	"sports-store/internal/cache"
	"sports-store/internal/database"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(db *database.Database, cache *cache.Cache) *Router {
	return &Router{
		mux: http.NewServeMux(),
		//Middleware and handlers initialization
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {

	return r.mux
}
