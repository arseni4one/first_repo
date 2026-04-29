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
	lenStep                    = 0.65 // средняя длина шага.    // зачем оно здесь?
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("ожидается 3 поля")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверное количество шагов: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверная продолжительность: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше 0")
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	return float64(steps) * (height * stepLengthCoefficient) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные данные")
	}
	calories := (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные данные")
	}

	calories := (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH * walkingCaloriesCoefficient
	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

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

	info := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainingType,
		duration.Hours(),
		dist,
		speed,
		calories,
	)

	return info, nil
}
