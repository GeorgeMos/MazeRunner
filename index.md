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
BFS is also a very simple algorithm. It traverses the tree scanning each time all the next available nodes from the current one. To do this it utilises a Queue(FIFO). The benefit of bfs is that it will always find the shortest path to the end node BUT it is a lot slower compared to dfs if there are more than one available paths. The problem with bfs is that when the algorithm finishes there is no path created, the only output is the end node. To solve this issue i use the `parrent` pointer in the `node` object. For each node visited the algorithm stores the pointer of the previous node as the parrent of the current one, this way a reverse path is created!
Because it is very easy to locate the end node we can create a path starting from the end and ending at the start (confusing...). Finaly we reverse it and the end result is the shortest route to the exit. Sadly the problems dont end here. For large Non-Perfect (more then 1 paths) mazes, the Non-Recursive implementation reached up to 7.5GB memmory usage before i stoped it! And the recursive type simply overflows. This indicates that the queue size is geting extremely large, and even if we had a lot of memmory the execution time is unpractical.
#### Non-Recursive
```go
func bfs(startNode *node) *node {
	var bfsQueue []*node
	bfsQueue = append(bfsQueue, startNode)
	currNode := startNode
	var endNode *node
	for currNode.nodeType != "end" {
		currNode.visited = true
		for i := 0; i < len(currNode.adjacentNodes); i++ {
			if !currNode.adjacentNodes[i].visited {
				bfsQueue = append(bfsQueue, currNode.adjacentNodes[i])
				currNode.adjacentNodes[i].parrent = currNode
			}
		}
		currNode = bfsQueue[1]
		bfsQueue = bfsQueue[1:]

		if currNode.nodeType == "end" {
			endNode = currNode
			break
		}
	}
	return endNode
}
```

#### Recursive
```go
func recbfs(currNode *node) *node {
	currNode.visited = true
	if len(bfsQueue) == 0 {
		bfsQueue = append(bfsQueue, currNode)
	}
	for i := 0; i < len(currNode.adjacentNodes); i++ {
		if !currNode.adjacentNodes[i].visited {
			bfsQueue = append(bfsQueue, currNode.adjacentNodes[i])
			currNode.adjacentNodes[i].parrent = currNode
		}
	}
	var nextNode *node
	if len(bfsQueue) > 1 {
		nextNode = bfsQueue[1]
		bfsQueue = bfsQueue[1:]
	} else {
		nextNode = nil
	}
	if currNode.nodeType != "end" {
		return recbfs(nextNode)
	} else {
		return currNode
	}
}
```
