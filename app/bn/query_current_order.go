package bn

import (
	bnresponse "bn_query_order/app/bn_response"
	handlerresponse "bn_query_order/app/handler_response"
	"bn_query_order/app/port"
	"encoding/json"
	"log"
	"time"

	bnrequest "bn_query_order/app/bn_request"

	bnpkg "github.com/non26/tradepkg/pkg/bn"

	"github.com/gorilla/websocket"
)

var (
	pingInterval = 180 * time.Second
	pongWait     = 20 * time.Second
)

type queryCurrentOrder struct {
	apiKey           string
	privateKey       string
	url              string
	method           string
	symbol           string
	bnws             *websocket.Conn
	positioninfo     chan []byte
	pingmessage      chan string
	positionresponse []byte
}

func NewQueryCurrentOrder() port.IBN {
	q := &queryCurrentOrder{
		positioninfo: make(chan []byte),
		pingmessage:  make(chan string),
	}
	return q
}

func (q *queryCurrentOrder) SetPrivateKey(privateKey string) {
	q.privateKey = privateKey
}

func (q *queryCurrentOrder) SetApiKey(apiKey string) {
	q.apiKey = apiKey
}

func (q *queryCurrentOrder) SetMethod(method string) {
	q.method = method
}

func (q *queryCurrentOrder) SetUrl(url string) {
	q.url = url
}

func (q *queryCurrentOrder) SetBnws(bnws *websocket.Conn) {
	q.bnws = bnws
}

func (q *queryCurrentOrder) SetSymbol(symbol string) {
	q.symbol = symbol
}

func (q *queryCurrentOrder) GetResponse() *handlerresponse.QueryOrderResponse {
	var positioninfo *bnresponse.QueryOrderResponse
	json.Unmarshal(q.positionresponse, positioninfo)
	return positioninfo.ToHandlerResponse()
}

func (q *queryCurrentOrder) GetQueryOrderRequest() bnrequest.QueryOrderRequest {
	req := bnrequest.QueryOrderRequest{
		ID:     "1",
		Method: q.method,
		Params: &bnrequest.QueryOrderRaram{
			APIKey:     q.apiKey,
			Timestamp:  time.Now().Unix() * 1000,
			RecvWindow: 30000,
		},
	}
	return req
}

func ReadPumpFromBNws(bnws *websocket.Conn, positioninfo chan []byte) {
	for {
		_, message, err := bnws.ReadMessage()
		if err != nil {
			log.Println("read bn ws error:", err)
			return
		}
		if message != nil {
			positioninfo <- message
		}
	}
}

// func (q *queryCurrentOrder) WritePumpToClientws() {
// 	for {
// 		select {
// 		case message := <-q.pingmessage:
// 			q.clientws.WriteMessage(websocket.PongMessage, []byte(message))
// 		default:
// 			continue
// 		}
// 	}
// }

func (q *queryCurrentOrder) Run() {
	bnDial, _, err := websocket.DefaultDialer.Dial(
		q.url,
		nil,
	)
	if err != nil {
		log.Println("dial bn ws error:", err)
		return
	}
	defer func() {
		println("defer close")
		bnDial.Close()
		close(q.positioninfo)
		close(q.pingmessage)
	}()
	q.SetBnws(bnDial)

	// bnDial.SetPingHandler(func(appData string) error {
	// 	log.Println("ping:", appData)
	// 	q.pingmessage <- appData
	// 	return nil
	// })

	req := q.GetQueryOrderRequest()
	bnsign := bnpkg.NewSignEd25519[bnrequest.QueryOrderRaram](q.privateKey)
	sig, err := bnsign.Sign(req.Params, "signature")
	if err != nil {
		log.Println("Ed25519 error:", err)
		return
	}
	req.Params.SetSignature(sig)

	// q.bnws.SetPingHandler(func(appData string) error {
	// 	log.Println("ping:", appData)
	// 	q.pingmessage <- appData
	// 	return nil
	// })
	j, _ := json.Marshal(req)
	err = q.bnws.WriteMessage(websocket.TextMessage, []byte(j))
	// err = q.bnws.WriteJSON(req)
	if err != nil {
		log.Println("write bn ws error:", err)
		return
	}

	go ReadPumpFromBNws(bnDial, q.positioninfo)

	counting := 0
	ticker := time.NewTicker(2 * time.Second)
	for {
		if counting == 10 {
			break
		}
		select {
		case <-ticker.C:
			counting++
		case message := <-q.positioninfo:
			log.Println("positioninfo:", string(message))
			counting++
		}
	}

}

// func signRequest2(payload string) (string, error) {
// 	keyData := "-----BEGIN PRIVATE KEY-----\nMC4CAQAwBQYDK2VwBCIEIKwzSGBlYbCE5cvr7ggQDU//tsiurzHbZSj7MM5ai8Aa\n-----END PRIVATE KEY-----"
// 	block, rest := pem.Decode([]byte(keyData))
// 	_ = rest
// 	if block == nil || block.Type != "PRIVATE KEY" {
// 		return "", fmt.Errorf("failed to decode PEM block containing the private key")
// 	}

// 	// Parse the RSA private key
// 	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
// 	if err != nil {
// 		return "", err
// 	}
// 	// Hash the payload using SHA-256
// 	// hash := sha256.New()
// 	// hash.Write([]byte(payload))
// 	// hashed := hash.Sum(nil)
// 	pk, ok := privateKey.(ed25519.PrivateKey)
// 	_ = pk
// 	if !ok {
// 		return "", fmt.Errorf("failed to convert private key to *rsa.PrivateKey")
// 	}

// 	// Sign the hashed payload using the private key
// 	signature := ed25519.Sign(pk, []byte(payload))
// 	// Encode the signature in base64
// 	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
// 	return signatureBase64, nil
// }
