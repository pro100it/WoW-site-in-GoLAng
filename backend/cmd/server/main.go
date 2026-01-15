package main

import (
    "log"
    "net/http"
    "time"
    "wow-registration/internal/config"
    "wow-registration/internal/database"
    "wow-registration/internal/handlers"
    "wow-registration/internal/middleware"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    // Загрузка конфигурации
    if err := config.Load(); err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Подключение к базе данных
    if err := database.Connect(); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer database.DB.Close()
    
    // Создание Echo инстанса
    e := echo.New()
    
    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.Gzip())
    e.Use(middleware.CORS())
    e.Use(middleware.Secure())
    e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(100)))
    
    // Статические файлы
    e.Static("/static", "./frontend/static")
    e.Static("/css", "./frontend/static/css")
    e.Static("/js", "./frontend/static/js")
    e.Static("/images", "./frontend/static/images")
    
    // Рендерер шаблонов
    e.Renderer = handlers.NewTemplateRenderer()
    
    // Роуты API
    api := e.Group("/api")
    {
        api.POST("/register", handlers.RegisterHandler)
        api.POST("/validate", handlers.RegisterHTMXHandler)
        api.POST("/login", handlers.LoginHandler)
        api.POST("/password/reset", handlers.ResetPasswordHandler)
        api.GET("/status", handlers.StatusHandler)
        api.GET("/stats/realtime", handlers.RealTimeStatsHandler)
    }
    
    // Web роуты
    e.GET("/", handlers.HomeHandler)
    e.GET("/register", handlers.RegistrationPageHandler)
    e.GET("/status", handlers.StatusPageHandler)
    e.GET("/rules", handlers.RulesPageHandler)
    e.GET("/players", handlers.OnlinePlayersHandler)
    
    // HTMX эндпоинты
    htmx := e.Group("/htmx")
    {
        htmx.POST("/validate/username", handlers.ValidateUsernameHandler)
        htmx.POST("/validate/email", handlers.ValidateEmailHandler)
        htmx.GET("/online-players", handlers.OnlinePlayersHTMXHandler)
        htmx.GET("/server-stats", handlers.ServerStatsHTMXHandler)
    }
    
    // Запуск сервера
    port := ":" + config.AppConfig.Server.Port
    s := &http.Server{
        Addr:         port,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    log.Printf("🚀 Server starting on http://localhost%s", port)
    log.Printf("📊 Environment: %s", config.AppConfig.Server.Environment)
    
    if err := e.StartServer(s); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
