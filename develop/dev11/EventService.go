package main

import (
	"errors"
	"sync"
	"time"
)

// интерфейс сервиса, бизнес-логика календаря
type EventService interface {
	//crud методы получения данных
	Create(Event) error
	Update(Event) error
	Delete(uint64) error
	//получаем по id
	Get(id uint64) (Event, error)
	//получаем по user_id и дате
	GetForDay(uid uint64, t time.Time) []Event
	GetForWeek(uid uint64, t time.Time) []Event
	GetForMonth(uid uint64, t time.Time) []Event
}

// данные будем хранить в памяти в мапе, без persist
type EventServiceImpl struct {
	m         *sync.RWMutex
	repo      map[uint64]Event
	idCounter uint64
}

func InitEventStorage() *EventServiceImpl {
	storage := &EventServiceImpl{
		m:         &sync.RWMutex{},
		repo:      make(map[uint64]Event),
		idCounter: 0,
	}
	return storage
}

//реализация сервиса, crud методов

func (es *EventServiceImpl) Create(e Event) error {
	es.m.Lock()
	defer es.m.Unlock()
	es.idCounter++
	if _, ok := es.repo[e.ID]; ok {
		return errors.New("event with this ID already exists")
	}
	e.ID = es.idCounter
	es.repo[es.idCounter] = e
	return nil
}

func (es *EventServiceImpl) Update(e Event) error {
	es.m.Lock()
	defer es.m.Unlock()
	if _, ok := es.repo[e.ID]; !ok {
		return errors.New("event with this ID not found")
	}
	es.repo[e.ID] = e
	return nil
}

func (es *EventServiceImpl) Delete(id uint64) error {
	es.m.Lock()
	defer es.m.Unlock()
	if _, ok := es.repo[id]; !ok {
		return errors.New("event with this ID not found")
	}
	delete(es.repo, id)
	return nil
}

func (es *EventServiceImpl) Get(id uint64) (Event, error) {
	es.m.RLock()
	defer es.m.RUnlock()
	event, ok := es.repo[id]
	if !ok {
		return Event{}, errors.New("event not found")
	}
	return event, nil
}

func (es *EventServiceImpl) GetForDay(uid uint64, t time.Time) []Event {
	es.m.RLock()
	defer es.m.RUnlock()
	out := make([]Event, 0)
	for _, e := range es.repo {
		if e.UserID == uid {
			if e.Date == t {
				out = append(out, e)
			}
		}
	}
	return out
}

func (es *EventServiceImpl) GetForWeek(uid uint64, t time.Time) []Event {
	es.m.RLock()
	defer es.m.RUnlock()
	out := make([]Event, 0)
	for _, e := range es.repo {
		if e.UserID == uid {
			if e.Date.After(t) && e.Date.Before(t.AddDate(0, 0, 7)) {
				out = append(out, e)
			}
		}
	}
	return out
}

func (es *EventServiceImpl) GetForMonth(uid uint64, t time.Time) []Event {
	es.m.RLock()
	defer es.m.RUnlock()
	out := make([]Event, 0)
	for _, e := range es.repo {
		if e.UserID == uid {
			if e.Date.After(t) && e.Date.Before(t.AddDate(0, 1, 0)) {
				out = append(out, e)
			}
		}
	}
	return out
}
