# genuuid

UUIDを生成する。


```sh
$ ./genuuid -h
Usage: genuuid [--gen-version NUM] [--ulid] [--xid] [--win] [--lower] [--upper] [--output OUT] [--v2-domain DOMAIN] [--v2-id ID] [--namespace UUID] [--data DATA] [
--version] [--debug] [NUM [NUM ...]]

Positional arguments:
  NUM                    生成個数。 [default: [1]]

Options:
  --gen-version NUM, -v NUM
                         生成するUUIDのバージョンを指定する。[1-7] [default: 4]
  --ulid                 ULIDを生成する。生成順でソートできるUUIDのようなもの。
  --xid                  xidを生成する。生成順でソートでき、かつURLエンコードが不要なUUIDのようなもの。
  --win, -w              Windowsのレジストリで使用される形式で出力する。{}で囲まれる。
  --lower                小文字で出力する。
  --upper                大文字で出力する。
  --output OUT, -o OUT   ファイル出力する。
  --v2-domain DOMAIN     (v2) 0,1,2のいずれかを指定する。(0:Person、1:Group、2:Org) [default: -1]
  --v2-id ID             (v2) ドメイン内でのID。PersonはUID、GroupはGIDである必要があります。Orgまたは非POSIXではidの意味は、サイトによって異なります。
  --namespace UUID       (v3/v5) 名前空間。UUID文字列を指定する。
  --data DATA            (v3/v5) データ。
  --version              このプログラムのバージョン情報を出力する。
  --debug, -d            デバッグ用。ログが詳細になる。
  --help, -h             ヘルプを出力する。
```

## 使い方

単に呼ぶと1つ生成する。

```sh
$ ./genguid.exe
73fb8ff8-7234-4e8b-9076-794fea4d8489
```

数値を与えるとその個数生成する。

```sh
$ ./genguid.exe 10
0b6b6f82-dd3d-41ee-a148-32534ba8b121
1234a65a-8b1b-42aa-ac10-e5f97f136c63
4b323d0e-433d-43b1-994d-fd7ccb108e99
f7a8b58e-2c2a-4d3a-8c0b-f403e211db9f
71934380-86e7-4f70-b6b3-0e0732536f45
30b3006d-c8ff-4ed2-9e7d-392ca1dc4df5
cff90b6d-889a-407c-b8a9-5fb1265d5a60
a10d0071-0a08-4643-b583-795f1d2176ac
3b48c35b-8b10-41a0-ae10-7ded891c1a2c
84dfc24e-4b93-47dc-a53e-4bb46ad6bb65
```

`-w`か`--win`をつけると、windowsのレジストリ形式で出力する。

```sh
$ ./genguid.exe -w
{7472910e-aaa2-48ea-9638-a1f44b18be2e}
```

`--upper`や`--lower`で大文字小文字を制御できる。

```sh
$ ./genguid.exe --upper
B1750630-3E89-4098-866D-F074C8CFD17D
$ ./genguid.exe --lower
1a10fa6e-5b88-45b3-b448-94260e701570
```

その他の機能として
[ULID](https://ja.wikipedia.org/wiki/UUID#ULID) と
[xid](https://github.com/rs/xid) の生成機能がある。

```sh
$ ./genguid.exe --ulid
01JEN25R8TJ9S2KZN7W103WWFJ
$ ./genguid.exe --xid
ctb9058aslujrt6u2c30
```

## 生成するUUIDのバージョン

`-v` で生成するUUIDのバージョンを1-7の値で指定できる。  
以下は3個ずつ生成させた例。  

注意点としてバージョン2,3,5は専用の引数で追加の情報を与える必要がある。  
詳細は <https://ja.wikipedia.org/wiki/UUID> 等を参照。  
下記はこのプログラムで与える必要がある情報についての仕様。
- v2
	- ドメイン: 0から2の数値。0:Person, 1:Group, 2:Orgを表す。
	- ID: ドメイン内でのID。ユーザーIDやグループIDに相当する。
- v3/v5
	- 名前空間: ハイフン付きの36文字UUID。
	- データ: 任意文字列。

```sh
$ ./genuuid 3 -v 1
dac53184-9cdb-11f0-b3a3-00155d303366
dac53390-9cdb-11f0-b3a3-00155d303366
dac5339a-9cdb-11f0-b3a3-00155d303366
$ ./genuuid 3 -v 2 --v2-domain 0 --v2-id 1
00000001-9cdc-21f0-8000-00155d303366
00000001-9cdc-21f0-8000-00155d303366
00000001-9cdc-21f0-8000-00155d303366
$ ./genuuid 3 -v 3 --namespace "d78946b5-bcc7-4dc5-9dcd-6843241d53cc" --data "hoge"
49484a0b-3559-3cb2-bb6d-e2d8b62f612e
49484a0b-3559-3cb2-bb6d-e2d8b62f612e
49484a0b-3559-3cb2-bb6d-e2d8b62f612e
$ ./genuuid 3 -v 4
b66f0dde-935b-474f-8ccc-b6a481841f63
23fef6d8-6e45-4e07-9921-51e84310dc9f
58a084a5-ff56-4f60-b06d-75daa445ea81
$ ./genuuid 3 -v 5 --namespace "d78946b5-bcc7-4dc5-9dcd-6843241d53cc" --data "hoge"
f135c5dc-ff74-5392-b32f-dcce26e9573e
f135c5dc-ff74-5392-b32f-dcce26e9573e
f135c5dc-ff74-5392-b32f-dcce26e9573e
$ ./genuuid 3 -v 6
01f09cdb-bba7-662b-8ef9-00155d303366
01f09cdb-bba7-6828-8ef9-00155d303366
01f09cdb-bba7-6831-8ef9-00155d303366
$ ./genuuid 3 -v 7
0199934b-3b18-7526-a97d-257336cb79bd
0199934b-3b18-7535-b6b8-487f8d51d284
0199934b-3b18-7549-ad6e-ebef9b873046
```
