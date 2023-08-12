package router

import (
	"net/http"

	"github.com/EdsonHTJ/stockfish-api/chess"
	"github.com/EdsonHTJ/stockfish-api/dto"
	"github.com/gin-gonic/gin"
)

//Chess Driver
var chessDriver *chess.Driver

func init() {
	chessDriver = chess.New()
}

func MoveHandler(c *gin.Context) {
	var req dto.MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	table := req.Table
	if !table.IsValid() {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	move, err := chessDriver.Move(req.Level, table)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(200, dto.MoveResponse{Move: move.Move, FenTable: string(move.Table)})
}

//Creates a new gintonic rote
func New() *gin.Engine {
	r := gin.Default()
	r.POST("/move", MoveHandler)
	return r
}
