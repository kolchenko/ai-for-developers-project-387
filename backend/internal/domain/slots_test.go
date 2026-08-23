package domain

import (
	"testing"
	"time"
)

func TestGenerateGridSlots(t *testing.T) {
	now := time.Date(2026, 8, 23, 0, 0, 0, 0, time.UTC)
	et := EventType{
		DurationMinutes: 30,
		AvailableFrom:   "09:00:00",
		AvailableTo:     "18:00:00",
	}

	slots := GenerateGridSlots(et, now)
	if len(slots) != 14*18 {
		t.Fatalf("expected %d slots, got %d", 14*18, len(slots))
	}

	// первый слот сегодня в 09:00, последний — 17:30 на 14-й день
	if !slots[0].StartsAt.Equal(time.Date(2026, 8, 23, 9, 0, 0, 0, time.UTC)) {
		t.Errorf("unexpected first slot: %v", slots[0])
	}
	last := slots[len(slots)-1]
	if !last.StartsAt.Equal(time.Date(2026, 9, 5, 17, 30, 0, 0, time.UTC)) ||
		!last.EndsAt.Equal(time.Date(2026, 9, 5, 18, 0, 0, 0, time.UTC)) {
		t.Errorf("unexpected last slot: %v", last)
	}

	// по 18 слотов на каждый из 14 дней
	days := map[string]int{}
	for _, s := range slots {
		days[s.StartsAt.Format("2006-01-02")]++
	}
	if len(days) != 14 {
		t.Fatalf("expected 14 days, got %d", len(days))
	}
	for d, n := range days {
		if n != 18 {
			t.Errorf("day %s: expected 18 slots, got %d", d, n)
		}
	}
}

func TestGenerateGridSlotsFiltersPast(t *testing.T) {
	now := time.Date(2026, 8, 23, 9, 15, 0, 0, time.UTC)
	et := EventType{DurationMinutes: 30, AvailableFrom: "09:00:00", AvailableTo: "10:30:00"}

	slots := GenerateGridSlots(et, now)
	// в день 1 остаются 09:30 и 10:00 (09:00 уже в прошлом), остальные 13 дней — по 3 слота
	if len(slots) != 2+13*3 {
		t.Fatalf("expected %d slots, got %d", 2+13*3, len(slots))
	}
	if !slots[0].StartsAt.Equal(time.Date(2026, 8, 23, 9, 30, 0, 0, time.UTC)) {
		t.Errorf("unexpected first slot: %v", slots[0])
	}
	for _, s := range slots {
		if s.StartsAt.Equal(time.Date(2026, 8, 23, 9, 0, 0, 0, time.UTC)) {
			t.Error("slot in the past not filtered")
		}
	}
}

func TestGenerateGridSlotsLastSlotFits(t *testing.T) {
	now := time.Date(2026, 8, 23, 0, 0, 0, 0, time.UTC)
	// 09:00..10:00 с шагом 45 мин: помещается только 09:00 (09:45+45=10:30 > 10:00)
	et := EventType{DurationMinutes: 45, AvailableFrom: "09:00:00", AvailableTo: "10:00:00"}
	slots := GenerateGridSlots(et, now)
	if len(slots) != 14 {
		t.Fatalf("expected 14 slots (1 per day), got %d", len(slots))
	}
}

func TestIsValidSlotStart(t *testing.T) {
	now := time.Date(2026, 8, 23, 0, 0, 0, 0, time.UTC)
	et := EventType{DurationMinutes: 30, AvailableFrom: "09:00:00", AvailableTo: "18:00:00"}

	valid := time.Date(2026, 8, 24, 9, 30, 0, 0, time.UTC)
	if !IsValidSlotStart(et, valid, now) {
		t.Error("expected valid slot start")
	}

	cases := []struct {
		name string
		at   time.Time
	}{
		{"off grid", time.Date(2026, 8, 24, 9, 10, 0, 0, time.UTC)},
		{"with seconds", time.Date(2026, 8, 24, 9, 30, 30, 0, time.UTC)},
		{"before open", time.Date(2026, 8, 24, 8, 30, 0, 0, time.UTC)},
		{"slot not fully fits", time.Date(2026, 8, 24, 17, 45, 0, 0, time.UTC)},
		{"in the past", time.Date(2026, 8, 22, 9, 0, 0, 0, time.UTC)},
		{"day before window", time.Date(2026, 8, 22, 9, 0, 0, 0, time.UTC)},
		{"outside window", time.Date(2026, 9, 6, 9, 0, 0, 0, time.UTC)},
	}
	for _, tc := range cases {
		if IsValidSlotStart(et, tc.at, now) {
			t.Errorf("%s: expected invalid", tc.name)
		}
	}
}

func TestFilterFree(t *testing.T) {
	day := func(h int) time.Time { return time.Date(2026, 8, 24, h, 0, 0, 0, time.UTC) }

	slots := []Slot{
		{day(9), day(10)},
		{day(10), day(11)},
		{day(11), day(12)},
	}
	bookings := []Booking{
		{StartsAt: day(9).Add(30 * time.Minute), EndsAt: day(10)}, // пересекает первый
		{StartsAt: day(11), EndsAt: day(12)},                       // ровно покрывает третий
	}

	free := FilterFree(slots, bookings)
	if len(free) != 1 || !free[0].StartsAt.Equal(day(10)) {
		t.Fatalf("expected only [10:00,11:00), got %+v", free)
	}
}
