package main

import (
	"bufio"
	"strings"

	"github.com/markbates/pkger"
)

//ParseLevel - processing a set of maps and splitting into levels
func ParseLevel(rawLevels []string) [][]string {
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

//LoadLevel - loading all the maps from a file
func LoadLevel(file string) ([]string, error) {
	var level []string

	f, err := pkger.Open(file)
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
