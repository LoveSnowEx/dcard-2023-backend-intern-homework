# page-list-service

這是一個使用 Go 語言實現的後端服務，提供了文章列表查詢和編輯的功能。

## 特色

本專案使用鏈結串列（Linked List）結構來管理文章列表。每個節點代表一篇文章（Page），使用者可以逐篇獲取文章，而不用一次載入整個列表。這種設計優化了查詢效率，能有效應對大量使用者的同時操作，從而增強系統的負載能力。而對於管理者，也提供了豐富的編輯操作，可以方便地進行列表的編輯，無需擔心分頁工程的問題。

> 專案想法參考了這篇[文章](https://medium.com/dcardlab/de07f45295f6)

## 介紹

- 包含 RESTful API、gRPC API 及 GUI。

### RESTful API

RESTful API 讓使用者可以透過 key 獲取列表及每一篇文章的資訊：

- `GET /head/<listkey>` 可獲得文章列表中第一篇文章的 key
- `GET /page/<pagekey>` 可獲得文章資訊及下一篇文章的 key

### GRPC API

gRPC API 提供管理者對列表進行編輯：

- New: 建立一個新的列表
- Delete: 刪除一個列表
- Begin: 獲取第一個元素的 Iterator
- End: 獲取最後一個元素之後的 Iterator
- Clear: 清空列表
- Insert: 於該 Iterator 之前插入新的元素
- Erase: 移除該 Iterator
- Set: 設定該 Iterator 儲存的 Page ID
- PushBack: 列表尾部新增一個元素
- PopBack: 移除列表尾部的元素
- PushFront: 列表首部新增一個元素
- PopFront: 移除列表首部的元素
- Clone: 複製一個列表

> TODO: 這部分可以使用 gRPC Bidirectional Streaming 進行優化，讓編輯操作的效率更高。

## GUI

可以用來測試 gRPC API，對文章列表進行編輯。

![grpcui](image/grpcui.png)

## 環境需求

- Go 1.20

## 如何使用

### 1. 初始化 env

請參考 `.env.example` 以對 `.env` 進行設定。

### 2. 啟動服務

執行以下指令以啟動所有服務：

```bash
docker-compose up -d
```

### 3. 服務端口

| service | port |
| --- | --- |
| RESTful API | 3000 |
| gRPC API | 50051 |
| GUI | 8080 |

## API 文件

### RESTful

#### GET `/head/<listkey>`

獲得列表的第一個元素的資訊，包含第一個元素的 key。

Response:

```json
{
  "nextPageKey": "xxx"
}
```

#### GET `/page/<pagekey>`

獲得文章的資訊，包含文章的標題、內容、網址 slug、是否發佈等等。

Response:

```json
{
  "article": {
    "title": "xxx",
    "content": "xxx",
    "slug": "xxx",
    "published": true
  },
  "nextPageKey": "xxx"
}
```

### gRPC

#### Message

```protobuf
message Empty {

}

message PageList {
  string key = 1;
}

message PageIterator {
  string key = 1;
  uint32 page_id = 2;
}

message DeleteRequest {
  string list_key = 1;
}

message BeginRequest {
  string list_key = 1;
}

message EndRequest {
  string list_key = 1;
}

message NextRequest {
  string iter_key = 1;
}

message PrevRequest {
  string iter_key = 1;
}

message ClearRequest {
  string list_key = 1;
}

message InsertRequest {
  string iter_key = 1;
  uint32 page_id = 2;
}

message EraseRequest {
  string iter_key = 1;
}

message SetRequest {
  string iter_key = 1;
  uint32 page_id = 2;
}

message PushRequest {
  string list_key = 1;
  uint32 page_id = 2;
}

message PopRequest {
  string list_key = 1;
}

message CloneRequest {
  string list_key = 1;
}
```

#### Service

##### New

建立一個新的列表。

```protobuf
rpc New(Empty) returns (PageList) {}
```

##### Delete

刪除一個列表。

```protobuf
rpc Delete(DeleteRequest) returns (Empty) {}
```

##### Begin

獲取第一個元素的 Iterator。

```protobuf
rpc Begin(BeginRequest) returns (PageIterator) {}
```

##### End

獲取最後一個元素之後的 Iterator。

```protobuf
rpc End(EndRequest) returns (PageIterator) {}
```

##### Next

獲取下一個元素的 Iterator。

```protobuf
rpc Next(NextRequest) returns (PageIterator) {}
```

##### Prev

獲取上一個元素的 Iterator。

```protobuf
rpc Prev(PrevRequest) returns (PageIterator) {}
```

##### Clear

清空列表。

```protobuf
  rpc Clear(ClearRequest) returns (Empty) {}
```

##### Insert

於該 Iterator 之前插入新的元素。

```protobuf
rpc Insert(InsertRequest) returns (PageIterator) {}
```

##### Erase

移除該 Iterator。

```protobuf
rpc Erase(EraseRequest) returns (PageIterator) {}
```

##### Set

設定該 Iterator 儲存的 Page ID。

```protobuf
rpc Set(SetRequest) returns (PageIterator) {}
```

##### PushBack

列表尾部新增一個元素。

```protobuf
rpc PushBack(PushRequest) returns (PageIterator) {}
```

##### PopBack

移除列表尾部的元素。

```protobuf
rpc PopBack(PopRequest) returns (Empty) {}
```

##### PushFront

列表首部新增一個元素。

```protobuf
rpc PushFront(PushRequest) returns (PageIterator) {}
```

##### PopFront

移除列表首部的元素。

```protobuf
rpc PopFront(PopRequest) returns (Empty) {}
```

##### Clone

複製一個列表。

```protobuf
rpc Clone(CloneRequest) returns (PageList) {}
```

## 資料庫

使用 PostgreSQL 作為資料庫，並透過 GORM 套件進行操作。

### 文章列表

```go
type PageList struct {
    gorm.Model
    Key uuid.UUID `gorm:"type:uuid;uniqueIndex"`
}

type PageNode struct {
    gorm.Model
    End     bool
    Key     uuid.UUID `gorm:"type:uuid;uniqueIndex"`
    PrevKey uuid.UUID `gorm:"type:uuid"`
    NextKey uuid.UUID `gorm:"type:uuid"`
    Page   *page.Page `gorm:"foreignkey:PageID;references:ID"`
    PageID  uint
    List    PageList  `gorm:"foreignkey:ListKey;references:Key"`
    ListKey uuid.UUID `gorm:"type:uuid;index"`
}
```

- `PageList`：文章列表，只記錄紀錄了該文章列表的 key。
- `PageNode`：文章列表節點，紀錄了文章列表中的每一篇文章，包含了上一篇文章的 key、下一篇文章的 key、所屬列表的 key 以及對應的文章 ID。

### 文章

```go
type Page struct {
    gorm.Model `json:"-"`
    Title      string `json:"title" gorm:"not null"`
    Content    string `json:"content" gorm:"type:text;not null"`
    Slug       string `json:"slug" gorm:"uniqueIndex;not null"`
    Published  bool   `json:"-" gorm:"default:false"`
}
```

`Page` 紀錄了文章的資訊，包含標題、內容、網址 slug 和發佈狀態。

## 測試

```bash
go test ./...
```

使用 `stretchr/testify` 和 `container/list` 套件來實現測試，透過比對 `list` 中的資料和資料庫中的資料是否相同，以確認操作是否正確。同時使用 in-memory sqlite 作為測試資料庫，以確保測試的可靠性。
