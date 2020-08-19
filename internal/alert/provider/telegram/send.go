package telegram

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"net/http"
	"os"
	"strings"
)

func (tg *Telegram) Send(mes *message.Message) error {
	tg.logger.Debug("tg send message")

	if mes.Image != "" {
		file, err := getPhotoFile(mes.Image)
		if err != nil {
			return fmt.Errorf("error get photo file, %w", err)
		}

		msg := tgbotapi.NewPhotoUpload(tg.chatID, file)
		tg.api.Send(msg)
	}

	textMsg := tgbotapi.NewMessage(tg.chatID, mes.Text)
	tg.api.Send(textMsg)

	return nil
}

func getPhotoFile(url string) (fileName string, err error) {
	tokens := strings.Split(url, "/")

	_ = os.Mkdir("tmp", os.ModePerm)

	fileName = "tmp/" + tokens[len(tokens)-1]

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	return fileName, nil
}
