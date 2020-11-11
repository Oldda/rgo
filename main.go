package main

import(
	"os"
	"rgo/rgo/http"
	"rgo/rgo/viper"
	"rgo/rgo/websocket"
	"log"
)

func main(){
	os.Setenv("RGO_RUNMOD","dev")
	config := viper.NewConfig("application","json")
	host := config.GetString("server.host")
	port := config.GetString("server.port")
	r := http.Default()

	r.StaticFile("/favicon.ico","./public/static/images/favicon.ico")
	r.GET("/",Index)
	user := r.Group("/user")
	user.Use(Middleware)
	{
		user.GET("/:name",User)
		user.GET("/ws",Ws)
	}
	r.Run(host+":"+port)
}

func Index(ctx *http.Context){
	ctx.JSON(200,http.H{
		"name":"ll",
	})
}

func User(ctx *http.Context){
	name := ctx.Param("name")
	ctx.JSON(200,name)
}

func Middleware(ctx *http.Context){
	log.Println("这是中间件，触发了")
}

func Ws (ctx *http.Context){
	wss := websocket.NewWsServer()
	go wss.Run()
	client := websocket.NewWsConn(wss)

	client.OnOpen(ctx.Writer,ctx.Req,func(cli *websocket.WsConn){
	    //绑定用户和链接
	    log.Println("open...")
	})

	client.OnMessage(func(cli *websocket.WsConn,msg []byte){
	    log.Println("message")
	})

	client.OnClose(func(cli *websocket.WsConn,err error){
	    log.Println("close")
	})
}