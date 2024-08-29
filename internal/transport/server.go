package transport

import (
	"context"
	"log"
	"net/http"
	"time"
	"wbl0/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type OrderService interface {
	Create(orderUID, data string) error
	Get(orderUID string) (string, error)
}

type StanService interface {
	Start()
	Stop()
}

type App struct {
	cfg *config.Config

	web      *http.Server
	router   *gin.Engine
	orderSvc OrderService
	stanSvc  StanService
}

func New(cfg *config.Config, orderSvc OrderService, stanSvc StanService) *App {
	router := gin.New()

	app := &App{
		cfg:      cfg,
		router:   router,
		orderSvc: orderSvc,
		stanSvc:  stanSvc,
	}

	app.web = &http.Server{
		Addr:    app.cfg.Web.Address(),
		Handler: router,
	}

	// Загрузка HTML-шаблонов

	router.LoadHTMLGlob("/app/templates/*")

	app.initRoutes()

	return app
}

func (app *App) initRoutes() {
	app.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Context-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	// Маршрут для отображения index.html
	app.router.GET("/", app.RenderIndex)

	// Маршрут для получения заказа по ID
	app.router.GET("/:id", app.GetById)
}

// RenderIndex рендерит HTML-шаблон index.html
func (app *App) RenderIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (app *App) Start() error {
	app.stanSvc.Start()
	log.Printf("start web server at http://%s\n", app.web.Addr)
	return app.web.ListenAndServe()
}

func (app *App) Stop(ctx context.Context) error {
	app.stanSvc.Stop()
	log.Println("stop web server")
	return app.web.Shutdown(ctx)
}

func (app *App) GetById(c *gin.Context) {
	uid := c.Param("id")
	log.Printf("[endpoint][0H] [GetById]: %s", uid)
	data, err := app.orderSvc.Get(uid)
	if err != nil {
		log.Printf("[endpoint][1H] err: %s", err)
		c.JSON(http.StatusNotFound, gin.H{
			"result": "not found",
			"data":   "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"data":   data,
	})
}
