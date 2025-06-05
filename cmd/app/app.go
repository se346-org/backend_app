package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chat-socio/backend/configuration"
	"github.com/chat-socio/backend/infrastructure/http"
	"github.com/chat-socio/backend/infrastructure/minio"
	"github.com/chat-socio/backend/infrastructure/nats"
	"github.com/chat-socio/backend/infrastructure/postgresql"
	"github.com/chat-socio/backend/infrastructure/redis"
	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/internal/handler"
	"github.com/chat-socio/backend/internal/middleware"
	"github.com/chat-socio/backend/internal/usecase"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	hertzSwagger "github.com/hertz-contrib/swagger"
	"github.com/hertz-contrib/websocket"
	natsjs "github.com/nats-io/nats.go"
	_ "github.com/chat-socio/backend/docs" // This is important!
)

// @title Chat Socio API
// @version 1.0
// @description This is the API documentation for Chat Socio backend service
// @host 100.70.60.105:8887
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

type Handler struct {
	UserHandler         *handler.UserHandler
	ConversationHandler *handler.ConversationHandler
	Middleware          *middleware.Middleware
	WebSocketHandler    *handler.WebSocketHandler
	UploadHandler       *handler.UploadHandler
	FriendHandler       *handler.FriendHandler
}

func CreateStream(js natsjs.JetStreamContext) error {
	_, err := js.AddStream(&natsjs.StreamConfig{
		Name:     domain.STREAM_NAME_CONVERSATION,
		Subjects: []string{domain.SUBJECT_WILDCARD_CONVERSATION},
	})

	if err != nil {
		return err
	}

	_, err = js.AddStream(&natsjs.StreamConfig{
		Name:     domain.STREAM_NAME_WS_MESSAGE,
		Subjects: []string{domain.SUBJECT_WILDCARD_MESSAGE},
	})
	if err != nil {
		return err
	}

	return nil

}

