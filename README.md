# V2bX-IpRecord
V2bX的外部记录器，用于同步和记录各节点的在线设备（IP）

PS：本人随手写出来自用的半成品，不保证任何可用性，有能力的建议参考源码自行实现。

## 配置文件说明
- `Addr` 监听地址，例：`127.0.0.1:1231`
- `Token` 与V2bX通信时的认证令牌
- `IpDb` IP数据库文件路径
- `MasterId` Telegram用户Id，用于推送历史连接Ip超限通知
- `BotToken` Telegram BotToken
- `HistoryIpLimit` 每日连接Ip数量限制，同城市算作一个Ip，超出将推送tg消息通知。为0不启用
- `OnlineIpLimit` 在线Ip数量限制，为0不启用
