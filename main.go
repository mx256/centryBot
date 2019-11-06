/*              Sentry Bot                */
/*   Protect your real telegramm account  */

package main

import (
	"fmt"
	"encoding/json"
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

	if msg.From.ID != ownerID { // Forward msg to owner
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
		searchResult, err := SearchInDb(msg.ReplyToMessage.MessageID) //Read usrId in db
		if err != nil {
			missMsg := tg.NewMessage(int64(ownerID), "Не найден автор сообщения")//en: missing autor of message
			missMsg.ReplyToMessageID = msg.MessageID
			Hub.tgApi.Send(missMsg)
			return
		}
		Hub.UniSender(searchResult, msg, 0)
	}
}

func (tlgrm *Bot)UniSender(usrId int64, msg *tg.Message, msgReply int){
	
	if msg.Sticker!=nil{
		stickerMsg:=tg.NewStickerShare(usrId, msg.Sticker.FileID)
		if msgReply != 0{
			stickerMsg.ReplyToMessageID = msgReply
		}
		tlgrm.tgApi.Send(stickerMsg)
	}
	if msg.Photo!=nil{	
		photoSize:=*msg.Photo
		photoMsg:=tg.NewPhotoShare(usrId, photoSize[0].FileID)
		photoMsg.Caption = msg.Caption
		if msgReply != 0{
			photoMsg.ReplyToMessageID = msgReply
		}
		tlgrm.tgApi.Send(photoMsg)
	}
	if msg.Video!=nil{
		videoMsg:=tg.NewVideoShare(usrId, msg.Video.FileID)
		videoMsg.Caption = msg.Caption
		if msgReply != 0{
			videoMsg.ReplyToMessageID = msgReply
		}
		tlgrm.tgApi.Send(videoMsg)
	}
	if msg.Animation!=nil{
		msgFromOwn := tg.NewMessage(int64(ownerID), "Animation is not supported yet") 
		tlgrm.tgApi.Send(msgFromOwn)
	}
	if msg.Audio !=nil{
		msgFromOwn := tg.NewMessage(int64(ownerID), "Audio is not supported yet") 
		tlgrm.tgApi.Send(msgFromOwn)
	}
	if msg.VideoNote!=nil{
		msgFromOwn := tg.NewMessage(int64(ownerID), "Videonote is not supported yet") 
		tlgrm.tgApi.Send(msgFromOwn)
	}
	if msg.Voice!=nil{
		msgFromOwn := tg.NewMessage(int64(ownerID), "Voice is not supported yet") 
		tlgrm.tgApi.Send(msgFromOwn)
	}
	if msg.Contact!=nil{
		msgFromOwn := tg.NewContact(usrId, msg.Contact.PhoneNumber,msg.Contact.FirstName)
		tlgrm.tgApi.Send(msgFromOwn)
	}
	if msg.Location!=nil{
		msgFromOwn := tg.NewLocation(usrId, msg.Location.Latitude, msg.Location.Longitude)
		tlgrm.tgApi.Send(msgFromOwn)	
	}
	if msg.Text!=""{
		msgFromOwn := tg.NewMessage(usrId, msg.Text) 
		tlgrm.tgApi.Send(msgFromOwn)
	}
}





	