// Package model 構造体定義
package model

type PostKajilabpayqrRequest struct {
	Barcode string `json:"barcode"`
}

type PostKajilabpayqrResponse struct {
	Barcode          string `json:"barcode"`
	Name             string `json:"name"`
	BalanceQRPayload string `json:"balance_qr_payload"`
}
