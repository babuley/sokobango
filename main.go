package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/danicat/simpleansi"
)

var keys = map[string]string{
	"A": "up",
	"B": "down",
	"C": "right",
	"D": "left",
}

var player Player

func parseLevel(rawLevels []string) [][]string {
	var maps [][]string
	var mapa []string
	var lidx int
	for _, line := range rawLevels {
		if strings.Contains(line, "Maze") {
			if len(mapa) > 0 {
				mapa = mapa[:len(mapa)]
				maps = append(maps, mapa)
				mapa = []string{}
			}
			lidx = 0
		}
		if lidx > 6 {
			mapa = append(mapa, line)
		}
		lidx++
	}
	return maps
}

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
	for y, line := range level {
		for x, char := range line {
			switch char {
			case '@':
				return Player{x, y}
			}
		}
	}
	return Player{}
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

func printMap(maps [][]string, idx int) {
	simpleansi.ClearScreen()
	for _, line := range maps[idx] {
		for _, chr := range line {
			switch chr {
			case 'X':
				fmt.Printf("%c", chr)

			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
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

func calculateMove(level []string, fromX int, fromY int, direction string) (toX int, toY int) {
	toX, toY = moves[direction](level, fromX, fromY)
	if level[toX][toY] == 'X' {
		toX = fromX
		toY = fromY
	}
	return
}

func movePlayer(level []string, dir string) {
	fmt.Println(player)
	fmt.Println(dir)
	player.X, player.Y = calculateMove(level, player.X, player.Y, dir)
}

func main() {
	initialise()
	allLevels, _ := loadLevel("levels/maps.txt")
	defer cleanup()
	startLevel := 0
	maps := parseLevel(allLevels)
	level := maps[startLevel]
	player = initPlayer(level)
	fmt.Println("Player reports for duty ", player)
	// game loop
	for {
		printMap(maps, startLevel)
		input, err := readInput()
		if err != nil {
			log.Println("Error reading input:", err)
			break
		}

		// process movement
		movePlayer(level, input)

		// process collisions

		// check game over
		if input == "ESC" {
			break
		}
		// repeat
	}
}
