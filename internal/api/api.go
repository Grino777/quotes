package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Grino777/quotes/internal/domain/models"
	"github.com/Grino777/quotes/internal/interfaces"
)

const (
	InternalError    = "internal server error"
	MethodNotAllowed = "method not allowed"
)

type API struct {
	logger  *slog.Logger
	service interfaces.Service
}

func NewApi(
	log *slog.Logger,
	service interfaces.Service,
) *API {
	return &API{logger: log, service: service}
}

func (a *API) HomeRoute(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(map[string]string{"result": "Quotes API"}); err != nil {
		http.Error(w, InternalError, http.StatusInternalServerError)
		return
	}
}

func (a *API) NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not found", http.StatusNotFound)
}

func (a *API) NotFoundFallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.NotFound(w, r)
	} else {
		a.HomeRoute(w, r)
	}

}

func (a *API) AllQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	if author != "" {
		data, err := a.service.FilterQuotes(r.Context(), author)
		if err != nil {
			http.Error(w, InternalError, http.StatusInternalServerError)
		}
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, InternalError, http.StatusInternalServerError)
			return
		}
	} else {
		data, err := a.service.GetQuotes(r.Context())
		if err != nil {
			http.Error(w, InternalError, http.StatusInternalServerError)
		}
		_, err = w.Write(data)
		if err != nil {
			http.Error(w, InternalError, http.StatusInternalServerError)
			return
		}
	}
}

func (a *API) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var q models.Quote

	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, InternalError, http.StatusInternalServerError)
		return
	}

	if err := q.Validate(); err != nil {
		http.Error(w, "fields author and quote is required", http.StatusBadRequest)
		return
	}

	res, err := a.service.CreateQuote(r.Context(), q)
	if err != nil {
		if err := json.NewEncoder(w).Encode(
			map[string][]byte{"result": res}); err != nil {
			http.Error(w, InternalError, http.StatusInternalServerError)
		}
	}

	_, err = w.Write(res)
	if err != nil {
		http.Error(w, InternalError, http.StatusInternalServerError)
		return
	}
}

func (a *API) RandomQuote(w http.ResponseWriter, r *http.Request) {
	data, err := a.service.GetRandomQuote(r.Context())
	if err != nil {
		http.Error(w, InternalError, http.StatusInternalServerError)
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, InternalError, http.StatusInternalServerError)
		return
	}
}

func (a *API) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 || pathParts[0] != "quotes" {
		http.Error(w, "invalid quote ID", http.StatusBadRequest)
		return
	}
	idStr := pathParts[1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid quote ID", http.StatusBadRequest)
		return
	}

	res, err := a.service.DeleteQuote(r.Context(), id)
	if err != nil {
		http.Error(w, InternalError, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		http.Error(w, InternalError, http.StatusInternalServerError)
		return
	}
}
