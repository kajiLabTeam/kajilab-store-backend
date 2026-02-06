package controller

import (
	"fmt"
	"net/http"

	"kajilab-store-backend/model"
	"kajilab-store-backend/service"
	"kajilab-store-backend/utils/qrutil"

	"github.com/gin-gonic/gin"
)

// バーコードからユーザ取得API
func GetUserByBarcode(c *gin.Context) {
	UserService := service.UserService{}
	barcode := c.Param("barcode")
	// ユーザ情報を取得
	user, err := UserService.GetUserByBarcode(barcode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get user by barcode")
		return
	}
	// レスポンスの型に変換
	userResponse := model.UserGetResponse{
		Id:      int64(user.ID),
		Name:    user.Name,
		Debt:    user.Debt,
		Barcode: user.Barcode,
	}

	c.JSON(http.StatusOK, userResponse)
}

func CreateUser(c *gin.Context) {
	UserService := service.UserService{}
	UserCreateRequest := model.UserCreateRequest{}
	err := c.Bind(&UserCreateRequest)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// バリデーション
	if UserCreateRequest.Barcode == nil || len(*UserCreateRequest.Barcode) != 13 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "need 13 characters barcode")
		return
	}
	if UserCreateRequest.Name == nil || *UserCreateRequest.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "need name")
		return
	}

	// 残高照会用QRコードのPayload生成
	// バーコードからQRペイロードを生成
	qrPayload := qrutil.HMACSHA256Short(*UserCreateRequest.Barcode)

	// リクエストの商品情報をデータベースの型へ変換
	user := model.User{
		Name:             *UserCreateRequest.Name,
		Debt:             0,
		Barcode:          *UserCreateRequest.Barcode,
		BalanceQrPayload: qrPayload,
	}

	err = UserService.CreateUser(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create user")
		return
	}

	// レスポンス用のデータ用意
	resUser := model.UserCreateResponse{
		Name:             *UserCreateRequest.Name,
		Barcode:          *UserCreateRequest.Barcode,
		BalanceQRPayload: qrPayload,
	}

	c.JSON(http.StatusOK, gin.H{
		"user": resUser,
	})
}

// ユーザの残高更新(主にチャージ)
func UpdateUserDebt(c *gin.Context) {
	UserService := service.UserService{}
	AssetService := service.AssetService{}
	UserUpdateDebtRequest := model.UserUpdateDebtRequest{}
	err := c.Bind(&UserUpdateDebtRequest)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// 変更前のユーザ
	beforeUser, err := UserService.GetUserById(UserUpdateDebtRequest.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get user")
		return
	}

	// ユーザ情報の更新
	user := model.User{
		Name:    beforeUser.Name,
		Debt:    beforeUser.Debt, // 変更しない
		Barcode: "",              // 変更しない
	}
	err = UserService.UpdateUser(UserUpdateDebtRequest.Id, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal update user")
		return
	}

	// 「変更後のユーザの残高　-　変更前のユーザの残高」を商店Debtに足す
	err = UserService.IncreaseKajilabpayDebt(int64(beforeUser.ID), -1, UserUpdateDebtRequest.Debt-beforeUser.Debt, "チャージ")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal decrease user debt")
		return
	}
	// 商店残高(money)を増やす
	err = AssetService.IncreaseMoney(UserUpdateDebtRequest.Debt - beforeUser.Debt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal decrease money")
		return
	}

	c.JSON(http.StatusOK, "success")
}

func UpdateUserBarcode(c *gin.Context) {
	UserService := service.UserService{}
	UserUpdateBarcodeRequest := model.UserUpdateBarcodeRequest{}
	err := c.Bind(&UserUpdateBarcodeRequest)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	user := model.User{
		Name:    "",                               // 変更しない
		Debt:    0,                                // 変更しない
		Barcode: UserUpdateBarcodeRequest.Barcode, // 変更しない
	}

	err = UserService.UpdateUser(UserUpdateBarcodeRequest.Id, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal update user")
		return
	}

	c.JSON(http.StatusOK, "success")
}

func CreateKajilabPayQR(c *gin.Context) {
	UserService := service.UserService{}
	KajilabPayQRRequest := model.PostKajilabpayqrRequest{}
	err := c.Bind(&KajilabPayQRRequest)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	if KajilabPayQRRequest.Barcode == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "need barcode")
		return
	}

	user, err := UserService.CreateKajilabPayQR(*KajilabPayQRRequest.Barcode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal update user")
		return
	}

	resUser := model.PostKajilabpayqrResponse{
		Barcode:          user.Barcode,
		Name:             user.Name,
		BalanceQRPayload: user.BalanceQrPayload,
	}

	c.JSON(http.StatusOK, resUser)
}
