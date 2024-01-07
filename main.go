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
	type response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Url     api    `json:"url"`
		Repo    string `json:"repo"`
		Postman string `json:"postman"`
	}
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		response := response{
			Code:    c.Response().StatusCode(),
			Message: "OK",
			Url: api{
				GetAll:  get{Method: "GET", Url: "/albums"},
				Create:  get{Method: "POST", Url: "/albums", Form: []string{"title", "artist", "price"}},
				GetById: get{Method: "GET", Url: "/albums/id", Param: "id integer"},
				Update:  get{Method: "PUT", Url: "/albums/id", Param: "id integer", Form: []string{"title", "artist", "price"}},
				Delete:  get{Method: "DELETE", Url: "/albums/id", Param: "id integer"},
			},
			Repo:    "https://github.com/WahidinAji/restfull-slice-fiber",
			Postman: "You guys can download the postman collection and environment *.json in that repository. You can run the testing with postman.",
		}
		c.Set("Content-type", "application/json; charset=UTF-8")
		return c.JSON(response)
	})
	setupRoute(app)
	log.Fatal(app.Listen(":8080"))
}

func getAlbums(c *fiber.Ctx) error {
	type response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  bool   `json:"status"`
		Data    []Album
	}
	res := response{
		Code:    c.Response().StatusCode(),
		Message: "Successfuly",
		Status:  true,
		Data:    albums,
	}
	c.Set("Content-type", "application/json; charset=UTF-8")
	return c.JSON(res)
}
func createAlbum(c *fiber.Ctx) error {
	// newAlbum := new(Album)
	var album Album
	if err := c.BodyParser(&album); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	album.ID = len(albums) + 1
	albums = append(albums, album)
	c.Status(fiber.StatusCreated)
	res := SingleResponse{
		Code:    c.Response().StatusCode(),
		Message: "Album created successfuly",
		Status:  true,
		Data:    album,
	}
	c.Set("Content-type", "application/json; charset=UTF-8")
	return c.JSON(res)
}
func getAlbumsById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	for _, album := range albums {
		if album.ID == id {
			res := SingleResponse{
				Code:    c.Response().StatusCode(),
				Message: "Successfuly",
				Status:  true,
				Data:    album,
			}
			c.Set("Content-type", "application/json; charset=UTF-8")
			return c.JSON(res)
		}
	}
	c.Status(404)
	resposne := ResponseDataString{
		Code:    c.Response().StatusCode(),
		Message: "Failed id was not found",
		Status:  false,
		Data:    "id was not found",
	}
	c.Set("Content-type", "application/json; charset=UTF-8")
	return c.JSON(resposne)
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
			response := SingleResponse{
				Code:    c.Response().StatusCode(),
				Message: "Updated successfuly",
				Status:  true,
				Data:    albums[i],
			}
			c.Set("Content-type", "application/json; charset=UTF-8")
			return c.JSON(response)
		}
	}
	c.Status(404)
	response := ResponseDataString{
		Code:    c.Response().StatusCode(),
		Message: "Failed id was not found",
		Status:  false,
		Data:    "id was not found",
	}
	c.Set("Content-type", "application/json; charset=UTF-8")
	return c.JSON(response)
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
			response := Response{
				Code:    c.Response().StatusCode(),
				Message: "Deleted successfuly",
				Status:  true,
			}
			c.Set("Content-type", "application/json; charset=UTF-8")
			return c.JSON(response)
		}
	}
	c.Status(404)
	response := ResponseDataString{
		Code:    c.Response().StatusCode(),
		Message: "Failed id was not found",
		Status:  false,
		Data:    "id was not found",
	}
	c.Set("Content-type", "application/json; charset=UTF-8")
	return c.JSON(response)
}

type SingleResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    Album
}

type Response struct {
	Code    int
	Message string
	Status  bool
}
type ResponseDataString struct {
	Code    int
	Message string
	Status  bool
	Data    string
}
