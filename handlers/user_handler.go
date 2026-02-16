package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"

	"go-microservice/models"
	"go-microservice/services"
)

var emailRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type UserHandler struct {
	svc    *services.UserService
	audit  *Audit
	notify *Notifier
}

func NewUserHandler(svc *services.UserService, audit *Audit, notify *Notifier) *UserHandler {
	return &UserHandler{svc: svc, audit: audit, notify: notify}
}

func (h *UserHandler) Register(r *mux.Router) {
	r.HandleFunc("/api/users", h.List).Methods(http.MethodGet)
	r.HandleFunc("/api/users/{id}", h.Get).Methods(http.MethodGet)
	r.HandleFunc("/api/users", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/api/users/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/api/users/{id}", h.Delete).Methods(http.MethodDelete)
}

// ListUsers godoc
// @Summary Get all users
// @Description Returns list of all users
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Router /api/users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users := h.svc.List()
	writeJSON(w, http.StatusOK, users)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Returns a single user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {string} string
// @Router /api/users/{id} [get]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	u, err := h.svc.Get(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, u)
}

// CreateUser godoc
// @Summary Create user
// @Description Creates a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User payload"
// @Success 201 {object} models.User
// @Failure 400 {string} string
// @Router /api/users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if u.Name == "" || !emailRe.MatchString(u.Email) {
		http.Error(w, "validation error", http.StatusBadRequest)
		return
	}

	saved := h.svc.Create(u)

	go h.audit.Log("CREATE", saved.ID)
	go h.notify.Send("created", saved.ID)

	writeJSON(w, http.StatusCreated, saved)
}

// UpdateUser godoc
// @Summary Update user
// @Description Updates user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User payload"
// @Success 200 {object} models.User
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /api/users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if u.Name == "" || !emailRe.MatchString(u.Email) {
		http.Error(w, "validation error", http.StatusBadRequest)
		return
	}
	updated, err := h.svc.Update(id, u)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	go h.audit.Log("UPDATE", updated.ID)
	go h.notify.Send("updated", updated.ID)

	writeJSON(w, http.StatusOK, updated)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Deletes user by ID
// @Tags users
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {string} string
// @Router /api/users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.svc.Delete(id); err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	go h.audit.Log("DELETE", id)
	go h.notify.Send("deleted", id)

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
