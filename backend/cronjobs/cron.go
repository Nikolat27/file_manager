package cronjobs

import (
	"file_manager/handlers"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"log/slog"
	"time"
)

const SixHoursInterval = "0 0 */6 * * *"

type CronJob struct {
	Cron *cron.Cron
}

func New() *CronJob {
	return &CronJob{
		Cron: cron.New(cron.WithSeconds()),
	}
}

func (cron *CronJob) Start(handler *handlers.Handler) {
	_, err := cron.Cron.AddFunc(SixHoursInterval, func() {
		if err := FilesCleanup(handler); err != nil {
			slog.Error("deleting file instance", "error", err)
		}
	})
	if err != nil {
		slog.Error("adding checkFileSettingsExpiration function", "error", err)
	}

	_, err = cron.Cron.AddFunc(SixHoursInterval, func() {
		if err := FileSettingsCleanup(handler); err != nil {
			slog.Error("deleting file setting instance", "error", err)
		}
	})
	if err != nil {
		slog.Error("adding checkFilesExpiration function", "error", err)
	}

	cron.Cron.Start()
}

func FilesCleanup(handler *handlers.Handler) error {
	filter := bson.M{
		"expire_at": bson.M{
			"$lte": time.Now(),
		},
	}

	deletedAmount, err := handler.Models.File.Delete(filter)
	if err != nil {
		return err
	}

	log.Printf("deleted %d files", deletedAmount)
	return nil
}

func FileSettingsCleanup(handler *handlers.Handler) error {
	filter := bson.M{
		"expiration_at": bson.M{
			"$lte": time.Now(),
		},
	}

	deletedAmount, err := handler.Models.FileSettings.Delete(filter)
	if err != nil {
		return err
	}

	log.Printf("deleted %d file settings", deletedAmount)
	return nil
}
