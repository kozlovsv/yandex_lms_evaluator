package grpc

import (
	"context"
	"log/slog"

	pbserver "evaluator/protos/gen/go"

	"github.com/kozlovsv/evaluator/server/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pbserver.UnimplementedExpressionServiceServer
	log             *slog.Logger
	storeAgent      *storage.AgentStore
	storeExpression *storage.ExpressionStore
	storeSettings   *storage.SettingsStore
}

func New(log *slog.Logger,
	storeAgent *storage.AgentStore,
	storeExpression *storage.ExpressionStore,
	storeSettings *storage.SettingsStore) *Server {
	return &Server{
		log:             log,
		storeAgent:      storeAgent,
		storeExpression: storeExpression,
		storeSettings:   storeSettings,
	}

}

func (s *Server) GetNewExpression(ctx context.Context, req *pbserver.ExpressionRequest) (*pbserver.ExpressionResponse, error) {
	agentName := req.AgentName
	//Регистрируем нового агента
	s.storeAgent.Add(agentName)

	exp, err := s.storeExpression.GetNewExpression()
	settings, err := s.storeSettings.Get()
	if err != nil {
		return &pbserver.ExpressionResponse{}, err
	}

	s.storeAgent.SetCurrentOp(agentName, exp.Value)

	task := &pbserver.Task{
		Expression: &pbserver.Expression{Value: exp.Value, Id: int32(exp.Id)},
		Settings: &pbserver.Settings{
			OpPlusTime:     int32(settings.OpPlusTime),
			OpMinusTime:    int32(settings.OpMinusTime),
			OpMultTime:     int32(settings.OpMultTime),
			OpDivTime:      int32(settings.OpDivTime),
			OpAgentTimeOut: int32(settings.OpAgentTimeOut),
		},
	}

	return &pbserver.ExpressionResponse{Task: task}, nil
}

func (s *Server) SendResult(ctx context.Context, in *pbserver.SendResultRequest) (*pbserver.SendResultResponse, error) {
	err := s.storeExpression.SetExpressionResult(int(in.ExpressionId), in.Result)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error while saving result to storage")
	}

	s.storeAgent.Add(in.AgentName)
	s.storeAgent.SetCurrentOp(in.AgentName, "")

	return &pbserver.SendResultResponse{Success: true}, nil
}

func (s *Server) SendError(ctx context.Context, in *pbserver.SendErrorRequest) (*pbserver.SendErrorResponse, error) {
	err := s.storeExpression.SetExpressionError(int(in.ExpressionId), in.Error)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error while saving expression error to storage")
	}

	s.storeAgent.Add(in.AgentName)
	s.storeAgent.SetCurrentOp(in.AgentName, "")
	return &pbserver.SendErrorResponse{Success: true}, nil
}
