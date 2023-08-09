package driver_test

import (
	"os"
	"testing"

	"github.com/EdsonHTJ/stockfish-api/driver"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_EXEC_PATH = "../stockfish/stockfish-ubuntu-x86-64-avx2"
)

//TestMain
func TestMain(m *testing.M) {
	os.Setenv(driver.EXE_FILE_PATH_ENV, TEST_EXEC_PATH)
	m.Run()
}

func TestNextMove(t *testing.T) {
	stockFish := driver.New()
	Move, err := stockFish.Move(20, driver.BASE_FEN)
	t.Log(Move)
	assert.NoError(t, err)
}
