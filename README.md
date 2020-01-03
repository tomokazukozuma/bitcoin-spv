# bitcoin spv
## SPVとは

## Bitcoinの受け取り
1. 秘密鍵と公開鍵（btcsuite/btcutil/btcd/btcecで生成）
2. Bitcoinアドレス（base58だけbtcsuite/btcutil/base58）
3. 秘密鍵と公開鍵のフォーマットとその意味（圧縮、非圧縮によるフォーマット）

## Bitocin Nodeとの通信
1. メッセージ構成
2. ハンドシェイク（Version/Verack）

## Bitcoinの残高
1. UTXO
2. Markle Tree
3. Bloom Filter
4. Filterload
5. Markle Block
6. Markle Path

## Bitcoinの送金
1. Bitcoin Script(P2PK/P2PKH/P2SH)
2. UTXOの構築
3. 署名処理
4. 手数料計算（satoshi/byte）
5. segwitの場合の送金と手数料
