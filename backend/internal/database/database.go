package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    "wow-registration/internal/config"
    _ "github.com/go-sql-driver/mysql"
    "github.com/redis/go-redis/v9"
)

var (
    DB    *sql.DB
    Redis *redis.Client
)

type Account struct {
    ID          int
    Username    string
    Email       string
    Password    string
    Salt        string
    Verifier    string
    Expansion   int
    CreatedAt   time.Time
    LastLogin   sql.NullTime
    IP          string
    Locked      bool
}

type Character struct {
    GUID    int
    Name    string
    Race    int
    Class   int
    Level   int
    Gender  int
    RealmID int
}

func Connect() error {
    cfg := config.AppConfig
    
    // MySQL connection
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.Name,
        cfg.Database.Charset,
    )
    
    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }
    
    // Test connection
    if err := DB.Ping(); err != nil {
        return fmt.Errorf("database ping failed: %w", err)
    }
    
    // Set connection pool settings
    DB.SetMaxOpenConns(25)
    DB.SetMaxIdleConns(5)
    DB.SetConnMaxLifetime(5 * time.Minute)
    
    // Redis connection
    Redis = redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,
    })
    
    log.Println("✅ Database connections established")
    return nil
}

func CreateAccount(account *Account) error {
    query := `
        INSERT INTO account (
            username, email, sha_pass_hash, expansion, 
            last_ip, joindate, sessionkey, v, s, locked
        ) VALUES (?, ?, ?, ?, ?, NOW(), '', ?, ?, ?)
    `
    
    _, err := DB.Exec(query,
        account.Username,
        account.Email,
        account.Password,
        account.Expansion,
        account.IP,
        account.Verifier,
        account.Salt,
        account.Locked,
    )
    
    return err
}

func AccountExists(username, email string) (bool, error) {
    query := `
        SELECT COUNT(*) FROM account 
        WHERE username = ? OR email = ?
    `
    
    var count int
    err := DB.QueryRow(query, username, email).Scan(&count)
    if err != nil {
        return false, err
    }
    
    return count > 0, nil
}

func GetAccountByUsername(username string) (*Account, error) {
    query := `
        SELECT id, username, email, expansion, joindate, last_login, last_ip, locked
        FROM account WHERE username = ?
    `
    
    account := &Account{}
    err := DB.QueryRow(query, username).Scan(
        &account.ID,
        &account.Username,
        &account.Email,
        &account.Expansion,
        &account.CreatedAt,
        &account.LastLogin,
        &account.IP,
        &account.Locked,
    )
    
    if err != nil {
        return nil, err
    }
    
    return account, nil
}

func UpdatePassword(username, newHash string) error {
    query := `
        UPDATE account 
        SET sha_pass_hash = ?, sessionkey = '', v = '', s = ''
        WHERE username = ?
    `
    
    _, err := DB.Exec(query, newHash, username)
    return err
}

func GetOnlinePlayers(realmID int) ([]Character, error) {
    query := `
        SELECT guid, name, race, class, level, gender
        FROM characters 
        WHERE online = 1 AND realm = ?
        ORDER BY level DESC
        LIMIT 20
    `
    
    rows, err := DB.Query(query, realmID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var characters []Character
    for rows.Next() {
        var c Character
        if err := rows.Scan(&c.GUID, &c.Name, &c.Race, &c.Class, &c.Level, &c.Gender); err != nil {
            return nil, err
        }
        c.RealmID = realmID
        characters = append(characters, c)
    }
    
    return characters, nil
}

func GetServerStats() (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    // Total accounts
    var totalAccounts int
    err := DB.QueryRow("SELECT COUNT(*) FROM account").Scan(&totalAccounts)
    if err != nil {
        return nil, err
    }
    stats["total_accounts"] = totalAccounts
    
    // Today's registrations
    var todayRegistrations int
    err = DB.QueryRow("SELECT COUNT(*) FROM account WHERE DATE(joindate) = CURDATE()").Scan(&todayRegistrations)
    if err != nil {
        return nil, err
    }
    stats["today_registrations"] = todayRegistrations
    
    // Online players (from characters DB - нужно отдельное соединение)
    stats["online_players"] = 0 // Заполняется отдельно
    
    return stats, nil
}
