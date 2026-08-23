---
title: ポストプロセッサーフック
weight: 55
---

Manael v3.1 では、画像変換の後段で独自処理を差し込める Go API フックが追加されました。Manael をライブラリとして組み込み、変換済みバイト列を送信前に確認、キャッシュ、置換、加工したい場合は `WithPostProcessor` を使用します。リクエスト単位の状態も必要な場合は `WithRequestPostProcessor` を使用します。

## フックが実行されるタイミング {#when-it-runs}

ポストプロセッサーフックは、Manael が画像変換に成功した後にのみ実行されます。

- 変換済み画像のバイト列を受け取ります。
- 何も変更せずそのまま返すこともできます。
- アプリケーションの都合に応じて別のバイト列に置き換えることもできます。
- アップストリームのレスポンスがそのまま通過した場合は呼び出されません。

## 基本例 {#basic-example}

```go
proxy := manael.NewServeProxy(upstreamURL,
	manael.WithPostProcessor(func(data []byte) ([]byte, error) {
		// 変換済みバイト列をキャッシュやオブジェクトストアへ保存します。
		return data, nil
	}),
)
```

## リクエスト対応の処理 {#request-aware-processing}

`WithRequestPostProcessor` は、変換済みバイト列とともにリクエストを受け取ります。リクエスト URL やヘッダーからキャッシュキーを作成する場合や、呼び出し元に応じた処理を行う場合に使用してください。

```go
proxy := manael.NewServeProxy(upstreamURL,
	manael.WithRequestPostProcessor(func(r *http.Request, data []byte) ([]byte, error) {
		cacheKey := r.URL.String() + ":" + r.Header.Get("Accept")
		_ = cacheKey // キャッシュなどのリクエスト単位の処理に使用します。
		return data, nil
	}),
)
```

両方のフックを設定した場合は `WithRequestPostProcessor` が優先され、`WithPostProcessor` は呼び出されません。

## 最終ペイロードのレスポンスヘッダー {#response-headers}

`WithResponseHeaderProcessor` は、受信したリクエスト、レスポンスヘッダー、および最終的な変換済みバイト列を受け取ります。`ETag`、キャッシュポリシー、監査用ヘッダーなど、最終ペイロードに依存するメタデータを追加する場合に使用します。このフックでペイロードを置き換えることはできません。

Manael は最終バイト列を決定し、それに対応する `Content-Type`、`Content-Length`、`Content-Disposition` を設定した後に、このフックを呼び出します。

```go
proxy := manael.NewServeProxy(upstreamURL,
	manael.WithResponseHeaderProcessor(func(r *http.Request, header http.Header, data []byte) error {
		header.Set("ETag", `"`+strconv.FormatInt(int64(len(data)), 10)+`"`)
		header.Set("Cache-Control", "private, max-age=60")
		return nil
	}),
)
```

## 主なユースケース {#use-cases}

- 変換済み画像を外部キャッシュへ保存する。
- フォーマット変換後にアプリケーション固有の後処理を行う。
- 最終的な変換結果を必要とする下流システムへ連携する。
- 最終的な変換済みペイロードから導出したレスポンスメタデータを設定する。

## エラーハンドリング {#error-handling}

フックがエラーを返した場合、Manael は失敗をログに記録し、加工途中の結果ではなく元のアップストリームレスポンスへフォールバックします。これにより、独自ロジックが失敗しても安全にリクエスト処理を継続できます。

## 挙動上の注意点 {#behavior-notes}

- このフックはコマンドラインオプションや環境変数ではなく、Go ライブラリ向け API です。
- フックはフォーマット変換後に実行されるため、最終的な変換済みペイロードを扱います。
- リクエストで変換が発生しなかった場合、フックは実行されません。
- リクエスト単位の状態が必要な場合は `WithRequestPostProcessor` を使用します。
- 最終ペイロードに依存する独自のレスポンスヘッダーには `WithResponseHeaderProcessor` を使用します。
