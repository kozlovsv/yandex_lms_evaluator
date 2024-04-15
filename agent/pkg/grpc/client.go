package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pbserver "evaluator/protos/gen/go"
)

type Client struct {
	conn    *grpc.ClientConn
	service pbserver.ExpressionServiceClient
}

func NewClient(grpcPort string) (*Client, error) {
	ctx := context.Background()
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return net.Dial("tcp", addr)
	}

	//Вынести хост в настройки
	//Сделать обработку случая когда нет строк
	//Протестировать
	//Поправить описание

	conn, err := grpc.DialContext(ctx, "server:"+grpcPort, grpc.WithContextDialer(dialer), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pbserver.NewExpressionServiceClient(conn)

	return &Client{conn, c}, nil
}

func (c *Client) GetNewExpression(agentName string) (*pbserver.ExpressionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.service.GetNewExpression(ctx, &pbserver.ExpressionRequest{AgentName: agentName})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) SendError(id int, errorMsg string, agentName string) (*pbserver.SendErrorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.service.SendError(ctx, &pbserver.SendErrorRequest{AgentName: agentName, ExpressionId: int32(id), Error: errorMsg})
	if err != nil {
		log.Fatalf("Error while calling SendError : %v", err)
	}

	return r, nil
}

func (c *Client) SendResult(id int, res float64, agentName string) (*pbserver.SendResultResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.service.SendResult(ctx, &pbserver.SendResultRequest{AgentName: agentName, ExpressionId: int32(id), Result: fmt.Sprintf("%.2f", res)})
	if err != nil {
		log.Fatalf("Error while calling SendResult : %v", err)
	}

	return r, err
}
