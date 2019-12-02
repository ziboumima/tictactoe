package main

import (
	"fmt"
	"math"
	"math/rand"
)

type MoveCoordinate struct {
	boardCoordinate int
	coordinate int
}


type UltimateState struct {
	self bool
	result float64
	ultimateBoard UltimateBoard
	lastMove MoveCoordinate // value equal -1 when there is no lastMove
	bestMove *UltimateState
}


type UltimateBoard = [9]Board

func emptyUltimateBoard() UltimateBoard {
	var board UltimateBoard
	for i := range board {
		for j := range board[i] {
			board[i][j] = EMPTY
		}
	}
	return board
}

func findWinner(board *Board) int {
	if checkWinner(SELF, board) {
		return SELF
	} else if checkWinner(OPPONENT, board) {
		return OPPONENT
	} else {
		return EMPTY
	}
}

/*
Return SELF, OPPONENT or EMPTY means no winner
 */
func findWinnerUltimate(ultimateBoard *UltimateBoard) int {
	board := emptyBoard()
	for i, b:= range ultimateBoard {
		board[i] = findWinner(&b)
	}
	return findWinner(&board)
}

func nextBoards(player int, board *Board) ([]Board, []int) {
	var boards []Board
	var coordinates []int
	for i, v := range board {
		if v == EMPTY {
			newBoard := *board
			newBoard[i] = player
			boards = append(boards, newBoard)
			coordinates = append(coordinates, i)
		}
	}
	return boards, coordinates
}

func player(self bool) int {
	if self {
		return SELF
	} else {
		return OPPONENT
	}
}

func possibleUltimateState(state *UltimateState, nextBoard int) []UltimateState {
	var result []UltimateState
	possibleBoards, coordinates := nextBoards(player(state.self), &state.ultimateBoard[nextBoard])
	for j, p := range possibleBoards {
		cloned := *state
		cloned.ultimateBoard[nextBoard] = p
		cloned.self = !cloned.self
		cloned.lastMove = MoveCoordinate{nextBoard, coordinates[j]}
		cloned.bestMove = nil
		result = append(result, cloned)
	}
	return result
}

func findNextPossibilitiesUltimate(state *UltimateState) []UltimateState {
	var result []UltimateState
	if state.lastMove.boardCoordinate == -1 {
		for i:= range state.ultimateBoard {
			for _, s:= range possibleUltimateState(state, i) {
				result = append(result, s)
			}
		}
	} else {
		nextBoard := state.lastMove.coordinate
		for _, s:= range possibleUltimateState(state, nextBoard) {
			result = append(result, s)
		}
	}
	return result
}

func scoreGameState(ultimateBoard *UltimateBoard) float64 {
	b := emptyBoard()
	for i, b:= range ultimateBoard {
		b[i] = findWinner(&b)
	}
	score := 0.0
	for _, v := range b {
		if v == SELF {
			score += 1
		} else if v == OPPONENT {
			score -= 1
		}
	}
	return score
}


func setResult(ultimateState *UltimateState, result float64) {
	ultimateState.result = result
	ultimateState.bestMove = nil
}


func alphaBeta(ultimateState *UltimateState, depth int, alpha float64, beta float64) {
	winner := findWinnerUltimate(&ultimateState.ultimateBoard)
	if winner == SELF {
		setResult(ultimateState, 10)
	} else if winner == OPPONENT {
		setResult(ultimateState, -10)
	} else if depth == 0 {
		setResult(ultimateState, scoreGameState(&ultimateState.ultimateBoard))
	} else {
		nextPossibilities := findNextPossibilitiesUltimate(ultimateState)
		if len(nextPossibilities) == 0 {
			setResult(ultimateState, 0)
		} else if ultimateState.self {
			value := math.Inf(-1)
			for i, s := range nextPossibilities {
				alphaBeta(&s, depth - 1, alpha, beta)
				value = math.Max(value, s.result)
				alpha = math.Max(alpha, value)
				if alpha >= beta {
					ultimateState.result = value
					ultimateState.bestMove = &nextPossibilities[i]
					break
				}
			}

		} else {
			value := math.Inf(1)
			for i, s := range nextPossibilities {
				alphaBeta(&s, depth - 1, alpha, beta)
				value = math.Min(value, s.result)
				beta = math.Min(beta, value)
				if alpha >= beta {
					ultimateState.result = value
					ultimateState.bestMove = &nextPossibilities[i]
					break
				}
			}
		}
	}
}

func play(ultimateState *UltimateState, depth int) *UltimateState {
	alphaBeta(ultimateState, depth, math.Inf(-1), math.Inf(1))
	if ultimateState.bestMove == nil {
		// create a random move
		possibleMoves := findNextPossibilitiesUltimate(ultimateState)
		nb := len(possibleMoves)
		var chosen UltimateState
		if nb != 0 {
			chosen = possibleMoves[rand.Int31n(int32(nb))]
			return &chosen
		} else {
			return nil
		}
	} else {
		return ultimateState.bestMove
	}
}

func moveUltimate(state *UltimateState, i int, j int) {
	boardi := i / 3
	boardj := j / 3
	boardCoordinate := boardj* 3 + boardi
	cellCoordinate := (j % 3) * 3 + i % 3
	lastMove := MoveCoordinate{boardCoordinate, cellCoordinate}
	state.ultimateBoard[boardCoordinate][cellCoordinate] = player(state.self)
	state.self = !state.self
	state.lastMove = lastMove
	state.bestMove = nil
}

func convertCoordinate(oneDim int) (int, int) {
	return oneDim % 3, oneDim / 3
}

func computeMove(move MoveCoordinate) (int, int) {
	boardX, boardY := convertCoordinate(move.boardCoordinate)
	i, j := convertCoordinate(move.coordinate)
	return boardX * 3 + i, boardY * 3 + j
}

func runUltimate() {
	state := UltimateState{true, 0.0,emptyUltimateBoard(),
		MoveCoordinate{-1, -1}, nil}
	for {
		var opponentRow, opponentCol int
		_, _ = fmt.Scan(&opponentRow, &opponentCol)

		var validActionCount int
		_, _ = fmt.Scan(&validActionCount)

		for i := 0; i < validActionCount; i++ {
			var row, col int
			_, _ = fmt.Scan(&row, &col)
		}

		if opponentRow == -1 && opponentCol == -1 {
			state.self = true
		} else {
			state.self = false
			moveUltimate(&state, opponentRow, opponentCol)
		}

		next := play(&state, 5)
		if next != nil {
			x, y := computeMove(next.lastMove)
			fmt.Println(x, y)// Write action to stdout
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")


	}
}