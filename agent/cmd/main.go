package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/kozlovsv/evaluator/agent/pkg/grpc"
	"github.com/kozlovsv/evaluator/agent/pkg/models"
	"github.com/kozlovsv/evaluator/agent/pkg/pool"
	"github.com/kozlovsv/evaluator/agent/pkg/rpn"
)

type ExpressionTask struct {
	exp        models.Expression
	settings   models.Settinsg
	agentName  string
	grpcClient grpc.Client
}

func (t *ExpressionTask) Execute() error {
	res, err := rpn.Evaluate(t.exp.Value, t.settings)
	if err == nil {
		_, err = t.grpcClient.SendResult(t.exp.Id, res, t.agentName)
		if err != nil {
			log.Println("[ERROR]", "Sending Result to server", err.Error())
		}
	} else {
		_, err = t.grpcClient.SendError(t.exp.Id, err.Error(), t.agentName)
		if err != nil {
			log.Println("[ERROR]", "Sending Error to server", err.Error())
		}
	}
	return nil
}

func (t *ExpressionTask) OnFailure(err error) {
	log.Println("[ERROR]", "Execution Task", err.Error())
}

func evaluateExpression(pool *pool.MyPool, grpcClient grpc.Client) {
	agentName, exists := os.LookupEnv("EVAL_GO_AGENT_NAME")
	if !exists {
		agentName = "Undefined"
	}

	for {
		er, err := grpcClient.GetNewExpression(agentName)

		if err == nil {
			exp := models.Expression{
				Value: er.Task.Expression.Value,
				Id:    int(er.Task.Expression.Id),
			}

			grpcSettings := er.Task.Settings

			settings := models.Settinsg{
				OpPlusTime:     int(grpcSettings.OpPlusTime),
				OpMinusTime:    int(grpcSettings.OpMinusTime),
				OpMultTime:     int(grpcSettings.OpMultTime),
				OpDivTime:      int(grpcSettings.OpDivTime),
				OpAgentTimeOut: int(grpcSettings.OpAgentTimeOut),
			}
			pool.AddWork(&ExpressionTask{exp: exp, settings: settings, agentName: agentName, grpcClient: grpcClient})
		} else {
			if err != sql.ErrNoRows {
				log.Println("[Error]", "get new expression", err.Error())
			}

		}

		//Ждем 2 секунды чтобы не долбить сильно сервер (этот параметр можно сделать настраиваемым)
		timer := time.NewTimer(time.Second * 2)
		<-timer.C
	}
}

func main() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found. ")
	}

	log.Println("Agent START!")

	poolSizeEnv, exists := os.LookupEnv("EVAL_GO_POOL_SIZE")
	poolSize := 5
	if exists {
		i, err := strconv.Atoi(poolSizeEnv)
		if err == nil {
			poolSize = i
		}
	}
	port, exists := os.LookupEnv("SERVER_GRPC_PORT")
	if !exists {
		port = "50051"
	}

	grpcClient, err := grpc.NewClient(port)

	if err != nil {
		panic(err)
	}

	pool, _ := pool.NewWorkerPool(poolSize, 0)
	pool.Start()
	evaluateExpression(pool, *grpcClient)
	pool.Stop()
	log.Println("Agent STOP!")
}
