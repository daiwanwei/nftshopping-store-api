# nftshopping-store-api
## Table of Contents

 * [專案描述](#專案描述)
 * [Demo](#Demo)
 * [執行專案](#執行專案)
 * [Optimization](#Optimization)

## 專案描述

### 商店管理
1. 會員管理
2. 藝術品管理
3. 交易管理(To-Do)
4. 品牌管理
5. 庫存管理
6. 收藏管理

## Demo
1.swagger文檔:http://storeapi.daiwanwei.xyz/swagger/index.html

## 執行專案

### 執行測試

```bash
$ make test
```

#### 執行應用程式

```bash
#到專案目錄下
$ cd path_to_dir/nftshopping-store-api

# 下載第三方套件
$ go mod download

# 生成swagger文檔
$ swag init 

# 編譯專案(輸出到當前目錄下,檔案名為main)
$ go build -o main . 

# 執行應用程式
$ ./main 

# 確認專案是否執行
$ curl localhost:8080/probe
```
#### API文檔(swagger)
網址打入
```bash
#網址打入(default host=>localhost:8080)
http://{host}/swagger/index.html
```

## Optimization
- [ ] 將訂購和交貨改成event driven