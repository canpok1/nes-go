# ROMS

テスト用のROMを`test/roms`に同梱しています。

* hello-world
    * `Hello World!`を表示する
    * 入手先: [NES研究室](http://hp.vector.co.jp/authors/VA042397/nes/sample.html)

その他のROMは、asmファイルからコンパイルすること

## asmファイルからの作成方法

nesasmでコンパイルすれば作れます。

詳しくは[ｷﾞｺ猫でもわかるファミコンプログラミング](https://github.com/thentenaar/nesasm)を参照。

ただしダウンロードできるnesasmはWindows用なので、Macの場合などは[ここ](https://github.com/thentenaar/nesasm)のソースコードからビルドすればOK。
sourceディレクトリ直下でmakeすればビルドできます。