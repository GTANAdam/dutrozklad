// Package weekday ..
package weekday

import "time"

// Time Base type for DateTime. Can be converted from time.Time
type Time time.Time

// Weekday A Weekday specifies a day of the week (Sunday = 0, ...).
type Weekday time.Weekday

// Weekday Returns specifies a day of the week (Sunday = 0, ...).
func (t Time) Weekday() Weekday {
	return Weekday(time.Time(t).Weekday())
}

func (w Weekday) String() (weekday string) {
	switch w {
	case 0:
		weekday = "Неділя"
	case 1:
		weekday = "Понеділок"
	case 2:
		weekday = "Вівторок"
	case 3:
		weekday = "Середа"
	case 4:
		weekday = "Четвер"
	case 5:
		weekday = "П'ятниця"
	case 6:
		weekday = "Субота"
	}
	return
}
