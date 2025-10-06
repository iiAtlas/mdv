# Markdown Demo

Welcome to **MDV** - a terminal-based markdown viewer!

## Text Formatting

You can format text in various ways:

- **Bold text** using `**double asterisks**`
- *Italic text* using `*single asterisks*`
- ~~Strikethrough text~~ using `~~tildes~~`
- `Inline code` using backticks

## Headers

# Header 1
## Header 2
### Header 3
#### Header 4
##### Header 5
###### Header 6

## Lists

### Unordered Lists

- Item 1
- Item 2
  - Nested item 2.1
  - Nested item 2.2
- Item 3

### Ordered Lists

1. First item
2. Second item
3. Third item
   1. Nested item 3.1
   2. Nested item 3.2

### Task Lists

- [x] Completed task
- [ ] Pending task
- [ ] Another pending task

## Code Blocks

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, MDV!")
}
```

```python
def greet(name):
    return f"Hello, {name}!"

print(greet("World"))
```

```javascript
function factorial(n) {
    return n <= 1 ? 1 : n * factorial(n - 1);
}

console.log(factorial(5));
```

## Tables

| Feature | Description | Status |
|---------|-------------|--------|
| TUI Mode | Terminal interface | ✅ |
| GUI Mode | Wails-based GUI | ✅ |
| Themes | Dark/Light themes | ✅ |
| Watch Mode | Auto-reload on changes | ✅ |

## Blockquotes

> This is a blockquote.
> It can span multiple lines.
>
> > Nested blockquotes are also supported.

## Links and Images

- [MDV on GitHub](https://github.com/iiatlas/mdv)
- [Markdown Guide](https://www.markdownguide.org/)

## Horizontal Rules

---

Above and below this text are horizontal rules.

---

## Inline HTML

<div style="color: blue;">
HTML elements are supported too!
</div>

## Emoji Support

:rocket: :star: :heart: :sparkles:

(Note: Emoji rendering depends on your terminal)

---

**Tip**: Press `q` to quit, `r` to reload, and `g`/`G` to jump to top/bottom!
