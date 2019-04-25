package otime

import (
	"fmt"
	"github.com/standup-raven/standup-raven/server/config"
	"strings"
	"time"
)

type OTime struct {
	time.Time
}

const (
	layoutTime            = "15:04"
	layoutTimeWithSeconds = "15:04:05"
	layoutDate            = "20060102"
	userTimeLayout        = "2006 15:04"
)

var nilTime = (time.Time{}).UnixNano()

func Parse(value string) (OTime, error) {
	argTime, err := time.Parse(layoutTime, value)
	if err != nil {
		return OTime{}, err
	}

	now := time.Now()
	argTime = time.Date(now.Year(), now.Month(), now.Day(), argTime.Hour(), argTime.Minute(), 0, 0, config.GetConfig().Location)
	return OTime{argTime}, nil
}

func Now() OTime {
	now := time.Now()
	return OTime{now.In(config.GetConfig().Location)}
}

func (ct OTime) GetTime() OTime {
	now, _ := time.Parse(layoutTime, ct.Format(layoutTime))
	return OTime{now.In(config.GetConfig().Location)}
}

func (ct OTime) GetTimeWithSeconds() OTime {
	now, _ := time.Parse(layoutTimeWithSeconds, ct.Format(layoutTimeWithSeconds))
	return OTime{now.In(config.GetConfig().Location)}
}

func (ct OTime) GetTimeString() string {
	return ct.Time.Format(layoutTime)
}

func (ct OTime) GetDate() OTime {
	now, _ := time.Parse(layoutDate, ct.Format(layoutDate))
	return OTime{now.In(config.GetConfig().Location)}
}

func (ct OTime) GetDateString() string {
	return ct.Time.Format(layoutDate)
}

func (ct *OTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	t, err := Parse(s)
	if err != nil {
		return err
	}

	ct.Time = t.Time
	return
}

func (ct OTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(layoutTime))), nil
}
