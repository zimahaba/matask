package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Project struct {
	Id            int
	Description   string
	Progress      int
	Task          Task
	DynamicFields DynamicFields
}

type DynamicFields map[string]interface{}

func (d DynamicFields) Value() (driver.Value, error) {
	if len(d) == 0 {
		return nil, nil
	}
	return json.Marshal(d)
}

func (a *DynamicFields) Scan(value interface{}) error {
	if value == nil {
		return json.Unmarshal([]byte("{}"), &a)
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type ProjectFilter struct {
	Name          string
	Progress1     int
	Progress2     int
	Page          int
	Size          int
	SortField     string
	SortDirection string
	UserId        int
}

type ProjectPageResult struct {
	Projects      []ProjectProjection
	Page          int
	Size          int
	TotalPages    int
	TotalElements int
}

type ProjectProjection struct {
	Id       int
	Name     string
	Progress int
}
