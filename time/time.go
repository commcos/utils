package time

import (
	"errors"
	"regexp"
	"time"
)

// WeekStartDay set week start day, default is Monday
var WeekStartDay = time.Monday

// TimeFormats default time formats will be parsed as
var TimeFormats = []string{"1/2/2006", "1/2/2006 15:4:5", "2006", "2006-1", "2006-1-2", "2006-1-2 15", "2006-1-2 15:4", "2006-1-2 15:4:5", "1-2", "15:4:5", "15:4", "15", "15:4:5 Jan 2, 2006 MST", "2006-01-02 15:04:05.999999999 -0700 MST", "2006-01-02T15:04:05-07:00"}

// BeginningOfMinute beginning of minute
func BeginningOfMinute(now time.Time) time.Time {
	return now.Truncate(time.Minute)
}

// BeginningOfHour beginning of hour
func BeginningOfHour(now time.Time) time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, now.Hour(), 0, 0, 0, now.Location())
}

// BeginningOfDay beginning of day
func BeginningOfDay(now time.Time) time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, now.Location())
}

// BeginningOfWeek beginning of week
func BeginningOfWeek(now time.Time) time.Time {
	t := BeginningOfDay(now)
	weekday := int(t.Weekday())

	if WeekStartDay != time.Sunday {
		weekStartDayInt := int(WeekStartDay)

		if weekday < weekStartDayInt {
			weekday = weekday + 7 - weekStartDayInt
		} else {
			weekday = weekday - weekStartDayInt
		}
	}
	return t.AddDate(0, 0, -weekday)
}

// BeginningOfMonth beginning of month
func BeginningOfMonth(now time.Time) time.Time {
	y, m, _ := now.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
}

// BeginningOfQuarter beginning of quarter
func BeginningOfQuarter(now time.Time) time.Time {
	month := BeginningOfMonth(now)
	offset := (int(month.Month()) - 1) % 3
	return month.AddDate(0, -offset, 0)
}

// BeginningOfYear BeginningOfYear beginning of year
func BeginningOfYear(now time.Time) time.Time {
	y, _, _ := now.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, now.Location())
}

// EndOfMinute end of minute
func EndOfMinute(now time.Time) time.Time {
	return BeginningOfMinute(now).Add(time.Minute - time.Nanosecond)
}

// EndOfHour end of hour
func EndOfHour(now time.Time) time.Time {
	return BeginningOfMinute(now).Add(time.Hour - time.Nanosecond)
}

// EndOfDay end of day
func EndOfDay(now time.Time) time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())
}

// EndOfWeek end of week
func EndOfWeek(now time.Time) time.Time {
	return BeginningOfWeek(now).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// EndOfMonth end of month
func EndOfMonth(now time.Time) time.Time {
	return BeginningOfMonth(now).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// EndOfQuarter end of quarter
func EndOfQuarter(now time.Time) time.Time {
	return BeginningOfQuarter(now).AddDate(0, 3, 0).Add(-time.Nanosecond)
}

// EndOfYear end of year
func EndOfYear(now time.Time) time.Time {
	return BeginningOfYear(now).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// Monday monday
func Monday(now time.Time) time.Time {
	t := BeginningOfDay(now)
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return t.AddDate(0, 0, -weekday+1)
}

// Sunday sunday
func Sunday(now time.Time) time.Time {
	t := BeginningOfDay(now)
	weekday := int(t.Weekday())
	if weekday == 0 {
		return t
	}
	return t.AddDate(0, 0, (7 - weekday))
}

// EndOfSunday end of sunday
func EndOfSunday(now time.Time) time.Time {
	return EndOfDay(Sunday(now))
}

func parseWithFormat(now time.Time, str string) (t time.Time, err error) {
	for _, format := range TimeFormats {
		t, err = time.Parse(format, str)
		if err == nil {
			return
		}
	}
	err = errors.New("Can't parse string as time: " + str)
	return
}

var hasTimeRegexp = regexp.MustCompile(`(\s+|^\s*)\d{1,2}((:\d{1,2})*|((:\d{1,2}){2}\.(\d{3}|\d{6}|\d{9})))\s*$`) // match 15:04:05, 15:04:05.000, 15:04:05.000000 15, 2017-01-01 15:04, etc
var onlyTimeRegexp = regexp.MustCompile(`^\s*\d{1,2}((:\d{1,2})*|((:\d{1,2}){2}\.(\d{3}|\d{6}|\d{9})))\s*$`)      // match 15:04:05, 15, 15:04:05.000, 15:04:05.000000, etc

// Parse parse string to time
func Parse(now time.Time, strs ...string) (t time.Time, err error) {
	var (
		setCurrentTime  bool
		parseTime       []int
		currentTime     = []int{now.Nanosecond(), now.Second(), now.Minute(), now.Hour(), now.Day(), int(now.Month()), now.Year()}
		currentLocation = now.Location()
		onlyTimeInStr   = true
	)

	for _, str := range strs {
		hasTimeInStr := hasTimeRegexp.MatchString(str) // match 15:04:05, 15
		onlyTimeInStr = hasTimeInStr && onlyTimeInStr && onlyTimeRegexp.MatchString(str)
		if t, err = parseWithFormat(now, str); err == nil {
			location := t.Location()
			if location.String() == "UTC" {
				location = currentLocation
			}

			parseTime = []int{t.Nanosecond(), t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}

			for i, v := range parseTime {
				// Don't reset hour, minute, second if current time str including time
				if hasTimeInStr && i <= 3 {
					continue
				}

				// If value is zero, replace it with current time
				if v == 0 {
					if setCurrentTime {
						parseTime[i] = currentTime[i]
					}
				} else {
					setCurrentTime = true
				}

				// if current time only includes time, should change day, month to current time
				if onlyTimeInStr {
					if i == 4 || i == 5 {
						parseTime[i] = currentTime[i]
						continue
					}
				}
			}

			t = time.Date(parseTime[6], time.Month(parseTime[5]), parseTime[4], parseTime[3], parseTime[2], parseTime[1], parseTime[0], location)
			currentTime = []int{t.Nanosecond(), t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), t.Year()}
		}
	}
	return
}

// MustParse must parse string to time or it will panic
func MustParse(now time.Time, strs ...string) (t time.Time) {
	t, err := Parse(now, strs...)
	if err != nil {
		panic(err)
	}
	return t
}

// Between check time between the begin, end time or not
func Between(now time.Time, begin, end string) bool {
	beginTime := MustParse(now, begin)
	endTime := MustParse(now, end)
	return now.After(beginTime) && now.Before(endTime)
}
