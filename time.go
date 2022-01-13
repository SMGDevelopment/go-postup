package postup

import (
	"encoding/json"
	"time"
)

// Time is a time struct that can be unmarshalled from PostUps time encoding.
type Time struct {
	t time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var (
		s   string
		err = json.Unmarshal(data, &s)
	)

	if err != nil {
		return err
	}

	t.t, err = time.Parse("2006-01-02T15:04:05Z", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("2006-01-02T15:04:05Z07:00", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("2006-01-02T15:04:05", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("2006-01-02 15:04:05", s)

	return err
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.t)
}
