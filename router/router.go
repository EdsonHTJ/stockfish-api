package router

import (
	"net/http"

	"github.com/EdsonHTJ/stockfish-api/chess"
	docs "github.com/EdsonHTJ/stockfish-api/docs"
	"github.com/EdsonHTJ/stockfish-api/dto"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/net/websocket"
)

// Chess Driver
var chessDriver *chess.Driver

// Move route godoc
// @Summary Ask for a move
// @Schemes
// @Description Ask for a move
// @Tags example
// @Accept json
// @Produce json
// @Param body body dto.MoveRequest true "Move Request Body"
// @Example {"table": "2k5/8/3b4/8/8/8/4R3/K1R5 b - - 0 1", "level": 20}
// @Success 200 {object} dto.MoveResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /move [post]
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
		var request dto.MoveWsRequest
		err := websocket.JSON.Receive(ws, &request)
		if err != nil {
			websocket.JSON.Send(ws, dto.ErrorResponse{Error: err.Error()})
			return
		}

		switch request.ReqType {
		case dto.PLAY_MOVE:
			PlayMoveWs(ws, request)
		case dto.PLAY_GAME:
			PlayGameWs(ws, request)
		}
	})

	handler.ServeHTTP(c.Writer, c.Request)
}

// Creates a new gintonic rote
func New() *gin.Engine {
	chessDriver = chess.New()

	docs.SwaggerInfo.BasePath = "/"

	r := gin.Default()
	r.POST("/move", MoveHandler)
	r.GET("/ws", HandleWebSocket)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
