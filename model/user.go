// Package model 構造体定義
package model

import "time"

type PostKajilabpayqrRequest struct {
	Barcode *string `json:"barcode"`
}

type PostKajilabpayqrResponse struct {
	Barcode          string `json:"barcode"`
	Name             string `json:"name"`
	BalanceQRPayload string `json:"balance_qr_payload"`
}

type UsersGetResponseUser struct {
	Name             string    `json:"name"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Barcode          string    `json:"barcode"`
	BalanceQRPayload string    `json:"balance_qr_payload"`
	Debt             int64     `json:"debt"`
	TotalPay         int64     `json:"total_pay"`
}

type UsersGetResponse struct {
	Users      []UsersGetResponseUser `json:"users"`
	TotalCount int64                  `json:"total_count"`
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
