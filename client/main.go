package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type drawBlock struct {
	x    int
	y    int
	opts *ebiten.DrawImageOptions
}

var snake []drawBlock
var score = 0

var square *ebiten.Image
var prevTime int64
var gameOver = false
var waitGame = true

const blockSize = 16
const boardSize = 10

var snakeMoveMSec int64 = 500

const (
	left = iota
	right
	up
	down
)

var food *drawBlock

var nextMoveDir = left

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func checkNewHead(newHead drawBlock) bool {
	if newHead.x < 0 || newHead.y < 0 || newHead.x >= boardSize || newHead.y >= boardSize {
		return true
	}

	for i := 0; i < len(snake); i++ {
		if snake[i].x == newHead.x && snake[i].y == newHead.y {
			return true
		}
	}

	return false
}

var (
	nameText = ""
	counter  = 0
)

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func update(screen *ebiten.Image) error {
	if waitGame {
		ebitenutil.DebugPrint(screen, "click to start")
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			gameOver = false
			waitGame = false
		}
		return nil
	}

	if gameOver {
		if len(nameText) < 3 {
			nameText += string(ebiten.InputChars())
		}
		if len(nameText) == 2 {
			gameOver = false
			waitGame = true
		}

		t := "input your name : " + nameText
		t += "_"
		ebitenutil.DebugPrintAt(screen, t, 0, 20)
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		nextMoveDir = up
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		nextMoveDir = left
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		nextMoveDir = down
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		nextMoveDir = right
	}

	now := makeTimestamp()
	if !gameOver && now-prevTime > snakeMoveMSec {
		snakeHead := &snake[0]
		var newHead drawBlock
		switch nextMoveDir {
		case left:
			newHead = drawBlock{x: snakeHead.x - 1, y: snakeHead.y}
		case right:
			newHead = drawBlock{x: snakeHead.x + 1, y: snakeHead.y}
		case up:
			newHead = drawBlock{x: snakeHead.x, y: snakeHead.y - 1}
		case down:
			newHead = drawBlock{x: snakeHead.x, y: snakeHead.y + 1}
		}

		if checkNewHead(newHead) {
			gameOver = true
			return nil
		}

		snake = append([]drawBlock{newHead}, snake...)

		if checkFood(newHead) {
			food = nil
			newFood()
			score++
		} else {
			_, snake = snake[len(snake)-1], snake[:len(snake)-1]
		}

		prevTime = now
	}

	for i := 0; i < len(snake); i++ {
		element := &snake[i]
		if element.opts == nil {
			element.opts = &ebiten.DrawImageOptions{}
			element.opts.GeoM.Translate(float64(element.x*blockSize), float64(element.y*blockSize))
		}
		square.Fill(color.White)
		screen.DrawImage(square, element.opts)
	}

	if food.opts == nil {
		food.opts = &ebiten.DrawImageOptions{}
		food.opts.GeoM.Translate(float64(food.x*blockSize), float64(food.y*blockSize))
	}
	square.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})
	screen.DrawImage(square, food.opts)

	if gameOver {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("score : %v / gameOver!", score))
	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("score : %v", score))
	}

	return nil
}

func newFood() {
	if food != nil {
		return
	}

	food = &drawBlock{x: rand.Intn(boardSize), y: rand.Intn(boardSize)}
}

func checkFood(newHead drawBlock) bool {
	if food.x == newHead.x && food.y == newHead.y {
		return true
	}
	return false
}

func sendScore(name string, score int) {
	// 간단한 http.PostForm 예제
	resp, err := http.PostForm("http://localhost:1323/api/rank", url.Values{"Name": {name}, "Score": {strconv.Itoa(score)}})
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	// // Response 체크.
	// respBody, err := ioutil.ReadAll(resp.Body)
	// if err == nil {
	// 	str := string(respBody)
	// 	println(str)
	// }
}

func main() {
	snake = append(snake, drawBlock{x: 5, y: 5})
	prevTime = makeTimestamp()
	gameOver = true
	waitGame = true
	newFood()

	square, _ = ebiten.NewImage(blockSize, blockSize, ebiten.FilterNearest)
	square.Fill(color.White)

	if err := ebiten.Run(update, 160, 160, 2, "SnakeClient"); err != nil {
		log.Fatal(err)
	}
}
