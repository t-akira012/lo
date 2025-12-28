# lo - List Oneline

`.md`, `.txt`, `.mkd` ファイルの1行目をテーブル表示

## Install

```bash
go install github.com/t-akira012/lo@latest
```

## Usage

```bash
lo              # テーブル表示
lo /path        # 指定ディレクトリ
lo -0           # NUL区切り（パイプライン向け）
lo --null       # 同上
```

## Output

### テーブル（デフォルト）

```
┌──────────────┬────────────────┐
│File Name     │1st line        │
├──────────────┼────────────────┤
│README.md     │# Project Title │
│notes.txt     │Meeting notes   │
└──────────────┴────────────────┘
```

### NUL区切り（-0）

```
README.md\t# Project Title\0notes.txt\tMeeting notes\0
```

## Pipeline

```bash
vim $(lo -0 | fzf --read0 -d'\t' --with-nth=2 | awk -F'\t' '{print $1}')
```
