package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/thebotguys/signalr"
)

func main()  {
	apiKey := "kkk"
	apiSecret := "sdsd"

	payload := "abc"
	// Payload SHA512 to hex encoding
	payloadSum := sha512.Sum512([]byte(payload))
	payloadHash := hex.EncodeToString(payloadSum[:])

	// Unix timestamp in milliseconds
	nonce := time.Now().Unix()*1000

	// All of the signature elemnts must be parsed as strings in this array.
	preSignatura := []string{strconv.Itoa(int(nonce)), payloadHash}
	signaturePayload := strings.Join(preSignatura, "")

	mac := hmac.New(sha512.New, []byte(apiSecret))
	_, err := mac.Write([]byte(signaturePayload))
	sig := hex.EncodeToString(mac.Sum(nil))

	client := signalr.NewWebsocketClient()

	client.OnClientMethod = func(hub, method string, arguments []json.RawMessage) {
		fmt.Println("Message Received: ")
		fmt.Println("HUB: ", hub)
		fmt.Println("METHOD: ", method)
		fmt.Println("ARGUMENTS: ", arguments)
	}
	client.OnMessageError = func (err error) {
		fmt.Println("ERROR OCCURRED: ", err)
	}

	err = client.Connect("https", "socket-v3.bittrex.com", []string{"c3"})
	if err != nil {
		fmt.Println("Connect fail:", err)
	}

	//fmt.Printf("client: %+v\n", client)

	raw, err := client.CallHub("c3", "Authenticate", apiKey, nonce, payload, sig)
	if err != nil {
		fmt.Println("CallHub fail:", err)
	}

	fmt.Println("Authenticate raw:", string(raw))

	raw, err = client.CallHub("c3", "Subscribe", "heartbeat")
	if err != nil {
		fmt.Println("CallHub fail:", err)
	}

	fmt.Println("CallHub raw:", string(raw))

	client.Close()

	fmt.Println("ok", err)
}
