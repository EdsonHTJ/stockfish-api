package router

import (
	"github.com/EdsonHTJ/stockfish-api/chess"
	"github.com/EdsonHTJ/stockfish-api/dto"
	"golang.org/x/net/websocket"
)

func PlayMoveWs(ws *websocket.Conn, request dto.MoveWsRequest) dto.MoveResponse {
	table := request.Table
	if !table.IsValid() {
		websocket.JSON.Send(ws, dto.ErrorResponse{Error: "Invalid table"})
		return dto.MoveResponse{}
	}

	move, err := chessDriver.Move(request.Level, table)
	if err != nil {
		websocket.JSON.Send(ws, dto.ErrorResponse{Error: err.Error()})
		return dto.MoveResponse{}
	}

	response := dto.MoveResponse{Move: move.Move, FenTable: string(move.Table)}
	websocket.JSON.Send(ws, response)
	return response
}

func PlayGameWs(ws *websocket.Conn, request dto.MoveWsRequest) {
	move := chess.Move{}
	move.Table = request.Table
	move.Move = ""

	moves := 0
	for (move.Move != "(none)") && (moves <= request.MoveLimits) {
		played := PlayMoveWs(ws, request)
		move.Table = chess.TableState(played.FenTable)
		move.Move = played.Move
		request.Table = move.Table
		moves++
	}
}
