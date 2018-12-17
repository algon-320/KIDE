# KIDE (Kyopro-Iikanjini-Dekiru-Environment)
*競プロ-いい感じに-出来る-環境*

----

## インストール
KIDEはGo言語で書かれているので、まずGo言語をインストールしてください。
その後、次のコマンドでダウンロード・ビルドを行います。
```sh
$ go get github.com/algon-320/KIDE
$ cd $GOPATH/src/github.com/algon-320/KIDE
$ go build
```
このコマンドを実行すると、ディレクトリに`KIDE`という実行ファイルが作成されるはずです。


## 仕様

### 対応しているオンラインジャッジ
- AtCoder
- Codeforces
- yukicoder
- AOJ


### 対応している言語
| 言語 | 指定するときの文字列 | ソースファイル拡張子 | デフォルトのコンパイルコマンド | デフォルトの実行コマンド |
|:----:|:----:|:----:|:----:|:----:|
| C++ | "C++" | ".cpp" | `g++ -std=c++11 -o a.out {SOURCEFILE_PATH}` | `./a.out` |
| Java | "Java" | ".java" | `javac {SOURCEFILE_PATH}` | `java Main` |
| Python2 | "Python2" | ".py" | 無し | `python {SOURCEFILE_PATH}` |
| Python3 | "Python3" | ".py" | 無し | `python {SOURCEFILE_PATH}` |

コンパイルコマンド・実行コマンドは`settings.json`で変更できる。
デフォルト言語も`setting.json`の`Language`->`DefaultLanguageName`で指定できる。

言語は比較的容易に追加できる。詳しくは`language/ADD_NEW_LANGUAGE.md`を参照。



### 対応しているエディタ(スニペット用)
- vscode


このプログラムは、実行ファイルの存在するディレクトリに各種ファイルを生成する。


## 主な機能
- `kide run`: コンパイル & 実行
- `kide dl {問題のURL}`: 問題のダウンロード
- `kide tester {問題id}`: テスト
- `kide submit {問題id}`: 提出
- `kide processer`: ソースコードを整形する（提出前にも適用される）
- `kide snippet`: スニペット管理
    - エディタ用のスニペット形式で出力
    - ライブラリ用のMarkdown出力

## サブコマンドの詳細
### `run`
カレントディレクトリにあるソースコードをコンパイル・実行する。

カレントディレクトリに指定した言語のソースファイルが1つしかない場合はそれをコンパイル・実行する。
複数ある場合はどれを実行するかの選択肢を表示する。
このソースファイルを決める仕組みは`tester`、`submit`でも同じ。

また、前回コンパイルした時とソースコードの内容が一致している場合、コンパイルはスキップされてそのまま実行される。(コンパイルの必要な言語)

コンパイルコマンド・実行コマンドは、settings.jsonで指定することが出来る。

設定項目が存在しない場合は基本的にデフォルトのものが使用される。

オプション
- `--language`、`-l`: コンパイル・実行したいソースコードの言語名を指定する(仕様の項目を参照)
    - 使える言語は language/language.go の`languageList`にあるもの
    - デフォルトは`setting.json`の`Language`->`DefaultLanguageName`で指定可能


### `dl {URL}`
指定した問題の情報とサンプル入出力をローカルに保存する。
AtCoder、Codeforces、yukicoder、AOJに対応しているが、正しく読み取れない問題もあるので注意。(ほとんど大丈夫なはず)

初めて使うときに、ログイン情報を要求されるので入力する。
(ユーザ名やパスワードを修正する場合は実行ファイルと同じディレクトリに作られる`settings.json`を変更する。)

ダウンロードするときに問題idが振られる。問題idは大文字小文字の区別なし。
- AtCoder、Codeforcesの場合はA問題なら"a"、B問題なら"b"といったようになる
- yukicoderは問題No.が問題idになる(No.001ならidは"001")
- AOJは問題のID(URLでid=XXXXのXXXX部分)が問題idとなる

AtCoder、Codeforces、yukicoderの場合、コンテストの問題一覧ページのURLを投げることで、一括して問題をダウンロードすることも出来る。

`view`で今保存されている問題一覧を表示出来る。引数で問題idを指定すると詳細を表示。

同じ問題IDの場合上書きされることに注意。
(例えば、あるコンテストのA問題をダウンロードして、別のコンテストのA問題をダウンロードすると
初めにダウンロードしたA問題は上書きされて消えてしまう。)

### `tester {問題id}`
指定された問題のサンプル入出力をテストする。`run`と同じようにコンパイルされた後に、自動でテストが行われる。
全て正解した場合は提出するか尋ねられ、そのまま提出できる。

オプション
`--case`、`-c`: 番号を指定すると特定のサンプルケースをテスト出来る


### `kide submit {問題id}`
指定した問題に対してソースコードを提出する。ジャッジ結果がACだった場合にソースコードを保存することも出来る。
（初回に保存するか尋ねられる。`settings.json`で変更可能。）


### `kide processer`
`settings.json`に`General->SourcecodeProcess->Command`で実行コマンドが設定されている場合にソースコードを整形することが出来る。設定しなければ、ソースコードがそのまま整形後のものとして扱われるため、意識する必要はない。

サブコマンド`kide processer`では、カレントディレクトリの対象ソースコードを整形し、その結果を出力する。
また、**`kide submit`や`kide tester`で提出する場合に、提出する直前にも整形が行われる**。

