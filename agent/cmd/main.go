package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/kozlovsv/evaluator/agent/pkg/models"
	"github.com/kozlovsv/evaluator/agent/pkg/pool"
	"github.com/kozlovsv/evaluator/agent/pkg/rpn"
)

const SERVER_URL = "http://server:8001/"

var ErrNoTask = errors.New("no tasks")

func getNewExpression(agentName string) (models.Expression, models.Settinsg, error) {
	exp := models.Expression{}
	setings := models.Settinsg{}

	resp, err := http.Get(SERVER_URL + "get-new-task?agent=" + agentName)

	if err != nil {
		return exp, setings, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return exp, setings, ErrNoTask
		} else {
			return exp, setings, errors.New("bad server responce: " + resp.Status)
		}
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return exp, setings, err
	}

	var task models.Task

	if err := json.Unmarshal(body, &task); err != nil {
		return exp, setings, err
	}
	return task.Expression, task.Settings, nil
}

func sendError(id int, message string, agentName string) error {
	resp, err := http.Get(fmt.Sprintf(SERVER_URL+"set-error?id=%d&err=%s&agent=%s", id, url.QueryEscape(message), agentName))
	if resp.StatusCode != http.StatusOK {
		return errors.New("bad server responce: " + resp.Status)
	}
	return err
}

func sendResult(id int, res float64, agentName string) error {
	resp, err := http.Get(fmt.Sprintf(SERVER_URL+"set-result?id=%d&res=%.2f&agent=%s", id, res, agentName))
	if resp.StatusCode != http.StatusOK {
		return errors.New("bad server responce: " + resp.Status)
	}
	return err
}

type ExpressionTask struct {
	exp       models.Expression
	settings  models.Settinsg
	agentName string
}

func (t *ExpressionTask) Execute() error {
	res, err := rpn.Evaluate(t.exp.Value, t.settings)
	if err == nil {
		err = sendResult(t.exp.Id, res, t.agentName)
		if err != nil {
			log.Println("[ERROR]", "Sending Result to server", err.Error())
		}
	} else {
		err = sendError(t.exp.Id, err.Error(), t.agentName)
		if err != nil {
			log.Println("[ERROR]", "Sending Error to server", err.Error())
		}
	}
	return nil
}

func (t *ExpressionTask) OnFailure(err error) {
	log.Println("[ERROR]", "Execution Task", err.Error())
}

func evaluateExpression(pool *pool.MyPool) {
	agentName, exists := os.LookupEnv("EVAL_GO_AGENT_NAME")
	if !exists {
		agentName = "Undefined"
	}

	for {
		exp, settings, err := getNewExpression(agentName)
		if err == nil {
			pool.AddWork(&ExpressionTask{exp: exp, settings: settings, agentName: agentName})
		} else {
			if err != ErrNoTask {
				log.Println("[Error]", "get new expression", err.Error())
			}
		}
		//Ждем 1 секунду чтобы не долбить сильно сервер (этот параметр можно сделать настраиваемым)
		timer := time.NewTimer(500 * time.Millisecond)
		<-timer.C
	}
}

func main() {
	log.Println("Agent START!")

	poolSizeEnv, exists := os.LookupEnv("EVAL_GO_POOL_SIZE")
	poolSize := 5
	if exists {
		i, err := strconv.Atoi(poolSizeEnv)
		if err == nil {
			poolSize = i
		}
	}

	pool, _ := pool.NewWorkerPool(poolSize, 0)
	pool.Start()
	evaluateExpression(pool)
	pool.Stop()
	log.Println("Agent STOP!")
}
