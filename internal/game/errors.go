package game

import "errors"

var (
	ErrGameNotFound        = errors.New("game not found")
	ErrGameFull            = errors.New("game is full")
	ErrGameAlreadyStarted  = errors.New("game already started")
	ErrGameNotStarted      = errors.New("game not started")
	ErrGameNotFinished     = errors.New("game not finished")
	ErrGameFinished        = errors.New("game already finished")
	ErrNotYourTurn         = errors.New("not your turn")
	ErrInvalidAction       = errors.New("invalid action")
	ErrPlayerNotFound      = errors.New("player not found")
	ErrPlayerAlreadyInGame = errors.New("player already in game")
	ErrInsufficientPlayers = errors.New("insufficient players to start game (minimum 2)")
	ErrTooManyPlayers      = errors.New("too many players (maximum 4)")
	ErrInvalidPlayerCount  = errors.New("player count must be between 2 and 4")
	ErrCorruptedState      = errors.New("corrupted game state")
	ErrSaveFailed          = errors.New("failed to save game")
	ErrLoadFailed          = errors.New("failed to load game")
	ErrTurnLimitReached    = errors.New("turn limit reached")
)
