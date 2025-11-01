# BOOTH VRChat Clothes Finder (Go)

A tiny Go CLI that fetches the latest or popular VRChat clothing items from BOOTH search results and prints them as JSON.

> Note: Please respect BOOTH's Terms of Service and avoid aggressive scraping. This tool fetches a single page per run and sets a browser-like User-Agent.

## Usage

- Default: popular items in Japanese locale, page 1
- You can tweak the search query, sort, page, and language.

### Build

```powershell
# from repo root
go build -o booth-vrc.exe ./
```

### Run

```powershell
# Popular (熱門) VRChat 衣装 on page 1 (Japanese locale)
./booth-vrc.exe --sort popular --page 1 --lang ja --query "vrchat 衣装"

# Latest (最新) results, Traditional Chinese locale
./booth-vrc.exe --sort new --page 1 --lang zh-tw --query "vrchat 衣服"
```

### Flags

- `--query` (string): search keywords (default: `vrchat 衣装`)
- `--sort` (string): `popular` or `new` (default: `popular`)
- `--page` (int): page number >= 1 (default: `1`)
- `--lang` (string): `ja` or `zh-tw` (default: `ja`)

### Output

Prints JSON array of items with fields:

```json
[
  {
    "title": "【しなの-Shinano-】Cat Hoodie",
    "url": "https://booth.pm/ja/items/7587937",
    "image": "https://.../path.jpg",
    "shop": "Micare Sewing",
    "price": "¥ 1,200"
  }
]
```

## Implementation notes

- Builds the BOOTH search URL like: `https://booth.pm/{lang}/search/{query}?sort={sort}&order=desc&page={page}`
- Parses products by locating links to `/items/{id}` and scraping nearby title, shop, price and image.
- Heuristics are used to be resilient to minor DOM changes; still, site updates may require selector adjustments.

## License

MIT
