package payload

import (
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(data []byte) error {
	date := string(data)
	// Ignore null, like in the main JSON package.
	if date == "null" || date == `""` {
		return nil
	}

	time, err := time.Parse(time.DateOnly, date[1:len(date)-1])
	*d = Date{time}
	return err
}
