package main

import(
	"os"
	"rgo/rgo/http"
	"rgo/rgo/viper"
	"rgo/rgo/websocket"
	"rgo/rgo/database"
	"rgo/rgo/images"
	"log"
)

func getServerAddr()string{
	config := viper.NewConfig("application","json")
	host := config.GetString("server.host")
	port := config.GetString("server.port")
	return host+":"+port
}

func initMysqlEngine(){
	err := database.NewMysqlEngine()
	if err != nil{
		log.Println("mysql init error:" + err.Error())
	}
}

type UserModel struct{
	ID uint `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name" json:"name"`
	CreateTime int `gorm:"autoCreateTime;column:createtime" json:"create_time"`
	UpdateTime int `gorm:"autoUpdateTime;column:updatetime" json:"update_time"`
	DeleteTime int `gorm:"autoDeleteTime;column:deletetime" json:"delete_time"`
}

func(user *UserModel)TableName()string{
	return "users"
}

func main(){
	
	//图片
	canvas := images.NewCanvas()
	canvas.AddImageElement("1.jpg",true)
	canvas.ToGray("11.jpg")
	//设置临时运行环境
	os.Setenv("RGO_RUNMOD","dev")
	//获取配置服务地址
	addr := getServerAddr()
	//初始化数据库连接池
	initMysqlEngine()
	//设置并启动http服务
	r := http.Default()
	r.StaticFile("/favicon.ico","./public/static/images/favicon.ico")
	r.GET("/",Index)
	user := r.Group("/user")
	user.Use(Middleware)
	{
		user.GET("/:name",User)
		user.GET("/ws",Ws)
	}
	r.Run(addr)
}

func Index(ctx *http.Context){
	user := new(UserModel)
	user.Name = "oldda"
	database.MysqlEngine.Create(user)
	database.MysqlEngine.Find(user)
	ctx.JSON(200,user)
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