package postup

import (
	"encoding/json"
	"strings"
	"time"
)

// Time is a time struct that can be unmarshalled from PostUps time encoding.
type Time struct {
	t time.Time
}

func NewTime(t time.Time) Time {
	return Time{t}
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var (
		s   = string(data)
		err error
	)

	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")

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
	if err == nil {
		return nil
	}

	// single and double digit variations on DD/MM/YYYY HH:mm
	t.t, err = time.Parse("01/02/2006 15:04", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("01/02/2006 15:4", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("1/02/2006 15:04", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("1/02/2006 15:4", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("01/2/2006 15:04", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("01/2/2006 15:4", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("1/2/2006 15:04", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("1/2/2006 15:4", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("01/02/06 15:04", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("01/02/06 15:4", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("1/02/06 15:04", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("1/02/06 15:4", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("01/2/06 15:04", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("01/2/06 15:4", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("1/2/06 15:04", s)
	if err == nil {
		return nil
	}

	t.t, err = time.Parse("1/2/06 15:4", s)

	return err
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.t)
}

func (t *Time) GetTime() time.Time {
	return t.t
}
