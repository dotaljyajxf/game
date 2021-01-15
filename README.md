# game

# 简介
一个TCPserver,可以接收客户端的消息中包含一个字符串描述的方法，类似“user.getInfo”，以及使用protobuf传递参数，实现获取用户信息的一种RPC调用。目前的server是一个单机服务，并没有做分布式。
