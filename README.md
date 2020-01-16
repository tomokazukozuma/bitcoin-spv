# bitcoin spv

## はじめに

## 1章 Bitcoinの現状

## 2章 P2P通信
1. Bitcoinノードの役割と種類
2. SPVとは

## 3章 Bitcoinの受け取り
1. 秘密鍵と公開鍵（btcsuite/btcutil/btcd/btcecで生成）
2. Bitcoinアドレス（base58だけbtcsuite/btcutil/base58）https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
3. 秘密鍵と公開鍵のフォーマットとその意味（圧縮、非圧縮によるフォーマット）

## 4章 Bitocin Nodeとの通信
1. Bitcoin Protocol(https://bitcoin.org/ja/developer-reference#protocol-versions)
2. メッセージ構成
3. ハンドシェイク（Version/Verack）

## 5章 Bitcoinの残高
1. UTXO
2. Markle Tree
3. Bloom Filter
4. Filterload
5. Markle Block
6. Markle Path

## 6章 Bitcoinの送金
1. Bitcoin Script(P2PK/P2PKH/P2SH)
2. UTXOの構築
3. 署名処理
4. 手数料計算（satoshi/byte）

## 7章 segwitについて
1. トランザクション構成
2. segwitの場合の送金と手数料

OP_CHECKSIG https://en.bitcoin.it/wiki/OP_CHECKSIG
BIP143 Segwitの署名検証 https://github.com/bitcoin/bips/blob/master/bip-0143.mediawiki
