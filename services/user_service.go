package services

import (
    "errors"
    "sync"
    "sync/atomic"

    "go-microservice/models"
)

var ErrNotFound = errors.New("user not found")

type UserService struct {
    mu     sync.RWMutex
    users  map[int]models.User
    nextID atomic.Int64
}

func NewUserService() *UserService {
    s := &UserService{
        users: make(map[int]models.User),
    }
    s.nextID.Store(0)
    return s
}

func (s *UserService) List() []models.User {
    s.mu.RLock()
    defer s.mu.RUnlock()
    res := make([]models.User, 0, len(s.users))
    for _, u := range s.users {
        res = append(res, u)
    }
    return res
}

func (s *UserService) Get(id int) (models.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    u, ok := s.users[id]
    if !ok {
        return models.User{}, ErrNotFound
    }
    return u, nil
}

func (s *UserService) Create(u models.User) models.User {
    id := int(s.nextID.Add(1))
    u.ID = id
    s.mu.Lock()
    s.users[id] = u
    s.mu.Unlock()
    return u
}

func (s *UserService) Update(id int, u models.User) (models.User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    if _, ok := s.users[id]; !ok {
        return models.User{}, ErrNotFound
    }
    u.ID = id
    s.users[id] = u
    return u, nil
}

func (s *UserService) Delete(id int) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    if _, ok := s.users[id]; !ok {
        return ErrNotFound
    }
    delete(s.users, id)
    return nil
}
