package game

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetScriptInfoEndpoint(ctx *gin.Context) {
	scriptID, ok := ctx.GetQuery("scriptID")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing scriptID"})
		return
	}

	if scriptID != "trouble_brewing" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Script not available, use 'trouble_brewing'"})
	}

	descriptions, err := GetAllRoleDescriptions()
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not load role descriptions"})
	} else {
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse role descriptions"})
		}
		ctx.JSON(http.StatusOK, gin.H{"roles": descriptions})
	}

}
