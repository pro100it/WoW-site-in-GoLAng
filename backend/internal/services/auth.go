package services

import (
    "crypto/rand"
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "math/big"
    "net"
    "strings"
    "time"
    "wow-registration/internal/config"
    "github.com/google/uuid"
)

type SRP6Verifier struct {
    Salt    string
    Verifier string
}

func GenerateSRP6(username, password string, coreType int) (*SRP6Verifier, error) {
    // Генерация соли
    salt := make([]byte, 32)
    if _, err := rand.Read(salt); err != nil {
        return nil, err
    }
    
    // Константы для SRP6
    g := big.NewInt(7)
    N := new(big.Int)
    N.SetString("894B645E89E1535BBDAD5B8B290650530801B18EBFBF5E8FAB3C82872A3E9BB7", 16)
    
    // Вычисление verifier
    h1 := sha1.Sum([]byte(strings.ToUpper(username + ":" + password)))
    
    var h2Input []byte
    if coreType == 5 { // CMangos
        // Reverse salt для CMangos
        revSalt := make([]byte, len(salt))
        for i, b := range salt {
            revSalt[len(salt)-1-i] = b
        }
        h2Input = append(revSalt, h1[:]...)
    } else { // TrinityCore
        h2Input = append(salt, h1[:]...)
    }
    
    h2 := sha1.Sum(h2Input)
    h2Int := new(big.Int).SetBytes(h2)
    
    verifier := new(big.Int).Exp(g, h2Int, N)
    verifierBytes := verifier.Bytes()
    
    // Pad to 32 bytes
    if len(verifierBytes) < 32 {
        padding := make([]byte, 32-len(verifierBytes))
        verifierBytes = append(padding, verifierBytes...)
    }
    
    if coreType == 5 { // CMangos
        // Reverse для CMangos
        revVerifier := make([]byte, len(verifierBytes))
        for i, b := range verifierBytes {
            revVerifier[len(verifierBytes)-1-i] = b
        }
        verifierBytes = revVerifier
        return &SRP6Verifier{
            Salt:     strings.ToUpper(hex.EncodeToString(salt)),
            Verifier: strings.ToUpper(hex.EncodeToString(verifierBytes)),
        }, nil
    }
    
    return &SRP6Verifier{
        Salt:     hex.EncodeToString(salt),
        Verifier: hex.EncodeToString(verifierBytes),
    }, nil
}

func GenerateSHA1Hash(username, password string) string {
    hash := sha1.Sum([]byte(strings.ToUpper(username + ":" + password)))
    return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func ValidatePassword(password string) error {
    cfg := config.AppConfig
    
    if len(password) < cfg.Security.PasswordMinLen {
        return fmt.Errorf("password must be at least %d characters", cfg.Security.PasswordMinLen)
    }
    
    if len(password) > cfg.Security.PasswordMaxLen {
        return fmt.Errorf("password cannot exceed %d characters", cfg.Security.PasswordMaxLen)
    }
    
    return nil
}

func ValidateUsername(username string) error {
    if len(username) < 3 {
        return fmt.Errorf("username must be at least 3 characters")
    }
    
    if len(username) > 16 {
        return fmt.Errorf("username cannot exceed 16 characters")
    }
    
    // Только латинские буквы, цифры и подчеркивания
    for _, r := range username {
        if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
            return fmt.Errorf("username can only contain letters, numbers and underscores")
        }
    }
    
    return nil
}

func ValidateEmail(email string) error {
    if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
        return fmt.Errorf("invalid email format")
    }
    
    if len(email) > 255 {
        return fmt.Errorf("email is too long")
    }
    
    return nil
}

func GenerateSessionToken() string {
    return uuid.New().String()
}

func GetClientIP(r interface{}) string {
    // Здесь нужно реализовать получение IP из запроса
    // В зависимости от используемого фреймворка
    return "127.0.0.1"
}

func RateLimitKey(ip string, action string) string {
    now := time.Now().Format("2006010215")
    return fmt.Sprintf("ratelimit:%s:%s:%s", ip, action, now)
}

func GenerateRandomString(length int) string {
    const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    result := make([]byte, length)
    
    for i := range result {
        n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
        result[i] = chars[n.Int64()]
    }
    
    return string(result)
}
