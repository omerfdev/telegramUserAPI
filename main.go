package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

const port = ":8080"
const botToken = "7105924273:AAHqk07jfhQrHyAbk1ppe_A3BrgPJOVaGas"

// Bot ve kullanıcı bilgilerini depolamak için kullanılan yapılar
type UserIDResponse struct {
	UserID int64 `json:"user_id"`
}

var (
	bot      *tb.Bot
	lastUserID int64
)

func main() {
	// Telegram botunu başlat
	b, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}
	bot = b

	// Bot mesajlarını dinleme
	bot.Handle(tb.OnText, handleMessage)

	// API endpoint'ini tanımla
	http.HandleFunc("/get_user_id", getUserID)

	// HTTP sunucusunu başlat
	fmt.Println("Server started at :8080")
	go http.ListenAndServe(port, nil)

	// Sonsuz döngüyü başlat
	select {}
}

// Kullanıcı ID'sini saklamak ve HTTP endpoint'ten döndürmek için kullanılan fonksiyon
func handleMessage(c *tb.Chat, m *tb.Message) {
	// Mesaj gönderen kullanıcı ID'sini al
	lastUserID = m.Sender.ID
	fmt.Printf("User ID for message: %d\n", lastUserID)
}

// API endpoint'ini tanımla
func getUserID(w http.ResponseWriter, r *http.Request) {
	// En son alınan kullanıcı ID'sini döndür
	response := UserIDResponse{
		UserID: lastUserID,
	}

	// JSON formatına dönüştür
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Yanıtı gönder
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}
