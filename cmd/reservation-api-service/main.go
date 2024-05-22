package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wac24-xbublavy-xskriba/xskriba-xbublavy-reservation-webapi/api"
	"github.com/wac24-xbublavy-xskriba/xskriba-xbublavy-reservation-webapi/internal/db_service"
	"github.com/wac24-xbublavy-xskriba/xskriba-xbublavy-reservation-webapi/internal/reservation"

	"time"

	"github.com/gin-contrib/cors"
)

func main() {
    log.Printf("Server started")
    port := os.Getenv("RESERVATION_API_PORT")
    if port == "" {
        port = "8080"
    }
    environment := os.Getenv("RESERVATION_API_ENVIRONMENT")
    if !strings.EqualFold(environment, "production") { // case insensitive comparison
        gin.SetMode(gin.DebugMode)
    }
    engine := gin.New()
    engine.Use(gin.Recovery())

		    corsMiddleware := cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
        AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
        ExposeHeaders:    []string{""},
        AllowCredentials: false,
        MaxAge: 12 * time.Hour,
    })
    engine.Use(corsMiddleware)

	// setup context update  middleware
    dbServiceAmbulance := db_service.NewMongoService[reservation.Ambulance](db_service.MongoServiceConfig{
        Collection: "ambulance",
    })
    dbServicePatient := db_service.NewMongoService[reservation.Patient](db_service.MongoServiceConfig{
        Collection: "patient",
    })
    defer dbServiceAmbulance.Disconnect(context.Background())
    engine.Use(func(ctx *gin.Context) {
        ctx.Set("db_service_ambulance", dbServiceAmbulance)
        ctx.Set("db_service_patient", dbServicePatient)
        ctx.Next()
    })

    // request routings
		reservation.AddRoutes(engine)

    engine.GET("/openapi", api.HandleOpenApi)
    engine.Run(":" + port)
}