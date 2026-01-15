package handlers

import (
    "encoding/json"
    "net/http"
    "strings"
    "time"
    "wow-registration/internal/database"
    "wow-registration/internal/services"
    "github.com/labstack/echo/v4"
)

type RegisterRequest struct {
    Username string `json:"username" validate:"required,min=3,max=16"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=4,max=16"`
    Captcha  string `json:"captcha,omitempty"`
}

type RegisterResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Account struct {
        ID       int    `json:"id,omitempty"`
        Username string `json:"username,omitempty"`
        Email    string `json:"email,omitempty"`
    } `json:"account,omitempty"`
}

func RegisterHandler(c echo.Context) error {
    var req RegisterRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, RegisterResponse{
            Success: false,
            Message: "Invalid request format",
        })
    }
    
    // Валидация
    if err := services.ValidateUsername(req.Username); err != nil {
        return c.JSON(http.StatusBadRequest, RegisterResponse{
            Success: false,
            Message: err.Error(),
        })
    }
    
    if err := services.ValidateEmail(req.Email); err != nil {
        return c.JSON(http.StatusBadRequest, RegisterResponse{
            Success: false,
            Message: err.Error(),
        })
    }
    
    if err := services.ValidatePassword(req.Password); err != nil {
        return c.JSON(http.StatusBadRequest, RegisterResponse{
            Success: false,
            Message: err.Error(),
        })
    }
    
    // Проверка капчи
    if config.AppConfig.Security.EnableCaptcha {
        if !verifyCaptcha(req.Captcha) {
            return c.JSON(http.StatusBadRequest, RegisterResponse{
                Success: false,
                Message: "Captcha verification failed",
            })
        }
    }
    
    // Проверка существования аккаунта
    exists, err := database.AccountExists(strings.ToUpper(req.Username), strings.ToUpper(req.Email))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, RegisterResponse{
            Success: false,
            Message: "Database error",
        })
    }
    
    if exists {
        return c.JSON(http.StatusConflict, RegisterResponse{
            Success: false,
            Message: "Username or email already exists",
        })
    }
    
    // Генерация SRP6 данных
    srp6, err := services.GenerateSRP6(
        strings.ToUpper(req.Username),
        req.Password,
        config.AppConfig.Game.ServerCore,
    )
    if err != nil {
        return c.JSON(http.StatusInternalServerError, RegisterResponse{
            Success: false,
            Message: "Failed to generate secure credentials",
        })
    }
    
    // Создание аккаунта
    account := &database.Account{
        Username:  strings.ToUpper(req.Username),
        Email:     strings.ToUpper(req.Email),
        Password:  services.GenerateSHA1Hash(req.Username, req.Password),
        Salt:      srp6.Salt,
        Verifier:  srp6.Verifier,
        Expansion: config.AppConfig.Game.Expansion,
        IP:        services.GetClientIP(c.Request()),
        CreatedAt: time.Now(),
        Locked:    false,
    }
    
    if err := database.CreateAccount(account); err != nil {
        return c.JSON(http.StatusInternalServerError, RegisterResponse{
            Success: false,
            Message: "Failed to create account",
        })
    }
    
    // Ответ
    resp := RegisterResponse{
        Success: true,
        Message: "Account created successfully",
    }
    resp.Account.ID = account.ID
    resp.Account.Username = req.Username
    resp.Account.Email = req.Email
    
    return c.JSON(http.StatusCreated, resp)
}

func verifyCaptcha(response string) bool {
    // Реализация проверки капчи (hCaptcha/ReCaptcha)
    return true // Заглушка
}

// HTMX версия регистрации
func RegisterHTMXHandler(c echo.Context) error {
    username := c.FormValue("username")
    email := c.FormValue("email")
    password := c.FormValue("password")
    confirmPassword := c.FormValue("confirm_password")
    
    // Быстрая валидация на стороне сервера
    if password != confirmPassword {
        return c.HTML(http.StatusOK, `
            <div class="text-red-500 text-sm mt-1" id="password-error">
                Passwords do not match
            </div>
        `)
    }
    
    if err := services.ValidateUsername(username); err != nil {
        return c.HTML(http.StatusOK, `
            <div class="text-red-500 text-sm mt-1" id="username-error">
                `+err.Error()+`
            </div>
        `)
    }
    
    // Проверка существования
    exists, _ := database.AccountExists(strings.ToUpper(username), strings.ToUpper(email))
    if exists {
        return c.HTML(http.StatusOK, `
            <div class="text-red-500 text-sm mt-1" id="username-error">
                Username or email already exists
            </div>
        `)
    }
    
    // Если все ок - скрываем ошибки
    return c.HTML(http.StatusOK, `
        <div class="hidden" id="username-error"></div>
        <div class="hidden" id="password-error"></div>
        <div class="text-green-500 text-sm mt-1">
            ✓ All checks passed
        </div>
    `)
}
