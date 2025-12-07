package question

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Rasulikus/qaservice/internal/api/handler"
	"github.com/go-playground/validator/v10"

	"github.com/Rasulikus/qaservice/internal/model"
	"github.com/Rasulikus/qaservice/internal/service"
)

// Handler реализует HTTP-обработчики для работы с вопросами и ответами.
type Handler struct {
	questions service.QuestionService
	answers   service.AnswerService
	validate  *validator.Validate
}

// New создаёт новый Handler с переданными сервисами вопросов и ответов.
func New(questionSvc service.QuestionService, answerSvc service.AnswerService) *Handler {
	return &Handler{
		questions: questionSvc,
		answers:   answerSvc,
		validate:  validator.New(),
	}
}

// Register регистрирует HTTP-маршруты обработчика на переданном ServeMux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.Handle("/questions", http.HandlerFunc(h.handleQuestions))
	mux.Handle("/questions/", http.HandlerFunc(h.handleQuestionPath))
}

// handleQuestions обрабатывает запросы по пути /questions для списка и создания вопросов.
func (h *Handler) handleQuestions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listQuestions(w, r)
	case http.MethodPost:
		h.createQuestion(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleQuestionPath обрабатывает запросы по путям /questions/{id} и /questions/{id}/answers.
func (h *Handler) handleQuestionPath(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/questions/")
	path = strings.Trim(path, "/")
	if path == "" {
		h.handleQuestions(w, r)
		return
	}

	parts := strings.Split(path, "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		handler.WriteError(w, model.ErrBadRequest)
		return
	}

	// Обработка /questions/{id}
	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			h.getQuestion(w, r, id)
		case http.MethodDelete:
			h.deleteQuestion(w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	// Обработка /questions/{id}/answers (только POST)
	if len(parts) == 2 && parts[1] == "answers" {
		if r.Method == http.MethodPost {
			h.createAnswer(w, r, id)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	handler.WriteError(w, model.ErrNotFound)
}

// createQuestionRequest описывает тело запроса на создание вопроса.
type createQuestionRequest struct {
	Text string `json:"text" validate:"required,min=1,max=10000"`
}

// createQuestion обрабатывает POST /questions и создаёт новый вопрос.
func (h *Handler) createQuestion(w http.ResponseWriter, r *http.Request) {
	var req createQuestionRequest
	if err := handler.DecodeJSON(r, &req); err != nil {
		handler.WriteError(w, err)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		if vErr, ok := model.AsValidationError(req, err); ok {
			handler.WriteError(w, vErr)
			return
		}
		handler.WriteError(w, model.ErrBadRequest)
		return
	}

	input := service.CreateQuestionInput{Text: req.Text}
	q, err := h.questions.Create(r.Context(), input)
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	handler.WriteJSON(w, http.StatusCreated, q)
}

// listQuestions обрабатывает GET /questions и возвращает список вопросов.
func (h *Handler) listQuestions(w http.ResponseWriter, r *http.Request) {
	qs, err := h.questions.List(r.Context())
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	handler.WriteJSON(w, http.StatusOK, qs)
}

// getQuestion обрабатывает GET /questions/{id} и возвращает вопрос с ответами.
func (h *Handler) getQuestion(w http.ResponseWriter, r *http.Request, id int) {
	q, err := h.questions.GetByID(r.Context(), id)
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	handler.WriteJSON(w, http.StatusOK, q)
}

// deleteQuestion обрабатывает DELETE /questions/{id} и удаляет вопрос.
func (h *Handler) deleteQuestion(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.questions.Delete(r.Context(), id); err != nil {
		handler.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// createAnswerRequest описывает тело запроса на создание ответа.
type createAnswerRequest struct {
	UserID string `json:"user_id" validate:"required,min=1,max=255"`
	Text   string `json:"text" validate:"required,min=1,max=10000"`
}

// createAnswer обрабатывает POST /questions/{id}/answers и создаёт новый ответ на вопрос.
func (h *Handler) createAnswer(w http.ResponseWriter, r *http.Request, questionID int) {
	var req createAnswerRequest
	if err := handler.DecodeJSON(r, &req); err != nil {
		handler.WriteError(w, err)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		if vErr, ok := model.AsValidationError(req, err); ok {
			handler.WriteError(w, vErr)
			return
		}
		handler.WriteError(w, model.ErrBadRequest)
		return
	}

	input := service.CreateAnswerInput{
		QuestionID: questionID,
		UserID:     req.UserID,
		Text:       req.Text,
	}
	a, err := h.answers.Create(r.Context(), input)
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	handler.WriteJSON(w, http.StatusCreated, a)
}
