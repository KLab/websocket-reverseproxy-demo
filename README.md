WebSocket の Dynamic Load Balancer
==================================

ステートレスな HTTP だと負荷分散ってアプリの事考えないでよかったけど、
WebSocket アプリを負荷分散させる場合はアプリの事考えたいこともあるよね.

リバースプロクシが redis からプロキシ先を動的に取得するようにしよう!

アプリのセットアップ
----------------------
複数のルームを持つ WebSocket chat を tornado で作りました。

```
$ git clone https://github.com/KLab/tornado-weboscket-sample.git
$ cd tornado-websocket-sample
$ pip install tornado
$ ./chatdemo.py --port=8888
```

これで http://127.0.0.1:8888/ で chat デモが動きます.
もう一アプリを立ち上げましょう

```
$ ./chatdemo.py --port=8889
```

http://127.0.0.1:8888/chat/room1 と http://127.0.0.1:8889/chat/room1
が別々に動いてることを確認してください.

redis
------

redis を適当にインストールして、ローカルホストで動かしてください。
リバースプロクシ先を手動で設定しておきます.

```
$ redis-cli
> set /chatsocket/room1 127.0.0.1:8888
ok
> set /chatsocket/room2 127.0.0.1:8889
ok
```

nginx + lua
---------------

WebSocket 対応が入ったバージョンの nginx を nginx-lua-module 入りでビルドしてください。
[lua-resty-redis](https://github.com/agentzh/lua-resty-redis) もインストールして、
lua から `require "resty.redis"` できるようにしておいてください。

```
$ nginx -p . -c nginx.conf
```

8080 番ポートで nginx が動き始めます.

http://127.0.0.1:8080/chat/room1 が http://127.0.0.1:8888/chat/room1 と連携し、
http://127.0.0.1:8080/chat/room2 が http://127.0.0.1:8889/chat/room2 と連携しています。

golang
--------
go の net/http の reverseproxy は WebSocket に対応していません。
fork して対応したものが [rproxy](https://github.com/methane/rproxy) です。

go の開発環境(注: hg tip で開発してるので go 1.0.3 では動かないかもしれません)を用意してください。

```
$ go get github.com/methane/rproxy
$ go get github.com/garyburd/redigo/redis
$ go run gorouter.go
```

http://127.0.0.1:8090/ で reverse proxy が動いています。あとは先程の nginx+lua と同じです.
