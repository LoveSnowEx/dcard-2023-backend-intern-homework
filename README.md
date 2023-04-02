# dcard-2023-backend-intern-homework

這是一個使用 Go 語言實現的後端應用程式，提供 REST API 和 gRPC API 給使用者進行操作。  
REST API 提供給一般使用者進行操作，而 gRPC API 提供給鏈結串列的編輯者進行操作。  
每個鏈結串列的元素紀錄了一個頁面（Page）的 ID。

## 設計

- 提供 REST API 給串列的使用者，讓使用者可以透過 key 獲取每一篇頁面的資訊
  - GET `/head/<listkey>` 可獲得 {nextPageKey: xxx}
  - GET `/page/<pagekey>` 可獲得 {article: {title: aaa, content: bbb, ...}, nextPageKey: xxx}
- 提供 gRPC API 給鏈結串列的編輯者進行操作
  - New: 建立一個新的串列
  - Delete: 刪除一個串列
  - Begin: 獲取第一個元素的 Iterator
  - End: 獲取最後一個元素之後的 Iterator
  - Clear: 清空串列
  - Insert: 於該 Iterator 之前插入新的元素
  - Erase: 移除該 Iterator
  - Set: 設定該 Iterator 儲存的 Page ID
  - PushBack: 串列尾部新增一個元素
  - PopBack: 移除串列尾部的元素
  - PushFront: 串列首部新增一個元素
  - PopFront: 移除串列首部的元素

## 環境需求

- Go 1.20

## 如何使用

### 執行資料庫

在執行伺服器之前，需要先執行資料庫，這裡使用的是 PostgreSQL，可以使用 Docker Compose 來執行。

```bash
docker-compose up -d
```

### 編譯伺服器

```bash
go build
```

### 執行伺服器

```bash
./dcard-2023-backend-intern-homework
```

- REST API 預設會在 port `:3000` 提供服務
- gRPC API 預設會在 port `:50051` 提供服務。
- 另外 gRPC 有 UI 可以使用，預設在 port `:8080`。

### 測試

```bash
go test ./...
```

### 自定義設定

可以透過修改 `.env` 檔案來自定義設定。`.env.example` 是範例檔案。

## REST API

### GET `/head/<listkey>`

獲得串列的第一個元素的資訊，包含第一個元素的 key。

Response:

```json
{
  "nextPageKey": "xxx"
}
```

### GET `/page/<pagekey>`

獲得頁面的資訊，包含頁面的標題、內容、網址 slug、是否發佈等等。

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

## gRPC API

### Message

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
```

### Service

#### New

建立一個新的串列。

```protobuf
rpc New(Empty) returns (PageList) {}
```

#### Delete

刪除一個串列。

```protobuf
rpc Delete(DeleteRequest) returns (Empty) {}
```

#### Begin

獲取第一個元素的 Iterator。

```protobuf
rpc Begin(BeginRequest) returns (PageIterator) {}
```

#### End

獲取最後一個元素之後的 Iterator。

```protobuf
rpc End(EndRequest) returns (PageIterator) {}
```

#### Next

獲取下一個元素的 Iterator。

```protobuf
rpc Next(NextRequest) returns (PageIterator) {}
```

#### Prev

獲取上一個元素的 Iterator。

```protobuf
rpc Prev(PrevRequest) returns (PageIterator) {}
```

#### Clear

清空串列。

```protobuf
  rpc Clear(ClearRequest) returns (Empty) {}
```

#### Insert

於該 Iterator 之前插入新的元素。

```protobuf
rpc Insert(InsertRequest) returns (PageIterator) {}
```

#### Erase

移除該 Iterator。

```protobuf
rpc Erase(EraseRequest) returns (PageIterator) {}
```

#### Set

設定該 Iterator 儲存的 Page ID。

```protobuf
rpc Set(SetRequest) returns (PageIterator) {}
```

#### PushBack

串列尾部新增一個元素。

```protobuf
rpc PushBack(PushRequest) returns (PageIterator) {}
```

#### PopBack

移除串列尾部的元素。

```protobuf
rpc PopBack(PopRequest) returns (Empty) {}
```

#### PushFront

串列首部新增一個元素。

```protobuf
rpc PushFront(PushRequest) returns (PageIterator) {}
```

#### PopFront

移除串列首部的元素。

```protobuf
rpc PopFront(PopRequest) returns (Empty) {}
```

## 資料庫設計

### 鏈結串列

鏈結串列的資料結構如下：

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
    // Page   *page.Page `gorm:"foreignkey:PageID;references:ID"`
    PageID  uint
    List    PageList  `gorm:"foreignkey:ListKey;references:Key"`
    ListKey uuid.UUID `gorm:"type:uuid"`
}
```

文章的資訊儲存在 `Page` 資料表中，而 `PageNode` 則是紀錄了鏈結串列的資訊，包含了頁面的 ID、頁面的前後頁面的 key、頁面所屬的串列的 key。

> 將 `Page` 註解掉的原因是因為 `Page` 相關的功能沒有實作，目前只有範例的 `Page` 資料表。

### 文章

文章的資料結構如下：

```go
type Page struct {
    gorm.Model `json:"-"`
    Title      string `json:"title" gorm:"not null"`
    Content    string `json:"content" gorm:"type:text;not null"`
    Slug       string `json:"slug" gorm:"uniqueIndex;not null"`
    Published  bool   `json:"-" gorm:"default:false"`
}
```

這邊只有簡單的紀錄了文章的標題、內容、網址 slug 和是否發佈。
