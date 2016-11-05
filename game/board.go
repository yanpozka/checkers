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
	PlayerUp   byte = 1
	PlayerDown      = 2
)

func InitGame(authorID string) []byte {
	b := flatbuffers.NewBuilder(0)

	pAPosition := b.CreateString(authorID)

	m.GameStart(b)
	m.GameAddPlayerA(b, pAPosition)
	m.GameAddNextPlayer(b, pAPosition)
	m.GameAddStatus(b, StatusWaitingOpponent)

	gamePosition := m.GameEnd(b)

	b.Finish(gamePosition)

	return b.Bytes[b.Head():]
}

func MakeGame() []byte {
	// TODO: everything
	return nil
}

func ReadGame(buf []byte) (string, string, int8) {
	game := m.GetRootAsGame(buf, 0)

	return string(game.PlayerA()), string(game.NextPlayer()), game.Status()
}

var firstBoard = [8][8]byte{
	{0, PlayerUp, 0, PlayerUp, 0, PlayerUp, 0, PlayerUp},
	{PlayerUp, 0, PlayerUp, 0, PlayerUp, 0, PlayerUp, 0},
	{0, PlayerUp, 0, PlayerUp, 0, PlayerUp, 0, PlayerUp},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{PlayerDown, 0, PlayerDown, 0, PlayerDown, 0, PlayerDown, 0},
	{0, PlayerDown, 0, PlayerDown, 0, PlayerDown, 0, PlayerDown},
	{PlayerDown, 0, PlayerDown, 0, PlayerDown, 0, PlayerDown, 0},
}
