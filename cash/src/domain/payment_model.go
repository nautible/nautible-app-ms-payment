package domain

type Payment struct {
	RequestId   string // 要求を一意に決めるキー（リクエストより取得）
	PaymentNo   string // 支払い処理を一意に決めるキー（本アプリ内で採番）
	AcceptNo    string // クレジット会社の受付キー（ダミーなので本アプリ内で採番）
	ReceiptDate string // クレジット会社の受付日時（ダミーなので本アプリ内で設定）
	OrderNo     string // 受注番号（リクエストより取得）
	OrderDate   string // 受注日時（リクエストより取得）
	CustomerId  int32  // 顧客番号（リクエストより取得）
	TotalPrice  int32  // 購入金額（リクエストより取得）
	OrderStatus string // 受注ステータス（本アプリ内で採番）
	DeleteFlag  bool   // 削除フラグ（false:有効 / true 削除）
}
