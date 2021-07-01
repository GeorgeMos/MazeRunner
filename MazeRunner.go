package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
)

//Globals
var sizeX int
var sizeY int
var path []*node
var suroundMap = [4][2]int{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
var nodeMap [][]node
var mazeFile string = "maze1000.png"
var nodes int = 0

//-------
//Opens the maze image
func initPxMap() [][]bool {
	maze, err := os.Open(mazeFile)
	if err != nil {
		log.Fatal(err)
	}
	defer maze.Close()

	colorMap, err := png.Decode(maze)
	if err != nil {
		log.Fatal(err)
	}

	sizeX = colorMap.Bounds().Max.X
	sizeY = colorMap.Bounds().Max.Y
	var pxMap [][]bool

	for x := 0; x < sizeX; x++ {
		var temp []bool
		for y := 0; y < sizeY; y++ {
			temp = append(temp, false)
		}
		pxMap = append(pxMap, temp)
	}

	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			if color.GrayModel.Convert(colorMap.At(x, y)).(color.Gray).Y == 255 {
				pxMap[x][y] = true
			} else {
				pxMap[x][y] = false
			}
		}
	}
	return pxMap
}

//Draws the path and saves to the output image
func saveImage(nodeMap [][]node) {
	solved := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			if nodeMap[x][y].nodeType == "wall" {
				solved.Set(x, y, color.Black)

			} else {
				solved.Set(x, y, color.White)
			}
		}
	}
	for i := 0; i < len(path)-1; i++ {
		solved.Set(path[i].pos[0], path[i].pos[1], color.RGBA{0, 255, 0, 255})

		//x++
		x := path[i].pos[0]
		y := path[i].pos[1]

		if path[i+1].pos[0] > path[i].pos[0] {
			for x < path[i+1].pos[0] {
				solved.Set(x, y, color.RGBA{0, 255, 0, 255})
				x++
			}
		}
		//x--
		x = path[i].pos[0]
		y = path[i].pos[1]

		if path[i+1].pos[0] < path[i].pos[0] {
			for x > path[i+1].pos[0] {
				solved.Set(x, y, color.RGBA{0, 255, 0, 255})
				x--
			}
		}
		//y++
		x = path[i].pos[0]
		y = path[i].pos[1]

		if path[i+1].pos[1] > path[i].pos[1] {
			for y < path[i+1].pos[1] {
				solved.Set(x, y, color.RGBA{0, 255, 0, 255})
				y++
			}
		}
		//y--
		x = path[i].pos[0]
		y = path[i].pos[1]

		if path[i+1].pos[1] < path[i].pos[1] {
			for y > path[i+1].pos[1] {
				solved.Set(x, y, color.RGBA{0, 255, 0, 255})
				y--
			}
		}

	}

	out, err := os.Create("solvedMaze.png")
	if err == nil {
		png.Encode(out, solved)
		out.Close()
	}
}

type node struct {
	nodeType      string
	adjacentNodes []*node
	openSideNum   int
	visited       bool
	pos           [2]int
}

//Initialises the node map
func initNodeMap(pxMap [][]bool) {
	for x := 0; x < sizeX; x++ {
		var temp []node
		for y := 0; y < sizeY; y++ {
			temp = append(temp, node{})
		}
		nodeMap = append(nodeMap, temp)
	}
}

//Initialises the nodes
func initNodes(pxMap [][]bool) {
	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeX; y++ {
			nodeMap[x][y].pos[0] = x
			nodeMap[x][y].pos[1] = y
			nodeMap[x][y].visited = false
			if pxMap[x][y] {
				//Scaning for type
				for i := 0; i < len(suroundMap); i++ {
					testX := x + suroundMap[i][0]
					testY := y + suroundMap[i][1]
					if testX >= 0 && testX < sizeX && testY >= 0 && testY < sizeY {
						if pxMap[testX][testY] {
							nodeMap[x][y].openSideNum++
						}
					}
				}

				if nodeMap[x][y].openSideNum >= 3 {
					nodeMap[x][y].nodeType = "junction"
					nodes++
				}
				if nodeMap[x][y].openSideNum == 1 {
					nodeMap[x][y].nodeType = "dEnd"
				}
				if nodeMap[x][y].openSideNum == 2 {
					if pxMap[x-1][y] && pxMap[x+1][y] || pxMap[x][y-1] && pxMap[x][y+1] {
						nodeMap[x][y].nodeType = "path"
					} else {
						nodeMap[x][y].nodeType = "corner"
						nodes++
					}
				}

				if y == 0 {
					nodeMap[x][y].nodeType = "start"
					nodes++
				}
				if y == sizeY-1 {
					nodeMap[x][y].nodeType = "end"
					nodes++
				}

			} else {
				nodeMap[x][y].nodeType = "wall"
			}
		}
	}
}

