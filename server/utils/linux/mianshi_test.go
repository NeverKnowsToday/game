package linux

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

//1.TestMakeAndAppend make默认值和append
func TestMakeAndAppend(t *testing.T) {
	s := make([]int, 3)
	s = append(s, 1, 2, 3)
	t.Log(s)
}

//2.TestDeferOrder defer调用顺序
func TestDeferOrder(t *testing.T) {
	defer func() { t.Log("打印前") }()
	defer func() { t.Log("打印中") }()
	defer func() { t.Log("打印后") }()

	panic("触发异常")
}

//3.TestGoClose 闭包
func TestGoClose(t *testing.T) {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	// for i := 0; i < 10; i++ {
	// 	go func() {
	// 		fmt.Println("A: ", i)
	// 		wg.Done()
	// 	}()
	// }

	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("B: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

//组合继承
type People struct{}

func (p *People) ShowA() {
	fmt.Println("showA")
	p.ShowB()
}
func (p *People) ShowB() {
	fmt.Println("showB")
}

type Teacher struct {
	People
}

func (t *Teacher) ShowB() {
	fmt.Println("teacher showB")
}

//4.TestGoInherit 组合继承
func TestGoInherit(t *testing.T) {
	teacher := Teacher{}
	teacher.ShowA()
}

//5.TestSelectRandom select测试
func TestSelectRandom(t *testing.T) {
	runtime.GOMAXPROCS(1)

	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)

	int_chan <- 1
	string_chan <- "hello"

	select {
	case value := <-int_chan:
		t.Log(value)
	case value := <-string_chan:
		panic(value)
	}
}

//6.map使用，怎么修改合适？
type UserAges struct {
	ages map[string]int
	sync.Mutex
}

func (ua *UserAges) Add(name string, age int) {
	ua.Lock()
	defer ua.Unlock()
	ua.ages[name] = age
}

func (ua *UserAges) Get(name string) int {
	if age, ok := ua.ages[name]; ok {
		return age
	}
	return -1
}

//7.nil interface
type Animal interface {
	Eat()
}

type Dog struct{}

func (stu *Dog) Eat() {}

func newAnimal() Animal {
	var dog *Dog

	return dog
}

func TestNilInterface(t *testing.T) {
	if newAnimal() == nil {
		t.Log("AAAAAAAAA")
	} else {
		t.Log("BBBBBBBBB")
	}
}

//8.channel bug，怎么改？
func TestChannelBug1(t *testing.T) {
	ch := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "finish job"
	}()

	select {
	case result := <-ch:
		t.Log(result)
	case <-time.After(time.Second):
		return
	}
}
