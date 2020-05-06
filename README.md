## 一、网络编程

注：目前只学习TCP版本

## 传输层 tcp

服务端使用Listen方法监听一个端口，当客户端需要通信时，使用Dial方法连接服务端。

服务端监听使用`net.listen("tcp", "ip:port")`

客户端连接`net.dial("tcp", "ip:port")`

## 应用层 socks5

1. 客户端请求服务端，发送socks5请求，包含本地支持的认证方法列表
2. 服务端收到请求，从认证列表挑选服务端支持的认证方法，并告知客户端
3. 若是无密码认证，直接开始发送请求连接。若是账号/密码认证，将账号/密码发送给服务端。
4. 服务端收到账号/密码，返回认证结果。
5. 若认证成功，开始发送请求连接。

客户端请求认证列表 -> 返回需要的认证方法 -> 需要密码则先认证 -> 发送请求连接

mermaid
graph TD
   A --> B
   

### 1. 建立连接

客户端与服务端建立TCP连接，向服务端发送连接请求，具体如下：

| version | nmethods | methods |
| --- | --- | --- |
|1Byte | 1Byte | 1-255Byte 由nmethods指定 |

`version` 代表socks代理版本，目前绝大多数情况使用`sock5`，即`0x05`

`nmethods` 客户端支持的认证方法数量

`methods` 客户端支持的方法列表，每个方法占用`1Byte`

#### METHOD定义

- `0x00` 不需要认证（常用）
- `0x01` GSSAPI认证
- `0x02` 账号密码认证（常用）
- `0x03` - 0x7F IANA分配
- `0x80` - 0xFE 私有方法保留
- `0xFF` 无支持的认证方法

如建立socks5，不需要认证，只支持不需要认认证协议，则数据格式为：
`0x05 0x01 0x00`

既支持无密码，又支持账号密码认证，数据格式为：
`0x05 0x02 0x00 0x02`

随后，服务端收到客户端支持的认证列表后，从中挑选一个支持的方法返回给客户端，

version | method
--- | --- 
1Byte | 1Byte

#### 1. 无需认证
 服务端返回：`0x05 0x00`，随后直接开始解包请求，代理客户端的流量。

#### 2. 密码认证 0x05 0x02
服务端返回：`0x05 0x02`，随后客户端按照以下格式发送账号密码：

VERSION| USERNAME_LENGTH | USERNAME | PASSWORD_LENGTH | PASSWORD
--- | --- | --- | --- | ---
1Byte|1Byte|1-255Byte|1Byte|1-255Byte

    version: 固定0x05
    user_length: 用户名长度
    username: 用户名
    password_length: 密码长度
    password: 密码

举例子：若账号是`0xFF`，密码是`0x0A 0x0B 0x0C`，则发送以下数据给服务端：
`0x05 0x01 0xFF 0x03 0x0A 0x0B 0x0C`。然后服务端收到并响应认证结果。

version | status
--- | ---
1Byte | 1Byte

    version: 固定0x05
    status: 0x00认证成功，其余认证失败（可由开发者自由定制）

服务端返回认证成功数据：0x05 0x00

随后客户端会发送连接命令给服务器，代理服务器会连接目标服务器，并返回结果。

### 客户端请求格式


VERSION|COMMAND|RSV|ADDRESS_TYPE|DST.ADDR|DST.PORT
---|---|---|---|---|---
1Byte|1Byte|1Byte|1Byte|1-255Byte|2Byte

    version: 固定0x05
    command: 
        0x01 CONNECT 连接上游服务器
        0x02 BND 绑定，被动模式
        0x03 UDP ACCOCIATE UDP中继
    RSV: 保留字段
    ADRESS_TYPE:
        0x01 IPv4
        0x03 域名
        0x04 IPv6

目前常用Connect模式：代理服务器收到请求后，代理访问该请请求，并将结果返回给客户端，完成整个连接过程。

#### golang实现socks5代理
见 [完整代码](./Proxy/Socks5/socks5.go)



#### 其它

浏览器访问socks5时，将源数据发给socks5服务端，服务端能够正常读取域名，此时实质上是http连接。

使用`curl --socks5 localhost:1080 -i baidu.com` 时，
发现并不能读取域名，实质上是curl将http请求在本地解析完成，
先解析域名服务器，端口，然后发送给代理服务器的是ip：port形式，
仅仅将纯tcp连接通过socks5传输。