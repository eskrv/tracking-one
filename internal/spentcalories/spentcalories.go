package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	part := strings.Split(data, ",")

	if len(part) != 3 {
				return 0, "", 0, errors.New("Не хватает значений для расчёта")
	}
	// парсим шаги в int
	steps, err := strconv.Atoi(part[0])
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("Число шагов равно или меньше 0")
	}

	// парсим время
	duration, err := time.ParseDuration(part[2])
	if err != nil {
		return 0, "", 0, err
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("Время равно или меньше 0")
	}

	return steps, part[1], duration, nil
}



func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	distanceKilometers := distanceMeters / mInKm

	return distanceKilometers
	
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	distanceKilometers := distance(steps, height)
	avgSpeed := distanceKilometers / duration.Hours()

	return avgSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64 

	switch trainingType {
		case "Бег":
			calories, err = RunningSpentCalories(steps, weight, height, duration)

		case "Ходьба":
			calories, err = WalkingSpentCalories(steps, weight, height, duration)

		default: 
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err

	}

	distanceKilometers := distance(steps, height)
	avgSpeed := meanSpeed(steps, height, duration)

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, duration.Hours(), distanceKilometers, avgSpeed, calories), nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Одно из значений меньше или равно нулю")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	return (weight * avgSpeed * durationMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Одно из значений меньше или равно нулю")
	}

	
	avgSpeed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	return ((weight * avgSpeed * durationMinutes) / minInH) * walkingCaloriesCoefficient, nil

}
