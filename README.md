# go_messaging

# Install

## Backend

1. cd ./cmd/chat_demo_backend
2. edit ./cmd/chat_demo_backend/config/app.toml
3. go build .
4. run bin file

## Frontend

1. cd ./cmd/chat_demo_frontend
2. edit ./cmd/chat_demo_frontend/config/app.toml ,set `WebsocketServer.WsUrl`
3. go build .
4. run bin file

# Consume Logic

u can modify `cmd/chat_demo_backend/consume_logic.go:19` to finish consume logic