package api

import (
	"T3B/bot_settings"
	"T3B/db"
	"T3B/types"
	"context"
	"encoding/json"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

type UserHandler struct {
	bot_settings.Bot
	userStore db.UserStore
}

var ()

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}

}

func (h *UserHandler) HandleGetUser(ctx tele.Context) error {
	var (
		id     = ctx.Text()
		botCtx = context.TODO()
	)

	user, err := h.userStore.GetUserByID(botCtx, id)
	if err != nil {
		return err
	}

	return ctx.Send(user)
}

func (h *UserHandler) HandleDeleteUser(ctx tele.Context) error {
	if err := h.userStore.DeleteUser(context.TODO(), ctx.Text()); err != nil {
		return ctx.Reply("❌ There is no account with this id ")
	}
	return ctx.Reply(fmt.Sprintf("⚙️ Done deleting account %s", ctx.Text()))

}

func (h *UserHandler) HandleCreateUser(ctx tele.Context) error {
	var params types.CreateUserParams
	params.TeleID = ctx.Chat().ID

	if err := json.Unmarshal([]byte(ctx.Text()), &params); err != nil {
		return ctx.Send("❌ There was an error with JSON format please try write it again")
	}

	ok, _ := h.userStore.CheckUsername(context.TODO(), params.Username)

	if !ok { // !ok == unvalid user
		return ctx.Send("❌ Username is taken ")
	}

	ok, _ = h.userStore.CheckTeleID(context.TODO(), ctx.Chat().ID)
	if !ok {
		return ctx.Send("❌ You have already registered")
	}

	if errors := params.Check(); len(errors) > 0 {
		errorMessage := "❌ There were errors in the parameters:\n"
		for _, err := range errors {
			errorMessage += "• " + err + "\n"
		}
		return ctx.Send(errorMessage)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ctx.Send("❌ There was an error creating the user")
	}
	_, err = h.userStore.CreateUser(context.TODO(), user)

	if err != nil {
		return ctx.Send("❌ There was an error creating the user")
	}

	userInserted := fmt.Sprintf(`
🎉 User created successfully! 🎉
_id : %v 
---------------------------
🆔 TeleID:  %d
👤 User:    %s
📧 Email:   %s
🔒 Password: %s
`, user.ID.Hex(), params.TeleID, params.Username, params.Email, types.MaskPassword(params.Password))

	// Send the userInserted data back to the user
	return ctx.Send(userInserted)

}
