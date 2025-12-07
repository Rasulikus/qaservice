package app

import (
	"net"
	"net/http"
	"path/filepath"

	answerHandler "github.com/Rasulikus/qaservice/internal/api/handler/answer"
	questionHandler "github.com/Rasulikus/qaservice/internal/api/handler/question"
	"github.com/Rasulikus/qaservice/internal/config"
	"github.com/Rasulikus/qaservice/internal/repository"
	answerRepository "github.com/Rasulikus/qaservice/internal/repository/answer"
	questionRepository "github.com/Rasulikus/qaservice/internal/repository/question"
	answerService "github.com/Rasulikus/qaservice/internal/service/answer"
	questionService "github.com/Rasulikus/qaservice/internal/service/question"
)

// App инициализирует все уровни приложения и возвращает готовый HTTP-сервер.
func App(cfg *config.Config) *http.Server {
	db, err := repository.NewDB(&cfg.DB, filepath.Join("migrations"))
	if err != nil {
		panic(err)
	}

	qRepo := questionRepository.NewRepository(db)
	aRepo := answerRepository.NewRepository(db)

	qService := questionService.NewService(qRepo)
	aService := answerService.NewService(aRepo, qRepo)

	mux := http.NewServeMux()

	qh := questionHandler.New(qService, aService)
	qh.Register(mux)

	ah := answerHandler.New(aService)
	ah.Register(mux)

	return &http.Server{
		Addr:    net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler: mux,
	}
}
