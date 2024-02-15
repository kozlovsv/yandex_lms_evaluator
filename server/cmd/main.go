package main

import (
	"log"
	"net/http"

	"github.com/kozlovsv/evaluator/server/pkg/checker"
	"github.com/kozlovsv/evaluator/server/pkg/handlers"
	"github.com/kozlovsv/evaluator/server/pkg/models"
)

func main() {
	db, err := models.OpenDB()
	if err != nil {
		log.Println("[ERROR]", "connection to DB", err.Error())
		panic(err)
	}
	defer db.Close()

	expressionStore := *models.NewExpressionStore(db)
	agentStore := *models.NewAgentStore(db)
	settingsStore := *models.NewSettingsStore(db)

	//Запускаем проверку задач, и агентов. Если задача долго висит, то возвращаем ее в статус новая, чтобы ее взал другой агент. Если агент долго не доступен то он сначала деактивируется, потом удаляется.
	settings, err := settingsStore.Get()
	if err != nil {
		log.Fatal(err)
	}
	checker := checker.NewChecker(expressionStore, agentStore, settings.OpAgentTimeOut, settings.OpAgentTimeOut, settings.OpAgentDeleteTimeOut)
	defer checker.Close()

	expressionHandler := handlers.NewExpressionHandler(expressionStore)

	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.Handle("/expressions", expressionHandler)
	mux.Handle("/expressions/", expressionHandler)
	mux.Handle("/get-new-task", handlers.NewGetNewExpressionHandler(expressionStore, settingsStore, agentStore))
	mux.Handle("/set-result", handlers.NewSetResultHandlerr(expressionStore, agentStore))
	mux.Handle("/set-error", handlers.NewSetErrorHandlerr(expressionStore, agentStore))
	mux.Handle("/settings", handlers.NewSettingsHandler(settingsStore))
	mux.Handle("/agents", handlers.NewAgentHandler(agentStore))
	log.Println("Server START!")
	http.ListenAndServe(":8001", mux)
	log.Println("Server END!")
}
