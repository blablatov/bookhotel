// Сервис бронирования номеров в отеле
// В предметной области выделены два понятия:
// Order — заказ, который включает в себя даты бронированияи контакты пользователя
// RoomAvailability — количество свободных номеров на конкретный день

package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	crtFile = filepath.Join(".", "certs", "server.crt")
	keyFile = filepath.Join(".", "certs", "server.key")
)

var Orders = []Order{}

var (
	mu    sync.Mutex
	cache = make(map[string]Order)
)

func main() {

	log.SetPrefix("Server event: ")
	log.SetFlags(log.Lshortfile)

	// Мультиплексор запросов. Router of http-requests.
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", createOrder)
	mux.HandleFunc("/other_func", otherFunc)

	LogInfo("Server listening on localhost:8443")

	err := http.ListenAndServeTLS(":8443", crtFile, keyFile, mux)
	//err := http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		LogInfo("Server closed")
	} else if err != nil {
		//log.Fatalf("Server failed: %v", err) // Функция log.Fatalf выводит сообщение и вызывает os.Exit(l)
		LogErrorf("Server failed: %s", err)
		os.Exit(1)
	}
}

// Обработчик запросов. Handler of requests
func createOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder Order
	json.NewDecoder(r.Body).Decode(&newOrder)

	daysToBook := daysBetween(newOrder.From, newOrder.To, newOrder.UserEmail)
	log.Printf("daysToBook: %v", daysToBook)
	if daysToBook == nil {
		LogErrorf("Hotel room is booked for selected dates:\n%v", newOrder)
		return
	}

	unavailableDays := make(map[time.Time]struct{})
	for _, day := range daysToBook {
		unavailableDays[day] = struct{}{}
	}

	for _, dayToBook := range daysToBook {
		for i, availability := range Availability {
			if !availability.Date.Equal(dayToBook) && availability.Quota < 1 {
				continue
			}
			availability.Quota -= 1
			log.Println("availability.Quota: ", availability.Quota)
			Availability[i] = availability
			delete(unavailableDays, dayToBook)
		}
	}

	if len(unavailableDays) != 0 {
		http.Error(w, "Hotel room is not available for selected dates", http.StatusInternalServerError)
		LogErrorf("Hotel room is not available for selected dates:\n%v\n%v", newOrder, unavailableDays)
		return
	}

	Orders = append(Orders, newOrder)
	log.Println("Orders: \n", Orders)

	// Мапа для хранения заказов в памяти сервиса. Map like dBase
	// Точка подключения внешней базы данных. Exchange point with dBase
	mu.Lock()
	for _, v := range Orders {
		cache[v.UserEmail] = v
		cache[v.From.GoString()] = v
		cache[v.To.GoString()] = v
		cache[v.HotelID] = v
		cache[v.RoomID] = v
	}
	mu.Unlock()
	log.Printf("dBase-cache all Orders: %v", cache)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)

	LogInfo("Order successfully created: %v\n", newOrder)
}

func otherFunc(w http.ResponseWriter, r *http.Request) {
	// ... // логика другой функции http мультиплексора
}
