package dbs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"marketplace/pkg/models"
)

type SightModel struct {
	DB *sql.DB
}

func (m *SightModel) Insert(sight *models.Sight) error {
	stmt := `
        INSERT INTO astana.sight
        (name, address, phone_number, content_info, bus_numbers, working_days, working_hours, visited, image_url) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := m.DB.Exec(stmt, sight.Name, sight.Address, sight.PhoneNumber, sight.ContentInfo, sight.BusNumbers, sight.WorkingDays, sight.WorkingHours, sight.Visited, sight.ImageUrl)
	if err != nil {
		return err
	}

	return nil
}

func (m *SightModel) GetSightById(id string) ([]byte, error) {
	stmt := `SELECT id, name, address, phone_number, content_info, bus_numbers, working_days, working_hours, visited, image_url FROM astana.sight WHERE id = ?`

	sightRow := m.DB.QueryRow(stmt, id)

	s := &models.Sight{}

	err := sightRow.Scan(&s.Id, &s.Name, &s.Address, &s.PhoneNumber, &s.ContentInfo, &s.BusNumbers, &s.WorkingDays, &s.WorkingHours, &s.Visited, &s.ImageUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	convertedSight, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return convertedSight, nil
}

func (m *SightModel) GetAllSights() ([]byte, error) {
	stmt := `SELECT id, name, address, phone_number, content_info, bus_numbers, working_days, working_hours, visited, image_url FROM astana.sight`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sights := []*models.Sight{}

	for rows.Next() {
		s := &models.Sight{}
		err = rows.Scan(&s.Id, &s.Name, &s.Address, &s.PhoneNumber, &s.ContentInfo, &s.BusNumbers, &s.WorkingDays, &s.WorkingHours, &s.Visited, &s.ImageUrl)
		if err != nil {
			return nil, err
		}
		sights = append(sights, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	convertedSights, err := json.Marshal(sights)
	if err != nil {
		return nil, err
	}
	return convertedSights, nil
}

func (m *SightModel) DeleteSightById(id int) error {
	stmt := `DELETE FROM astana.sight WHERE id = ?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
