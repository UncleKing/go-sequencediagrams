package sequence

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// write api handlers
// write logic to place the participants
// extract participants and store them in a list
// precompute lengths of each message and if necessary increase the distance between two participants
// precompute heights of each line and if necessary boxes that go below them -->*
// draw each sequence
// save the image.

func RegisterSequenceHandler(router *gin.Engine, aph *jwt.GinJWTMiddleware) {
	u := GinSequenceHandler{}

	router.POST("/api/v1/sequence/", u.Sequence)
}

type GinSequenceHandler struct {
}

func (u *GinSequenceHandler) Sequence(c *gin.Context) {
	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fullText := string(rawData)
	//f, err := os.Create("proflile.pprof")
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//
	//}
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()
	responseBytes, err := CreateDiagram(fullText)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "image/png", responseBytes)
}
