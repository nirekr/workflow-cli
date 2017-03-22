package mock

import (
	"fmt"
	"net/http"

	"github.com/braintree/manners"
	"github.com/dellemc-symphony/workflow-cli/models"
	"github.com/gin-gonic/gin"
)

// StopMock stops the Server
func StopMock() {
	manners.Close()
}

// CreateMock starts a mock REST Endpoint for the FRU workflow
func CreateMock() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	id := "123abc-456def-789ghi"
	retry := true

	router.GET("/fru/api/about", func(c *gin.Context) {

		c.String(http.StatusOK, "up and running")

	})

	router.GET("/fru/api/data/:taskid", func(c *gin.Context) {
		taskid := c.Param("taskid")

		response := models.MockDataResponse{
			Response: "YooHoo!",
			TaskID:   taskid,
		}

		c.JSON(http.StatusOK, response)

	})

	router.POST("/fru/api/workflow", func(c *gin.Context) {
		stepNext := models.Link{
			Rel:    "step-next",
			Href:   fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/vcenter-endpoint", "http://", id),
			Type:   "application/vnd.dellemc.vcenter.endpoint+json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "Initiate Workflow",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)

	})

	router.GET("/fru/api/workflow", func(c *gin.Context) {
		workflow := models.Workflow{
			URI: fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s", "http://", id),
		}

		workflows := models.Workflows{workflow}

		c.JSON(http.StatusOK, workflows)

	})

	router.GET("/fru/api/workflow/:trackingid", func(c *gin.Context) {
		stepNext := models.Link{
			Rel:    "step-next",
			Href:   fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/vcenter-endpoint", "http://", id),
			Type:   "application/vnd.dellemc.vcenter.endpoint+json",
			Method: "POST",
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "Initiate Workflow",
			Links:       links,
		}

		c.JSON(http.StatusOK, response)

	})

	// Step 2
	router.POST("/fru/api/workflow/:trackingid/vcenter-endpoint", func(c *gin.Context) {
		id = c.Param("trackingid")

		// Validate JSON Body
		var vcenterCreds models.Endpoint
		if c.BindJSON(&vcenterCreds) == nil {

			stepNext := models.Link{
				Rel:    "step-next",
				Href:   fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/scaleio-endpoint", "http://", id),
				Type:   "application/vnd.dellemc.scaleio.endpoint+json",
				Method: "POST",
			}

			links := models.Links{stepNext}

			response := models.Response{
				ID:          id,
				Workflow:    "quanta-replacement-d51b-esxi",
				CurrentStep: "capturevCenterEndpoint",
				Links:       links,
			}

			c.JSON(http.StatusCreated, response)
		}
	})

	// Step 3
	router.POST("/fru/api/workflow/:trackingid/scaleio-endpoint", func(c *gin.Context) {
		id = c.Param("trackingid")

		// Validate JSON Body
		var scaleioCreds models.Endpoint
		if c.BindJSON(&scaleioCreds) == nil {

			stepNext := models.Link{
				Rel:    "step-next",
				Href:   fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/start-data-collection", "http://", id),
				Type:   "application/json",
				Method: "POST",
			}

			links := models.Links{stepNext}

			response := models.Response{
				ID:          id,
				Workflow:    "quanta-replacement-d51b-esxi",
				CurrentStep: "capturevScaleioEndpoint",
				Links:       links,
			}

			c.JSON(http.StatusCreated, response)
		}
	})

	// Step 4
	router.POST("/fru/api/workflow/:trackingid/start-data-collection", func(c *gin.Context) {
		var stepNext models.Link

		if retry == true {
			retry = false
			id = c.Param("trackingid")

			stepNext = models.Link{
				Rel:    "step-retry",
				Href:   fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/start-data-collection", "http://", id),
				Type:   "application/json",
				Method: "POST",
			}

		} else {
			id = c.Param("trackingid")

			stepNext = models.Link{
				Rel:    "step-next",
				Href:   fmt.Sprintf("%slocalhost:8080/fru/api/workflow/%s/verify-data-collection", "http://", id),
				Type:   "application/json",
				Method: "GET",
			}
		}

		links := models.Links{stepNext}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "startDataCollection",
			Links:       links,
		}

		c.JSON(http.StatusCreated, response)

	})

	// Step 5
	router.GET("/fru/api/workflow/:trackingid/verify-data-collection", func(c *gin.Context) {
		id = c.Param("trackingid")

		returnedNode := models.Node{
			ID:           "node01",
			SerialNumber: "123456789",
		}
		returnedNode2 := models.Node{
			ID:           "node02",
			SerialNumber: "987654321",
		}
		nodes := models.Nodes{returnedNode, returnedNode2}

		response := models.Response{
			ID:          id,
			Workflow:    "quanta-replacement-d51b-esxi",
			CurrentStep: "verifyDataCollection",
			Nodes:       nodes,
		}

		c.JSON(http.StatusCreated, response)

	})

	manners.ListenAndServe(":8080", router)
}
