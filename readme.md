curl http://127.0.0.1:8001/settings
curl http://127.0.0.1:8001/expressions
curl -X POST http://127.0.0.1:8001/expressions -H 'Content-Type: application/json' -d '{"op_plus":100,"op_minus":200,"op_mult":300,"op_div":400,"op_agent_timeout":500}'


