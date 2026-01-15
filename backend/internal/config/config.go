package config

import (
    "os"
    "strconv"
    "github.com/joho/godotenv"
)

type Config struct {
    Server struct {
        Port         string
        SecretKey    string
        Environment  string
        RateLimit    int
    }
    Database struct {
        Host     string
        Port     string
        User     string
        Password string
        Name     string
        Charset  string
    }
    Redis struct {
        Host     string
        Port     string
        Password string
        DB       int
    }
    Game struct {
        Expansion    int
        RealmList    string
        ServerCore   int
        MaxAccounts  int
        AllowMultiIP bool
    }
    Security struct {
        EnableCaptcha   bool
        CaptchaSecret   string
        CaptchaSiteKey  string
        Enable2FA       bool
        PasswordMinLen  int
        PasswordMaxLen  int
    }
}

var AppConfig *Config

func Load() error {
    _ = godotenv.Load()
    
    AppConfig = &Config{}
    
    // Server
    AppConfig.Server.Port = getEnv("PORT", "8080")
    AppConfig.Server.SecretKey = getEnv("SECRET_KEY", "supersecretkey")
    AppConfig.Server.Environment = getEnv("ENVIRONMENT", "development")
    AppConfig.Server.RateLimit, _ = strconv.Atoi(getEnv("RATE_LIMIT", "100"))
    
    // Database
    AppConfig.Database.Host = getEnv("DB_HOST", "localhost")
    AppConfig.Database.Port = getEnv("DB_PORT", "3306")
    AppConfig.Database.User = getEnv("DB_USER", "root")
    AppConfig.Database.Password = getEnv("DB_PASSWORD", "")
    AppConfig.Database.Name = getEnv("DB_NAME", "auth")
    AppConfig.Database.Charset = getEnv("DB_CHARSET", "utf8mb4")
    
    // Redis
    AppConfig.Redis.Host = getEnv("REDIS_HOST", "localhost")
    AppConfig.Redis.Port = getEnv("REDIS_PORT", "6379")
    AppConfig.Redis.Password = getEnv("REDIS_PASSWORD", "")
    AppConfig.Redis.DB, _ = strconv.Atoi(getEnv("REDIS_DB", "0"))
    
    // Game
    AppConfig.Game.Expansion, _ = strconv.Atoi(getEnv("GAME_EXPANSION", "2"))
    AppConfig.Game.RealmList = getEnv("REALMLIST", "logon.yourserver.com")
    AppConfig.Game.ServerCore, _ = strconv.Atoi(getEnv("SERVER_CORE", "0"))
    AppConfig.Game.MaxAccounts, _ = strconv.Atoi(getEnv("MAX_ACCOUNTS", "5"))
    AppConfig.Game.AllowMultiIP, _ = strconv.ParseBool(getEnv("ALLOW_MULTI_IP", "false"))
    
    // Security
    AppConfig.Security.EnableCaptcha, _ = strconv.ParseBool(getEnv("ENABLE_CAPTCHA", "true"))
    AppConfig.Security.CaptchaSecret = getEnv("CAPTCHA_SECRET", "")
    AppConfig.Security.CaptchaSiteKey = getEnv("CAPTCHA_SITEKEY", "")
    AppConfig.Security.Enable2FA, _ = strconv.ParseBool(getEnv("ENABLE_2FA", "false"))
    AppConfig.Security.PasswordMinLen, _ = strconv.Atoi(getEnv("PASSWORD_MIN_LEN", "4"))
    AppConfig.Security.PasswordMaxLen, _ = strconv.Atoi(getEnv("PASSWORD_MAX_LEN", "16"))
    
    return nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
