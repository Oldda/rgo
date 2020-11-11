package websocket

import(
	"errors"
)

type WsServer struct {

	clients map[*WsConn]bool //链接

	broadcast chan []byte //广播通道

	register chan *WsConn //注册通道

	unregister chan *WsConn //注销通道
}

func NewWsServer() *WsServer {
	return &WsServer{
		broadcast:  make(chan []byte),
		register:   make(chan *WsConn),
		unregister: make(chan *WsConn),
		clients:    make(map[*WsConn]bool),
	}
}

//获取链接总数
func(this *WsServer)Statistic()int{
	return len(this.clients)
}

//广播消息
func (this *WsServer)Broadcast(msg []byte){
	this.broadcast<-msg
}

//单发消息
func(this *WsServer)SendTo(cli *WsConn,msg []byte)error{
	if has,ok := this.clients[cli];has && ok{
		cli.Send(msg)
	}
	return errors.New("不存在的链接")
}

func (this *WsServer) Run() {
	for {
		select {
		case client := <-this.register: //注册
			this.clients[client] = true

		case client := <-this.unregister: //注销
			if _, ok := this.clients[client]; ok {
				delete(this.clients, client) //移除此条链接
				close(client.send)//关闭此条链接通道
			}

		case message := <-this.broadcast: //广播
			for client := range this.clients {
				select {
				case client.send <- message:
				default:
					delete(this.clients, client)
					close(client.send)
				}
			}
		}
	}
}
