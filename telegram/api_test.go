package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"www.samermurad.com/piBot/telegram/models"
)

func Test(t *testing.T) {
	h := models.BotMessage{
		ChatId: 68386493,
		Text:   "Testing New Telegram Api",
	}
	ch := make(chan *models.Message)
	go SendMessage(h, ch)
	data := <-ch
	assert.NotNil(t, data, "Call failed")
	assert.NotNil(t, data.MessageId, "Call failed")
}

func TestFormattedMessage(t *testing.T) {
	msg := models.BotMessage{
		ChatId:    68386493,
		ParseMode: models.MarkdownV2,
		Text: `
		*Testing New Telegram Api MarkdownV2*
*bold \*text*
_italic \*text_
__underline__
~strikethrough~
*bold _italic bold ~italic bold strikethrough~ __underline italic bold___ bold*
[inline URL](http://www.example.com/)
[@soom](tg://user?id=68386493)
` +
			"`inline fixed-width code`" +
			"```\npre-formatted fixed-width code block```\n\n" +
			"```python\n" +
			`def sum(a,b):
				return a + b

			` +
			"```",
	}
	ch := make(chan *models.Message)
	go SendMessage(msg, ch)
	data := <-ch
	assert.NotNil(t, data, "Call failed")
	assert.NotNil(t, data.MessageId, "Call failed")
}
