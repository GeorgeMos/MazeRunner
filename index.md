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
	weights []int //The length to go from the current node to each of the adjacent nodes
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
`initNodes(pxMap [][]bool)` assigns the `nodeType` in each node and `initAdjMap(pxMap [][]bool)` links all the nodes and creates a tree data structure. Now everything is ready for us to start solving!

## The algorithms
This is a WIP project and more algorithms will be added in the future so keep checking this section out for updates!
### Depth First Search(DFS)
DFS is pretty simple. Starting from the root of the tree (In this case the start of the maze) DFS travels as deep in to the tree as it can. If it comes to a dead end it backtracks to a node with an available path. To achive this DFS uses a Stack(LIFO). If the current node is marked as the "end" then the loop stops. The array `path` is used as the stack and it stores pointers to the nodes we have to follow in order to solve the maze. I implemented DFS with a while loop (Non-Recursive) and with Recursion in order to compare.
#### Non-Recursive:
```go
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
```
#### Recursive:
```go
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
```
### Breadth First Search(BFS)
BFS is also a very simple algorithm.
