# Archimedes
## go 企业级开发手脚架

gin构建企业级脚手架，代码简洁易读，可快速进行高效web开发。
主要功能有：
1. 请求链路日志打印，涵盖mysql/redis/request
2. 支持多语言错误信息提示及自定义错误提示。
3. 支持了多配置环境
4. 封装了 log/redis/mysql/http.client 常用方法
5. 请求响应加解密
6. IP限制
7. jwt认证
8. 限流
9. token自动刷新
10. 参数校验
11. 微服务（待完成）
11. 容器化部署（待完成）

项目地址：https://github.com/fehu-asia/Archimedes
### 现在开始
- 安装软件依赖

```
git clone git@github.com:fehu-asia/Archimedes.git
cd fehu
go mod tidy
```
- 确保正确配置了 conf/mysql_map.toml、conf/redis_map.toml：

- 运行脚本

```
go run main.go
```
- 测试mysql与请求链路

创建测试表：
```
CREATE TABLE `area` (
 `id` bigint(20) NOT NULL AUTO_INCREMENT,
 `area_name` varchar(255) NOT NULL,
 `city_id` int(11) NOT NULL,
 `user_id` int(11) NOT NULL,
 `update_at` datetime NOT NULL,
 `create_at` datetime NOT NULL,
 `delete_at` datetime NOT NULL,
 PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='area';
INSERT INTO `area` (`id`, `area_name`, `city_id`, `user_id`, `update_at`, `create_at`, `delete_at`) VALUES (NULL, 'area_name', '1', '2', '2019-06-15 00:00:00', '2019-06-15 00:00:00', '2019-06-15 00:00:00');
```

```
curl 'http://127.0.0.1:8880/base/dao?id=1'
{
    "errno": 0,
    "errmsg": "",
    "data": "[{\"id\":1,\"area_name\":\"area_name\",\"city_id\":1,\"user_id\":2,\"update_at\":\"2019-06-15T00:00:00+08:00\",\"create_at\":\"2019-06-15T00:00:00+08:00\",\"delete_at\":\"2019-06-15T00:00:00+08:00\"}]",
    "trace_id": "c0a8fe445d05b9eeee780f9f5a8581b0"
}

查看链路日志（确认是不是一次请求查询，都带有相同trace_id）：
tail -f gin_scaffold.inf.log

```
- 测试参数绑定与多语言验证(目前不可与加解密共用)

```
curl 'http://127.0.0.1:8880/demo/bind?name=name&locale=zh'
{
    "errno": 500,
    "errmsg": "年龄为必填字段,密码为必填字段",
    "data": "",
    "trace_id": "c0a8fe445d05badae8c00f9fb62158b0"
}

curl 'http://127.0.0.1:8880/demo/bind?name=name&locale=en'
{
    "errno": 500,
    "errmsg": "Age is a required field,Passwd is a required field",
    "data": "",
    "trace_id": "c0a8fe445d05bb4cd3b00f9f3a768bb0"
}
```

### 文件分层
```
├── README.md
├── conf            配置文件夹
├── controller      控制器
│   └── demo.go
├── constant        常量
├── model             输入输出结构层
│   └── form         输出
│   └── param.go     请求参数
├── go.mod
├── go.sum
├── main.go         入口文件
├── middleware      中间件层
├── public          公共文件
└── router          路由层
│    ├── httpserver.go
│    └── route.go
├── util            工具包
```