//Links all the nodes and creates a tree structure
func initAdjMap(pxMap [][]bool) {
	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			if pxMap[x][y] {
				//scanning x++
				ix := x
				iy := y
				for pxMap[ix][y] {
					if ix+1 < sizeX {
						ix++
					} else {
						break
					}

					if nodeMap[ix][iy].nodeType != "dEnd" && nodeMap[ix][iy].nodeType != "wall" && nodeMap[ix][iy].nodeType != "path" {
						nodeMap[x][y].adjacentNodes = append(nodeMap[x][y].adjacentNodes, &nodeMap[ix][iy])
						break
					}
				}
				//scanning y++
				ix = x
				iy = y
				for pxMap[x][iy] {
					if iy+1 < sizeY {
						iy++
					} else {
						break
					}
					if nodeMap[ix][iy].nodeType != "dEnd" && nodeMap[ix][iy].nodeType != "wall" && nodeMap[ix][iy].nodeType != "path" {
						nodeMap[x][y].adjacentNodes = append(nodeMap[x][y].adjacentNodes, &nodeMap[ix][iy])
						break
					}
				}
				//scanning x--
				ix = x
				iy = y
				for pxMap[ix][y] {
					if ix-1 >= 0 {
						ix--
					} else {
						break
					}

					if nodeMap[ix][iy].nodeType != "dEnd" && nodeMap[ix][iy].nodeType != "wall" && nodeMap[ix][iy].nodeType != "path" {
						nodeMap[x][y].adjacentNodes = append(nodeMap[x][y].adjacentNodes, &nodeMap[ix][iy])
						break
					}
				}
				//scanning y--
				ix = x
				iy = y
				for pxMap[x][iy] {
					if iy-1 >= 0 {
						iy--
					} else {
						break
					}
					if nodeMap[ix][iy].nodeType != "dEnd" && nodeMap[ix][iy].nodeType != "wall" && nodeMap[ix][iy].nodeType != "path" {
						nodeMap[x][y].adjacentNodes = append(nodeMap[x][y].adjacentNodes, &nodeMap[ix][iy])
						break
					}
				}
			}
		}
	}
}

func rstVisited() {
	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			nodeMap[x][y].visited = false
		}
	}
	for {
		if len(path) > 0 {
			path = path[:len(path)-1]
		} else {
			break
		}
	}
}

//DFS (Non Recursive)
func dfs(pxMap [][]bool) {
	var currNode *node
	var nextNode *node
	//find start
	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			if nodeMap[x][y].nodeType == "start" {
				currNode = &nodeMap[x][y]
				break
			}
		}
	}
	path = append(path, currNode)
	for {
		currNode.visited = true
		if currNode.nodeType == "end" {
			path = append(path, currNode)
			break
		}
		validPath := false
		for i := 0; i < len(currNode.adjacentNodes); i++ {
			if !currNode.adjacentNodes[i].visited {
				nextNode = currNode.adjacentNodes[i]
				validPath = true
				break
			}
		}
		if !validPath {
			//If all next nodes are visited go back 1 node and remove the last one from the path
			if len(path) > 0 {
				nextNode = path[len(path)-1]
				path = path[:len(path)-1]
			} else {
				break
			}
			currNode = nextNode
		} else {
			currNode = nextNode
			path = append(path, currNode)
		}

	}

}

//-------------
func getStart() *node {
	var startNode *node = nil
	for x := 0; x < sizeX; x++ {
		for y := 0; y < sizeY; y++ {
			if nodeMap[x][y].nodeType == "start" {
				startNode = &nodeMap[x][y]
				break
			}
		}
	}
	return startNode
}

//BFS (Recursive)
func delPath(target *node) {
	var temp []*node
	for i := 0; i < len(path); i++ {
		if !(path[i].pos[0] == target.pos[0] && path[i].pos[1] == target.pos[1]) {
			temp = append(temp, path[i])
		}
	}
	path = temp
}
func bfs(currNode *node) {
	fmt.Println(currNode.pos)
	if !currNode.visited {
		currNode.visited = true
		path = append(path, currNode)
	}
	hasPath := false
	if currNode.nodeType != "end" {
		for i := 0; i < len(currNode.adjacentNodes); i++ {
			if !currNode.adjacentNodes[i].visited {
				hasPath = true
				nextNode := currNode.adjacentNodes[i]
				bfs(nextNode)
			}
		}

		if !hasPath {
			delPath(currNode)
		}

	}
}

//--------------

//DFS (Recursive)
func recDfs(currNode *node) {
	if !currNode.visited {
		path = append(path, currNode)
		currNode.visited = true
	}
	hasPath := false
	var nextNode *node
	if currNode.nodeType != "end" {
		//The next node is the next available adjacent node
		for i := 0; i < len(currNode.adjacentNodes); i++ {
			if !currNode.adjacentNodes[i].visited {
				hasPath = true
				nextNode = currNode.adjacentNodes[i]
				break
			}
		}
		if hasPath { //If there is an available path go there
			recDfs(nextNode)
		} else { //Else go bact to the previous node and remove the current from the path
			if len(path) > 0 {
				nextNode = path[len(path)-1]
				path = path[:len(path)-1]
				recDfs(nextNode)
			}
		}
	} else { //If currNode = the end add it to the path
		path = append(path, currNode)
	}
}

func main() {
	var pxMap = initPxMap()
	initNodeMap(pxMap)
	initNodes(pxMap)
	initAdjMap(pxMap)
	fmt.Println("Nodes:", nodes)
	//RecDFS
	sTime := time.Now()
	recDfs(getStart())
	elapsed := time.Since(sTime)

	fmt.Println("Finished Recursive D.F.S in: " + elapsed.String())
	saveImage(nodeMap)
	rstVisited()

	//NonRecDFS
	sTime = time.Now()
	dfs(pxMap)
	elapsed = time.Since(sTime)
	fmt.Println("Finished Non-Recursive D.F.S in: " + elapsed.String())
	saveImage(nodeMap)
	rstVisited()

}
