package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	part := strings.Split(data, ",")

	if len(part) != 2 {
		return 0, 0, errors.New("Не хватает значений для расчёта")
	}
	// парсим шаги в int
	steps, err := strconv.Atoi(part[0])
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, errors.New("Число шагов равно или меньше 0")
	}
	// парсим время
	duration, err := time.ParseDuration(part[1])
	if err != nil {
		return 0, 0, err
	}

	if duration <= 0 {
		return 0, 0, errors.New("Время равно или меньше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)

	// выводим оишибку и далее возвращаем пустую строку
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distanceMeters := float64(steps) * stepLength
	distanceKilometers := distanceMeters / mInKm

	// расчет калорий
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKilometers, calories)
}
