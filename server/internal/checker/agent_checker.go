package checker

import (
	"log/slog"
	"time"

	"github.com/kozlovsv/evaluator/server/internal/lib/logger/sl"
	"github.com/kozlovsv/evaluator/server/internal/storage"
)

type Checker struct {
	expressionStore    *storage.ExpressionStore
	agentStore         *storage.AgentStore
	timeOutExpression  int
	timeOutAgent       int
	timeOutDeleteAgent int
	log                *slog.Logger
	stop               chan struct{}
}

func NewChecker(expressionStore *storage.ExpressionStore, agentStore *storage.AgentStore, settingsStore *storage.SettingsStore, log *slog.Logger) *Checker {
	//Запускаем проверку задач, и агентов. Если задача долго висит, то возвращаем ее в статус новая, чтобы ее взал другой агент. Если агент долго не доступен то он сначала деактивируется, потом удаляется.
	settings, err := settingsStore.Get()
	if err != nil {
		panic(err)
	}

	return &Checker{
		expressionStore:    expressionStore,
		agentStore:         agentStore,
		timeOutExpression:  settings.OpAgentTimeOut,
		timeOutAgent:       settings.OpAgentTimeOut,
		timeOutDeleteAgent: settings.OpAgentDeleteTimeOut,
		log:                log,
		stop:               make(chan struct{}),
	}
}

func (c *Checker) Run() {
	go c.check()
}

func (c *Checker) check() {
	log := c.log.With(slog.String("op", "Checker.Check"))
	ticker := time.NewTicker(time.Second)
	log.Info("checker was started")
	for {
		select {
		case <-ticker.C:
			c.updateStorage()
		case <-c.stop:
			ticker.Stop()
			log.Info("checker was stoped")
			return
		}
	}
}

func (c *Checker) updateStorage() {
	log := c.log.With(slog.String("op", "Checker.UpdateStorage"))
	err := c.expressionStore.UpFrozenExpressions(c.timeOutExpression)
	if err != nil {
		log.Error("UpFrozenExpressions", sl.Err(err))
	}
	c.agentStore.SetNotAvailable(c.timeOutAgent)
	if err != nil {
		log.Error("SetNotAvailable", sl.Err(err))
	}
	c.agentStore.DeleteNotAvailable(c.timeOutDeleteAgent)
	if err != nil {
		log.Error("DeleteNotAvailable", sl.Err(err))
	}
}

func (c *Checker) Stop() {
	close(c.stop)
}
