package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []Album{
	{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func setupRoute(app *fiber.App) {
	app.Get("/albums", getAlbums)
	app.Post("/albums", createAlbum)
	app.Get("/albums/:id", getAlbumsById)
	app.Put("/albums/:id", updateAlbumsById)
	app.Delete("/albums/:id", deleteAlbumsById)
}

func main() {
	type get struct {
		Method string
		Url    string
		Form   []string
		Param  string
	}
	type api struct {
		GetAll  get
		Create  get
		GetById get
		Update  get
		Delete  get
	}
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		response := c.JSON(&fiber.Map{
			"code":    c.Response().StatusCode(),
			"message": "OK",
			"url": api{
				GetAll:  get{Method: "GET", Url: "/albums"},
				Create:  get{Method: "POST", Url: "/albums", Form: []string{"title", "artist", "price"}},
				GetById: get{Method: "GET", Url: "/albums/id", Param: "id integer"},
				Update:  get{Method: "PUT", Url: "/albums/id", Param: "id integer", Form: []string{"title", "artist", "price"}},
				Delete:  get{Method: "DELETE", Url: "/albums/id", Param: "id integer"},
			},
			"repo":    "https://github.com/WahidinAji/fiber-example/tree/main/restapi-slice",
			"postman": "You guys can download the postman collection and environment *.json in that repository. You can run the testing with postman.",
		})
		c.Set("Content-type", "application/json; charset=UTF-8")
		return response
	})
	setupRoute(app)
	log.Fatal(app.Listen(":8080"))
}

func getAlbums(c *fiber.Ctx) error {
	response := c.JSON(&fiber.Map{
		"code":    c.Response().StatusCode(),
		"message": "Successfuly",
		"status":  true,
		"data":    albums,
	})
	c.Set("Content-type", "application/json; charset=UTF-8")
	return response
}
func createAlbum(c *fiber.Ctx) error {
	newAlbum := new(Album)
	if err := c.BodyParser(&newAlbum); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	newAlbum.ID = len(albums) + 1
	// albums = append(albums, Album{ID: newAlbum.ID, Title: newAlbum.Title, Artist: newAlbum.Artist, Price: newAlbum.Price})
	albums = append(albums, *newAlbum)
	c.Status(fiber.StatusCreated)
	response := c.JSON(&fiber.Map{
		"code":    c.Response().StatusCode(),
		"message": "Album created successfuly",
		"status":  true,
		"data":    newAlbum,
	})
	c.Set("Content-type", "application/json; charset=UTF-8")
	return response
}
func getAlbumsById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	for _, album := range albums {
		if album.ID == id {
			response := c.JSON(&fiber.Map{
				"code":    c.Response().StatusCode(),
				"message": "Successfuly",
				"status":  true,
				"data":    album,
			})
			c.Set("Content-type", "application/json; charset=UTF-8")
			return response
		}
	}
	c.Status(404)
	resposne := c.JSON(&fiber.Map{
		"code":    c.Response().StatusCode(),
		"message": "Failed id was not found",
		"status":  false,
		"data":    "",
	})
	c.Set("Content-type", "application/json; charset=UTF-8")
	return resposne
}
func updateAlbumsById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	album := new(Album)
	if err := c.BodyParser(album); err != nil {
		return err
	}
	i := 0
	for i = 0; i < len(albums); i++ {
		if albums[i].ID == id {
			albums[i] = Album{ID: id, Title: album.Title, Artist: album.Artist, Price: album.Price}
			response := c.JSON(&fiber.Map{
				"code":    c.Response().StatusCode(),
				"message": "Updated successfuly",
				"status":  true,
				"data":    albums[i],
			})
			c.Set("Content-type", "application/json; charset=UTF-8")
			return response
		}
	}
	c.Status(404)
	response := c.JSON(&fiber.Map{
		"code":    c.Response().StatusCode(),
		"message": "Failed id was not found",
		"status":  false,
		"data":    "id was not found",
	})
	c.Set("Content-type", "application/json; charset=UTF-8")
	return response
}

func deleteAlbumsById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return err
	}

	for i := 0; i < len(albums); i++ {
		url := albums[i]
		if url.ID == id {
			albums = append(albums[:i], albums[i+1:]...)
			i--
			response := c.JSON(&fiber.Map{
				"code":    c.Response().StatusCode(),
				"message": "Deleted Successfuly",
				"status":  true,
			})
			c.Set("Content-type", "application/json; charset=UTF-8")
			return response
		}
	}
	c.Status(404)
	response := c.JSON(&fiber.Map{
		"code":    c.Response().StatusCode(),
		"message": "Failed id was not found",
		"status":  false,
		"data":    "id was not found",
	})
	c.Set("Content-type", "application/json; charset=UTF-8")
	return response
}
