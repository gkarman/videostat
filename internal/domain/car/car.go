package car

import (
	"time"
)

type Car struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	events []any
}

func New(id string, name string) *Car {
	car := &Car{
		ID:   id,
		Name: name,
	}

	event := &Created{
		ID:   car.ID,
		Name: car.Name,
		At:   time.Now(),
	}
	car.addEvent(event)
	return car
}

func (c *Car) addEvent(e any) {
	c.events = append(c.events, e)
}

func (c *Car) PullEvents() []any {
	evs := c.events
	c.events = nil
	return evs
}
