# 卡片和标签接口

所有接口都在 `/api/v1` 下，需要登录后的 `jwt` cookie。

统一响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 卡片

### 创建卡片

`POST /api/v1/cards`

```json
{
  "content": "卡片内容",
  "tags": ["原始标签数组"],
  "tag_ids": [1, 2],
  "external_links": ["https://example.com"],
  "internal_links": ["card:12"]
}
```

### 查询卡片

`POST /api/v1/cards/query`

```json
{
  "content": "关键词",
  "tags": ["原始标签数组"],
  "tag_ids": [1],
  "page": 1,
  "page_size": 10
}
```

### 获取卡片详情

`GET /api/v1/cards/{id}`

### 更新卡片

`PUT /api/v1/cards/{id}`

只传需要修改的字段。`tag_ids` 会整体替换卡片的规范化标签关系。

```json
{
  "content": "新的卡片内容",
  "tags": ["原始标签数组"],
  "tag_ids": [2, 3],
  "external_links": ["https://example.com"],
  "internal_links": ["card:12"]
}
```

### 删除卡片

`DELETE /api/v1/cards/{id}`

删除卡片时会同步软删除它的 `card_tag` 关系。

## 标签

### 创建标签

`POST /api/v1/tags`

```json
{
  "name": "项目"
}
```

同一用户下标签名称唯一。

### 查询标签

`POST /api/v1/tags/query`

```json
{
  "name": "项",
  "page": 1,
  "page_size": 10
}
```

### 获取标签详情

`GET /api/v1/tags/{id}`

### 更新标签

`PUT /api/v1/tags/{id}`

```json
{
  "name": "新名称"
}
```

### 删除标签

`DELETE /api/v1/tags/{id}`

删除标签时会同步软删除它的 `card_tag` 关系。

## 卡片标签关系

### 查询某张卡片的标签

`GET /api/v1/cards/{id}/tags`

### 给卡片添加一个标签

`POST /api/v1/cards/{id}/tags`

```json
{
  "tag_id": 1
}
```

重复添加会返回已有关系，不会重复创建。

### 替换卡片的全部标签

`PUT /api/v1/cards/{id}/tags`

```json
{
  "tag_ids": [1, 2, 3]
}
```

传空数组会清空该卡片的全部标签关系。

### 删除卡片上的某个标签

`DELETE /api/v1/cards/{id}/tags/{tag_id}`

### 查询映射关系

`POST /api/v1/card-tags/query`

```json
{
  "card_id": 1,
  "tag_id": 2,
  "page": 1,
  "page_size": 10
}
```

`card_id` 和 `tag_id` 都可选；同时传时查询某个卡片和某个标签的关系。
