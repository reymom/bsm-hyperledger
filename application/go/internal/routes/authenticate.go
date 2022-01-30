package routes

import (
	"fmt"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[generateBuyerRoutes]")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[generateBuyerRoutes]")
}
