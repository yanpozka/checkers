package game

import (
	"bytes"
	"testing"
)

const playerID = "player-123-456"

func TestInitGame(t *testing.T) {
	buf := InitGame([]byte(playerID))

	pA, nP, _, status := ReadGame(buf)

	if !bytes.Equal(pA, []byte(playerID)) || !bytes.Equal(nP, []byte(playerID)) {
		t.Fatalf("Expected %q, got playerA: %q and nextPlayer: %q\n", playerID, pA, nP)
	}
	if status != StatusWaitingOpponent {
		t.Fatalf("Expected: %d, got: %d", StatusWaitingOpponent, status)
	}
}

func TestMakeGame(t *testing.T) {
	buf := MakeGame([]byte(playerID), firstBoard)

	pA, nP, board, status := ReadGame(buf)

	if !bytes.Equal(pA, []byte(playerID)) || !bytes.Equal(nP, []byte(playerID)) {
		t.Fatalf("Expected %q, got playerA: %q and nextPlayer: %q\n", playerID, pA, nP)
	}
	if status != StatusPlaying {
		t.Fatalf("Expected: %d, got: %d", StatusWaitingOpponent, status)
	}

	if board != firstBoard {
		t.Fatalf("Boards aren't equals: %v\n", board)
	}
}

//
// Based in benchmarks from https://github.com/rw/go-flatbuffers-example
// all rights reserved for the author
//

func BenchmarkWrite(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := MakeGame([]byte(playerID), firstBoard)
		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}
	}
}

func BenchmarkRoundtrip(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := MakeGame([]byte(playerID), firstBoard)

		pA, nP, board, _ := ReadGame(buf)
		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}

		if !bytes.Equal(pA, []byte(playerID)) || !bytes.Equal(nP, []byte(playerID)) {
			b.Fatalf("Expected %q, got playerA: %q and nextPlayer: %q\n", playerID, pA, nP)
		}

		if board != firstBoard {
			b.Fatalf("Boards aren't equals: %v\n", board)
		}
	}
}
