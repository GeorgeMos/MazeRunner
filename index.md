# Welcome to MazeRunner

This project was inspired by [Computerphile's video](https://www.youtube.com/watch?v=rop0W4QDOUI&t=535s) on maze solving algorithms. I originaly planned and began to implement dfs in Python but soon realised that trying to solve a 1000x1000 maze in Python is painfully slow and gave up on the idea for a while. Recently thow i discovered GoLang! The combination of blazing fast run time, easy syntax and garbage collection made me try again.

## How it works 
At first the maze is loaded in to memmory as an array of RGBA colors. Then in linear time a function scans the color array and creates a boolean array where
false = wall and true = path.<br>
With the help of the following struct, again in linear time an array of nodes is created.
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

