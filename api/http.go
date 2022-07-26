package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/navisot/go-url-shortener/serializer/json"
	"github.com/navisot/go-url-shortener/serializer/msgpack"
	"github.com/navisot/go-url-shortener/shortener"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

// RedirectHandler is the interface that holds the available methods
type RedirectHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
}

// handler inits the handler
type handler struct {
	redirectService shortener.RedirectService
}

// NewHandler returns a RedirectHandler interface
func NewHandler(redirectService shortener.RedirectService) RedirectHandler {
	return &handler{redirectService: redirectService}
}

// Get redirects to the original url based on the provided URL code
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect.URL, http.StatusMovedPermanently)
}

// Post stores the short URL into database
func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	redirect, err := h.serializer(contentType).Decode(requestBody)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.redirectService.Store(redirect)

	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responseBody, err := h.serializer(contentType).Encode(redirect)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	setupResponse(w, contentType, responseBody, http.StatusCreated)
}

// serializer is responsible to serialize to msgpack or json
func (h *handler) serializer(contentType string) shortener.RedirectSerializer {
	if contentType == "aplication/x-msgpack" {
		return &msgpack.Redirect{}
	}

	return &json.Redirect{}
}

// setupResponse is responsible for setting the response based on the chosen serializer
func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)

	if err != nil {
		log.Println(err)
	}
}