func RunApp() {
	ctx, cancel := context.WithCancel(context.Background())
	// Initialize the database connection
	db, err := postgresql.Connect(ctx, configuration.ConfigInstance.Postgres)
	if err != nil {
		panic(err)
	}

	redisClient := redis.Connect(configuration.ConfigInstance.Redis)

	natsClient := nats.Connect(configuration.ConfigInstance.Nats.Address)
	js, err := natsClient.JetStream()
	if err != nil {
		panic(err)
	}
	//Init websocket
	domain.InitWebSocket()

	// Create stream
	err = CreateStream(js)
	if err != nil {
		panic(err)
	}

	observability, err := observability.New(observability.Config{
		TracingEnabled: configuration.ConfigInstance.Observability.TracingEnabled,
		JaegerEndpoint: configuration.ConfigInstance.Observability.JaegerEndpoint,
		ServiceName:    configuration.ConfigInstance.Observability.JaegerService,
	})
	if err != nil {
		panic(err)
	}
	// Initialize storage
	storage, err := minio.NewMinioClient(configuration.ConfigInstance.Minio, observability)
	if err != nil {
		panic(err)
	}

	// Initialize repositories
	accountRepository := postgresql.NewAccountRepository(db)
	userRepository := postgresql.NewUserRepository(db, observability)
	sessionRepository := postgresql.NewSessionRepository(db)
	sessionCacheRepository := redis.NewSessionCacheRepository(redisClient)
	userCacheRepository := redis.NewUserCacheRepository(redisClient)
	conversationRepository := postgresql.NewConversationRepository(db)
	messageRepository := postgresql.NewMessageRepository(db)
	userOnlineRepository := postgresql.NewUserOnlineRepository(db)
	seenMessageRepository := postgresql.NewSeenMessageRepository(db, observability)
	friendRepository := postgresql.NewFriendRepository(db)

	// Initialize publisher
	messagePublisher := nats.NewPublisher(js)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(accountRepository, userRepository, sessionRepository, sessionCacheRepository, userCacheRepository, observability)
	conversationUseCase := usecase.NewConversationUseCase(conversationRepository, messageRepository, messagePublisher, userOnlineRepository, userRepository, seenMessageRepository, observability)
	userOnlineUseCase := usecase.NewUserOnlineUsecase(userOnlineRepository)
	friendUseCase := usecase.NewFriendUseCase(friendRepository, userRepository, observability)

	// Initialize the handler
	handler := &Handler{
		UserHandler: &handler.UserHandler{
			UserUseCase: userUseCase,
			Obs:         observability,
		},

		Middleware: middleware.NewMiddleware(sessionCacheRepository, sessionRepository),
		WebSocketHandler: handler.NewWebSocketHandler(&websocket.HertzUpgrader{
			CheckOrigin: func(c *app.RequestContext) bool {
				return true
			},
		}, userOnlineUseCase, userUseCase, observability),
		ConversationHandler: &handler.ConversationHandler{
			ConversationUseCase: conversationUseCase,
			UserUseCase:         userUseCase,
			Obs:                 observability,
		},
		UploadHandler: &handler.UploadHandler{
			Storage: storage,
			Obs:     observability,
		},
		FriendHandler: handler.NewFriendHandler(friendUseCase, observability, userCacheRepository),
	}

	// Init subscriber
	WsNewMessageSubscriber := nats.NewSubscriber(js, domain.CONSUMER_NAME_WS_MESSAGE_NEW)
	err = WsNewMessageSubscriber.Subscribe(ctx, domain.SUBJECT_NEW_MESSAGE, nats.WrapHandler(conversationUseCase.HandleNewMessage))
	if err != nil {
		panic(err)
	}

	UpdateLastMessageSubscriber := nats.NewQueueSubscriber(js, domain.QUEUE_NAME_WS_MESSAGE_UPDATE_LAST_MESSAGE, domain.CONSUMER_NAME_WS_MESSAGE_UPDATE_LAST_MESSAGE)
	err = UpdateLastMessageSubscriber.Subscribe(ctx, domain.SUBJECT_UPDATE_LAST_MESSAGE_ID, nats.WrapHandler(conversationUseCase.HandleUpdateLastMessageID))
	if err != nil {
		panic(err)
	}

	SeenMessageSubscriber := nats.NewQueueSubscriber(js, domain.QUEUE_NAME_SEEN_MESSAGE, domain.CONSUMER_NAME_SEEN_MESSAGE)
	err = SeenMessageSubscriber.Subscribe(ctx, domain.SUBJECT_SEEN_MESSAGE, nats.WrapHandler(conversationUseCase.HandleSeenMessage))
	if err != nil {
		panic(err)
	}

	// Initialize the server
	s := http.NewServer(configuration.ConfigInstance.Server)
	s.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

	// Set up routes
	SetUpRoutes(s, handler)

	//graceful shutdown
	var signalChan = make(chan os.Signal, 1)
	go func() {
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		fmt.Println("Received shutdown signal, shutting down gracefully...")
		WsNewMessageSubscriber.Unsubscribe()
		db.Close()
		redisClient.Close()
		natsClient.Drain()
		cancel()
	}()

	// Start the server
	s.Spin()
}

