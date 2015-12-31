package telegram

// Update is https://core.telegram.org/bots/api#update
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is https://core.telegram.org/bots/api#message
type Message struct {
	MessageID             int64       `json:"message_id"`
	From                  *User       `json:"from"`
	Date                  int64       `json:"date"`
	Chat                  *Chat       `json:"chat"`
	ForwardFrom           *User       `json:"forward_from"`
	ForwardDate           int64       `json:"forward_date"`
	ReplyToMessage        *Message    `json:"reply_to_message"`
	Text                  string      `json:"text"`
	Audio                 *Audio      `json:"audio"`
	Document              *Document   `json:"document"`
	Photo                 []PhotoSize `json:"photo"`
	Sticker               *Sticker    `json:"sticker"`
	Video                 *Video      `json:"video"`
	Voice                 *Voice      `json:"voice"`
	Caption               string      `json:"caption"`
	Contact               *Contact    `json:"contact"`
	Location              *Location   `json:"location"`
	NewChatParticipant    *User       `json:"new_chat_participant"`
	LeftChatParticipant   *User       `json:"left_chat_participant"`
	NewChatTtitle         string      `json:"new_chat_title"`
	NewChatPhoto          []PhotoSize `json:"new_chat_photo"`
	DeleteChatPhoto       bool        `json:"delete_chat_photo"`
	GroupChatCreated      bool        `json:"group_chat_created"`
	SupergroupChatCreated bool        `json:"supergroup_chat_created"`
	ChannelChatCreated    bool        `json:"channel_chat_created"`
	MigrateToChatID       int64       `json:"migrate_to_chat_id"`
	MigrateFromChatID     int64       `json:"migrate_from_chat_id"`
}

// User is https://core.telegram.org/bots/api#user
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Chat is https://core.telegram.org/bots/api#chat
type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Username string `json:"username"`
	FistName string `json:"first_name"`
	LastName string `json:"last_name"`
}

// Audio is https://core.telegram.org/bots/api#audio
type Audio struct{}

// Document is https://core.telegram.org/bots/api#document
type Document struct{}

// PhotoSize is https://core.telegram.org/bots/api#photosize
type PhotoSize struct{}

// Sticker is https://core.telegram.org/bots/api#sticker
type Sticker struct{}

// Video is https://core.telegram.org/bots/api#video
type Video struct{}

// Voice is https://core.telegram.org/bots/api#voice
type Voice struct{}

// Contact is https://core.telegram.org/bots/api#contact
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FistName    string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      int64  `json:"user_id"`
}

// Location is https://core.telegram.org/bots/api#location
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ReplyKeyboardMarkup is https://core.telegram.org/bots/api#replykeyboardmarkup
type ReplyKeyboardMarkup struct {
	Keyboard        []string `json:"keyboard"`
	ResizeKeyboard  bool     `json:"resize_keyboard"`
	OneTimeKeyboard bool     `json:"one_time_keyboard"`
	Selective       bool     `json:"selective"`
}

// ReplyKeyboardHide is https://core.telegram.org/bots/api#replykeyboardhide
type ReplyKeyboardHide struct {
	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"`
}

// ForceReply is https://core.telegram.org/bots/api#forcereply
type ForceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"selective"`
}

// SetWebhook holds the data necessary for the setWebhook API call.
type SetWebhook struct {
	URL string `json:"url"`
}

// SendMessage holds the data necessary for the sendMessage API call.
type SendMessage struct {
	ChatID                int64       `json:"chat_id"`
	Text                  string      `json:"text"`
	ParseMode             string      `json:"parse_mode"`
	DisableWebPagePreview string      `json:"disable_web_page_preview"`
	ReplyToMessageID      int64       `json:"reply_to_message_id"`
	ReplyMarkup           interface{} `json:"reply_markup"`
}
