package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"strings"
)

// SMTPConfig 配置 SMTP 伺服器資訊
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// SendEmail 發送郵件，夾帶檔案附件
func (s *SMTPConfig) SendEmail(to, subject string, attachments map[string][]string) error {
	from := s.Username
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	var emailBody bytes.Buffer
	writer := multipart.NewWriter(&emailBody)

	// 設定標頭
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=\"%s\"", writer.Boundary())

	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	emailBody.WriteString(message.String())

	// 寫入郵件正文，直接以 UTF-8 發送
	part, _ := writer.CreatePart(map[string][]string{"Content-Type": {"text/plain; charset=UTF-8"}})
	part.Write([]byte("今日命中告誡名單清單"))

	// 添加附件
	for fileName, fileLines := range attachments {
		attachment, _ := writer.CreateFormFile("attachment", fileName)
		// 將 []string 中的每一行內容分開並寫入檔案
		attachmentContent := strings.Join(fileLines, "\r\n") // 用換行符號連接每行
		attachment.Write([]byte(attachmentContent))
	}

	writer.Close()

	// 設定 TLS 連線
	serverAddr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	tlsConfig := &tls.Config{
		ServerName: s.Host,
	}

	conn, err := tls.Dial("tcp", serverAddr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(emailBody.Bytes())
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}
