package main

import (
	"backend_job_assistant/internal/config"
	"backend_job_assistant/internal/controllers"
	"backend_job_assistant/internal/middlewares"
	"backend_job_assistant/internal/repositories/memory"
	"backend_job_assistant/internal/services"
	"log"
	"net/http"
)

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
	mux.HandleFunc("POST /auth/telegram", authController.TelegramLogin)
	mux.HandleFunc("GET /me", authMiddleware.Require(userController.Me))

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

/*TODO это то что жду от фронта для авторизации

await fetch("/auth/telegram", {
	method: "POST",
	headers: { "Content-Type": "application/json" },
	credentials: "include",
	body: JSON.stringify({
		init_data: window.Telegram.WebApp.initData})});
*/
