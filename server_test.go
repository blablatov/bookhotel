// Выполнить тестовый запрос go test .

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestServer(t *testing.T) {

	// URL тестового метода локально. Test method
	apiUrl := "https://localhost:8443/orders"

	// Формирование параметров запроса. JSON params of request
	payload, _ := json.Marshal(struct {
		HotelID   string `json:"hotel_id"`
		RoomID    string `json:"room_id"`
		UserEmail string `json:"uemail"`
		From      string `json:"from"`
		To        string `json:"to"`
	}{
		HotelID:   "reddison",
		RoomID:    "lux",
		UserEmail: "guest@mail.ru",
		From:      "2024-01-03T00:00:00Z",
		To:        "2024-01-05T00:00:00Z",
	})

	// Подгрузка сертификата и ключа. Loads the certs
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("Сертификат и ключ не получены: %v\n", err)
	}

	// Logs CLIENT_SERVER_HANDSHAKE_TRAFFIC_SECRETS
	var w io.Writer
	w = os.Stdout

	// Форматирование запроса. Formatting of the request
	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(payload))
	// Формирование заголовков запроса. Headers of request
	req.Header.Set("Content-Type", "application/json")

	// Формирование метаданных структуры запроса. Struct of request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				KeyLogWriter:       w,
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req) // Выполнение запроса. Send of request
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	// Отложеное выполнение закрытия запроса, до выполнения метода и получения ответа
	// Defer to finished the method and got response
	defer resp.Body.Close()

	fmt.Printf("Status = %v ", resp.Status) // Статус ответа сервера. Status of response

	// Чтение данных сервера, обработка ошибок. Reads data from server, check errors
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println("\nResponse of server: \n", string([]byte(body)))
}
