package resource

import (
	"encoding/json"
	"fmt"
	"time"
)

type IdResource struct {
	Id int
}

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	fmt.Printf("date %v", d.Format(time.DateOnly))
	if d.Time.IsZero() {
		return json.Marshal("")
	}

	return json.Marshal(d.Format(time.DateOnly))
}
