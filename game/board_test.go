package game

import "testing"

func TestInitGame(t *testing.T) {
	const playerID = "player-123-456"

	buf := InitGame(playerID)

	pA, nP, status := ReadGame(buf)

	if pA != playerID || nP != playerID {
		t.Fatalf("Expected %q, got playerA: %q and nextPlayer: %q\n", playerID, pA, nP)
	}
	if status != StatusWaitingOpponent {
		t.Fatalf("Expected: %d, got: %d", StatusWaitingOpponent, status)
	}
}
