package utils

import (
	"bytes"
	"strconv"
	"time"
)

func ToArray(in map[string]string) []byte {
	buf := bytes.NewBuffer([]byte("[")) //nolint:mirror //.
	for _, v := range in {
		buf.WriteString(v[1 : len(v)-1])
		buf.WriteString(",")
	}
	res := buf.Bytes()[:buf.Len()-1]
	if len(res) > 0 {
		res = append(res, ']')
	}
	return res
}

func ToStringsSlice(ints []int64) []string {
	res := make([]string, 0, len(ints))
	for _, v := range ints {
		res = append(res, strconv.FormatInt(v, 10))
	}
	return res
}

func InterfacesToArrayObj(in []interface{}) []byte {
	if len(in) == 1 && in[0] == nil {
		return []byte(`[]`)
	}
	buf := bytes.NewBuffer([]byte("[")) //nolint:mirror //.
	for _, v := range in {
		if s, ok := v.(string); ok {
			buf.WriteString(s)
			buf.WriteString(",")
		}
	}
	if buf.Len() == 1 {
		return []byte(`[]`)
	}

	res := buf.Bytes()[:buf.Len()-1]
	return append(res, ']')
}

func NextDay(date string) (string, error) {
	dt, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return "", err
	}
	dt = dt.AddDate(0, 0, 1)
	return dt.Format(time.DateOnly), nil
}

func GetMonthStart(dt time.Time) time.Time {
	year, month, _ := dt.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}
func GetMonthEnd(dt time.Time) time.Time {
	year, month, _ := dt.Date()
	return time.Date(year, month, DaysIn(dt.Month(), dt.Year()), 0, 0, 0, 0, time.Local)
}

func DaysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func StartOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	daysToSubtract := int(weekday - time.Monday)
	return t.AddDate(0, 0, -daysToSubtract).Truncate(24 * time.Hour) //nolint:mnd //hours per day
}
func EndOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	daysToAdd := 7 - int(weekday)                              //nolint:mnd //days per week
	return t.AddDate(0, 0, daysToAdd).Truncate(24 * time.Hour) //nolint:mnd //hours per day
}
