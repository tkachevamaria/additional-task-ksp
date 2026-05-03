package server

import (
	"fmt"
	"time"
)

// ZodiacSign определяет знак зодиака по дате рождения
func ZodiacSign(birthDate string) (int, string, error) {
	date, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 0, "", fmt.Errorf("invalid date format: %w", err)
	}

	month := date.Month()
	day := date.Day()

	switch {
	case (month == time.March && day >= 21) || (month == time.April && day <= 19):
		return 1, "Овен", nil
	case (month == time.April && day >= 20) || (month == time.May && day <= 20):
		return 2, "Телец", nil
	case (month == time.May && day >= 21) || (month == time.June && day <= 20):
		return 3, "Близнецы", nil
	case (month == time.June && day >= 21) || (month == time.July && day <= 22):
		return 4, "Рак", nil
	case (month == time.July && day >= 23) || (month == time.August && day <= 22):
		return 5, "Лев", nil
	case (month == time.August && day >= 23) || (month == time.September && day <= 22):
		return 6, "Дева", nil
	case (month == time.September && day >= 23) || (month == time.October && day <= 22):
		return 7, "Весы", nil
	case (month == time.October && day >= 23) || (month == time.November && day <= 21):
		return 8, "Скорпион", nil
	case (month == time.November && day >= 22) || (month == time.December && day <= 21):
		return 9, "Стрелец", nil
	case (month == time.December && day >= 22) || (month == time.January && day <= 19):
		return 10, "Козерог", nil
	case (month == time.January && day >= 20) || (month == time.February && day <= 18):
		return 11, "Водолей", nil
	case (month == time.February && day >= 19) || (month == time.March && day <= 20):
		return 12, "Рыбы", nil
	default:
		return 0, "", fmt.Errorf("could not determine zodiac sign")
	}
}
