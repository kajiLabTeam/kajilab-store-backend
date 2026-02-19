package barcodeutil

func MaskBarcode(barcode string) string {
	if len(barcode) <= 6 {
		return barcode // 短すぎる場合はそのまま返す
	}

	head := barcode[:3]
	tail := barcode[len(barcode)-3:]
	return head + "*******" + tail
}
