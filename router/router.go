package router

import (
	"go-rest-api/controller"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, tc controller.ITaskController) *echo.Echo {
	e := echo.New()
	// Echoに対し CORS（Cross-Origin Resource Sharing）を適応、設定を行う
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// クライアントからのrequestを許可するリスト
		AllowOrigins: []string{"thhp://localhost:3000", os.Getenv("FE_URL")},
		// 許可されるHTTPヘッダーのリスト
		AllowHeaders: []string{
			// HTTPリクエストが発信されたウェブページのドメイン
			echo.HeaderOrigin,
			// クエストまたはレスポンスのコンテンツタイプ（データの種類や形式）
			echo.HeaderContentType,
			// クライアントが受け入れるレスポンスのコンテンツタイプ(データの形式)を指定するヘッダー
			echo.HeaderAccept,
			// サーバーがクライアントから受け入れる追加のHTTPヘッダーのリストを指定するヘッダー
			echo.HeaderAccessControlAllowHeaders,
			// クロスサイトリクエストフォージェリ（CSRF）トークンを含むヘッダー
			echo.HeaderXCSRFToken},
			// 許可されるHTTPメソッド（HTTPリクエストの種類）のリスト
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE"},
		// クレデンシャル（認証情報）の送信を許可
		AllowCredentials: true,
	}))
	// Echoに対し CSRF（Cross-Site Request Forgery）を適応、保護の設定
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		// CSRFトークンが設定されるクッキーのパスを指定
		// "/"はアプリケーション全体で使用可能
		CookiePath: "/",
		// CSRFトークンが設定されるクッキーのドメインを指定
		// 指定されたドメインの下でのみクッキーが送信される
		CookieDomain: os.Getenv("API_DOMAIN"),
		// CSRFトークンを含むクッキーをJavaScriptからアクセスできないようにする
		CookieHTTPOnly: true,
		// CSRFトークンを含むクッキーのSameSite属性
		// すべてのリクエストにクッキーを送信することを許可
		CookieSameSite: http.SameSiteNoneMode,
		// CookieSameSite: http.SameSiteDefaultMode,
		// CookirMaxAge: 60,
	}))
	// ルートエンドポイントの設定
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)
	// EchoでJWT（JSON Web Token）認証を使用するための設定
	// /tasks グループの設定
	t := e.Group("/tasks")
	// JWT認証ミドルウェアをtに対して適応
	// ユーザーが認証済みであることを確認するためにトークンを使用
	t.Use(echojwt.WithConfig(echojwt.Config{
		// フィールド: JWT署名の鍵を指定
		// os.Getenv("SECRET") から環境変数を読み込み、JWTの署名鍵として使用
		SigningKey: []byte(os.Getenv("SECRET")),
		// フィールド: JWTトークンをどこで検索するかを指定
		// クッキー内でJWTトークンを検索
		TokenLookup: "cookie:token",
	}))
	// タスク関連のエンドポイント設定 controller.ITaskControllerで処理
	t.GET("", tc.GetAllTasks)
	t.GET("/:taskId", tc.GetTaskById)
	t.POST("", tc.CreateTask)
	t.PUT("/:taskId", tc.UpdateTask)
	t.DELETE("/:taskId", tc.DeleteTask)
	// リクエストのハンドリングとルーティングを担当するinstance
	return e
}