package router

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/yervsil/auth_service/internal/transport/http"
	"github.com/yervsil/auth_service/internal/transport/websocket"
)

func Routes(httpHandler *http.Handler, wsHandler *websocket.Handler) *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/sign-up", httpHandler.Sign_up()).Methods("POST")
	r.HandleFunc("/sign-in", httpHandler.Sign_in()).Methods("POST")
	r.HandleFunc("/refresh_token", httpHandler.Refresh_token()).Methods("POST")

	protectedRoutes := r.PathPrefix("/api").Subrouter()
	protectedRoutes.Use(
		httpHandler.JWTMiddleware(),
	)

	protectedRoutes.HandleFunc("/chat/{roomId}", wsHandler.Chat())

	return r
}