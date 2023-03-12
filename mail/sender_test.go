package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weichunnn/neobank/util"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	to := []string{"neobankdev@gmail.com"}
	attachFiles := []string{"../README.md"}
	subject := "A test email"

	content := `
  <h1>Hi there</h1>
  <p>This is a test email</p>
  <p>Regards</p>
  <p>NeoBank Dev</p>
  `

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
