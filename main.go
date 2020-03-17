package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/danicat/simpleansi"
)

var keys = map[string]string{
	"A": "up",
	"B": "down",
	"C": "right",
	"D": "left",
}

var player Player
var boulders []*Boulder
var targets []*Target

func loadLevel(file string) ([]string, error) {
	var level []string
	f, err := os.Open(file)
	if err != nil {
		return level, err
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		level = append(level, line)
	}
	return level, nil
}

func initPlayer(level []string) Player {
	for x, line := range level {
		for y, char := range line {
			switch char {
			case '@':
				return Player{x, y}
			}
		}
	}
	return Player{}
}

func initBoulder(level []string) []*Boulder {
	var boulders []*Boulder
	for x, line := range level {
		for y, char := range line {
			switch char {
			case '*':
				boulders = append(boulders, &Boulder{x, y})
			}
		}
	}
	return boulders
}

func initTarget(level []string) []*Target {
	var targets []*Target
	for x, line := range level {
		for y, char := range line {
			switch char {
			case '.':
				targets = append(targets, &Target{x, y})
			}
		}
	}
	return targets
}

func printLevel(level []string) {
	for _, line := range level {
		fmt.Println(line)
	}
}

func printMaps(maps [][]string) {
	for _, line := range maps {
		for _, s := range line {
			fmt.Println(s)
		}
	}
}

var reset = "\x1b[0m"
var oo = "\x1b[42m" + " " + reset

func printMap(maps [][]string, idx int) {
	simpleansi.ClearScreen()
	for _, line := range maps[idx] {
		for _, chr := range line {
			switch chr {
			case 'X':
				fmt.Print(simpleansi.WithBackground(" ", "GREEN"))
			case '.':
				fmt.Printf("%c", chr)
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	for _, b := range boulders {
		simpleansi.MoveCursor(b.X, b.Y)
		fmt.Print("*")
	}
	simpleansi.MoveCursor(player.X, player.Y)
	fmt.Print("@")
	simpleansi.MoveCursor(len(maps[idx])+1, 0)
}

func recoverFatal(msg string, err error) {
	if err != nil {
		log.Fatalln("Error activating cbreak mode:", err)
	}
}

func runTerminal(term *exec.Cmd) error {
	term.Stdin = os.Stdin
	return term.Run()
}

func initCooked() func() *exec.Cmd {
	return func() *exec.Cmd {
		return exec.Command("stty", "cbreak", "-echo")
	}
}

func initCBreak() func() *exec.Cmd {
	return func() *exec.Cmd {
		return exec.Command("stty", "-cbreak", "echo")
	}
}

func initialise() {
	cbTerm := initCooked()()
	err := runTerminal(cbTerm)
	recoverFatal("Error activating cbreak mode:", err)
}

func cleanup() {
	cookedTerm := initCBreak()()
	err := runTerminal(cookedTerm)
	recoverFatal("Error activating cooked mode:", err)
}

func readInput() (string, error) {
	buffer := make([]byte, 100)
	cnt, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", err

	}
	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	}

	if cnt >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' {
			return keys[string(buffer[2])], nil
		}
	}
	return "", nil
}

func moveUp(level []string, x int, y int) (int, int) {
	return x - 1, y
}

func moveDown(level []string, x int, y int) (int, int) {
	return x + 1, y
}

func moveRight(level []string, x int, y int) (int, int) {
	return x, y + 1
}

func moveLeft(level []string, x int, y int) (int, int) {
	return x, y - 1
}

var moves = map[string]func(level []string, x int, y int) (int, int){
	"up":    moveUp,
	"down":  moveDown,
	"left":  moveLeft,
	"right": moveRight,
}

func isPositionOfOccupied(level []string, x int, y int, positionMarker byte) bool {
	if positionMarker == '*' {
		b := getBoulderAtPosition(x, y)
		if b != nil {
			return true
		}
		return false
	}
	return level[x][y] == positionMarker

}

func hitWall(level []string, x int, y int) bool {
	return isPositionOfOccupied(level, x, y, 'X')
}

func isPositionOccupied(level []string, x int, y int) bool {
	return hitWall(level, x, y) || isPositionOfOccupied(level, x, y, '*')
}

func calculateMove(level []string, fromX int, fromY int, direction string) (toX int, toY int) {
	if level == nil {
		return
	}
	if moveFunc, ok := moves[direction]; ok {
		toX, toY = moveFunc(level, fromX, fromY)
		if hitWall(level, toX, toY) {
			toX = fromX
			toY = fromY
		}
		//Move boulder
		b := getBoulderAtPosition(toX, toY)
		if b != nil {
			candX, candY := moveFunc(level, toX, toY)
			if isPositionOccupied(level, candX, candY) {
				b.X, b.Y = toX, toY
				toX, toY = fromX, fromY
			} else {
				b.X, b.Y = candX, candY
			}
		}
	}
	return
}

func getBoulderAtPosition(x int, y int) *Boulder {
	for _, cand := range boulders {
		if cand.X == x && cand.Y == y {
			return cand
		}
	}
	return nil
}

func boulderAtPosition(x int, y int) bool {
	for _, cand := range boulders {
		if cand.X == x && cand.Y == y {
			return true
		}
	}
	return false
}

func movePlayer(level []string, dir string) {
	if level != nil {
		player.X, player.Y = calculateMove(level, player.X, player.Y, dir)
	}
}

func matchBoulderToTarget(b *Boulder, t *Target) bool {
	return b.X == t.X && b.Y == t.Y
}

func isLevelCompleted() bool {
	c := 0
	for _, b := range boulders {
		for _, t := range targets {
			if matchBoulderToTarget(b, t) {
				c++
			}
		}
	}
	return c == len(boulders)
}

func initLevel(allLevels []string, startLevel int) ([][]string, []string) {
	maps := ParseLevel(allLevels)
	level := maps[startLevel]
	player = initPlayer(level)
	targets = initTarget(level)
	boulders = initBoulder(level)
	return maps, level
}

func main() {
	initialise()
	allLevels, _ := LoadLevel("levels/maps.txt")
	defer cleanup()
	startLevel := 0
	maps, level := initLevel(allLevels, startLevel)

	input := make(chan string)
	go func(ch chan<- string) {
		for {
			input, err := readInput()
			if err != nil {
				log.Println("Error reading input:", err)
				ch <- "ESC"
			}
			ch <- input
		}
	}(input)
	exit := false
	// game loop
	for {

		// process movement
		select {
		case evt := <-input:
			if evt == "ESC" {
				exit = true
			}
			movePlayer(level, evt)
		default:
		}

		printMap(maps, startLevel)

		if exit {
			break
		}

		// is completed
		if isLevelCompleted() {
			fmt.Println("Level completed")
			startLevel++
			maps, level = initLevel(allLevels, startLevel)
		}

		// repeat
		time.Sleep(100 * time.Millisecond)
	}
}
