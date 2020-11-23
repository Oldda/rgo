package websocket

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"rgo/rgo/viper"
)

var (
	wsConfig    = viper.NewConfig("ws", "json")
	heartBeat   = time.Duration(wsConfig.GetInt64("heartbeat_frequency"))
	msgSize     = wsConfig.GetInt64("max_message_size")
	readBuffer  = wsConfig.GetInt("read_buffer_size")
	writeBuffer = wsConfig.GetInt("write_buffer_size")
	origin, _   = wsConfig.Get("origin_allowed").([]string)
)

var (
	// 写入超时时间s
	writeWait = heartBeat / 2 * time.Second

	// 读取pong消息的超时时间
	pongWait = (heartBeat + 3) * time.Second

	// 发送ping的时间频率
	pingPeriod = heartBeat * time.Second

	// 消息大小设置
	maxMessageSize = msgSize
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readBuffer,
	WriteBufferSize: writeBuffer,
	CheckOrigin: func(r *http.Request) bool {
		if origin == nil || len(origin) < 1 {
			return true
		}
		for _, v := range origin {
			if r.URL.Path == v {
				return true
			}
		}
		return false
	},
}

type openHandler func(*WsConn)

type messageHandler func(*WsConn, []byte)

type closeHandler func(*WsConn, error)

// 接入的单条链接
type WsConn struct {
	svr *WsServer

	conn *websocket.Conn

	send chan []byte

	messageHandler messageHandler

	closeHandler closeHandler
}

//初始化一个client对象
func NewWsConn(svr *WsServer) *WsConn {
	return &WsConn{
		svr:  svr,
		send: make(chan []byte, 256),
	}
}

//设置OPEN回调函数
func (this *WsConn) OnOpen(w http.ResponseWriter, r *http.Request, handler openHandler) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	this.conn = conn
	this.svr.register <- this

	go this.writePump()
	go this.readPump()
	go handler(this) //调用回调函数
	return
}

//设置MESSAGE回调函数
func (this *WsConn) OnMessage(handler messageHandler) {
	this.messageHandler = handler
}

//设置CLOSE回调函数
func (this *WsConn) OnClose(handler closeHandler) {
	this.closeHandler = handler
}

//发送消息
func (this *WsConn) Send(msg []byte) {
	this.send <- msg
}

//读取客户端消息
func (this *WsConn) readPump() {
	defer func() {
		this.svr.unregister <- this
		this.conn.Close()
	}()
	this.conn.SetReadLimit(maxMessageSize)
	this.conn.SetReadDeadline(time.Now().Add(pongWait)) //now+30s
	this.conn.SetPongHandler(func(string) error { this.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		//ReadMessage()该操作会阻塞线程所以建议运行在其他协程上
		_, message, err := this.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				this.closeHandler(this, err)
			}
			log.Println(err)
			return
		}
		//客户端心跳
		if string(message) == "ping" {
			this.send <- []byte("pong")
			continue
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//this.svr.broadcast <- message
		this.messageHandler(this, message)
	}
}

//客户端发送消息
func (this *WsConn) writePump() {
	ticker := time.NewTicker(pingPeriod) //27s
	defer func() {
		ticker.Stop()
		this.conn.Close()
	}()
	for {
		select {
		case message, ok := <-this.send:
			this.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 发送一条close消息给对端
				this.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := this.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			//Add queued chat messages to the current websocket message.
			n := len(this.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-this.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C: //后端心跳
			this.conn.SetWriteDeadline(time.Now().Add(writeWait)) //now+10s
			if err := this.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
