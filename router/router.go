package router

import (
	"net/http"

	"github.com/EdsonHTJ/stockfish-api/chess"
	"github.com/EdsonHTJ/stockfish-api/dto"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

// Chess Driver
var chessDriver *chess.Driver

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

func HandleWebSocket(c *gin.Context) {
	handler := websocket.Handler(func(ws *websocket.Conn) {
		var request dto.MoveRequest
		err := websocket.JSON.Receive(ws, &request)
		if err != nil {
			websocket.JSON.Send(ws, dto.ErrorResponse{Error: err.Error()})
			return
		}

		table := request.Table
		if !table.IsValid() {
			websocket.JSON.Send(ws, dto.ErrorResponse{Error: "Invalid table"})
			return
		}

		move, err := chessDriver.Move(request.Level, table)
		if err != nil {
			websocket.JSON.Send(ws, dto.ErrorResponse{Error: err.Error()})
			return
		}

		websocket.JSON.Send(ws, dto.MoveResponse{Move: move.Move, FenTable: string(move.Table)})
	})

	handler.ServeHTTP(c.Writer, c.Request)
}

// Creates a new gintonic rote
func New() *gin.Engine {
	chessDriver = chess.New()

	r := gin.Default()
	r.POST("/move", MoveHandler)
	r.GET("/ws", HandleWebSocket)
	return r
}
