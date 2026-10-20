package meshio

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func readFile(name string) ([]byte, error) { return os.ReadFile(filepath.FromSlash(name)) }

type signatureVectors struct {
	Tolerance int `json:"tolerance_seconds"`
	Vectors   []struct {
		Name      string `json:"name"`
		Secret    string `json:"secret"`
		Timestamp int64  `json:"timestamp"`
		Now       int64  `json:"now"`
		Body      string `json:"body"`
		Signature string `json:"signature"`
		Valid     bool   `json:"valid"`
	} `json:"vectors"`
}

func TestWebhookSignatureVectors(t *testing.T) {
	raw, err := readFile("testdata/webhook-signatures.json")
	if err != nil {
		t.Fatal(err)
	}
	var sv signatureVectors
	if err := json.Unmarshal(raw, &sv); err != nil {
		t.Fatal(err)
	}
	if len(sv.Vectors) == 0 {
		t.Fatal("sem vetores")
	}
	for _, v := range sv.Vectors {
		v := v
		t.Run(v.Name, func(t *testing.T) {
			h := http.Header{}
			h.Set("Meshio-Timestamp", itoa(v.Timestamp))
			h.Set("Meshio-Signature", v.Signature)
			h.Set("Meshio-Event-Id", "evt_01test")
			ev, err := VerifyWebhook([]byte(v.Body), h, v.Secret,
				WithTolerance(time.Duration(sv.Tolerance)*time.Second),
				WithClock(func() time.Time { return time.Unix(v.Now, 0) }))
			if v.Valid {
				if err != nil {
					t.Fatalf("esperava válido: %v", err)
				}
				if ev.EventID != "evt_01test" || ev.EventType != "webhook.test" {
					t.Fatalf("evento inesperado: %+v", ev)
				}
				return
			}
			var wv *WebhookVerificationError
			if !errors.As(err, &wv) || !errors.Is(err, ErrWebhookVerification) {
				t.Fatalf("esperava WebhookVerificationError, veio %v", err)
			}
		})
	}
}

func TestWebhookHeadersCaseInsensitive(t *testing.T) {
	raw, _ := readFile("testdata/webhook-signatures.json")
	var sv signatureVectors
	_ = json.Unmarshal(raw, &sv)
	v := sv.Vectors[0]
	// mapa montado à mão, com nomes em minúsculas (sem canonicalização)
	h := map[string][]string{"meshio-timestamp": {itoa(v.Timestamp)}, "MESHIO-SIGNATURE": {v.Signature}}
	if _, err := ConstructEvent([]byte(v.Body), h, v.Secret, WithClock(func() time.Time { return time.Unix(v.Now, 0) })); err != nil {
		t.Fatal(err)
	}
	// assinaturas em várias linhas de cabeçalho
	h2 := map[string][]string{"Meshio-Timestamp": {itoa(v.Timestamp)}, "Meshio-Signature": {"v1=00", v.Signature}}
	if err := VerifyWebhookSignature([]byte(v.Body), h2, v.Secret, WithClock(func() time.Time { return time.Unix(v.Now, 0) })); err != nil {
		t.Fatal(err)
	}
	if err := VerifyWebhookSignature([]byte(v.Body), map[string][]string{}, v.Secret); !errors.Is(err, ErrWebhookVerification) {
		t.Fatalf("cabeçalhos ausentes devem falhar: %v", err)
	}
}

func TestEventExamplesDecode(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("api", "events", "examples", "*.json"))
	if err != nil || len(files) == 0 {
		t.Fatalf("sem exemplos: %v", err)
	}
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		ev, err := ParseWebhookEvent(raw)
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		d, err := ev.ParsedData()
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if _, isMap := d.(map[string]any); isMap {
			t.Errorf("%s: tipo %s sem estrutura tipada", f, ev.EventType)
		}
		if ev.EventID == "" || ev.OccurredAt.IsZero() {
			t.Errorf("%s: envelope incompleto", f)
		}
	}
	raw, _ := os.ReadFile(filepath.Join("api", "events", "examples", "message.received.json"))
	ev, _ := ParseWebhookEvent(raw)
	d, _ := ev.ParsedData()
	mr, ok := d.(*MessageReceivedData)
	if !ok || !strings.HasPrefix(mr.MessageID, "msg_") || mr.From == "" {
		t.Fatalf("message.received mal decodificado: %+v", d)
	}
}

func TestUnknownEventType(t *testing.T) {
	ev, err := ParseWebhookEvent([]byte(`{"event_id":"evt_1","event_type":"future.thing","occurred_at":"2026-10-20T12:00:00Z","business_id":null,"channel_id":null,"data":{"a":1},"extra":true}`))
	if err != nil {
		t.Fatal(err)
	}
	d, err := ev.ParsedData()
	m, ok := d.(map[string]any)
	if err != nil || !ok || m["a"] != float64(1) {
		t.Fatalf("tipo desconhecido deveria virar mapa: %v %v", d, err)
	}
}

func itoa(n int64) string { return strconv.FormatInt(n, 10) }
