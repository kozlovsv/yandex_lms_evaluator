package checker

import (
	"log"
	"time"

	"github.com/kozlovsv/evaluator/server/pkg/models"
)

type Checker struct {
	expressionStore    models.ExpressionStore
	agentStore         models.AgentStore
	timeOutExpression  int
	timeOutAgent       int
	timeOutDeleteAgent int
	stop               chan struct{}
}

func NewChecker(expressionStore models.ExpressionStore, agentStore models.AgentStore, timeOutExpression int, timeOutAgent int, timeOutDeleteAgent int) *Checker {
	checker := &Checker{
		expressionStore:    expressionStore,
		agentStore:         agentStore,
		timeOutExpression:  timeOutExpression,
		timeOutAgent:       timeOutAgent,
		timeOutDeleteAgent: timeOutDeleteAgent,
		stop:               make(chan struct{}),
	}

	go checker.check()
	return checker
}

func (c *Checker) check() {
	ticker := time.NewTicker(time.Second)
	log.Println("[INFO] checker was started")
	for {
		select {
		case <-ticker.C:
			c.updateStorage()
		case <-c.stop:
			ticker.Stop()
			log.Println("[INFO] checker was stoped")
			return
		}
	}
}

func (c *Checker) updateStorage() {
	err := c.expressionStore.UpFrozenExpressions(c.timeOutExpression)
	if err != nil {
		log.Println("[ERROR]", "[Agent Checker]", "UpFrozenExpressions", err.Error())
	}
	c.agentStore.SetNotAvailable(c.timeOutAgent)
	if err != nil {
		log.Println("[ERROR]", "[Agent Checker]", "SetNotAvailable", err.Error())
	}
	c.agentStore.DeleteNotAvailable(c.timeOutDeleteAgent)
	if err != nil {
		log.Println("[ERROR]", "[Agent Checker]", "DeleteNotAvailable", err.Error())
	}
}

func (c *Checker) Close() {
	close(c.stop)
}
