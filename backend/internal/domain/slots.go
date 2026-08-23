package domain

import (
	"fmt"
	"time"
)

const BookingWindowDays = 14

var slotDurations = map[int]bool{15: true, 30: true, 45: true, 60: true}

func IsValidDuration(minutes int) bool {
	return slotDurations[minutes]
}

func durationOf(et EventType) time.Duration {
	return time.Duration(et.DurationMinutes) * time.Minute
}

func dayStartUTC(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func ParseTimeOfDay(value string) (time.Duration, error) {
	t, err := time.Parse("15:04:05", value)
	if err != nil {
		return 0, err
	}
	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second, nil
}

func windowEnd(now time.Time) time.Time {
	return dayStartUTC(now).Add(BookingWindowDays * 24 * time.Hour)
}

// GenerateGridSlots возвращает все слоты сетки типа события на окно в 14 дней,
// начиная с текущей даты включительно. Слоты в прошлом отбрасываются.
func GenerateGridSlots(et EventType, now time.Time) []Slot {
	today := dayStartUTC(now)
	dur := durationOf(et)
	var slots []Slot
	for day := 0; day < BookingWindowDays; day++ {
		d := today.Add(time.Duration(day) * 24 * time.Hour)
		from := d.Add(mustTimeOfDay(et.AvailableFrom))
		to := d.Add(mustTimeOfDay(et.AvailableTo))
		for end := from.Add(dur); !end.After(to); end = end.Add(dur) {
			s := Slot{StartsAt: end.Add(-dur), EndsAt: end}
			if !s.StartsAt.Before(now) {
				slots = append(slots, s)
			}
		}
	}
	return slots
}

// IsValidSlotStart проверяет, что startsAt совпадает с началом допустимого слота:
// не в прошлом, внутри окна 14 дней, совпадает с сеткой типа события и целиком
// умещается в ежедневный диапазон.
func IsValidSlotStart(et EventType, startsAt, now time.Time) bool {
	if startsAt.Before(now) {
		return false
	}
	day := dayStartUTC(startsAt)
	today := dayStartUTC(now)
	if day.Before(today) || !day.Before(windowEnd(now)) {
		return false
	}

	from := day.Add(mustTimeOfDay(et.AvailableFrom))
	to := day.Add(mustTimeOfDay(et.AvailableTo))
	dur := durationOf(et)

	offset := startsAt.Sub(from)
	if offset < 0 || offset%dur != 0 {
		return false
	}
	return !startsAt.Add(dur).After(to)
}

// FilterFree оставляет слоты, не пересекающиеся ни с одной бронью.
func FilterFree(slots []Slot, bookings []Booking) []Slot {
	free := make([]Slot, 0, len(slots))
	for _, s := range slots {
		busy := false
		for _, b := range bookings {
			if s.StartsAt.Before(b.EndsAt) && b.StartsAt.Before(s.EndsAt) {
				busy = true
				break
			}
		}
		if !busy {
			free = append(free, s)
		}
	}
	return free
}

func mustTimeOfDay(value string) time.Duration {
	d, err := ParseTimeOfDay(value)
	if err != nil {
		panic(fmt.Sprintf("invalid time of day %q", value))
	}
	return d
}
