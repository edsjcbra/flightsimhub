package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/edsjcbra/flightsimhub/config"
	"github.com/edsjcbra/flightsimhub/internal/database"
	"github.com/edsjcbra/flightsimhub/internal/routes"
)

func main() {
	// 1️⃣ Carregar configuração
	config.LoadConfig()

	// 2️⃣ Conectar ao banco
	database.Connect()
	defer database.Close()

	// 3️⃣ Iniciar o router Gin
	router := gin.Default()

	// 4️⃣ Registrar rotas
	routes.RegisterRoutes(router)

	// 5️⃣ Iniciar servidor
	port := config.AppConfig.Port
	log.Printf("🌍 Server running on http://localhost:%s\n", port)
	router.Run(fmt.Sprintf(":%s", port))
}
