package handlers

import (
    "html/template"
    "net/http"
    "os"
    "path/filepath"
    "wow-registration/internal/database"
    "github.com/labstack/echo/v4"
)

type Template struct {
    templates *template.Template
}

func NewTemplateRenderer() *Template {
    tmpl := template.New("")
    
    // Автоматически загружаем все шаблоны
    templateDir := "./frontend/templates"
    tmpl = template.Must(tmpl.ParseGlob(filepath.Join(templateDir, "*.html")))
    tmpl = template.Must(tmpl.ParseGlob(filepath.Join(templateDir, "partials/*.html")))
    
    return &Template{templates: tmpl}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

type PageData struct {
    Title       string
    Description string
    Config      interface{}
    Stats       map[string]interface{}
    OnlinePlayers []database.Character
}

func HomeHandler(c echo.Context) error {
    stats, _ := database.GetServerStats()
    onlinePlayers, _ := database.GetOnlinePlayers(1) // realmID = 1
    
    data := PageData{
        Title:       "WoW Server Registration",
        Description: "Register your World of Warcraft account",
        Config:      config.AppConfig,
        Stats:       stats,
        OnlinePlayers: onlinePlayers,
    }
    
    return c.Render(http.StatusOK, "index.html", data)
}

func StatusHandler(c echo.Context) error {
    stats, err := database.GetServerStats()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    
    onlinePlayers, _ := database.GetOnlinePlayers(1)
    
    return c.JSON(http.StatusOK, map[string]interface{}{
        "stats": stats,
        "online_players": onlinePlayers,
    })
}

func RealTimeStatsHandler(c echo.Context) error {
    stats, _ := database.GetServerStats()
    onlinePlayers, _ := database.GetOnlinePlayers(1)
    
    return c.Render(http.StatusOK, "partials/stats.html", map[string]interface{}{
        "stats": stats,
        "online_players": onlinePlayers,
    })
}
