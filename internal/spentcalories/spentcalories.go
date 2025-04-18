package spentcalories

import (
	"errors"
	"fmt"
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

func splitData(data string) ([]string, error) {
	splitedData := strings.Split(data, ",")
	if len(splitedData) != 3 {
		return []string{}, errors.New("wrong data string format")
	}

	return splitedData, nil
}

func getStepsCount(steps string) (int, error) {
	stepsCount, err := strconv.Atoi(steps)
	if stepsCount <= 0 && err == nil {
		return 0, errors.New("steps count smaller or equil 0")
	}

	return stepsCount, err
}

func getWalkDuration(duration string) (time.Duration, error) {
	walkDuration, err := time.ParseDuration(duration)
	if walkDuration <= 0 && err == nil {
		return 0, errors.New("walk duration smaller or equil 0")
	}
	return walkDuration, err
}

func parseTraining(data string) (int, string, time.Duration, error) {
	splitedData, err := splitData(data)
	if err != nil {
		return 0, "", 0, err
	}

	// Validate steps count
	stepsCount, err := getStepsCount(splitedData[0])
	if err != nil {
		return 0, "", 0, err
	}

	// Validate walkDuration
	walkDuration, err := getWalkDuration(splitedData[2])
	if err != nil {
		return 0, "", 0, err
	}

	return stepsCount, splitedData[1], walkDuration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	dist := float64(steps) * stepLength
	distKm := dist / mInKm
	return distKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	averageSpeed := dist / duration.Hours()
	return averageSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepsCount, activityType, walkingDuration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	// Calculate distance for activity type
	var dist, averageSpeed, calories float64

	dist = distance(stepsCount, height)
	averageSpeed = meanSpeed(stepsCount, height, walkingDuration)

	switch activityType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(stepsCount, weight, height, walkingDuration)
		if err != nil {
			return "", err
		}
	case "Бег":
		calories, err = RunningSpentCalories(stepsCount, weight, height, walkingDuration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	// Create result string
	res := fmt.Sprintf("Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f\n",
		activityType, walkingDuration.Hours(), dist, averageSpeed, calories)

	return res, nil
}

func validateInfo(steps int, weight, height float64, duration time.Duration) error {
	if steps <= 0 {
		return errors.New("steps are smaller or equil 0")
	}

	if weight <= 0 {
		return errors.New("weight is smaller or equil 0")
	}

	if height <= 0 {
		return errors.New("height is smaller or equil 0")
	}

	if duration <= 0 {
		return errors.New("duration is smaller or equil 0")
	}

	return nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Validate all paramaters
	err := validateInfo(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	// Caclulate calories
	averageSpeed := meanSpeed(steps, height, duration)

	calories := (weight * averageSpeed * duration.Minutes()) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Validate all paramaters
	err := validateInfo(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	// Caclulate calories
	averageSpeed := meanSpeed(steps, height, duration)

	calories := (weight * averageSpeed * duration.Minutes()) / minInH * walkingCaloriesCoefficient

	return calories, nil
}
