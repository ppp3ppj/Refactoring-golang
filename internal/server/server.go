package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"github.com/ppp3ppj/go-refactoring-workshop/config"
	"github.com/ppp3ppj/go-refactoring-workshop/db"
	server_middlewares "github.com/ppp3ppj/go-refactoring-workshop/internal/middlewares"
)

type echoServer struct {
    app *echo.Echo
    conf *config.Config
    db db.IDatabase
}

var (
    server *echoServer
    once sync.Once
)

func NewEchoServer(conf *config.Config, db db.IDatabase) *echoServer {
    echoApp := echo.New()
    echoApp.Logger.SetLevel(log.DEBUG)

    once.Do(func() {
        server = &echoServer{
            app: echoApp,
            conf: config.ConfigGetting(),
            db: db,
        }
    })

    return server
}

type Person struct {
    Key         string          `json:"key"`
    Name        string          `json:"name"`
    Description string          `json:"description"`
    Image       string          `json:"image"`
    Traits      []Trait `json:"traits"` // Use RawMessage to hold JSON data
    Tags        []string `json:"tags"`   // Use RawMessage to hold JSON data
}

type Trait struct {
    Personality string `json:"personality"`
    Like string `json:"like"`
    ZodiacSign string `json:"Zodiac Sign"`
    Emoji string `json:"emoji"`
    Color string `json:"color"`
}



