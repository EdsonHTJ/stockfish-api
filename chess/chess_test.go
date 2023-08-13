package chess_test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/EdsonHTJ/stockfish-api/chess"
	"github.com/stretchr/testify/require"
)

const (
	TEST_EXEC_PATH      = "../stockfish/stockfish-ubuntu-x86-64-avx2"
	TEST_SCRIPTS_FOLDER = "../scripts"
)

// TestMain
func TestMain(m *testing.M) {
	os.Setenv(chess.EXE_FILE_PATH_ENV, TEST_EXEC_PATH)
	os.Setenv(chess.SCRIPTS_FILE_PATH_ENV, TEST_SCRIPTS_FOLDER)
	m.Run()
}

func TestNextMove(t *testing.T) {
	stockFish := chess.New()
	Move, err := stockFish.Move(20, chess.BASE_FEN)
	require.NoError(t, err)

	t.Log(Move)
}

func TestFullGame(t *testing.T) {
	stockFish := chess.New()

	//Move 1
	move := &chess.Move{Table: "2k5/8/3b4/8/8/8/4R3/K1R5 b - - 0 1"}
	for {
		var err error
		tableCpy := move.Table
		fmt.Println(tableCpy)
		fmt.Println(move.Move)
		move, err = stockFish.Move(20, tableCpy)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestParalelGames(t *testing.T) {
	stockFish := chess.New()

	wg := sync.WaitGroup{}
	wg.Add(10)
	createGame := func(index int, t *testing.T) {
		move := &chess.Move{Table: "2k5/8/3b4/8/8/8/4R3/K1R5 b - - 0 1"}
		for {
			var err error
			tableCpy := move.Table
			fmt.Println("Game: ", index, "Stauts: ", move)
			move, err = stockFish.Move(20, tableCpy)
			if err != nil {
				t.Error(err)
				return
			}

			if move.Move == "(none)" {
				fmt.Println("Game: ", index, "Stauts: ", move)
				break
			}
		}

		wg.Done()
	}

	for i := 0; i < 10; i++ {
		go createGame(i, t)
	}

	wg.Wait()
}
