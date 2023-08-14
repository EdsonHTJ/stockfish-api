package router_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/EdsonHTJ/stockfish-api/chess"
	"github.com/EdsonHTJ/stockfish-api/dto"
	"github.com/EdsonHTJ/stockfish-api/router"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
)

const (
	TEST_EXEC_PATH      = "../stockfish/stockfish-ubuntu-x86-64-avx2"
	TEST_SCRIPTS_FOLDER = "../scripts"
	URL                 = "http://localhost:8080"
	MOVE_PATH           = "/move"
	WEBSOCKET_PATH      = "/ws"
)

// TestMain
func TestMain(m *testing.M) {
	os.Setenv(chess.EXE_FILE_PATH_ENV, TEST_EXEC_PATH)
	os.Setenv(chess.SCRIPTS_FILE_PATH_ENV, TEST_SCRIPTS_FOLDER)
	go router.New().Run(":8080")
	m.Run()
}

func TestMoveHandler(t *testing.T) {
	moveRequest := dto.MoveRequest{Table: chess.BASE_FEN, Level: 20}
	b, err := json.Marshal(moveRequest)
	require.NoError(t, err)

	bytesBody := bytes.NewReader(b)
	response, err := http.Post(URL+MOVE_PATH, "application/json", bytesBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)

	var moveResponse dto.MoveResponse
	responseBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	err = json.Unmarshal(responseBytes, &moveResponse)
	require.NoError(t, err)

	require.NotEmpty(t, moveResponse.Move)
	require.NotEmpty(t, moveResponse.FenTable)
}

func TestWebSocketHandler(t *testing.T) {
	moveWsRequest := dto.MoveWsRequest{
		ReqType: dto.PLAY_MOVE,
		MoveRequest: dto.MoveRequest{
			Table: chess.BASE_FEN,
			Level: 20,
		},
	}

	ws, err := websocket.Dial("ws://localhost:8080/ws", "", URL)
	require.NoError(t, err)

	err = websocket.JSON.Send(ws, moveWsRequest)
	require.NoError(t, err)

	var moveResponse dto.MoveResponse
	err = websocket.JSON.Receive(ws, &moveResponse)

	require.NoError(t, err)
	require.NotEmpty(t, moveResponse.Move)
	require.NotEmpty(t, moveResponse.FenTable)
}

func TestPlayGameWs(t *testing.T) {
	moveWsRequest := dto.MoveWsRequest{
		ReqType:    dto.PLAY_GAME,
		MoveLimits: 5,
		MoveRequest: dto.MoveRequest{
			Table: chess.BASE_FEN,
			Level: 20,
		},
	}

	ws, err := websocket.Dial("ws://localhost:8080/ws", "", URL)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(5)
	go func() {
		for {
			var moveResponse dto.MoveResponse
			err = websocket.JSON.Receive(ws, &moveResponse)
			require.NoError(t, err)
			require.NotEmpty(t, moveResponse.Move)
			require.NotEmpty(t, moveResponse.FenTable)
			t.Log(moveResponse)
			wg.Done()
		}
	}()

	err = websocket.JSON.Send(ws, moveWsRequest)
	require.NoError(t, err)

	wg.Wait()
}
