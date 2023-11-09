package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	// HTTPリクエストとレスポンスの情報にアクセスするためのメソッドやフィールドを提供
)

// requestを処理、対応した操作を実行

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

func (uc *userController) SignUp(c echo.Context) error {
	user := model.User{}
	// HTTPリクエストのボディデータを読み取る
	// そのデータを&userの構造体のフィールドにセットする
	if err := c.Bind(&user); err != nil {
		// HTTP 400 Bad Request、エラーメッセージを含んだJSONレスポンスをクライアントに返す
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// userRes = model.UserResponse{ID:newUser.ID,Email: newUser.Email,}
	// エラー時 userRes = model.UserResponse{} と エラー
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// tokenString = 署名つきJWTトークン string
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		// HTTPステータスコード 500（Internal Server Error）
		// errに格納されているエラーメッセージをJSON形式でクライアントに返却
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// http.Cookie 構造体の新しいインスタンスをメモリに割り当てポインタを返す
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	// クッキーを同じオリジンおよびクロスオリジンのリクエストに送信できる
	// サードパーティの認証情報やセッションクッキーを使用する際に必要
	cookie.SameSite = http.SameSiteNoneMode
	// クライアントにクッキーを送信
	// セッション管理、認証情報、ユーザ設定などの情報をクライアント側に保持
	c.SetCookie(cookie)
	// エラーなしで成功のステータスコードを持つレスポンスを返す
	return c.NoContent(http.StatusOK)
}

// クッキーに格納されたトークンを削除し、ユーザーをログアウト
func (uc *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	// tokenString を空白に
	cookie.Value = ""
	// cookieの有効期限を現在時刻に
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly  = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

func (uc *userController) CsrfToken(c echo.Context) error {
	// EchoコンテキストからCSRFトークンを取得し、文字列として使用
	token := c.Get("csrf").(string)
	// 200 (OK)と、CSRFトークンを含むJSONレスポンスを返す
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}