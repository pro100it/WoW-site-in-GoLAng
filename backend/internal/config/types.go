package config

import (
    "time"
)

type Config struct {
    Server    ServerConfig
    Database  DatabaseConfig
    Redis     RedisConfig
    Game      GameConfig
    Security  SecurityConfig
    Email     EmailConfig
    Cache     CacheConfig
    Logging   LoggingConfig
    Monitoring MonitoringConfig
    Integrations IntegrationsConfig
    Frontend  FrontendConfig
    SecurityHeaders SecurityHeadersConfig
    Maintenance MaintenanceConfig
    Debug     DebugConfig
    Custom    CustomConfig
}

type ServerConfig struct {
    Port                string
    SecretKey           string
    Environment         string
    Domain              string
    BaseURL             string
    LogLevel            string
    RateLimit           int
    RateLimitWindow     int
    RegistrationCooldown int
}

type DatabaseConfig struct {
    Host              string
    Port              string
    User              string
    Password          string
    Name              string
    Charset           string
    
    // Character database
    CharsHost         string
    CharsPort         string
    CharsUser         string
    CharsPassword     string
    CharsName         string
    
    // World database (optional)
    WorldHost         string
    WorldPort         string
    WorldUser         string
    WorldPassword     string
    WorldName         string
    
    // Connection pool
    MaxOpenConns      int
    MaxIdleConns      int
    ConnMaxLifetime   time.Duration
}

type RedisConfig struct {
    Host         string
    Port         string
    Password     string
    DB           int
    PoolSize     int
    MinIdleConns int
}

type GameConfig struct {
    ServerCore       int
    Expansion        int
    RealmList        string
    ServerName       string
    ServerMOTD       string
    MaxAccountsPerIP int
    AllowMultiIP     bool
    BattlenetSupport bool
    SRP6Version      int
}

type SecurityConfig struct {
    PasswordMinLen               int
    PasswordMaxLen               int
    RequireComplexPassword       bool
    PasswordHistoryCount         int
    
    UsernameMinLen               int
    UsernameMaxLen               int
    AllowSpecialChars            bool
    SpecialCharsAllowed          string
    
    RequireEmailVerification     bool
    AllowMultipleAccountsPerEmail bool
    EmailDomainsBlacklist        []string
    
    EnableCaptcha                bool
    CaptchaProvider              string
    CaptchaSecret                string
    CaptchaSiteKey               string
    
    Enable2FA                    bool
    TwoFAProvider                string
    TwoFAIssuer                  string
}

type EmailConfig struct {
    SMTPHost       string
    SMTPPort       string
    SMTPUser       string
    SMTPPassword   string
    SMTPFrom       string
    SMTPFromName   string
    SMTPSecure     bool
    TemplatePath   string
}

type CacheConfig struct {
    Enabled         bool
    Type            string
    Duration        int
    StatsDuration   int
    PlayersDuration int
}

type LoggingConfig struct {
    FilePath    string
    MaxSize     int
    MaxAge      int
    MaxBackups  int
    Compress    bool
}

type MonitoringConfig struct {
    EnableMetrics       bool
    MetricsPort         string
    EnableHealthChecks  bool
    HealthCheckInterval int
    PrometheusEnabled   bool
    PrometheusPath      string
}

type IntegrationsConfig struct {
    DiscordWebhookURL        string
    DiscordRegistrationNotify bool
    DiscordChannelID         string
    
    TelegramBotToken         string
    TelegramChatID           string
    
    SOAPEnabled              bool
    SOAPHost                 string
    SOAPPort                 string
    SOAPUser                 string
    SOAPPassword             string
}

type FrontendConfig struct {
    Theme                 string
    PrimaryColor          string
    SecondaryColor        string
    AccentColor           string
    
    EnableWebSockets      bool
    EnablePWA             bool
    EnableOfflineMode     bool
    EnableServiceWorker   bool
    
    GoogleAnalyticsID     string
    CloudflareAnalyticsToken string
}

type SecurityHeadersConfig struct {
    EnableCSP              bool
    EnableHSTS             bool
    EnableXSSProtection    bool
    EnableContentTypeOptions bool
    EnableFrameOptions     bool
    CSPDirectives          string
}

type MaintenanceConfig struct {
    Enabled   bool
    Message   string
    Start     string
    End       string
}

type DebugConfig struct {
    Enabled                   bool
    EnableSwagger             bool
    SwaggerPath               string
    AllowInsecureRegistrations bool
    SkipCaptchaInDev         bool
}

type CustomConfig struct {
    RegistrationBonusEnabled  bool
    StartGold                 int
    StartItems                string
    StartSpells               string
    
    VoteSystemEnabled         bool
    VoteSites                 string
    VoteRewardItem            string
    VoteRewardCount           int
    
    ReferralSystemEnabled     bool
    ReferralReward            string
    MaxReferralsPerAccount    int
}
