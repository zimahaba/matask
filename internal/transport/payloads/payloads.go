package payload

import (
	"time"
)

type MataskTime struct {
	time.Time
}

func (t *MataskTime) UnmarshalJSON(data []byte) error {
	date := string(data)
	// Ignore null, like in the main JSON package.
	if date == "null" || date == `""` {
		return nil
	}

	time, err := time.Parse(time.DateOnly, date[1:len(date)-1])
	*t = MataskTime{time}
	return err
}
