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

func parsePackage(data string) (int, time.Duration, error) {
	// First validating
	values := strings.Split(data, ",")

	if len(values) != 2 {
		return 0, 0, errors.New("wrong string format")
	}

	// Get steps Count
	stepsCount, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, 0, err
	}

	if stepsCount <= 0 {
		return 0, 0, errors.New("amount of steps equil 0")
	}

	// Get walk duration
	walkDuration, err := time.ParseDuration(values[1])
	if err != nil {
		return 0, 0, err
	}

	if walkDuration <= 0 {
		return 0, 0, errors.New("too small time")
	}

	return stepsCount, walkDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Extra validating
	stepsCount, walkDuration, err := parsePackage(data)
	if err != nil {
		log.Print(err)
		return ""
	}
	if stepsCount <= 0 {
		return ""
	}

	// Count result data
	distance := stepLength * float64(stepsCount)
	distanceKm := distance / mInKm
	kkall, err := spentcalories.WalkingSpentCalories(stepsCount, weight, height, walkDuration)
	if err != nil {
		log.Print(err)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", stepsCount, distanceKm, kkall)

	return result
}
