package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

const telegramAPI = "https://api.telegram.org/bot"

type GetChatResponse struct {
    Ok     bool `json:"ok"`
    Result struct {
        ID int64 `json:"id"`
    } `json:"result"`
}

func main() {
    http.HandleFunc("/get_user_id", getUserIDHandler)
    port := ":8080"
    fmt.Printf("Server listening on port %s\n", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        fmt.Println(err)
    }
}

func getUserIDHandler(w http.ResponseWriter, r *http.Request) {
    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "Username is required", http.StatusBadRequest)
        return
    }

    token := os.Getenv("TELEGRAM_BOT_TOKEN")
    if token == "" {
        http.Error(w, "Telegram bot token is required", http.StatusInternalServerError)
        return
    }

    url := fmt.Sprintf("%s%s/getChat?chat_id=@%s", telegramAPI, token, username)

    resp, err := http.Get(url)
    if err != nil {
        http.Error(w, "Failed to make request to Telegram API", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        http.Error(w, "Error from Telegram API", http.StatusInternalServerError)
        return
    }

    var getChatResponse GetChatResponse
    if err := json.NewDecoder(resp.Body).Decode(&getChatResponse); err != nil {
        http.Error(w, "Failed to parse response from Telegram API", http.StatusInternalServerError)
        return
    }

    if !getChatResponse.Ok {
        http.Error(w, "Invalid response from Telegram API", http.StatusInternalServerError)
        return
    }

    userID := getChatResponse.Result.ID

    response := map[string]int64{
        "user_id": userID,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