KIDEと実行コマンドとの間では標準入力と標準出力でソースコードをやり取りする。
標準入力から読み取ったソースコードを整形し、標準出力に出力するプログラムを作成し、実行コマンドとして登録することで、ソースコードを整形することが出来る。

また、コマンドの文字列の中の`{EXE_DIR}`はKIDEの実行ファイルのあるディレクトリのパスに置き換えられる。


### `cf-mysubmissions {コンテストid}`
Codeforcesのコンテストidを指定し、そのコンテストにおける自分の提出のジャッジ結果を表示する。

「In Queue」の提出が存在する場合、1分毎に確認し、ジャッジ結果が更新されていた場合その結果を表示する。


### `snippet`
対応しているエディタ用のスニペットを出力したり、ライブラリ用にMarkdownを出力する。
これらは標準出力に書き込まれるため、必要に応じてリダイレクションなどでファイルに出力する。

スニペットは項目ごとに1つのファイルを作る。形式は次の通りで、ファイル名は`{filename}.snip`にする必要がある。
```
<NAME> {名前}
<TRIGGER> {スニペットを発動させる文字列}
<TAG> {タグ}
<*NOTE>
複数行で説明などを記述できる。
Markdown出力する際にそのまま出力されるので、Markdown記法を用いることが出来る。

- pandocなどでHTMLやTeX形式に変換することが出来て便利そう
- `snippet`コマンド
- $F_i = F_{i-1} + F_{i-2}$

<*CODE>
// ここにコードを書く（複数行）
printf("Hello,World\n");
printf("sample snippet\n");
```

初回実行時にスニペットのある親ディレクトリを尋ねられるので絶対パスで入力する。
再帰的に`.snip`ファイルを検索するため、親ディレクトリ以下の構成については自由。

コマンドを実行するとどのエディタ用のスニペットを出力するか尋ねられる。
ここで、markdown出力を指定することができる。


## 使い方
KIDEをダウンロードコンパイルしてあり、`KIDE`という実行ファイルをパスの通ったディレクトリに配置してあるという
前提で、[AtCoder Practice Contest](https://beta.atcoder.jp/contests/practice)を例に説明します。

0. 適当なディレクトリに入る
1. ウェブブラウザでBeta版AtCoderの「A - はじめてのあっとこーだー（Welcome to AtCoder）」を開く [link](https://beta.atcoder.jp/contests/practice/tasks/practice_1)
2. 問題のURLをコピーする
3. `$ KIDE dl https://beta.atcoder.jp/contests/practice/tasks/practice_1`を実行する
4. AtCoderアカウントのユーザ名・パスワードを聞かれるので入力
5. サンプルケースが保存される。
6. ソースコードを作成する(例としてC++を想定)
    ここでディレクトリに`.cpp`のファイルが複数ある場合、`KIDE run`、`KIDE tester`、`KIDE submit`でコンパイルする対象を選択する画面が出る。
7. テストする
    - C++以外の言語を使う場合は`$ KIDE run --language Python`などと指定すること
    - `$ KIDE run`を実行すると先ほど書いたソースコードがコンパイルされて実行される（実行するだけなので、入力などは自分で書く）
    - `$ KIDE tester A --case 1`を実行するとサンプル1がテストされる
    - `$ KIDE tester A`を実行するとサンプル1、サンプル2がテストされ、どちらとも正解した場合はこのまま提出するかを尋ねられる
8. 提出する
    - `tester`コマンドで自動提出しない場合は`submit`コマンドを用いる
    - `$ submit A`で提出できる。この場合、サンプルのテストを行わずに提出するので注意。
9. ジャッジされるのを待つ
10. ジャッジ結果が表示される
    - ACした場合、ソースコードを保存するか尋ねられる。保存したい場合は保存するディレクトリを指定する。（`settings.json`で変更可能）
    - 保存を有効にするとこれ以降ACしたときに自動的に指定したディレクトリに保存される。
    - このとき、問題URL・提出URL・提出日時・ステータス（Accepted）がソースコードの先頭にコメントアウトされて追加されたものが保存される。
11. あとは精進するだけ！


## `settings.json`の例
```json
{
  "General": {
    "SaveSourceFileAfterAccepted": true,
    "SaveSourceFileDirectory": "{EXE_DIR}/ac_sources"
  },
  "Language": {
    "DefaultLanguageName": "C++",
    "C++": {
      "CompileCommand": "g++ -std=c++14 -O0 -g -o a.out {SOURCEFILE_PATH}",
      "RunningCommand": "./a.out"
    },
    "Java": {
      "CompileCommand": "javac {SOURCEFILE_PATH}",
      "RunningCommand": "java Main"
    },
    "Python": {
      "CompileCommand": "",
      "RunningCommand": "python {SOURCEFILE_PATH}"
    },
    "Python2": {
      "CompileCommand": "",
      "RunningCommand": "python {SOURCEFILE_PATH}"
    },
    "Python3": {
      "CompileCommand": "",
      "RunningCommand": "python {SOURCEFILE_PATH}"
    }
  },
  "OnlineJudge": {
    "AOJ": {
      "Handle": "aoj_handle",
      "Password": "********"
    },
    "AtCoder": {
      "Handle": "atcoder_handle",
      "Password": "********"
    },
    "Codeforces": {
      "Handle": "codeforces_handle",
      "Password": "********"
    },
    "yukicoder": {
      "Handle": "yukicoder_handle",
      "Password": "********"
    }
  },
  "snippet_manager": {
    "root_dir": "/home/ユーザ名/competitive_programming/libraries(snippets)"
  }
}
```