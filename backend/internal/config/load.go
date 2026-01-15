package config

import (
    "os"
    "strconv"
    "strings"
    "time"
    "github.com/joho/godotenv"
)

// LoadConfig загружает конфигурацию из .env файла и переменных окружения
func LoadConfig() error {
    // Загружаем .env файл если существует
    _ = godotenv.Load()
    
    // Загружаем конфигурацию в структуру
    cfg := &Config{}
    
    // Server
    cfg.Server.Port = getEnv("PORT", "8080")
    cfg.Server.SecretKey = getEnv("SECRET_KEY", "change-this-in-production")
    cfg.Server.Environment = getEnv("ENVIRONMENT", "production")
    cfg.Server.Domain = getEnv("DOMAIN", "localhost")
    cfg.Server.BaseURL = getEnv("BASE_URL", "http://localhost:8080")
    cfg.Server.LogLevel = getEnv("LOG_LEVEL", "info")
    
    // Rate Limiting
    cfg.Server.RateLimit, _ = strconv.Atoi(getEnv("RATE_LIMIT", "100"))
    cfg.Server.RateLimitWindow, _ = strconv.Atoi(getEnv("RATE_LIMIT_WINDOW", "60"))
    cfg.Server.RegistrationCooldown, _ = strconv.Atoi(getEnv("REGISTRATION_COOLDOWN", "300"))
    
    // Database
    cfg.Database.Host = getEnv("DB_HOST", "localhost")
    cfg.Database.Port = getEnv("DB_PORT", "3306")
    cfg.Database.User = getEnv("DB_USER", "root")
    cfg.Database.Password = getEnv("DB_PASSWORD", "")
    cfg.Database.Name = getEnv("DB_NAME", "auth")
    cfg.Database.Charset = getEnv("DB_CHARSET", "utf8mb4")
    
    cfg.Database.CharsHost = getEnv("DB_CHARS_HOST", cfg.Database.Host)
    cfg.Database.CharsPort = getEnv("DB_CHARS_PORT", cfg.Database.Port)
    cfg.Database.CharsUser = getEnv("DB_CHARS_USER", cfg.Database.User)
    cfg.Database.CharsPassword = getEnv("DB_CHARS_PASSWORD", cfg.Database.Password)
    cfg.Database.CharsName = getEnv("DB_CHARS_NAME", "characters")
    
    cfg.Database.MaxOpenConns, _ = strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
    cfg.Database.MaxIdleConns, _ = strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
    cfg.Database.ConnMaxLifetime, _ = time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "300") + "s")
    
    // Redis
    cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
    cfg.Redis.Port = getEnv("REDIS_PORT", "6379")
    cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
    cfg.Redis.DB, _ = strconv.Atoi(getEnv("REDIS_DB", "0"))
    cfg.Redis.PoolSize, _ = strconv.Atoi(getEnv("REDIS_POOL_SIZE", "10"))
    cfg.Redis.MinIdleConns, _ = strconv.Atoi(getEnv("REDIS_MIN_IDLE_CONNS", "2"))
    
    // Game
    cfg.Game.ServerCore, _ = strconv.Atoi(getEnv("SERVER_CORE", "0"))
    cfg.Game.Expansion, _ = strconv.Atoi(getEnv("EXPANSION", "2"))
    cfg.Game.RealmList = getEnv("REALMLIST", "logon.yourserver.com")
    cfg.Game.ServerName = getEnv("SERVER_NAME", "WoW WotLK Server")
    cfg.Game.ServerMOTD = getEnv("SERVER_MOTD", "Welcome to our WoW Server!")
    cfg.Game.MaxAccountsPerIP, _ = strconv.Atoi(getEnv("MAX_ACCOUNTS_PER_IP", "5"))
    cfg.Game.AllowMultiIP, _ = strconv.ParseBool(getEnv("ALLOW_MULTI_IP", "false"))
    
    cfg.Game.BattlenetSupport, _ = strconv.ParseBool(getEnv("BATTLENET_SUPPORT", "false"))
    cfg.Game.SRP6Version, _ = strconv.Atoi(getEnv("SRP6_VERSION", "0"))
    
    // Security
    cfg.Security.PasswordMinLen, _ = strconv.Atoi(getEnv("PASSWORD_MIN_LEN", "4"))
    cfg.Security.PasswordMaxLen, _ = strconv.Atoi(getEnv("PASSWORD_MAX_LEN", "16"))
    cfg.Security.RequireComplexPassword, _ = strconv.ParseBool(getEnv("REQUIRE_COMPLEX_PASSWORD", "false"))
    cfg.Security.PasswordHistoryCount, _ = strconv.Atoi(getEnv("PASSWORD_HISTORY_COUNT", "3"))
    
    cfg.Security.UsernameMinLen, _ = strconv.Atoi(getEnv("USERNAME_MIN_LEN", "3"))
    cfg.Security.UsernameMaxLen, _ = strconv.Atoi(getEnv("USERNAME_MAX_LEN", "16"))
    cfg.Security.AllowSpecialChars, _ = strconv.ParseBool(getEnv("ALLOW_SPECIAL_CHARS", "true"))
    cfg.Security.SpecialCharsAllowed = getEnv("SPECIAL_CHARS_ALLOWED", "_-.")
    
    cfg.Security.RequireEmailVerification, _ = strconv.ParseBool(getEnv("REQUIRE_EMAIL_VERIFICATION", "true"))
    cfg.Security.AllowMultipleAccountsPerEmail, _ = strconv.ParseBool(getEnv("ALLOW_MULTIPLE_ACCOUNTS_PER_EMAIL", "false"))
    cfg.Security.EmailDomainsBlacklist = strings.Split(getEnv("EMAIL_DOMAINS_BLACKLIST", "tempmail.com,10minutemail.com"), ",")
    
    cfg.Security.EnableCaptcha, _ = strconv.ParseBool(getEnv("ENABLE_CAPTCHA", "true"))
    cfg.Security.CaptchaProvider = getEnv("CAPTCHA_PROVIDER", "hcaptcha")
    cfg.Security.CaptchaSecret = getEnv("CAPTCHA_SECRET", "")
    cfg.Security.CaptchaSiteKey = getEnv("CAPTCHA_SITEKEY", "")
    
    cfg.Security.Enable2FA, _ = strconv.ParseBool(getEnv("ENABLE_2FA", "false"))
    cfg.Security.TwoFAProvider = getEnv("2FA_PROVIDER", "totp")
    cfg.Security.TwoFAIssuer = getEnv("2FA_ISSUER", "WoW Server")
    
    // Email
    cfg.Email.SMTPHost = getEnv("SMTP_HOST", "smtp.gmail.com")
    cfg.Email.SMTPPort = getEnv("SMTP_PORT", "587")
    cfg.Email.SMTPUser = getEnv("SMTP_USER", "")
    cfg.Email.SMTPPassword = getEnv("SMTP_PASSWORD", "")
    cfg.Email.SMTPFrom = getEnv("SMTP_FROM", "noreply@wowserver.com")
    cfg.Email.SMTPFromName = getEnv("SMTP_FROM_NAME", "WoW Server")
    cfg.Email.SMTPSecure, _ = strconv.ParseBool(getEnv("SMTP_SECURE", "true"))
    cfg.Email.TemplatePath = getEnv("EMAIL_TEMPLATE_PATH", "./templates/email/")
    
    // Cache
    cfg.Cache.Enabled, _ = strconv.ParseBool(getEnv("CACHE_ENABLED", "true"))
    cfg.Cache.Type = getEnv("CACHE_TYPE", "redis")
    cfg.Cache.Duration, _ = strconv.Atoi(getEnv("CACHE_DURATION", "300"))
    cfg.Cache.StatsDuration, _ = strconv.Atoi(getEnv("STATS_CACHE_DURATION", "60"))
    cfg.Cache.PlayersDuration, _ = strconv.Atoi(getEnv("PLAYERS_CACHE_DURATION", "30"))
    
    // Logging
    cfg.Logging.FilePath = getEnv("LOG_FILE_PATH", "./logs/app.log")
    cfg.Logging.MaxSize, _ = strconv.Atoi(getEnv("LOG_MAX_SIZE", "100"))
    cfg.Logging.MaxAge, _ = strconv.Atoi(getEnv("LOG_MAX_AGE", "30"))
    cfg.Logging.MaxBackups, _ = strconv.Atoi(getEnv("LOG_MAX_BACKUPS", "7"))
    cfg.Logging.Compress, _ = strconv.ParseBool(getEnv("LOG_COMPRESS", "true"))
    
    // Monitoring
    cfg.Monitoring.EnableMetrics, _ = strconv.ParseBool(getEnv("ENABLE_METRICS", "true"))
    cfg.Monitoring.MetricsPort = getEnv("METRICS_PORT", "9090")
    cfg.Monitoring.EnableHealthChecks, _ = strconv.ParseBool(getEnv("ENABLE_HEALTH_CHECKS", "true"))
    cfg.Monitoring.HealthCheckInterval, _ = strconv.Atoi(getEnv("HEALTH_CHECK_INTERVAL", "30"))
    
    cfg.Monitoring.PrometheusEnabled, _ = strconv.ParseBool(getEnv("PROMETHEUS_ENABLED", "true"))
    cfg.Monitoring.PrometheusPath = getEnv("PROMETHEUS_PATH", "/metrics")
    
    // Third-party Integrations
    cfg.Integrations.DiscordWebhookURL = getEnv("DISCORD_WEBHOOK_URL", "")
    cfg.Integrations.DiscordRegistrationNotify, _ = strconv.ParseBool(getEnv("DISCORD_REGISTRATION_NOTIFY", "true"))
    cfg.Integrations.DiscordChannelID = getEnv("DISCORD_CHANNEL_ID", "")
    
    cfg.Integrations.TelegramBotToken = getEnv("TELEGRAM_BOT_TOKEN", "")
    cfg.Integrations.TelegramChatID = getEnv("TELEGRAM_CHAT_ID", "")
    
    cfg.Integrations.SOAPEnabled, _ = strconv.ParseBool(getEnv("SOAP_ENABLED", "false"))
    cfg.Integrations.SOAPHost = getEnv("SOAP_HOST", "localhost")
    cfg.Integrations.SOAPPort = getEnv("SOAP_PORT", "7878")
    cfg.Integrations.SOAPUser = getEnv("SOAP_USER", "admin")
    cfg.Integrations.SOAPPassword = getEnv("SOAP_PASSWORD", "admin")
    
    // Frontend
    cfg.Frontend.Theme = getEnv("THEME", "dark")
    cfg.Frontend.PrimaryColor = getEnv("PRIMARY_COLOR", "#d4af37")
    cfg.Frontend.SecondaryColor = getEnv("SECONDARY_COLOR", "#c41f3b")
    cfg.Frontend.AccentColor = getEnv("ACCENT_COLOR", "#0078ff")
    
    cfg.Frontend.EnableWebSockets, _ = strconv.ParseBool(getEnv("ENABLE_WEBSOCKETS", "true"))
    cfg.Frontend.EnablePWA, _ = strconv.ParseBool(getEnv("ENABLE_PWA", "true"))
    cfg.Frontend.EnableOfflineMode, _ = strconv.ParseBool(getEnv("ENABLE_OFFLINE_MODE", "false"))
    cfg.Frontend.EnableServiceWorker, _ = strconv.ParseBool(getEnv("ENABLE_SERVICE_WORKER", "true"))
    
    cfg.Frontend.GoogleAnalyticsID = getEnv("GOOGLE_ANALYTICS_ID", "")
    cfg.Frontend.CloudflareAnalyticsToken = getEnv("CLOUDFLARE_ANALYTICS_TOKEN", "")
    
    // Security Headers
    cfg.SecurityHeaders.EnableCSP, _ = strconv.ParseBool(getEnv("ENABLE_CSP", "true"))
    cfg.SecurityHeaders.EnableHSTS, _ = strconv.ParseBool(getEnv("ENABLE_HSTS", "true"))
    cfg.SecurityHeaders.EnableXSSProtection, _ = strconv.ParseBool(getEnv("ENABLE_XSS_PROTECTION", "true"))
    cfg.SecurityHeaders.EnableContentTypeOptions, _ = strconv.ParseBool(getEnv("ENABLE_CONTENT_TYPE_OPTIONS", "true"))
    cfg.SecurityHeaders.EnableFrameOptions, _ = strconv.ParseBool(getEnv("ENABLE_FRAME_OPTIONS", "true"))
    cfg.SecurityHeaders.CSPDirectives = getEnv("CSP_DIRECTIVES", "default-src 'self'; script-src 'self' 'unsafe-inline' https://hcaptcha.com https://*.hcaptcha.com; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self' https://hcaptcha.com https://*.hcaptcha.com")
    
    // Maintenance
    cfg.Maintenance.Enabled, _ = strconv.ParseBool(getEnv("MAINTENANCE_MODE", "false"))
    cfg.Maintenance.Message = getEnv("MAINTENANCE_MESSAGE", "Server is under maintenance. Please check back later.")
    cfg.Maintenance.Start = getEnv("MAINTENANCE_START", "")
    cfg.Maintenance.End = getEnv("MAINTENANCE_END", "")
    
    // Debugging
    cfg.Debug.Enabled, _ = strconv.ParseBool(getEnv("DEBUG", "false"))
    cfg.Debug.EnableSwagger, _ = strconv.ParseBool(getEnv("ENABLE_SWAGGER", "true"))
    cfg.Debug.SwaggerPath = getEnv("SWAGGER_PATH", "/api/docs")
    cfg.Debug.AllowInsecureRegistrations, _ = strconv.ParseBool(getEnv("ALLOW_INSECURE_REGISTRATIONS", "false"))
    cfg.Debug.SkipCaptchaInDev, _ = strconv.ParseBool(getEnv("SKIP_CAPTCHA_IN_DEV", "true"))
    
    // Custom Settings
    cfg.Custom.RegistrationBonusEnabled, _ = strconv.ParseBool(getEnv("REGISTRATION_BONUS_ENABLED", "false"))
    cfg.Custom.StartGold, _ = strconv.Atoi(getEnv("START_GOLD", "0"))
    cfg.Custom.StartItems = getEnv("START_ITEMS", "")
    cfg.Custom.StartSpells = getEnv("START_SPELLS", "")
    
    cfg.Custom.VoteSystemEnabled, _ = strconv.ParseBool(getEnv("VOTE_SYSTEM_ENABLED", "false"))
    cfg.Custom.VoteSites = getEnv("VOTE_SITES", "")
    cfg.Custom.VoteRewardItem = getEnv("VOTE_REWARD_ITEM", "")
    cfg.Custom.VoteRewardCount, _ = strconv.Atoi(getEnv("VOTE_REWARD_COUNT", "1"))
    
    cfg.Custom.ReferralSystemEnabled, _ = strconv.ParseBool(getEnv("REFERRAL_SYSTEM_ENABLED", "false"))
    cfg.Custom.ReferralReward = getEnv("REFERRAL_REWARD", "")
    cfg.Custom.MaxReferralsPerAccount, _ = strconv.Atoi(getEnv("MAX_REFERRALS_PER_ACCOUNT", "10"))
    
    // Сохраняем конфигурацию в глобальную переменную
    AppConfig = cfg
    
    return nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
