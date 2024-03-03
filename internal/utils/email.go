package handlers

import (
	"SantaWeb/internal/db"
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/smtp"

	"go.mongodb.org/mongo-driver/bson"
)

func GenerateConfirmationCode() (string, error) {
    b := make([]byte, 4)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("%x", b), nil
}
func SendConfirmationEmail(to, code string) error {
    from := "220680@astanait.edu.kz"
    password := "1MToT3pTDm0Ah" 
    smtpHost := "smtp.office365.com"
    smtpPort := "587"

    
    message := []byte("To: " + to + "\r\n" +
        "Subject: Подтверждение регистрации\r\n" +
        "\r\n" +
        "Ваш код подтверждения: " + code + "\r\n")

    // Аутентификация
    auth := smtp.PlainAuth("", from, password, smtpHost)

 
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
    if err != nil {
        return err
    }

    return nil
}

func ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "Confirmation code is required", http.StatusBadRequest)
        return
    }

    collection := db.Client.Database("SantaWeb").Collection("volunteers")
    filter := bson.M{"confirmationCode": code}
    update := bson.M{"$set": bson.M{"isEmailConfirmed": true}}

    result, err := collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        http.Error(w, "Failed to confirm email", http.StatusInternalServerError)
        return
    }

    if result.ModifiedCount == 0 {
        http.Error(w, "Invalid confirmation code", http.StatusBadRequest)
        return
    }


    http.Redirect(w, r, "/email-confirmed-successfully", http.StatusSeeOther)
}
