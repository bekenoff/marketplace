package dbs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"marketplace/pkg/models"
)

type EventModel struct {
	DB *sql.DB
}

func (m *EventModel) Insert(event *models.Events) error {
	stmt := `
        INSERT INTO astana.event
        (name, schedule, bus,price, description, image_url, category) 
        VALUES (?, ?, ?, ?, ?, ?, ?);`

	_, err := m.DB.Exec(stmt, event)
	if err != nil {
		return err
	}

	return nil
}

func (m *EventModel) GetEventById(id string) ([]byte, error) {
	stmt := `SELECT id, name,schedule,bus,price,description,image_url,category FROM astana.event WHERE id = ?`

	eventRow := m.DB.QueryRow(stmt, id)

	s := &models.Events{}

	err := eventRow.Scan(&s.Id, &s.Name, &s.Schedule, &s.Bus, &s.Price, &s.Description, &s.ImageUrl, &s.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	convertedEvent, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return convertedEvent, nil
}

func (m *EventModel) GetAllEvents() ([]byte, error) {
	stmt := `SELECT id, name, schedule, bus, price, description, image_url, category  FROM astana.event`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*models.Events{}

	for rows.Next() {
		s := &models.Events{}
		err = rows.Scan(&s.Id, &s.Name, &s.Schedule, &s.Bus, &s.Price, &s.Description, &s.ImageUrl, &s.Category)
		if err != nil {
			return nil, err
		}
		events = append(events, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	convertedEvents, err := json.Marshal(events)
	if err != nil {
		return nil, err
	}
	return convertedEvents, nil
}

func (m *SightModel) DeleteEventById(id int) error {
	stmt := `DELETE FROM astana.event WHERE id = ?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
