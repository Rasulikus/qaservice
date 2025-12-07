package answer

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Rasulikus/qaservice/internal/model"
	"github.com/Rasulikus/qaservice/internal/service"
)

// Handler реализует HTTP-обработчики для работы с ответами.
type Handler struct {
	answers service.AnswerService
}

// New создаёт новый Handler с переданным сервисом ответов.
func New(answerSvc service.AnswerService) *Handler {
	return &Handler{answers: answerSvc}
}

// Register регистрирует HTTP-маршруты обработчика ответов на переданном ServeMux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.Handle("/answers/", http.HandlerFunc(h.handleAnswerPath))
}

// handleAnswerPath обрабатывает запросы по пути /answers/{id}.
func (h *Handler) handleAnswerPath(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/answers/")
	path = strings.Trim(path, "/")
	if path == "" {
		writeError(w, model.ErrNotFound)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, model.ErrBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getAnswer(w, r, id)
	case http.MethodDelete:
		h.deleteAnswer(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getAnswer обрабатывает GET /answers/{id} и возвращает ответ по его идентификатору.
func (h *Handler) getAnswer(w http.ResponseWriter, r *http.Request, id int) {
	a, err := h.answers.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// deleteAnswer обрабатывает DELETE /answers/{id} и удаляет ответ по его идентификатору.
func (h *Handler) deleteAnswer(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.answers.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// writeJSON отправляет JSON-ответ с указанным HTTP-статусом и данными.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError конвертирует ошибку в публичное представление и отправляет её в виде JSON-ответа.
func writeError(w http.ResponseWriter, err error) {
	status, pub := model.ToHTTP(err)
	writeJSON(w, status, pub)
}
