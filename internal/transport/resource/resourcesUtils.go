package resource

import (
	"encoding/json"
	"time"
)

type IdResource struct {
	Id int
}

type Page struct {
	Number        int
	Size          int
	TotalPages    int
	TotalElements int
}

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return json.Marshal("")
	}

	return json.Marshal(d.Format(time.DateOnly))
}
