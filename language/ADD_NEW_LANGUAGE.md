## 新しい言語を追加する手順

1. languageディレクトリに新しい言語のソースファイルを作る。
2. Languageインターフェースを実装したstructを書く。
3. language.go の`GetLanguage`に新しい言語を追加する。
4. OnlineJudgeに言語を追加する(各オンラインジャッジのソースコードの`getLangID`に追加する)