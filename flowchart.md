# Main 函数流程图

这是 `main.go` 程序的执行流程图，展示了 Booth.pm 商品搜索爬虫的主要处理步骤。

```mermaid
flowchart TD
    Start([开始 main]) --> ParseFlags[解析命令行参数<br/>query, sort, page, lang]
    ParseFlags --> BuildURL[构建搜索 URL<br/>buildSearchURL]
    BuildURL --> CheckURLError{URL 构建<br/>是否成功?}
    CheckURLError -->|失败| ErrorExit1[记录错误并退出]
    CheckURLError -->|成功| FetchDoc[发送 HTTP 请求<br/>fetchDocument]
    FetchDoc --> CheckFetchError{请求是否<br/>成功?}
    CheckFetchError -->|失败| ErrorExit2[记录错误并退出]
    CheckFetchError -->|成功| ExtractItems[提取商品信息<br/>extractItems]
    ExtractItems --> FindLinks[查找所有商品链接<br/>匹配 /items/\d+]
    FindLinks --> ForEach[遍历每个链接]
    ForEach --> GetContainer[定位商品容器元素]
    GetContainer --> ExtractFields[提取字段<br/>Title, Image, Shop, Price]
    ExtractFields --> CheckDuplicate{URL 是否<br/>重复?}
    CheckDuplicate -->|是| ForEach
    CheckDuplicate -->|否| AddItem[添加到结果列表]
    AddItem --> MoreLinks{还有更多<br/>链接?}
    MoreLinks -->|是| ForEach
    MoreLinks -->|否| EncodeJSON[编码为 JSON 格式]
    EncodeJSON --> CheckEncodeError{编码是否<br/>成功?}
    CheckEncodeError -->|失败| ErrorExit3[记录错误并退出]
    CheckEncodeError -->|成功| OutputJSON[输出 JSON 到标准输出]
    OutputJSON --> End([结束])
    
    style Start fill:#e1f5e1
    style End fill:#e1f5e1
    style ErrorExit1 fill:#ffe1e1
    style ErrorExit2 fill:#ffe1e1
    style ErrorExit3 fill:#ffe1e1
```

## 主要流程说明

1. **参数解析** - 使用 `flag` 包解析命令行参数（查询关键字、排序方式、页码、语言）
2. **URL 构建** - 调用 `buildSearchURL()` 生成 Booth.pm 搜索 URL
3. **HTTP 请求** - 调用 `fetchDocument()` 获取搜索结果页面
4. **数据提取** - 调用 `extractItems()` 解析 HTML 并提取商品信息
5. **去重处理** - 使用 map 确保每个 URL 只出现一次
6. **JSON 输出** - 将结果编码为 JSON 格式并输出到标准输出

## 关键函数

- `buildSearchURL()` - 构建搜索 URL，支持语言、排序和分页
- `fetchDocument()` - 发送 HTTP GET 请求，返回 goquery 文档对象
- `extractItems()` - 从 HTML 中提取商品标题、链接、图片、店铺名和价格
- `absolute()` - 将相对 URL 转换为绝对 URL
