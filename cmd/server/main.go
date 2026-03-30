package main

import (
	_ "backend_job_assistant/cmd/server/docs" // твои сгенерированные swagger файлы
	"backend_job_assistant/internal/config"
	"backend_job_assistant/internal/controllers"
	"backend_job_assistant/internal/middlewares"
	"backend_job_assistant/internal/repositories/memory"
	"backend_job_assistant/internal/services"

	httpSwagger "github.com/swaggo/http-swagger" // swagger UI для net/http

	"log"
	"net/http"
)

// @title Job Assistant API
// @version 1.0
// @description API для твоего проекта
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()

	users := memory.NewUserMemoryRepository()

	telegramAuth := services.TelegramAuthService{
		BotToken: cfg.BotToken,
		MaxAge:   cfg.TelegramMaxAge,
	}

	sessions := services.SessionService{
		Secret: []byte(cfg.SessionSecret),
		TTL:    cfg.SessionTTL,
	}

	authController := &controllers.AuthController{
		TelegramAuth: telegramAuth,
		Sessions:     sessions,
		Users:        users,
	}

	userController := &controllers.UserController{
		Users: users,
	}

	authMiddleware := middlewares.AuthMiddleware{Sessions: sessions}

	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("POST /auth/telegram", authController.TelegramLogin)
	mux.HandleFunc("GET /me", authMiddleware.Require(userController.Me))

	log.Println("Это порт что я жду:", cfg.ServerPort)
	log.Println("Это токен что я жду:", cfg.BotToken)
	log.Println("server listening on", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, mux))

}

/*TODO это то что жду от фронта для авторизации

await fetch("/auth/telegram", {
	method: "POST",
	headers: { "Content-Type": "application/json" },
	credentials: "include",
	body: JSON.stringify({
		init_data: window.Telegram.WebApp.initData})});
*/
