package meshio_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/meshio-br/meshio-go"
)

func ExampleMessagesService_Send() {
	client := meshio.NewClient(meshio.WithAPIKey(os.Getenv("MESHIO_API_KEY")))
	msg, err := client.Messages.Send(context.Background(), &meshio.SendMessageParams{
		ChannelID: "chn_…",
		To:        "+5511988887777",
		Type:      "text",
		Text:      &meshio.TextParams{Body: "Olá!"},
	}, meshio.WithIdempotencyKey("pedido-1234"))
	var apiErr *meshio.APIError
	if errors.As(err, &apiErr) {
		log.Println(apiErr.Code, apiErr.Message)
		return
	}
	if err == nil {
		fmt.Println(msg.ID)
	}
}

func ExampleBusinessesService_ListAll() {
	client := meshio.NewClient()
	it := client.Businesses.ListAll(context.Background(), &meshio.BusinessesListParams{Limit: 200})
	for it.Next() {
		fmt.Println(it.Current().ID)
	}
	if err := it.Err(); err != nil {
		log.Println(err)
	}
}

func ExampleVerifyWebhook() {
	secret := os.Getenv("MESHIO_WEBHOOK_SECRET")
	http.HandleFunc("/webhooks/meshio", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		event, err := meshio.VerifyWebhook(body, r.Header, secret)
		if err != nil {
			http.Error(w, "assinatura inválida", http.StatusBadRequest)
			return
		}
		data, _ := event.ParsedData()
		switch d := data.(type) {
		case *meshio.MessageReceivedData:
			log.Printf("mensagem de %s (%s/%s)", d.From, d.Direction, d.Origin)
		case *meshio.MessageStatusChangedData:
			log.Printf("mensagem %s: %s", d.MessageID, d.Status)
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func ExampleMicrosToReais() {
	reais, _ := meshio.MicrosToReais("12500000")
	fmt.Println(reais)
	// Output: 12.50
}
