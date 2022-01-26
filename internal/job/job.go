package job

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	"github.com/mispon/stewart-bot/internal/commands"
	"github.com/mispon/stewart-bot/internal/config"
)

type cronJob struct {
	session *discordgo.Session
	config  *config.Config
	jobs    []commands.JobCommand
}

// New creates new job instance
func New(session *discordgo.Session, cfg *config.Config) *cronJob {
	return &cronJob{
		session: session,
		config:  cfg,
		jobs: []commands.JobCommand{
			commands.NewHoroscopeV2Cmd(cfg),
		},
	}
}

// Run schedule periodic job
func (cj cronJob) Run() chan<- struct{} {
	doneCh := make(chan struct{})

	go func(done <-chan struct{}) {
		for {
			select {
			case <-done:
				logrus.Debug("terminating cron job...")
				return
			case <-time.After(time.Minute):
				cj.triggerJobs()
			}
		}
	}(doneCh)

	return doneCh
}

func (cj cronJob) triggerJobs() {
	for _, job := range cj.jobs {
		now := time.Now()
		jh, jm := job.TriggerTime()

		if now.Hour() == jh && now.Minute() == jm {
			err := job.Run(cj.session)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}
