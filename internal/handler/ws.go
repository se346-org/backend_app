package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/chat-socio/backend/configuration"
	"github.com/chat-socio/backend/infrastructure/websocket"
	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/internal/usecase"
	"github.com/chat-socio/backend/pkg/jwt"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/cloudwego/hertz/pkg/app"
	ws "github.com/hertz-contrib/websocket"
)

type WebSocketHandler struct {
	upgrader          *ws.HertzUpgrader
	UserOnlineUsecase usecase.UserOnlineUsecase
	UserUsecase       usecase.UserUseCase
	obs               *observability.Observability
}

func (wsh *WebSocketHandler) HandleWebsocket(ctx context.Context, c *app.RequestContext) {
	ctx, span := wsh.obs.StartSpan(ctx, "WebSocketHandler.HandleWebsocket")
	defer span()

	// Upgrade the connection to a WebSocket connection
	err := wsh.upgrader.Upgrade(c, func(conn *ws.Conn) {
		wsConn, err := websocket.NewWSConnection(conn)
		if err != nil {
			conn.WriteMessage(ws.TextMessage, fmt.Appendf(nil, "Failed to create WebSocket connection: %v", err))
			conn.Close()
			return
		}

		msg, err := wsConn.ReceiveMessage()
		if err != nil {
			wsConn.SendMessage(fmt.Appendf(nil, "Failed to read message: %v", err))
			wsConn.Close()
			return
		}

		// fmt.Println("Received message:", string(msg))

		// wsConn.SendMessage(msg)

		// Process the message (this is where you would handle incoming messages)
		var wsMessage domain.WebSocketMessage

		err = json.Unmarshal(msg, &wsMessage)
		if err != nil {
			wsConn.SendMessage(fmt.Appendf(nil, "Failed to unmarshal message: %v", err))
			wsConn.Close()
			return
		}

		// Handle the message based on its type
		switch wsMessage.Type {
		case domain.WsAuthorization:
			// Handle authorization message
			token, ok := wsMessage.Payload["token"]
			if !ok {
				wsConn.SendMessage(fmt.Appendf(nil, "Token not found in message"))
				wsConn.Close()
				return
			}
			// Validate the token and get the user ID
			jwtToken, err := jwt.ParseHS256JWT(token.(string), configuration.ConfigInstance.JWT.SecretKey)
			if err != nil {
				wsConn.SendMessage(fmt.Appendf(nil, "Failed to parse token: %v", err))
				wsConn.Close()
				return
			}

			// validate
			b, err := jwt.ValidateHS256JWT(jwtToken)
			if err != nil {
				wsConn.SendMessage(fmt.Appendf(nil, "Failed to validate token: %v", err))
				wsConn.Close()
				return
			}

			if !b {
				wsConn.SendMessage(fmt.Appendf(nil, "Token is invalid"))
				wsConn.Close()
				return
			}

			// Get claims from the token
			claims, err := jwt.ExtractClaims(jwtToken)
			if err != nil {
				wsConn.SendMessage(fmt.Appendf(nil, "Failed to extract claims: %v", err))
				wsConn.Close()
				return
			}

			// Get claims from the token
			claims, err = jwt.ExtractClaims(jwtToken)
			if err != nil {
				wsConn.SendMessage(fmt.Appendf(nil, "Failed to extract claims: %v", err))
				wsConn.Close()
				return
			}

			domain.WebSocket.AddWrapConnection(wsConn)
			userID, err := wsh.UserUsecase.GetUserIDByAccountID(ctx, claims.Sub)
			if err != nil {
				wsConn.SendMessage(fmt.Appendf(nil, "Failed to get user ID: %v", err))
				wsConn.Close()
				return
			}
			userOnline := &domain.UserOnline{
				UserID:       userID,
				ConnectionID: wsConn.GetID(),
			}
			err = wsh.UserOnlineUsecase.CreateUserOnline(ctx, userOnline)
			if err != nil {
				wsConn.SendMessage(fmt.Appendf(nil, "Failed to create user online: %v", err))
				wsConn.Close()
				return
			}
			wsResonse := domain.NewWebSocketMessage(domain.WsAuthorization, map[string]any{
				"account_id":     claims.Sub,
				"user_id":        userID,
				"user_online_id": userOnline.ID,
			})
			wsConn.SendMessage([]byte(wsResonse.String()))
		default:
			// Close the connection if the message type is not recognized
			wsConn.SendMessage(fmt.Appendf(nil, "Unknown message type"))
			wsConn.Close()
			return
		}

		for {
			// Read messages from the WebSocket connection
			msg, err := wsConn.ReceiveMessage()
			if err != nil {
				wsConn.SendMessage(fmt.Appendf(nil, "Failed to read message: %v", err))
				domain.WebSocket.RemoveConnection(wsConn.GetID())
				break
			}

			// wsConn.SendMessage(fmt.Appendf(nil, "Received message: %s", msg))
			// Process the message (this is where you would handle incoming messages)
			var wsMessage domain.WebSocketMessage
			err = json.Unmarshal(msg, &wsMessage)
			if err != nil {
				wsConn.SendMessage(fmt.Appendf(nil, "Failed to unmarshal message: %v", err))
				domain.WebSocket.RemoveConnection(wsConn.GetID())
				break
			}
			// Handle the message based on its type
			switch wsMessage.Type {
			case domain.WsPing:
				// Handle ping message
				// Send a pong message back
				pongMessage := domain.NewWebSocketMessage(domain.WsPong, nil)
				err = wsConn.SendMessage([]byte(pongMessage.String()))
			default:
				// Handle other message types
				// wsConn.SendMessage([]byte(fmt.Sprintf("Unknown message type: %s", wsMessage.Type)))
				// wsConn.Close()

			}

			// Close the connection if an error occurs
			if err != nil {
				log.Println("Failed to read message:", err)
				domain.WebSocket.RemoveConnection(wsConn.GetID())
				break
			}
		}
	})

	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
}

func NewWebSocketHandler(upgrader *ws.HertzUpgrader, userOnlineUsecase usecase.UserOnlineUsecase, userUsecase usecase.UserUseCase, obs *observability.Observability) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader:          upgrader,
		UserOnlineUsecase: userOnlineUsecase,
		UserUsecase:       userUsecase,
		obs:               obs,
	}
}
