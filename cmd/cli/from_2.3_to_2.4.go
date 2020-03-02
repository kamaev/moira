package main

import (
	moira2 "github.com/moira-alert/moira/internal/moira"
)

func updateFrom23(logger moira2.Logger, dataBase moira2.Database) error {
	logger.Info("Update 2.3 -> 2.4 start")

	logger.Info("Start marking unused triggers")
	if err := resaveTriggers(dataBase); err != nil {
		return err
	}

	logger.Info("Update 2.3 -> 2.4 finish")
	return nil
}

func downgradeTo23(logger moira2.Logger, dataBase moira2.Database) error {
	return nil
}

func resaveTriggers(database moira2.Database) error {
	allTriggerIDs, err := database.GetAllTriggerIDs()
	if err != nil {
		return err
	}
	allTriggers, err := database.GetTriggers(allTriggerIDs)
	if err != nil {
		return err
	}
	for _, trigger := range allTriggers {
		if trigger != nil {
			if err = database.SaveTrigger(trigger.ID, trigger); err != nil {
				return err
			}
		}
	}
	return nil
}
