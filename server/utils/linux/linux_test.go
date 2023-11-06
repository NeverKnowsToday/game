package linux

import (
	"container/list"
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestChannel1(t *testing.T) {
	var ch = make(chan int, 10)
	for i := 0; i < 10; i++ {
		//t.Log(i)
		select {
		case ch <- i:
		case v := <-ch:
			t.Log(v)
		}
	}
	// for v := range ch {
	// 	t.Log(v)
	// }
}

const (
	FullBackupSwitchOnName    = 1 << iota
	ArchiveBackupSwitchOnName = 1 << iota
)

type User struct {
	Name string `json:"name"` // `json:"xxx"`：在结构体和json字符串字段顺序不一致的情况下：unmarshal根据tag去寻找对应字段的内容
	Age  int    `json:"age"`
}

// func sortByAge(u []User) {
// 	sort.Slice(u, func(i, j int) bool { //
// 		return u[i].Age < u[j].Age
// 	})
// }
type Student struct {
	Name     string  //名称
	IsOnline int64   //是否在线
	Integral int64   //积分
	Gold     float64 //金币
}

type Students []*Student

type haCheckReceive struct {
	result    string
	startTime string
	endTime   string
}

type student struct {
	Name string
	Age  int
}

type TaskQueue struct {
	lock sync.Mutex

	taskList *list.List

	parallelNum       int
	controlChan       chan struct{}
	unfinishedTaskCnt int

	task TaskInterface
}

type TaskInterface interface {
	Run() error
}

type op1Object struct {
}

func Newop1Object(t *TaskQueue) *op1Object {
	return &op1Object{}
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		unfinishedTaskCnt: 0,
		taskList:          list.New(),
	}
}

func (o *op1Object) Run() error {
	// clusterID := q.Get().(int)

	clusterID := 1
	fmt.Println("work", clusterID)
	time.Sleep(1 * time.Second)
	fmt.Println("work end", clusterID)

	return nil
}

func (t *TaskQueue) Append(item interface{}) {
	t.taskList.PushBack(item)
	t.unfinishedTaskCnt += 1
}

func (t *TaskQueue) Init(parallelNum int) {
	t.parallelNum = parallelNum
	t.controlChan = make(chan struct{}, parallelNum)
	// 注册任务
	t.task = Newop1Object(t)

	for i := 0; i < t.parallelNum; i++ {
		t.controlChan <- struct{}{}
	}
}

func (t *TaskQueue) RunTask() {
	{
		if t.Add() {
			defer func() {
				if err := recover(); err != nil {
					stack := make([]byte, 1<<20)
					stack = stack[:runtime.Stack(stack, false)]
				}
			}()
			//t.task.Run()
			t.Done()
		}
	}
}

func (t *TaskQueue) Add() bool {
	_, ok := <-t.controlChan
	return ok
}

func (t *TaskQueue) Get() interface{} {
	t.lock.Lock()
	item := t.taskList.Front()
	t.taskList.Remove(item)
	t.lock.Unlock()

	return item.Value
}

func (t *TaskQueue) Done() {
	t.controlChan <- struct{}{}

	t.unfinishedTaskCnt = t.unfinishedTaskCnt - 1
	if t.unfinishedTaskCnt == 0 {
		close(t.controlChan)
	}
}

func TestControlTask(t *testing.T) {
	q := NewTaskQueue()
	clusterID := []int{1, 2, 3, 4, 5, 6, 7, 8}
	for _, id := range clusterID {
		q.Append(id)
	}
	q.Init(2)
	t.Log("start...")
	for i := 0; i < q.unfinishedTaskCnt; i++ {
		go q.RunTask()
	}

	time.Sleep(10 * time.Second)
}

