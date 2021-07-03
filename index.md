# Welcome to MazeRunner

This project was inspired by [Computerphile's video](https://www.youtube.com/watch?v=rop0W4QDOUI&t=535s) on maze solving algorithms. I originaly planned and began to implement dfs in Python but soon realised that trying to solve a 1000x1000 maze in Python is painfully slow and gave up on the idea for a while. Recently thow i discovered GoLang! The combination of blazing fast run time, easy syntax and garbage collection made me try again.

## How it works 
At first the maze is loaded in to memmory as an array of RGBA colors. Then in linear time a function scans the color array `colorMap` and creates a boolean array `pxMap` where
false = wall and true = path.<br>
The input png image must follow a set of rules.

- The maze must be black and white (Black for wall and white for path)
- Wall and path are 1 pixel
- Start is at the top and exit at the bottom

With the help of the following struct, again in linear time an array of nodes is created for each pixel of the maze.
```go
type node struct {
	nodeType      string
	adjacentNodes []*node
	openSideNum   int
	visited       bool
	pos           [2]int
	//Used in bfs
	parrent *node
}
```
The two main variables that we use are `nodeType` and `adjacentNodes`. The rest will be explained later. <br>
### Linking the nodes
`nodeType` can have the following values
- "wall"
- "path"
- "dEnd"
- "corner"
- "junction"

Instead of making an adjacency matrix which does not scale very well, `adjacentNodes` like a linked list, stores pointers to all the nodes that are next to the current node.
The following function assigns types to every node.
```go
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
```
