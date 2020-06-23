package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"time"
)

// board size
const size int = 8

// pos structure for board
type pos struct {
	x, y int
}

// check does pos with the same position is already on the desk
func doesPosExist(p []pos, x, y int) bool {
	for _, v := range p {
		if x == v.x && y == v.y {
			return true
		}
	}
	return false
}

// check intersection between new pos and others
func doesPosIntersect(p []pos, q pos) bool {
	for _, v := range p {
		if checkIntersection(v, q) {
			return true
		}
	}
	return false
}

// check intersection between two poses
func checkIntersection(p pos, q pos) bool {
	// check x axis
	if p.x == q.x {
		return true
	}
	// check y axis
	if p.y == q.y {
		return true
	}

	// check bot right diagonal
	x, y := q.x, q.y
	for x < size && y < size {
		if p.x == x && p.y == y {
			return true
		}
		x++
		y++
	}

	// check top left diagonal
	x, y = q.x, q.y
	for x >= 0 && y >= 0 {
		if p.x == x && p.y == y {
			return true
		}
		x--
		y--
	}

	// check top right diagonal
	x, y = q.x, q.y
	for x < size && y >= 0 {
		if p.x == x && p.y == y {
			return true
		}
		x++
		y--
	}

	// check bot left diagonal
	x, y = q.x, q.y
	for x >= 0 && y < size {
		if p.x == x && p.y == y {
			return true
		}
		x--
		y++
	}

	return false
}

// generate new original pos
func generatePos(p []pos) pos {
	for {
		x := rand.Intn(size)
		y := rand.Intn(size)
		q := pos{x, y}
		if !doesPosIntersect(p, q) {
			return q
		}
	}
}

// sort poses by x and then compare new pos scope with existed
func checkDoesPosesEqual(a [][]pos, b []pos) bool {
	sort.SliceStable(b, func(i, j int) bool {
		return b[i].x < b[i].x && b[i].y < b[i].y
	})
	for c := 0; c < len(a); c++ {
		sort.SliceStable(a[c], func(i, j int) bool {
			return a[c][i].x < a[c][i].x && a[c][i].y < a[c][i].y
		})
		if reflect.DeepEqual(a[c], b) {
			return true
		}
	}
	return false
}

func main() {
	// set count
	count := 10
	// 92 is maximum count of ways
	if count > 92 {
		count = 92
	}
	// chan for getting new valid position scope
	send := make(chan bool)
	// create list of position scopes
	posss := make([][]pos, count)

	// run async func
	go func() {
		for {
			// chan for getting new valid position scope
			done := make(chan []pos)
			rand.Seed(time.Now().UnixNano())

			// run 100 async funcs that trying to generate new pos scope
			for i := 0; i < 100; i++ {
				go func() {
					poss := make([]pos, 0, 8)
					for i := 0; i < 8; i++ {
						// generate pos and check intersection
						pos := generatePos(poss)
						poss = append(poss, pos)
					}
					// send succed pos scope
					done <- poss
				}()
			}

			select {
			case poss := <-done:
				// check does same scope of poses exist
				if !checkDoesPosesEqual(posss, poss) {
					fmt.Println(poss)
					posss = append(posss, poss)
					send <- true
				}
			}
		}
	}()

	// get count of poses scope
	for j := 0; j < count; {
		select {
		case <-send:
			j++
		}
	}

}
