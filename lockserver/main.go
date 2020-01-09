package main

import lock "zklock/lockserver"

import "fmt"

func main() {
	lockserver, err := lock.GetLockServer("redis")
	if err != nil {
		panic(err)
	}
	lockserver.Lock()
}

func num1() {
	for i := 0; i < 50; i++ {
		fmt.Println(i)
	}
}
func num2() {
	for i := 51; i < 100; i++ {
		fmt.Println(i)
	}
}
