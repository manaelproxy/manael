---
title: "Manael v3.2 リリース！"
date: 2026-08-23
description: "Manael v3.2 では、リクエスト情報を利用した後処理と、最終ペイロードに基づくレスポンスヘッダーのカスタマイズを追加しました。"
---

## 概要

Manael v3.2 では、Go ライブラリとして組み込むアプリケーション向けのポストプロセッサーフックを拡張しました。変換済み画像の後処理で受信リクエストを参照できるようになり、最終的な変換結果に基づいてレスポンスヘッダーをカスタマイズすることもできます。

## リクエスト情報を利用する後処理

新しい `WithRequestPostProcessor` は、変換済み画像のバイト列に加えて、受信した `*http.Request` を受け取ります。リクエスト URL や `Accept` ヘッダーを使ってキャッシュキーを生成したり、呼び出し元に応じて処理を変えたりする場合に便利です。

```go
proxy := manael.NewServeProxy(upstreamURL,
	manael.WithRequestPostProcessor(func(r *http.Request, data []byte) ([]byte, error) {
		cacheKey := r.URL.String() + ":" + r.Header.Get("Accept")
		_ = cacheKey // キャッシュなどのリクエスト単位の処理に使用します。
		return data, nil
	}),
)
```

`WithRequestPostProcessor` と従来の `WithPostProcessor` を両方設定した場合は、前者が優先されます。

## 最終ペイロードに基づくレスポンスヘッダー

`WithResponseHeaderProcessor` は、受信リクエスト、レスポンスヘッダー、最終的な変換済みバイト列を受け取るフックです。`ETag`、キャッシュポリシー、監査用ヘッダーなど、最終ペイロードから導出するメタデータを設定できます。このフックでペイロード自体を置き換えることはできません。

Manael は、最終バイト列を決定して `Content-Type`、`Content-Length`、`Content-Disposition` を設定した後に、このフックを呼び出します。

```go
proxy := manael.NewServeProxy(upstreamURL,
	manael.WithResponseHeaderProcessor(func(r *http.Request, header http.Header, data []byte) error {
		header.Set("ETag", `"`+strconv.FormatInt(int64(len(data)), 10)+`"`)
		header.Set("Cache-Control", "private, max-age=60")
		return nil
	}),
)
```

## 安全なフォールバック動作

これらのフックは、Manael が実際に画像を変換した場合にのみ実行されます。変換が発生しなければフックはスキップされます。いずれかのフックがエラーを返した場合は、Manael が失敗をログ出力し、加工途中の結果ではなく元のアップストリームレスポンスへ安全にフォールバックします。

## 詳細

ポストプロセッサーフックのドキュメントと、v3.2 の変更内容をあわせて確認してください。
