package dto

import "github.com/EdsonHTJ/stockfish-api/chess"

const (
	PLAY_MOVE = iota
	PLAY_GAME
)

type MoveRequest struct {
	Table chess.TableState `json:"table" example:"2k5/8/3b4/8/8/8/4R3/K1R5 b - - 0 1"`
	Level uint16           `json:"level" example:"20"`
}

type MoveWsRequest struct {
	ReqType    int `json:"reqType" example:"0"`
	MoveLimits int `json:"movesLimits" example:"0"`
	MoveRequest
}
