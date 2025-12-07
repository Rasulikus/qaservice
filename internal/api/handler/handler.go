package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Rasulikus/qaservice/internal/model"
)

// DecodeJSON декодирует JSON-тело запроса в dst.
// Запрещает неизвестные поля, а при любой ошибке парсинга
// возвращает model.ErrBadRequest.
func DecodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return model.ErrBadRequest
	}
	return nil
}

// WriteJSON отправляет HTTP-ответ с указанным статусом и JSON-телом.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// WriteError преобразует доменную ошибку в HTTP-статус и публичное сообщение
// и отправляет их в виде JSON-ответа.
func WriteError(w http.ResponseWriter, err error) {
	status, pub := model.ToHTTP(err)
	WriteJSON(w, status, pub)
}
