package main

import (
	"T3B/api"
	"T3B/bot_settings"
	"T3B/db"
	"context"
	"log"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	tele "gopkg.in/telebot.v3"
)

const (
	dburi  = "mongodb://localhost:27017"
	dbname = "EDUNT-TOOLS"

	usersColl = "users" // coll = group
)

func main() {
	app := fiber.New()
	_ = app

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		panic(err)
	}

	bot, err := tele.NewBot(bot_settings.Settings)
	if err != nil {
		log.Fatalf("Error connecting to bot: %v\n", err)
	}

	_ = bot_settings.Bot{
		Menu: &tele.ReplyMarkup{ResizeKeyboard: true},
		Bot:  bot,
		Run:  true,
	}

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	bot.Handle("/createuser", func(ctx tele.Context) error {
		ctx.Reply("📋 Please provide user details in JSON format: \n\n" +
			`{
    "user":"your_user",
    "email":"your_email@gmail.com",
    "password":"yourPassword"
}`)
		bot.Handle(tele.OnText, userHandler.HandleCreateUser)
		return nil
	})

	bot.Handle("/deleteuser", func(ctx tele.Context) error {
		ctx.Reply("Please Enter The ID Of your Account")

		bot.Handle(tele.OnText, userHandler.HandleDeleteUser)
		return nil
	})

	// Handle incoming text messages for user creation

	bot.Handle("/start", func(ctx tele.Context) error {
		return ctx.Reply("Hello! Type /createuser to create a new user.\nType /deleteuser to delete account")
	})

	bot.Start()
}
