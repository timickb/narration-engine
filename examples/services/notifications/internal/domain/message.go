package domain

// Message Структура email сообщения.
type Message struct {
	MailFrom string `json:"mail_from"`
	MailTo   string `json:"mail_to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}
