package telegram

import(
  "errors"
  "crypto/rand"
  "encoding/base64"
  "encoding/json"
  "log"
  "fmt"
  "net/http"
  "net/url"

  "github.com/gorilla/mux"
)

type apiResponse struct {
  Ok bool
  Description string
}

// WebhookServer is a HTTP server implementing the Telegram bot webhook.
type WebhookServer struct {
  Token string
  BasePath string

  webhookToken string
}

// NewWebhookServer creates a new server.
func NewWebhookServer(config Config) (*WebhookServer) {
  ws := &WebhookServer{
    Token: config.Token,
    BasePath: config.BasePath,
  }
  ws.setWebhookToken()
  return ws
}

// Run runs the server. This method is long-running.
func (ws *WebhookServer) Run() {
  ws.setWebhook()
  r := mux.NewRouter()
	r.HandleFunc(fmt.Sprintf(webhookPath, ws.webhookToken), ws.WebhookHandler)
	http.Handle("/", r)
}

func (ws *WebhookServer) setWebhookToken() {
  b := make([]byte, 64)
  _, err := rand.Read(b)
  dieOnError(err, "Unable to get randomness: %v")
  ws.webhookToken = base64.StdEncoding.EncodeToString(b)
}

// Request makes a request to the Telegram API.
func (ws *WebhookServer) Request(method string, parameters url.Values) error {
  u, err := url.Parse(fmt.Sprintf(telegramAPI, ws.Token, method))
  if err != nil {
    log.Fatalf("Unable to create request url: %v", err)
  }
  u.RawQuery = parameters.Encode()
  resp, err := http.Get(u.String())
  if err != nil {
    log.Fatalf("Unable to set webhook: %v", err)
  }
  decoder := json.NewDecoder(resp.Body)
  apiResp := &apiResponse{}
  decoder.Decode(apiResp)
  if !apiResp.Ok {
    return errors.New(apiResp.Description)
  }
  return nil
}

func (ws *WebhookServer) setWebhook() {
  v := url.Values{}
  v.Add("url", fmt.Sprintf(ws.BasePath + webhookPath, ws.webhookToken))
  if err := ws.Request("setWebhook", v); err != nil {
    log.Fatalf("Unable to set the webhook: %v", err)
  }
}

// WebhookHandler handles a request from the Telegram bot API.
func (ws *WebhookServer) WebhookHandler(w http.ResponseWriter,
  req *http.Request) {

}
