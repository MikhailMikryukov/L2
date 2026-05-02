package service

import (
	"L2.18/entities"
	"errors"
	"strconv"
	"sync"
	"time"
)

var (
	// ErrInvalidDate некорректная дата
	ErrInvalidDate = errors.New("invalid date")
	// ErrNoSuchEvent несуществующее событие
	ErrNoSuchEvent = errors.New("no such event")
	// ErrInvalidID некорректный ID
	ErrInvalidID = errors.New("invalid user id")
	// ErrUserNotFound несуществующий пользователь
	ErrUserNotFound = errors.New("user not found")
)

// UserService сервис юзеров
type UserService struct {
	users map[int]entities.User
	mu    sync.Mutex
}

// NewUserService Создание экземпляра UserService
func NewUserService() *UserService {
	return &UserService{users: make(map[int]entities.User)}
}

// DeleteEvent удаление события
func (us *UserService) DeleteEvent(idStr string, dateStr string, event string) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	user, err := us.getUserByID(idStr)
	if err != nil {
		return err
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ErrInvalidDate
	}

	events := user.Events[date]
	var idx = -1

	for i := 0; i < len(events); i++ {
		if events[i] == event {
			idx = i
			break
		}
	}

	if idx == -1 {
		return ErrNoSuchEvent
	}

	events = append(events[0:idx], events[idx+1:]...)
	return nil
}

// UpdateEvent обновление события
func (us *UserService) UpdateEvent(idStr string, dateStr string, event string) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	user, err := us.getUserByID(idStr)
	if err != nil {
		return err
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ErrInvalidDate
	}

	events := user.Events[date]

	for i := 0; i < len(events); i++ {
		if events[i] == event {
			events[i] = event
		}
	}

	return nil
}

// GetEventsForUserID получение событий по ID
func (us *UserService) GetEventsForUserID(idStr string, dateStr string, days int) ([]string, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	user, err := us.getUserByID(idStr)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	var events []string

	for i := 0; i < days; i++ {
		event := user.Events[date.Add(time.Hour*24*time.Duration(i))]
		events = append(events, event...)
	}

	return events, nil
}

// CreateEvent создание события
func (us *UserService) CreateEvent(idStr string, dateStr string, event string) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	
	user, err := us.getUserByID(idStr)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			user = us.addUser(idStr)
		} else {
			return err
		}
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ErrInvalidDate
	}

	user.Events[date] = append(user.Events[date], event)
	return nil
}

func (us *UserService) getUserByID(idStr string) (entities.User, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return entities.User{}, ErrInvalidID
	}

	if user, ok := us.users[id]; ok {
		return user, nil
	}
	return entities.User{}, ErrUserNotFound
}

func (us *UserService) addUser(idStr string) entities.User {
	id, _ := strconv.Atoi(idStr)
	user := entities.User{
		ID:     id,
		Events: make(map[time.Time][]string),
	}
	us.users[id] = user
	return user
}
