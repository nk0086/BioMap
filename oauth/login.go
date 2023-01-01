package oauth

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuthConfigはOAuth 2.0の設定を保持する構造体です
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// LoginHandlerはGoogleアカウントでのログインを行うハンドラーです
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// OAuth 2.0の設定を生成する
	oauthConfig := &oauth2.Config{
		ClientID:     "218183045866-9rblkruk7ievr17t3jk44svit0qug09b.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// 認証用のURLを生成する
	authURL := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)

	// 認証画面にリダイレクトする
	http.Redirect(w, r, authURL, http.StatusFound)
}

// CallbackHandlerはOAuth 2.0の認証コールバックを処理するハンドラーです
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// OAuth 2.0の設定を生成する
	oauthConfig := &oauth2.Config{
		ClientID:     "218183045866-9rblkruk7ievr17t3jk44svit0qug09b.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// クエリパラメータからコードを取得する
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Invalid code", http.StatusBadRequest)
		return
	}

	// コードを使用してトークンを取得する
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// トークンを使用してユーザー情報を取得する
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// ユーザー情報を表示する
	fmt.Fprintf(w, "User Info: %#v", resp)
}
