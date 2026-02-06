// Package model 構造体定義
package model

type PostKajilabpayqrRequest struct {
	Barcode *string `json:"barcode"`
}

type PostKajilabpayqrResponse struct {
	Barcode          string `json:"barcode"`
	Name             string `json:"name"`
	BalanceQRPayload string `json:"balance_qr_payload"`
}

type UserCreateRequest struct {
	Name    *string `json:"name"`
	Barcode *string `json:"barcode"`
}

type UserCreateResponse struct {
	Name             string `json:"name"`
	Barcode          string `json:"barcode"`
	BalanceQRPayload string `json:"balance_qr_payload"`
}
