package driver

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/notnil/chess"
)

const (
	EXE_FILE_PATH_ENV = "PATH_TO_EXECUTABLE"
	DEFAULT_PATH      = "./stockfish/stockfish-ubuntu-x86-64-avx2"
)

type TableState string

type Move struct {
	Move  string
	Table TableState
}

type Driver struct {
	exePath string
}

const (
	BASE_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

func New() *Driver {
	exepath := os.Getenv(EXE_FILE_PATH_ENV)
	if exepath == "" {
		exepath = DEFAULT_PATH
	}

	return &Driver{exePath: exepath}
}

func (d *Driver) Move(skillLevel uint16, state TableState) (*Move, error) {
	if skillLevel > 20 {
		skillLevel = 20
	}

	commandList := createMoveCommandString(skillLevel, state)
	output, err := d.ExecStockfishCommand(commandList)
	if err != nil {
		return nil, fmt.Errorf("error on stockfish engine")
	}

	moveTxt := parseOutput(output)
	if moveTxt == "" {
		return nil, errors.New("stockfish: couldn't parse stockfish output - " + output)
	}

	newFEN, err := updateFENWithMove(string(state), moveTxt)
	if err != nil {
		return nil, errors.New("stockfish: couldn't update FEN with move - " + err.Error())
	}

	return &Move{Move: moveTxt, Table: newFEN}, nil
}

func (d *Driver) ExecStockfishCommand(commandList []string) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	cmd := exec.Command(d.exePath)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	stdi, err := cmd.StdinPipe()
	if err != nil {
		return "", errors.New("stockfish: error occured when creating stdin pipe " + err.Error())
	}

	if err := cmd.Start(); err != nil {
		return "", errors.New("stockfish: error occured when running stockfish command " + err.Error())
	}

	for _, s := range commandList {
		stdi.Write([]byte(s))
	}

	stdi.Close()
	if err := cmd.Wait(); err != nil {
		return "", errors.New("stockfish: error occured when waiting for stockfish command to finish " + err.Error())
	}

	output := buf.String()
	return output, nil
}

func parseOutput(output string) string {
	output = strings.Replace(output, "\n", " ", -1)
	words := strings.Split(output, " ")
	next := false
	for _, word := range words {
		if next {
			return word
		}
		if word == "bestmove" {
			next = true
		}
	}
	return ""
}

func updateFENWithMove(fenStr string, moveStr string) (TableState, error) {
	fenOpts, err := chess.FEN(fenStr)
	if err != nil {
		return "", err
	}

	game := chess.NewGame(fenOpts)
	err = game.MoveStr(moveStr)
	if err != nil {
		return "", err
	}

	return TableState(game.FEN()), nil
}

func createMoveCommandString(skilllevel uint16, state TableState) []string {
	//Refactor this for each string
	s := make([]string, 0)
	s = append(s, fmt.Sprintf("setoption name Skill Level value %d\n", skilllevel))
	s = append(s, "position fen "+string(state)+"\n")
	s = append(s, "go movetime 150\n")
	s = append(s, "d\n")
	return s
}
