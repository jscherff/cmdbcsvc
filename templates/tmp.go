package main

import `time`

func main() {
	println(time.Now().Format(time.RFC3339))
}
