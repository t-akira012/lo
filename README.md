# lo - List Oneline

- `.md`, `.txt`, `.mkd` ファイルの1行目をテーブル表示

## Install

```bash
go install github.com/t-akira012/lo@latest
```

## Usage

```bash
lo          # カレントディレクトリ
lo /path    # 指定ディレクトリ
```

## Output

```
┌──────────────┬────────────────┐
│File Name     │1st line        │
├──────────────┼────────────────┤
│README.md     │# Project Title │
│notes.txt     │Meeting notes   │
└──────────────┴────────────────┘
```
