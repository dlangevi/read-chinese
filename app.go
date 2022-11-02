package main

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/jmoiron/sqlx"

	"read-chinese/backend"
)

//go:embed src-node/build/read-chinese.node
var program []byte

func startNode() *exec.Cmd {

	userConfigDir := backend.ConfigDir()
	userProgram := path.Join(userConfigDir, "read-chinese.node")
	err := os.WriteFile(userProgram, program, 0777)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(userProgram, userConfigDir)
	pipe, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func(p io.ReadCloser) {
		reader := bufio.NewReader(pipe)
		line, err := reader.ReadString('\n')
		for err == nil {
			fmt.Println(line)
			line, err = reader.ReadString('\n')
		}
	}(pipe)

	return cmd

}

// App struct
type App struct {
	ctx context.Context
	cmd *exec.Cmd
	db  *sqlx.DB
}

// NewApp creates a new App application struct
// NewApp 创建一个新的 App 应用程序
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
// startup 在应用程序启动时调用
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	// 在这里执行初始化设置
	a.ctx = ctx

	// While doing dev I will just run the node server myself,
	// No need for launching a bundle
	// a.cmd = startNode()

}

func (a *App) NodeIpc(function string, argsJson string) string {
	postBody, _ := json.Marshal(map[string]string{
		"function": function,
		"args":     argsJson,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:3451/ipc", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	return sb
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}

// domReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (a *App) domReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
// beforeClose在单击窗口关闭按钮或调用runtime.Quit即将退出应用程序时被调用.
// 返回 true 将导致应用程序继续，false 将继续正常关闭。
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
// 在应用程序终止时被调用
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
	// Kill it:
	// if err := a.cmd.Process.Kill(); err != nil {
	// 	log.Fatal("failed to kill process: ", err)
	// }
}