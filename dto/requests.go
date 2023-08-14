package dto

import "github.com/EdsonHTJ/stockfish-api/chess"

type MoveRequest struct {
	Table chess.TableState `json:"table" example:"2k5/8/3b4/8/8/8/4R3/K1R5 b - - 0 1"`
	Level uint16           `json:"level" example:"20"`
}
