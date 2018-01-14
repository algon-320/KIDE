## 新しい言語を追加する手順

1. languageディレクトリに新しい言語のソースファイルを作る。
2. Languageインターフェースを実装したstructを書く。
    - languageBase structを埋め込むと楽
        - name, compliceCommand, runningCommand, commentBegin, commentEndを決めるとある程度自動でやってくれる
        - 細かい動作を変更する際はそれぞれのメソッドを実装する
3. language.go の`languageList`に新しい言語を追加する。
4. OnlineJudgeに言語を追加する(各オンラインジャッジのソースコードの`getLangID`に追加する)