func (s * echoServer) Start() {
    timeOutMiddleware := server_middlewares.GetTimeOutMiddleware(s.conf.Server.Timeout)
    corsMiddleware := server_middlewares.GetCORSMiddleware(s.conf.Server.AllowOrigins)

    s.app.Use(middleware.Recover())
    s.app.Use(middleware.Logger())

    s.app.Use(timeOutMiddleware)
    s.app.Use(corsMiddleware)

    // set format color
    s.app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
        Output: s.app.Logger.Output(),
    }))

    s.app.HideBanner = true
    s.app.HidePort = true
    asciiArt := fmt.Sprintf(`
        WWWWWWWWWWWWWMWWWWWWWWWWWWWWWWWWWWWWWWWW
        WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW
        WWWWW0k0NWWWWNNWWWWNKOkXWWWWWWWWWWWWWWWW
        WWWWWO:coxOOkkkkkkxlc;c0WWWWWWWWWWWWWWWW
        WWWWW0lcdxkxxOOkkkxddlcxKWWWWWWWWWWWWWWW
        WWWW0xlcx0x:oOdcdOOOOOkxxONWWWWWWWWWWWWW
        WWW0odd:lkd:okd:okkkkkxd:,oNWWWWWWWWWWWW
        WWWx,.':ldddxdl:,''cdkOc  .OWWWWWWWWWWWW
        WWWO;..:::OXOdl;.  ..''.  .l0NWWWWWWWWWW
        WWWXd'  .c0K0kdl:'.      .:ddkXWWWWWWWWW
        WWWWKc,,l00O0X0kdl:;,,..'lkkdox0NWWWWWWW
        WWWWNkodoxOO00KKOxdl:;'. ,xXOlcokKNWWWWW
        WWWWWKxodO0kooxxd:'.... ..;ol;;odoOWWWWW
        WWWWWW0ooxOOddkdl:,,,,,;;::;:dO0koo0WWWW
        WWWWWWWOddooddx0KOl:cccccloxk0KXk;.'OWWW
        WWWWWWWWX00OOdclc:,lO00KXNWWWNKx;..;0WWW
        WWWWWWWWWWWWWXd;..,ONWWWWWWWWWWWX0KNWWWW
        WWWWWWWWWWWWWWWX00XWWWWWMWWWWWWWWWWWWWWW
        WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW
        WWWWWWWWWWWWWWWWWWWWWWWWWWWMWWWWWWWWWWWW

        Haerin,
        Name: %s Port: %d
    `,
     s.conf.AppInfo.Name,
     s.conf.Server.Port,
    )

    fmt.Print(asciiArt)
    s.app.GET("/v1/health", s.healthCheck)

    s.app.GET("/person/:key", func(c echo.Context) error {
        key := c.Param("key")
        row := s.db.Connect().QueryRow(`SELECT key, name, description, image, traits, tags FROM "Person" WHERE key = $1`, key)

        var person Person
        var traits []byte
        var tags pq.StringArray
        if err := row.Scan(&person.Key, &person.Name, &person.Description, &person.Image, &traits, &tags); err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }

        if err := json.Unmarshal(traits, &person.Traits); err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }

        person.Tags = tags
        return c.JSON(http.StatusOK, person)

    })

    s.app.GET("/persons", func(c echo.Context) error {
        row, err := s.db.Connect().Query(`SELECT key, name, description, image, traits, tags FROM "Person"`)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }

        persons := []Person{}
        for row.Next() {
            var person Person
            var traits []byte
            var tags pq.StringArray
            if err := row.Scan(&person.Key, &person.Name, &person.Description, &person.Image, &traits, &tags); err != nil {
                return c.JSON(http.StatusInternalServerError, err)
            }

            if err := json.Unmarshal(traits, &person.Traits); err != nil {
                return c.JSON(http.StatusInternalServerError, err)
            }

            person.Tags = tags
            persons = append(persons, person)
        }

        return c.JSON(http.StatusOK, persons)
    })

    s.app.POST("/persons", func(c echo.Context) error {
        p := new(Person)
        if err := c.Bind(p); err != nil {
            return c.String(http.StatusBadRequest, "bad request")
        }

        stmt, err := s.db.Connect().Prepare(`INSERT INTO "Person" (key, name, description, image, traits, tags) VALUES ($1, $2, $3, $4, $5, $6)`)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }

        traits, err := json.Marshal(p.Traits)
        if err != nil {
            c.JSON(http.StatusInternalServerError, err)
        }

        tags := pq.StringArray(p.Tags)
        _, err = stmt.Exec(p.Key, p.Name, p.Description, p.Image, traits, tags)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }

        return c.JSON(http.StatusOK, "OK")

    })

    s.app.PUT("/person/:key", func(c echo.Context) error {
        key := c.Param("key")
        p := new(Person)
        if err := c.Bind(p); err != nil {
            return c.String(http.StatusBadRequest, "bad request")
        }

        stmt, err := s.db.Connect().Prepare(`
            UPDATE "Person"
            SET name = $1, description = $2, image = $3, traits = $4, tags = $5
            WHERE key = $6
        `)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }

        traits, err := json.Marshal(p.Traits)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }

        tags := pq.StringArray(p.Tags)
        _, err = stmt.Exec(p.Name, p.Description, p.Image, traits, tags, key)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }

        return c.JSON(http.StatusOK, "Updated successfully")
    })

    s.app.DELETE("/person/:key", func(c echo.Context) error {
        key := c.Param("key")

        stmt, err := s.db.Connect().Prepare(`DELETE FROM "Person" WHERE key = $1`)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }

        result, err := stmt.Exec(key)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }

        if rowsAffected == 0 {
            return c.JSON(http.StatusNotFound, "Person not found")
        }

        return c.JSON(http.StatusOK, "Deleted successfully")
    })

    // Graceful Shutdown
    quitCh := make(chan os.Signal, 1)
    signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
    go s.gracefullyShutdown(quitCh)

    s.httpListening()
}

func (s *echoServer) gracefullyShutdown(quitCh <-chan os.Signal) {
    ctx := context.Background()

    <-quitCh
    s.app.Logger.Info("Shutting down the service...")

    if err := s.app.Shutdown(ctx); err != nil {
        s.app.Logger.Fatalf("Error: %s", err.Error())
    }
}

func (s *echoServer) httpListening() {
    url := fmt.Sprintf(":%d", s.conf.Server.Port)
    if err := s.app.Start(url); err != nil && err != http.ErrServerClosed {
        s.app.Logger.Fatalf("shutting down the server: %v", err)
    }
}

func (s *echoServer) healthCheck(c echo.Context) error {
    return c.String(http.StatusOK, "OK")
}
