package www

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TradfriFields struct {
	Gateway string `label:"Gateway"`
	Key     string `label:"Key"`
}

func Interface_Server() {
	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"testFunction": test_html,
	})

	router.Static("/bootstrap", "www/static/bootstrap")
	router.LoadHTMLGlob("www/templates/*.tmpl")

	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":        "Tradfri2MQTT",
			"active_page":  "tradfri",
			"mqtt_active":  "",
			"form_tradfri": QuickForm("/tradfri", TradfriFields{}),
		})

	})

	router.POST("/tradfri", func(c *gin.Context) {
		var tradfri TradfriFields

		// fmt.Println(c.PostForm("Gateway"))

		if err := c.ShouldBind(&tradfri); err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "Ok", "Gateway": tradfri.Gateway})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	router.Run(":8321")
}

func test_html(lbl string, msg string) template.HTML {
	return template.HTML("<h2>test</h2>")
}
