/*              Sentry Bot                */
/*   Protect your real telegramm account  */

package main

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api" //for compact
)

type Bot struct {
	tgApi *tg.BotAPI
}

const botToken = ""      // YOUR TELEGRAM BOT TOCKEN; SEE @BotFather
const ownerID = 00000000 // YOUR ACCOUNT ID; SEE @userinfobot
// (dont confuse with same adress crypto chanell)

func main() {

	InitDb() // SQLite
	bot, err := tg.NewBotAPI(botToken)
	if err != nil {
		fmt.Println(err)
	}
	msgHub := Bot{
		tgApi: bot,
	}

	bot.Debug = false
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		msgHub.msgHub(update.Message)
	}
}

func (Hub *Bot) msgHub(msg *tg.Message) {

	if msg.From.ID != ownerID { //Send msg to owner
		msgToOwn := tg.NewForward(int64(ownerID), msg.Chat.ID, msg.MessageID)
		msgInBot, _ := Hub.tgApi.Send(msgToOwn)
		/*go */ SaveToDb(msgInBot.MessageID, msg.Chat.ID) // Database write

	} else { //Recive msg from owner

		if msg.ReplyToMessage == nil { // quote is reqired
			missReply := tg.NewMessage(int64(ownerID), "Пропущено цитирование сообщения") //en: missing quote for message
			missReply.ReplyToMessageID = msg.MessageID
			Hub.tgApi.Send(missReply)
			return
		}
		searchResult := SearchInDb(msg.ReplyToMessage.MessageID)
		if searchResult == int64(ownerID) {
			missMsg := tg.NewMessage(int64(ownerID), "Не найден автор сообщения")
			missMsg.ReplyToMessageID = msg.MessageID
			Hub.tgApi.Send(missMsg)
			return
		}
		msgFromOwn := tg.NewMessage(SearchInDb(msg.ReplyToMessage.MessageID), msg.Text) //Read usrID in db
		Hub.tgApi.Send(msgFromOwn)
	}
}
