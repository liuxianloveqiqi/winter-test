package main

//
//import (
//	"encoding/json"
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt"
//	"net/http"
//	"net/url"
//	"strings"
//	"time"
//)
//
//type User struct {
//	Username string `json:"username"`
//}
//
//func main() {
//	r := gin.Default()
//
//	// 配置客户端
//	clientID := "client-id"
//	clientSecret := "client-secret"
//	authServerURL := "https://auth-server-url"
//	redirectURI := "http://localhost:8080/callback"
//	scopes := []string{"openid", "profile", "email"}
//
//	// 构造认证请求
//	r.GET("/login", func(c *gin.Context) {
//		state := "state"
//		nonce := "nonce"
//		u, _ := url.Parse(authServerURL)
//		u.Path = "/oauth2/auth"
//		q := u.Query()
//		q.Set("client_id", clientID)
//		q.Set("response_type", "code")
//		q.Set("redirect_uri", redirectURI)
//		q.Set("scope", strings.Join(scopes, " "))
//		q.Set("state", state)
//		q.Set("nonce", nonce)
//		u.RawQuery = q.Encode()
//		c.Redirect(http.StatusTemporaryRedirect, u.String())
//	})
//
//	// 处理认证服务器的响应
//	r.GET("/callback", func(c *gin.Context) {
//		code := c.Query("code")
//		state := c.Query("state")
//		u, _ := url.Parse(authServerURL)
//		u.Path = "/oauth2/token"
//		q := u.Query()
//		q.Set("grant_type", "authorization_code")
//		q.Set("code", code)
//		q.Set("redirect_uri", redirectURI)
//		u.RawQuery = q.Encode()
//
//		req, _ := http.NewRequest(http.MethodPost, u.String(), nil)
//		req.SetBasicAuth(clientID, clientSecret)
//		resp, err := http.DefaultClient.Do(req)
//		if err != nil {
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//		defer resp.Body.Close()
//
//		if resp.StatusCode != http.StatusOK {
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//
//		var tokenResponsestruct {
//			AccessToken  string `json:"access_token"`
//			TokenType    string `json:"token_type"`
//			ExpiresIn    int    `json:"expires_in"`
//			RefreshToken string `json:"refresh_token"`
//		}
//		err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
//		if err != nil {
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//

//
//		// 获取用户信息
//		u, _ = url.Parse(userInfoURL)
//		req, _ = http.NewRequest(http.MethodGet, u.String(), nil)
//		req.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)
//		resp, err = http.DefaultClient.Do(req)
//		if err != nil {
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//		defer resp.Body.Close()
//
//		if resp.StatusCode != http.StatusOK {
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//
//		var userInfostruct {
//			Sub       string `json:"sub"`
//			Email     string `json:"email"`
//			FirstName string `json:"given_name"`
//			LastName  string `json:"family_name"`
//		}
//		err = json.NewDecoder(resp.Body).Decode(&userInfo)
//		if err != nil {
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//
//		// 验证用户是否已经注册过
//		// 如果是第一次登录，则进行用户注册流程
//		// 如果已经注册过，则直接登录
//

//		c.JSON(http.StatusOK, gin.H{
//			"message": "login success",
//			"user":    userInfo,
//		})
//	})
//}