func SetUpRoutes(s *server.Hertz, handler *Handler) {
	// Swagger documentation
	url := hertzSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	s.GET("/swagger/*any", hertzSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Route not use auth middleware
	// @Summary Register a new user
	// @Description Register a new user with email and password
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param request body presenter.RegisterRequest true "Register Request"
	// @Success 200 {object} presenter.BaseResponse[presenter.RegisterResponse]
	// @Router /user/register [post]
	s.POST(("/user/register"), handler.UserHandler.Register)

	// @Summary Login user
	// @Description Login with email and password
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param request body presenter.LoginRequest true "Login Request"
	// @Success 200 {object} presenter.BaseResponse[presenter.LoginResponse]
	// @Router /user/login [post]
	s.POST(("/user/login"), handler.UserHandler.Login)

	// Route use auth middleware
	authGroup := s.Group("/auth")
	authGroup.Use(handler.Middleware.AuthMiddleware())

	// @Summary Get user info
	// @Description Get current user information
	// @Tags user
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Success 200 {object} presenter.BaseResponse[presenter.GetUserInfoResponse]
	// @Router /auth/user/info [get]
	authGroup.GET("/user/info", handler.UserHandler.GetMyInfo)

	// @Summary Search users
	// @Description Search for users by query
	// @Tags user
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param keyword query string false "Search keyword"
	// @Param limit query int false "Limit" default(10)
	// @Param last_id query string false "Last ID for pagination"
	// @Success 200 {object} presenter.BaseResponse[[]presenter.GetUserInfoResponse]
	// @Router /auth/user/search [get]
	authGroup.GET("/user/search", handler.UserHandler.GetListUser)

	// @Summary Get conversations
	// @Description Get list of user's conversations
	// @Tags conversation
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param conversation_id query string false "Conversation ID"
	// @Param last_message_id query string false "Last message ID for pagination"
	// @Param limit query int false "Limit" default(20)
	// @Success 200 {object} presenter.BaseResponse[[]presenter.GetListConversationResponse]
	// @Router /auth/conversation [get]
	authGroup.GET("/conversation", handler.ConversationHandler.GetListConversation)

	// @Summary Create conversation
	// @Description Create a new conversation
	// @Tags conversation
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body presenter.CreateConversationRequest true "Create Conversation Request"
	// @Success 200 {object} presenter.BaseResponse[presenter.ConversationResponse]
	// @Router /auth/conversation [post]
	authGroup.POST("/conversation", handler.ConversationHandler.CreateConversation)

	// @Summary Send message
	// @Description Send a new message in a conversation
	// @Tags message
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body presenter.SendMessageRequest true "Send Message Request"
	// @Success 200 {object} presenter.BaseResponse[presenter.MessageResponse]
	// @Router /auth/message [post]
	authGroup.POST("/message", handler.ConversationHandler.SendMessage)

	// @Summary Get messages
	// @Description Get messages from a conversation
	// @Tags message
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param conversation_id query string true "Conversation ID"
	// @Param last_message_id query string false "Last message ID for pagination"
	// @Param limit query int false "Limit" default(20)
	// @Success 200 {object} presenter.BaseResponse[[]presenter.MessageResponse]
	// @Router /auth/message [get]
	authGroup.GET("/message", handler.ConversationHandler.GetListMessage)

	// @Summary Upload file
	// @Description Upload a file (avatar)
	// @Tags upload
	// @Accept multipart/form-data
	// @Produce json
	// @Security BearerAuth
	// @Param file formData file true "File to upload"
	// @Param bucket_name formData string true "Bucket name"
	// @Param object_name formData string true "Object name"
	// @Success 200 {object} presenter.BaseResponse[presenter.UploadResponse]
	// @Router /auth/upload [post]
	authGroup.POST("/upload", handler.UploadHandler.UploadFile)

	// @Summary Mark message as seen
	// @Description Mark a message as seen
	// @Tags message
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body presenter.SeenMessageRequest true "Seen Message Request"
	// @Success 200 {object} presenter.BaseResponse[presenter.SeenMessageResponse]
	// @Router /auth/seen-message [post]
	authGroup.POST("/seen-message", handler.ConversationHandler.SeenMessage)

	// @Summary WebSocket connection
	// @Description WebSocket endpoint for real-time communication
	// @Tags websocket
	// @Accept json
	// @Produce json
	// @Success 101 {string} string "Switching Protocols"
	// @Router /ws [get]
	s.GET("/ws", handler.WebSocketHandler.HandleWebsocket)

	// Friend routes
	authGroup.POST("/friend/:friend_id", handler.FriendHandler.SendFriendRequest)
	authGroup.POST("/friend/:friend_id/accept", handler.FriendHandler.AcceptFriendRequest)
	authGroup.POST("/friend/:friend_id/reject", handler.FriendHandler.RejectFriendRequest)
	authGroup.GET("/friend", handler.FriendHandler.GetFriends)
	authGroup.GET("/friend/requests", handler.FriendHandler.GetFriendRequests)
	authGroup.GET("/friend/received", handler.FriendHandler.GetFriendRequestsReceived)
	authGroup.DELETE("/friend/:friend_id", handler.FriendHandler.Unfriend)
}
