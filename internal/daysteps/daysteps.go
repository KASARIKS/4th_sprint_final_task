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
	hInd, mInd := strings.IndexByte(values[1], 'h'), strings.IndexByte(values[1], 'm')

	// Empty or invalid time
	if hInd == -1 && mInd == -1 {
		return 0, 0, errors.New("invalid time")
	}

	var hours, minutes float64

	// Get hours and minutes
	if hInd != -1 {
		hours, err = strconv.ParseFloat(values[1][0:hInd], 64)
		if err != nil {
			return 0, 0, err
		}
	} else {
		hours = 0
	}

	if mInd != -1 {
		minutes, err = strconv.ParseFloat(values[1][hInd+1:mInd], 64)
		if err != nil {
			return 0, 0, err
		}
	} else {
		minutes = 0
	}

	// Wrong numbers in duration
	if hours <= 0 && minutes <= 0 || hours < 0 || minutes < 0 {
		return 0, 0, errors.New("Duration smaller or equil 0")
	}

	walkDuration := time.Duration(hours*float64(time.Hour) + minutes*float64(time.Minute))

	return stepsCount, walkDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Extra validating
	stepsCount, walkDuration, err := parsePackage(data)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", stepsCount, distanceKm, kkall)

	return result
}
