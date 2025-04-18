package daysteps

import (
	"errors"
	"fmt"
	"log"
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

func splitData(data string) ([]string, error) {
	splitedData := strings.Split(data, ",")
	if len(splitedData) != 2 {
		return []string{}, errors.New("wrong string format")
	}

	return splitedData, nil
}

func getStepsCount(steps string) (int, error) {
	stepsCount, err := strconv.Atoi(steps)
	if stepsCount <= 0 && err == nil {
		return 0, errors.New("amount of steps equil 0")
	}

	return stepsCount, err
}

func getWalkDuration(duration string) (time.Duration, error) {
	walkDuration, err := time.ParseDuration(duration)
	if walkDuration <= 0 && err == nil {
		return 0, errors.New("too small time")
	}

	return walkDuration, err
}

func parsePackage(data string) (int, time.Duration, error) {
	splitedData, err := splitData(data)
	if err != nil {
		return 0, 0, err
	}

	stepsCount, err := getStepsCount(splitedData[0])
	if err != nil {
		return 0, 0, err
	}

	walkDuration, err := getWalkDuration(splitedData[1])
	if err != nil {
		return 0, 0, err
	}

	return stepsCount, walkDuration, nil
}

func distanceInKilometres(stepsCount int) float64 {
	distance := stepLength * float64(stepsCount)
	distanceKm := distance / mInKm
	return distanceKm
}

func DayActionInfo(data string, weight, height float64) string {
	// Get parsed statistics
	stepsCount, walkDuration, err := parsePackage(data)
	if err != nil {
		log.Print(err)
		return ""
	}

	// Count result data
	distanceKm := distanceInKilometres(stepsCount)
	kkall, err := spentcalories.WalkingSpentCalories(stepsCount, weight, height, walkDuration)
	if err != nil {
		log.Print(err)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", stepsCount, distanceKm, kkall)

	return result
}
