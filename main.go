package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
	"github.com/fermyon/spin-go-sdk/sqlite"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		cors(w)

		if r.Method == http.MethodOptions {
			return
		}

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var contact Contact
		err = json.Unmarshal(raw, &contact)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = createContact(r.Context(), contact)
		if err != nil {
			fmt.Println("ERROR: failed to add contact to db", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	})
}

func cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")
}

type Contact struct {
	Id        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	Msg       string `json:"msg" db:"msg"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}

func createContact(ctx context.Context, contact Contact) error {
	conn := db()
	defer conn.Close()

	contact.CreatedAt = time.Now().Format(time.RFC3339)
	contact.Id = uuid.NewString()

	_, err := conn.QueryxContext(ctx, "INSERT INTO contact values (?, ?, ?, ?, ?)", contact.Id, contact.Name, contact.Email, contact.Msg, contact.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func db() *sqlx.DB {
	conn := sqlite.Open("default")
	return sqlx.NewDb(conn, "sqlite")
}

func main() {}
