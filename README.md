zenbun
====

簡易全文検索エンジン

雑なので速度とか精度とかアテにしないでほしい。

Usage
----
```go
db := zenbun.New()

db.Index("neko", "吾輩は猫である")
db.Index("ningen", "人間失格")

candidates := db.Find("猫")
for _, c := range candidates {
    fmt.Printf("%s: %.4f\n", c.DocumentName, c.Score)
}
```

Installation
----

```bash
go get github.com/asparlose/zenbun
```

License
----

MIT

ひとこと
----

英語できないので誰か助けて.
