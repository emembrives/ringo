package telegram

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

const (
	telegramAPI = "https://api.telegram.org/bot%s/%s"
	webhookPath = "/webhook/%s/"
)

type apiResponse struct {
	Ok          bool
	Description string
}

// WebhookBot is a HTTP server implementing the Telegram bot webhook.
type WebhookBot struct {
	Token    string
	BasePath string

	webhookToken string
	logic        BotLogic

	updateChan  chan Update
	messageChan chan SendMessage
}

// NewWebhookBot creates a new server.
func NewWebhookBot(config Config, logic BotLogic) *WebhookBot {
	ws := &WebhookBot{
		Token:       config.Token,
		BasePath:    config.BasePath,
		logic:       logic,
		updateChan:  make(chan Update),
		messageChan: make(chan SendMessage),
	}
	ws.setWebhookToken()
	logic.SetIncomingChannel(ws.updateChan)
	logic.SetOutgoingChannel(ws.messageChan)
	return ws
}

// Run runs the server. This method is long-running.
func (ws *WebhookBot) Run() {
	ws.setWebhook()
	r := mux.NewRouter()
	r.HandleFunc("/webhook/{webhookToken}/", ws.WebhookHandler)
	http.Handle("/", r)
	go ws.logic.Run()
	go ws.sendIncomingMessages()
	log.Fatal(http.ListenAndServe("127.0.0.1:9100", nil))
}

func (ws *WebhookBot) setWebhookToken() {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("Unable to get randomness: %v", err)
	}

	ws.webhookToken = base32.StdEncoding.EncodeToString(b)
}

// Request makes a request to the Telegram API.
func (ws *WebhookBot) Request(method string, data interface{}) error {
	u, err := url.Parse(fmt.Sprintf(telegramAPI, ws.Token, method))
	if err != nil {
		log.Fatalf("Unable to create request url: %v", err)
	}
	log.Printf("Sending request %s", u.String())
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Error while serializing the request data: %s", err)
	}

	resp, err := http.Post(u.String(), "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("Error while sending the request: %s", err)
	}
	decoder := json.NewDecoder(resp.Body)
	apiResp := &apiResponse{}
	decoder.Decode(apiResp)
	if !apiResp.Ok {
		return fmt.Errorf("%+v", apiResp)
	}
	return nil
}

func (ws *WebhookBot) setWebhook() {
	d := SetWebhook{
		URL: fmt.Sprintf(ws.BasePath+webhookPath, ws.webhookToken),
	}
	if err := ws.Request("setWebhook", d); err != nil {
		log.Printf("Unable to set the webhook: %v", err)
	}
}

// WebhookHandler handles a request from the Telegram bot API.
func (ws *WebhookBot) WebhookHandler(w http.ResponseWriter,
	req *http.Request) {
	log.Printf("New request: %s", req.URL.String())
	vars := mux.Vars(req)
	webhookToken := vars["webhookToken"]
	if webhookToken != ws.webhookToken {
		log.Printf("Wrong webhook token: %s vs %s", webhookToken, ws.webhookToken)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	decoder := json.NewDecoder(req.Body)
	var update Update
	err := decoder.Decode(&update)
	if err != nil {
		log.Printf("Error while decoding update: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("%+v", update)
	log.Printf("From: %+v", update.Message.From)
	log.Printf("Chat: %+v", update.Message.Chat)

	w.WriteHeader(http.StatusAccepted)
}

func (ws *WebhookBot) sendIncomingMessages() {
	for {
		m := <-ws.messageChan
		ws.Request("sendMessage", m)
	}
}
