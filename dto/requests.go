package dto

import "github.com/EdsonHTJ/stockfish-api/chess"

type MoveRequest struct {
	Table chess.TableState `json:"table"`
	Level uint16           `json:"level"`
}
