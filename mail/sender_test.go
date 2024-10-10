package mail

import (
	"github.com/WooDMaNbtw/BankApp/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `<h1>おはよ、やさし　えまいぉ　です　ね</h1>`

	to := []string{"nikita.profatilov5@gmail.com"}
	attachFiles := []string{"../README.md"}
	err = sender.SendSubject(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
