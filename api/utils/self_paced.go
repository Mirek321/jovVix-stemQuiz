package utils

import (
	"fmt"
	"time"

	"github.com/Improwised/jovvix/api/constants"
)

// ParseSessionMode validates the optional session mode query value.
// An empty value defaults to live mode.
func ParseSessionMode(mode string) (string, error) {
	switch mode {
	case "", constants.ModeLive:
		return constants.ModeLive, nil
	case constants.ModeSelfPaced:
		return constants.ModeSelfPaced, nil
	default:
		return "", fmt.Errorf("invalid session mode %q", mode)
	}
}

// SelfPacedDeadline returns the last moment an answer to a self-paced question
// is still accepted: delivery time + question duration + grace period.
func SelfPacedDeadline(deliveredAt time.Time, durationInSeconds int) time.Time {
	return deliveredAt.Add(time.Duration(durationInSeconds+constants.SelfPacedGracePeriodSeconds) * time.Second)
}

// IsSelfPacedTimedOut reports whether a self-paced question delivered at
// deliveredAt can no longer be answered at now.
func IsSelfPacedTimedOut(now, deliveredAt time.Time, durationInSeconds int) bool {
	return now.After(SelfPacedDeadline(deliveredAt, durationInSeconds))
}

// SelfPacedResponseTimeMs computes the server-side response time in
// milliseconds, clamped to [0, duration*1000].
func SelfPacedResponseTimeMs(now, deliveredAt time.Time, durationInSeconds int) int {
	elapsed := now.Sub(deliveredAt).Milliseconds()
	maxMs := int64(durationInSeconds) * 1000
	if elapsed < 0 {
		return 0
	}
	if maxMs < 0 {
		return 0
	}
	if elapsed > maxMs {
		return int(maxMs)
	}
	return int(elapsed)
}

// IsAnswerCorrect reports whether the selected original option keys are a
// correct answer. Survey questions are correct whenever any option was chosen.
func IsAnswerCorrect(selected, correct []int, questionType int) bool {
	if len(selected) == 0 {
		return false
	}
	if questionType == constants.Survey {
		return true
	}
	correctSet := make(map[int]struct{}, len(correct))
	for _, key := range correct {
		correctSet[key] = struct{}{}
	}
	for _, key := range selected {
		if _, ok := correctSet[key]; !ok {
			return false
		}
	}
	return true
}
