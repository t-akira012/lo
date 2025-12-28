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
lo -s           # シンプル出力（awk/fzf向け）
lo --simple     # 同上
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

### シンプル（-s）

```
"README.md"	# Project Title
"notes.txt"	Meeting notes
```

### パイプライン例

```bash
lo -s | fzf | awk -F'\t' '{print $1}' | xargs vim
```
