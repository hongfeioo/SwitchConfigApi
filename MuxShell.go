package main

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

//MuxShell  执行命令和输出结果的引擎，核心的执行动作函数
func MuxShell(w io.Writer, r, e io.Reader) (chan<- string, <-chan string) {
	in := make(chan string, 3)
	out := make(chan string, 5)
	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself ,登陆命令是第一条队列任务

	go func() {
		for cmd := range in { // 循环输入命令
			wg.Add(1) // 添加队列任务
			w.Write([]byte(cmd + "\n"))
			fmt.Printf("%s\n", cmd)
			// //由于netconf客户端回复hello的时候，server端没有任何回应，所以直接不用等待out
			// if strings.Contains(cmd, "</hello>]]>]]>") {
			// 	fmt.Println("contains")
			// 	out <- string("ccont")
			// 	wg.Done()
			// }
			wg.Wait() // 等待队列任务完成,执行完一条命令后，再执行下一条命令
		}
	}()

	go func() { // 进入死循环，一直读取数据；
		var (
			buf [128 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])

			if err != nil {
				fmt.Println(err.Error())
				close(in)
				close(out)
				return
			}
			t += n

			result := string(buf[:t])

			if strings.Contains(result, "- More -") { //  遇到超常输出，需要加空格;  这种空格不会影响数据的完整性，亲测dis arp; dis mac-address ;
				w.Write([]byte(" "))
			}

			if strings.Contains(result, "name:") || strings.Contains(result, "word:") ||
				strings.Contains(result, ">") || strings.Contains(result, "]") || strings.Contains(result, "]]>]]>") {
				out <- string(buf[:t])
				t = 0

				wg.Done() //读到头一次，完成一个队列任务；
			}
		}
	}()
	return in, out
}
