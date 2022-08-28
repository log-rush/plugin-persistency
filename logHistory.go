package pluginPersistency

import (
	"github.com/gofiber/fiber/v2"
	"github.com/log-rush/distribution-server/pkg/app"
)

type routerPlugin struct {
	config Config
}

func newRouterPlugin(config Config) *routerPlugin {
	return &routerPlugin{
		config: config,
	}
}

func (p *routerPlugin) SetupRouter(router fiber.Router, context *app.Context) {
	router.Get("/files/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")
		return c.Status(200).JSON(struct {
			Files []string `json:"files"`
		}{
			Files: p.config.Adapter.ListLogFiles(key),
		})
	})

	router.Get("/logs/:key/:file", func(c *fiber.Ctx) error {
		key := c.Params("key")
		logFile := c.Params("file")
		logs, err := p.config.Adapter.GetLogs(key, logFile)
		if err != nil {
			return c.Status(400).Send([]byte(err.Error()))
		}
		return c.Status(200).Send(logs)
	})
}
