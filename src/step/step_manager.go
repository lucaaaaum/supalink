package step

import (
	"fmt"
	"supalink/src/settings"
)

type StepManager struct {
	currentStep      int
	currentStepCount int
}

func (sm *StepManager) NextStep(settings settings.Settings) (int, int, error) {
	if len(settings.Steps) == 0 {
		return 0, 0, fmt.Errorf("no steps defined")
	}

	if sm.currentStep == 0 {
		sm.currentStep = 1
		sm.currentStepCount = 1
		return sm.currentStep, sm.currentStepCount, nil
	}

	if sm.currentStepCount >= settings.Steps[sm.currentStep-1] {
		if sm.currentStep >= len(settings.Steps) {
			return 0, 0, fmt.Errorf("exceeded the number of defined steps")
		}
		sm.currentStep++
		sm.currentStepCount = 1
	} else {
		sm.currentStepCount++
	}

	return sm.currentStep, sm.currentStepCount, nil
}
