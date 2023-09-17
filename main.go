package main

import (
    "log"
    "os"
	"net/http"
	"math/rand"

	"github.com/pocketbase/dbx"
	"github.com/labstack/echo/v5"
	"github.com/google/uuid"
	
    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/apis"
    "github.com/pocketbase/pocketbase/core"	
)

func main() {
    app := pocketbase.New()

    // serves static files from the provided public dir (if exists)
    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
        return nil
    })

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("api-ext/hello/:name", func(c echo.Context) error {
			name := c.PathParam("name")
			newUUID := uuid.New()

			cities := []string{"New York", "Los Angeles", "London", "Paris", "Tokyo", "Sydney"}
			randomCity := cities[rand.Intn(len(cities))]
	
			return c.JSON(http.StatusOK, map[string]string{
				"message": "Hello " + name + "!",
				"uuid": newUUID.String(),
				"city": randomCity,
			})
		})
	
		return nil
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("api-ext/characters", func(c echo.Context) error {
			characters:= []struct {
				Id string `db:"id" json:"id"`
				Name string `db:"name" json:"name"`
				Age int `db:"age" json:"age"`
				CreatedAt string `db:"created" json:"created"`
				UpdatedAt string `db:"updated" json:"updated"`
			}{}

			err := app.Dao().DB().
				Select("id", "name", "age", "created", "updated").
				From("characters").
				All(&characters)
			
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
			}
			
			return c.JSON(http.StatusOK, characters)
		})
		
		e.Router.POST("api-ext/characters", func(c echo.Context) error {
			type Character struct {
				Name string `db:"name" json:"name"`
				Age  int    `db:"age" json:"age"`
			}

			character := &Character{}

			bind_err := c.Bind(character)
			if bind_err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid character"})
			}
			
			_, err := app.Dao().DB().
				NewQuery("INSERT INTO characters (name, age) VALUES ({:name}, {:age})").
				Bind(dbx.Params{
					"name": character.Name,
					"age": character.Age,
					}).
				Execute()

			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to insert character"})
			}

			return c.JSON(http.StatusOK, map[string]string{"message": "Character inserted successfully"})
		})
		
		return nil
	})

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
