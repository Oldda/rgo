# rgo
写个简易的框架试试

我的第一个基于golang开发的web开发脚手架。
# http服务器
根据标准库net/http包封装基础框架。包含动态路由，路由分组，中间件，上下文，日志和异常处理等。

# websocket服务器
根据gorilla/websocket封装的脚手架，支持单聊，群聊，广播等。并封装为事件触发模式。如果您用过swoole或者workerman那么您将找到熟悉的感觉。

# viper
简单封装，根据viper13实现的配置处理。

# database
第一版本仅根据gorm封装了mysql引擎。后续会有mongodb，redis，elasticsearch等。

# images
图像处理包。 该包的实现类似一张工作桌，将图片和文字作为素材管理。
功能包括：图片合成，水印，裁剪，缩放，灰度。

# 待续...
