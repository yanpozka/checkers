package game

import (
	flatbuffers "github.com/google/flatbuffers/go"
	m "github.com/yanpozka/checkers/game/messages"
)

const (
	StatusWaitingOpponent int8 = iota + 1
	StatusPlaying
	StatusEnded
)

const (
	PlayerA int8 = 1
	PlayerB      = 2
)

var gameBuilder = flatbuffers.NewBuilder(0)

func InitGame(playerID []byte) []byte {
	gameBuilder.Reset()

	pAPosition := gameBuilder.CreateByteString(playerID)

	m.GameStart(gameBuilder)
	m.GameAddPlayerA(gameBuilder, pAPosition)
	m.GameAddNextPlayer(gameBuilder, pAPosition)
	m.GameAddStatus(gameBuilder, StatusWaitingOpponent)

	gamePosition := m.GameEnd(gameBuilder)

	gameBuilder.Finish(gamePosition)

	return gameBuilder.Bytes[gameBuilder.Head():]
}

func MakeGame(playerID []byte, board [8][8]int8) []byte {
	gameBuilder.Reset()

	pAPosition := gameBuilder.CreateByteString(playerID)

	var boardPosition flatbuffers.UOffsetT
	{
		var rows [8]flatbuffers.UOffsetT

		for rx, countRows := 7, 0; rx >= 0; rx-- { // start allocation from the last cell

			m.RowStartCellsVector(gameBuilder, 8)
			for cx := 7; cx >= 0; cx-- {
				gameBuilder.PlaceInt8(board[rx][cx])
			}
			cellsPosition := gameBuilder.EndVector(8)

			m.RowStart(gameBuilder)
			m.RowAddCells(gameBuilder, cellsPosition)
			rowPosition := m.RowEnd(gameBuilder)

			rows[countRows] = rowPosition
			countRows++
		}

		m.GameStartBoardVector(gameBuilder, 8)
		for _, rowPosition := range rows {
			gameBuilder.PrependUOffsetT(rowPosition)
		}
		boardPosition = gameBuilder.EndVector(8)
	}

	m.GameStart(gameBuilder)
	m.GameAddBoard(gameBuilder, boardPosition)
	m.GameAddPlayerA(gameBuilder, pAPosition)
	m.GameAddNextPlayer(gameBuilder, pAPosition)
	m.GameAddStatus(gameBuilder, StatusPlaying)

	gamePosition := m.GameEnd(gameBuilder)

	gameBuilder.Finish(gamePosition)

	return gameBuilder.Bytes[gameBuilder.Head():]
}

func ReadGame(buf []byte) ([]byte, []byte, [8][8]int8, int8) {
	game := m.GetRootAsGame(buf, 0)

	var arr [8][8]int8

	for rx, rlen := 0, game.BoardLength(); rx < rlen; rx++ {
		var rowObj m.Row
		game.Board(&rowObj, rx)

		for cx, clen := 0, rowObj.CellsLength(); cx < clen; cx++ {
			arr[rx][cx] = rowObj.Cells(cx)
		}
	}

	return game.PlayerA(), game.NextPlayer(), arr, game.Status()
}

var firstBoard = [8][8]int8{
	{0, PlayerA, 0, PlayerA, 0, PlayerA, 0, PlayerA},
	{PlayerA, 0, PlayerA, 0, PlayerA, 0, PlayerA, 0},
	{0, PlayerA, 0, PlayerA, 0, PlayerA, 0, PlayerA},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{PlayerB, 0, PlayerB, 0, PlayerB, 0, PlayerB, 0},
	{0, PlayerB, 0, PlayerB, 0, PlayerB, 0, PlayerB},
	{PlayerB, 0, PlayerB, 0, PlayerB, 0, PlayerB, 0},
}
