package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Movie struct {
	Id         int
	Synopsis   string
	Comments   string
	Year       string
	Rate       int
	Actors     Actors
	Director   string
	PosterPath string
	Task       Task
}

type Actors struct {
	Actors []string
}

func (a Actors) Value() (driver.Value, error) {
	if len(a.Actors) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

func (a *Actors) Scan(value interface{}) error {
	if value == nil {
		return json.Unmarshal([]byte("{\"Actors\": []}"), &a)
	}
	v, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(v, &a)
}