func TestExecTimeout5(t *testing.T) {
	intZeroSetList := make([]string, 0)
	for _, i := range intZeroSetList {
		t.Log(i)
	}
	t.Log("f")
	// l := list.New()
	// clusterID := []int{3, 2, 3, 4, 5, 6, 7, 8}
	// for _, id := range clusterID {
	// 	l.PushBack(id)
	// 	t.Log(l.Front().Value.(int))
	// }

	//t.Log(l.Remove(l.Front()))

	// go func() {
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			stack := make([]byte, 1<<20)
	// 			stack = stack[:runtime.Stack(stack, false)]
	// 			t.Logf("go routine panic, err: %s, stack: %s", err, stack)
	// 		}
	// 	}()

	// 	err := op()
	// 	if err != nil {
	// 		t.Log(err.Error())
	// 	}
	// }()
	// time.Sleep(2 * time.Second)
	// ins := make([]string, 0)
	// ins = append(ins, "a", "b", "c", "d")

	// parallelNum := 2
	// var ch = make(chan struct{}, parallelNum)
	// wg := sync.WaitGroup{}
	// wg.Add(len(ins))

	// go func() {
	// 	for i := 0; i < parallelNum; i++ {
	// 		ch <- struct{}{}
	// 	}
	// }()

	// for _, slave := range ins {

	// 	go func(slave string) {
	// 		{
	// 			_, ok := <-ch
	// 			if ok {
	// 				t.Log("start")
	// 				time.Sleep(2 * time.Second)
	// 				ch <- struct{}{}
	// 				wg.Done()
	// 				t.Log("end")
	// 			}
	// 		}
	// 	}(slave)
	// }
	// wg.Wait()
	// close(ch)

	// m := make(map[string]*student)
	// stus := []student{
	// 	{Name: "zhou", Age: 24},
	// 	{Name: "li", Age: 23},
	// 	{Name: "wang", Age: 22},
	// }
	// for _, stu := range stus {
	// 	m[stu.Name] = &stu
	// }

	// for k, v := range m {
	// 	println(k, "=>", v.Name)
	// }
	// updateClusterID := ""
	// clusterIDSlice := []string{"1", "2", "1", "2", "3"}
	// _, clusterIDSliceInter := array.Unique(clusterIDSlice)

	// for _, clusterID := range clusterIDSliceInter.([]interface{}) {
	// 	updateClusterID += clusterID.(string) + ","
	// }
	// t.Log(updateClusterID)

	// terminationChan := make(chan struct{})
	// receiveChan := make(chan haCheckReceive, 1)

	// go func() {
	// 	go func() {
	// 		result := haCheckReceive{}
	// 		for {
	// 			select {
	// 			case <-terminationChan:
	// 				t.Log("termination ha work")
	// 				receiveChan <- result
	// 				return
	// 			default:
	// 				time.Sleep(1 * time.Second)
	// 				t.Log("subtask")
	// 				result.result += "subtask 111"
	// 				result.startTime = "1"
	// 				result.endTime = "2"
	// 			}
	// 		}
	// 	}()

	// 	t.Log("start ha work")
	// 	time.Sleep(5 * time.Second)
	// 	t.Log("finish ha work")
	// 	terminationChan <- struct{}{}
	// 	res := <-receiveChan
	// 	t.Log(res)
	// }()

	// time.Sleep(20 * time.Second)

	// // // diff := SortedSliceStringMaxDiffItem(a1, a2)

	// diff := []string{}
	// findFlag := -1
	// c := 1
	// for _, v := range base {
	// 	if v == currentLast {
	// 		findFlag = 1
	// 	}
	// 	if findFlag >= 0 {
	// 		if c > 1 {
	// 			diff = append(diff, v)
	// 		}
	// 		c++
	// 	}
	// }

	// t.Log(diff)
	// for i := 0; i <= 23; i++ {
	// 	t.Log(i/2 + 1)
	// }

	// var	 users []User
	// var u User
	// u.Name = "test1"
	// u.Age = 22
	// users = append(users, u)
	// u.Name = "test2"
	// u.Age = 20
	// users = append(users, u)
	// u.Name = "test3"
	// u.Age = 26
	// users = append(users, u)

	// sortByAge(users)
	// t.Log(users)
	// a := 1
	// t.Log(FullBackupSwitchOnName + ArchiveBackupSwitchOnName)
	// t.Log(a & FullBackupSwitchOnName)
	// t.Log(a & ArchiveBackupSwitchOnName)
	// a := "/app/opengauss/brmsoft"
	// aa := strings.Split(a, "/")

	// t.Log(aa)
	// res := ""
	// for c, v := range aa {
	// 	//t.Log(v)
	// 	if c+1 != len(aa) {
	// 		res += v + "/"
	// 	}
	// }
	// t.Log(res)
	// s := "/app/appuser/dbpaas/console/package/expansion/opengauss/og_expansion/"
	// sList := strings.Split(s, "/")
	// t.Log(sList[0:3])
	//t.Log(strings.Trim(sList[4], ")"))

	// currentHour := time.Now().Hour()
	// windowStartHour := 15
	// windowStart, err := time.Parse("15:05", "01:00")
	// if err == nil {
	// 	windowStartHour = windowStart.Hour()
	// }
	// //t.Log(windowStartHour)
	// windowEndHour := 15
	// windowEnd, err := time.Parse("15:05", "15")
	// if err == nil {
	// 	windowEndHour = windowEnd.Hour()
	// }

	// if currentHour >= windowEndHour || currentHour <= windowStartHour {
	// 	t.Log("no")
	// }
	// t.Log("f")

}

// func concatStr(s ...string) string {
// 	r := ""
// 	for _, v := range s {
// 		r += v
// 	}
// 	return r
// }

// type testObject struct {
// 	Agent string `json:"agent"`
// }

func TestBuildSSHTrust(t *testing.T) {
	err := BuildSSHTrust(context.Background(), true, "zhouxin28", "zhouxin28", "10.133.83.13", "10.133.83.13")
	if err != nil {
		t.Fatal(err.Error())
	}
}

type CommonInterface interface {
	func1(string)
}

type CommonObject struct{}

func (c *CommonObject) func1(message string) {
	print(message)
}

type Task struct {
	commonInterface CommonInterface
}

func Test01(t *testing.T) {
	common := &CommonObject{}

	task := &Task{}
	task.commonInterface = common

}
