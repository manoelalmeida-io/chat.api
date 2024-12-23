package main

import (
	"chat_api/internal/amqp"
	"chat_api/internal/configuration"
	"chat_api/internal/event"
	"chat_api/internal/handler"
	"chat_api/internal/jwt"
	"chat_api/internal/persistence"
	"chat_api/internal/repository"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	config := viper.AllSettings()
	configuration.ResolveAndUpdateAllSettings(config)

	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("mysql://%v:%v@tcp(%v:3306)/%v",
			viper.Get("mysql.user"), viper.Get("mysql.password"), viper.Get("mysql.host"), viper.Get("mysql.database")),
	)

	if err != nil {
		log.Fatalf("Error connecting to database while running migrations: %v", err)
	}

	log.Print("Running database migrations.")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error accuried appling migrations to database: %v", err)
	}

	log.Println("Migrations applied successfully.")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fmt.Println(viper.GetString("client-base-url"))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{viper.GetString("client-base-url")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	db := persistence.CreateConnection()
	amqpChannel := amqp.CreateConnectionRabbitmq()

	userRepository := repository.NewUserRepository(db)
	userContactRepository := repository.NewUserContactRepository(db)
	chatRepository := repository.NewChatRepository(db)
	chatMessageRepository := repository.NewChatMessageRepository(db)

	userTokenConverter := jwt.NewUserTokenConverter(userRepository)

	e.Use(echojwt.WithConfig(echojwt.Config{
		KeyFunc:       jwt.KeyFunc,
		SigningMethod: "RS256",
		TokenLookup:   "header:Authorization:Bearer ",
	}))
	e.Use(userTokenConverter.UserTokenConverterMiddleware)

	eventPublisher := event.NewEventPublisher(amqpChannel)
	eventConsumer := event.NewEventConsumer(chatRepository, chatMessageRepository)

	q, err := amqpChannel.QueueDeclare(
		"chat.message.sent.queue", // name
		false,                     // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		log.Fatal("failed to declare queue", err)
	}

	msgs, err := amqpChannel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal("failed to register consumer", err)
	}

	go func() {
		for d := range msgs {
			eventConsumer.ReceiveMessageSent(d)
		}
	}()

	userHandler := handler.NewUserHandler(userRepository, userContactRepository)
	chatHandler := handler.NewChatHandler(eventPublisher, userRepository)

	e.POST("/users/sign-in", userHandler.SignInHandler)
	e.GET("/users/contacts", userHandler.FindContactsByUserHandler)
	e.POST("/users/contacts", userHandler.AddContactHandler)
	e.GET("/users/contacts/:id", userHandler.GetContactByIdHandler)
	e.PUT("/users/contacts/:id", userHandler.UpdateContactHandler)
	e.DELETE("/users/contacts/:id", userHandler.DeleteContactHandler)

	e.POST("/chats/send", chatHandler.SendMessageHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
