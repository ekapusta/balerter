package telegram

import (
	"github.com/balerter/balerter/internal/alert/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

type apiMock struct {
	mock.Mock
}

func (m *apiMock) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	args := m.Called(c)

	msg := tgbotapi.Message{}
	err := args.Error(1)
	return msg, err
}

func TestSend(t *testing.T) {
	m := &apiMock{}
	m.On("Send", mock.Anything).Return(nil)

	tg := &Telegram{
		api:    m,
		logger: zap.NewNop(),
		chatID: 42,
	}

	mes := &message.Message{
		Level:     "foo",
		AlertName: "bar",
		Text:      "baz",
		Fields:    []string{"f1", "f2"},
		Image:     "",
	}

	err := tg.Send(mes)
	require.NoError(t, err)
}
