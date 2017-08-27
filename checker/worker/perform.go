package worker

import (
	"github.com/moira-alert/moira-alert/checker"
	"sync"
	"time"
)

func (worker *Checker) perform(triggerIDs []string, noCache bool, cacheTTL int64, wg *sync.WaitGroup) {
	if noCache {
		for _, id := range triggerIDs {
			wg.Add(1)
			go func(triggerID string) {
				defer wg.Done()
				if err := worker.handleTriggerToCheck(triggerID); err != nil {
					worker.Logger.Errorf("Failed to perform trigger: %s error: %s", triggerID, err.Error())
				}
			}(id)
		}
	} else {
		for _, id := range triggerIDs {
			wg.Add(1)
			go func(triggerID string) {
				defer wg.Done()
				//todo triggerId add check cache cacheTTL seconds
				if err := worker.handleTriggerToCheck(triggerID); err != nil {
					worker.Logger.Errorf("Failed to perform trigger: %s error: %s", triggerID, err.Error())
				}
			}(id)
		}
	}
}

func (worker *Checker) handleTriggerToCheck(triggerID string) error {
	acquired, err := worker.Database.SetTriggerCheckLock(triggerID)
	if err != nil {
		return err
	}
	if acquired {
		start := time.Now()
		if err := worker.checkTrigger(triggerID); err != nil {
			return err
		}
		end := time.Now()
		worker.Metrics.TriggerCheckTime.UpdateSince(start)
		worker.Metrics.TriggerCheckGauge.Update(worker.Metrics.TriggerCheckGauge.Value() + int64(end.Sub(start)))
	}
	return nil
}

func (worker *Checker) checkTrigger(triggerID string) error {
	defer worker.Database.DeleteTriggerCheckLock(triggerID)
	triggerChecker := checker.TriggerChecker{
		TriggerID: triggerID,
		Database:  worker.Database,
		Logger:    worker.Logger,
		Config:    worker.Config,
	}

	err := triggerChecker.InitTriggerChecker()
	if err != nil {
		if err == checker.ErrTriggerNotExists {
			return nil
		}
		return err
	}
	//todo cacheTTL
	return triggerChecker.Check()
}