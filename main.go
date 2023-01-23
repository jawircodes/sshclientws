package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	Code    int    `json:"code"`
	Details string `json:"details"`
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var interval time.Duration
var ORIGIN string
var PORT int
var HOST string

//var Don_update bool

type Host struct {
	Hostname string `form:"hostname" json:"hostname" binding:"required"`
	Port     string `form:"port,default=22" json:"port" `
	User     string `form:"user,default=root" json:"user"`
	Password string `form:"password" json:"password"` //base64 encoded password
	Cols     int    `form:"cols,default=120" json:"cols"`
	Rows     int    `form:"rows,default=32" json:"rows"`
}

func main() {
	Init()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	gin.SetMode(gin.ReleaseMode)
	e := gin.Default()
	e.Use(cors.Default())
	e.GET("/ws", WsSsh)
	e.Run(fmt.Sprintf("%s:%d", HOST, PORT))

}
func Init() {
	flag.DurationVar(&interval, "interval", 20*time.Second, "set ping pong frequency")
	flag.StringVar(&ORIGIN, "origin", "*", "set origins ,like \"http://127.0.0.1:8080,http://localhost:8080\"")
	flag.IntVar(&PORT, "port", 9018, "set port")
	flag.StringVar(&HOST, "host", "", "set bind host (default all)")
	/*var v = flag.Bool("V", false, "show version")
	if *v {
		showVersion()
	}*/
	flag.Parse()
	log.Printf("interval is %s", interval)
}

func wshandleError(ctx *gin.Context, err error) bool {
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Message{http.StatusInternalServerError,
			fmt.Sprintf("%s", err)})
	}
	return (err != nil)
}
func WsSsh(c *gin.Context) {
	var host Host
	err := c.ShouldBindQuery(&host)
	if wshandleError(c, err) {
		fmt.Println("Rrr", err.Error())
		return
	}

	client, err := NewSshClient(host.User, host.Hostname, host.Port, host.Password, nil)

	if wshandleError(c, err) {
		fmt.Println("ERRRR", err.Error())

		return
	}
	defer client.Close()

	ssConn, err := NewSshConn(host.Cols, host.Rows, client)
	if wshandleError(c, err) {
		fmt.Println("Rrr")
		return
	}
	defer ssConn.Close()
	// after configure, the WebSocket is ok.
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if wshandleError(c, err) {
		return
	}
	defer wsConn.Close()

	quitChan := make(chan bool, 3)

	var logBuff = new(bytes.Buffer)

	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(wsConn, logBuff, quitChan)
	go ssConn.SendComboOutput(wsConn, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan
}
