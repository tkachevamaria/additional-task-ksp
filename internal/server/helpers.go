package server

import (
	"fmt"
	"log"
	"time"
)

// ZodiacSign определяет знак зодиака по дате рождения
func ZodiacSign(birthDate string) (int, string, error) {
	log.Printf("🔮 [ZodiacSign] Начинаю определение знака зодиака для даты: %s", birthDate)

	date, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		log.Printf("❌ [ZodiacSign] Ошибка парсинга даты '%s': %v", birthDate, err)
		return 0, "", fmt.Errorf("invalid date format: %w", err)
	}

	month := date.Month()
	day := date.Day()

	log.Printf("📅 [ZodiacSign] Распаршено: месяц=%d (%s), день=%d", month, month.String(), day)

	var zodiacID int
	var zodiacName string

	switch {
	case (month == time.March && day >= 21) || (month == time.April && day <= 19):
		zodiacID, zodiacName = 1, "Овен"
	case (month == time.April && day >= 20) || (month == time.May && day <= 20):
		zodiacID, zodiacName = 2, "Телец"
	case (month == time.May && day >= 21) || (month == time.June && day <= 20):
		zodiacID, zodiacName = 3, "Близнецы"
	case (month == time.June && day >= 21) || (month == time.July && day <= 22):
		zodiacID, zodiacName = 4, "Рак"
	case (month == time.July && day >= 23) || (month == time.August && day <= 22):
		zodiacID, zodiacName = 5, "Лев"
	case (month == time.August && day >= 23) || (month == time.September && day <= 22):
		zodiacID, zodiacName = 6, "Дева"
	case (month == time.September && day >= 23) || (month == time.October && day <= 22):
		zodiacID, zodiacName = 7, "Весы"
	case (month == time.October && day >= 23) || (month == time.November && day <= 21):
		zodiacID, zodiacName = 8, "Скорпион"
	case (month == time.November && day >= 22) || (month == time.December && day <= 21):
		zodiacID, zodiacName = 9, "Стрелец"
	case (month == time.December && day >= 22) || (month == time.January && day <= 19):
		zodiacID, zodiacName = 10, "Козерог"
	case (month == time.January && day >= 20) || (month == time.February && day <= 18):
		zodiacID, zodiacName = 11, "Водолей"
	case (month == time.February && day >= 19) || (month == time.March && day <= 20):
		zodiacID, zodiacName = 12, "Рыбы"
	default:
		log.Printf("❌ [ZodiacSign] Не удалось определить знак зодиака для %d-%d", month, day)
		return 0, "", fmt.Errorf("could not determine zodiac sign")
	}

	log.Printf("✅ [ZodiacSign] Определён знак зодиака: %s (id=%d)", zodiacName, zodiacID)
	return zodiacID, zodiacName, nil
}
