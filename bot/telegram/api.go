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

type apiResponse struct {
	Ok          bool
	Description string
}

// WebhookServer is a HTTP server implementing the Telegram bot webhook.
type WebhookServer struct {
	Token    string
	BasePath string

	webhookToken string
}

// NewWebhookServer creates a new server.
func NewWebhookServer(config Config) *WebhookServer {
	ws := &WebhookServer{
		Token:    config.Token,
		BasePath: config.BasePath,
	}
	ws.setWebhookToken()
	return ws
}

// Run runs the server. This method is long-running.
func (ws *WebhookServer) Run() {
	ws.setWebhook()
	r := mux.NewRouter()
	r.HandleFunc("/webhook/{webhookToken}/", ws.WebhookHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("127.0.0.1:9100", nil))
}

func (ws *WebhookServer) setWebhookToken() {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	dieOnError(err, "Unable to get randomness: %v")
	ws.webhookToken = base32.StdEncoding.EncodeToString(b)
}

// Request makes a request to the Telegram API.
func (ws *WebhookServer) Request(method string, data interface{}) error {
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

func (ws *WebhookServer) setWebhook() {
	d := SetWebhook{
		URL: fmt.Sprintf(ws.BasePath+webhookPath, ws.webhookToken),
	}
	if err := ws.Request("setWebhook", d); err != nil {
		log.Printf("Unable to set the webhook: %v", err)
	}
}

// WebhookHandler handles a request from the Telegram bot API.
func (ws *WebhookServer) WebhookHandler(w http.ResponseWriter,
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
	w.WriteHeader(http.StatusAccepted)
}
