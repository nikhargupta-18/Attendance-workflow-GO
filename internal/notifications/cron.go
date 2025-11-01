package notifications

import (
	"attendance-workflow/pkg/config"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type CronService struct {
	scheduler *asynq.Scheduler
}

func NewCronService() *CronService {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}

	redisOpt := asynq.RedisClientOpt{
		Addr:     config.AppConfig.Redis.Addr,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	}

	scheduler := asynq.NewScheduler(
		redisOpt,
		&asynq.SchedulerOpts{
			Location: loc,
		},
	)

	return &CronService{scheduler: scheduler}
}

func (s *CronService) Start() error {
	// Schedule leave reminder notifications every day at 9 AM
	if _, err := s.scheduler.Register("0 9 * * *", asynq.NewTask(TypeReminderEmail, nil)); err != nil {
		return err
	}

	// Schedule absentee report generation every Sunday at 6 PM
	if _, err := s.scheduler.Register("0 18 * * 0", asynq.NewTask(
		"report:absentee",
		nil,
		asynq.Queue("reports"),
		asynq.MaxRetry(3),
	)); err != nil {
		return err
	}

	return s.scheduler.Run()
}

func (s *CronService) Stop() {
	s.scheduler.Shutdown()
}
