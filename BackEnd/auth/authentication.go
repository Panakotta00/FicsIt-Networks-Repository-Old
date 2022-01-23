package auth

import (
	"FINRepository/Database"
	"FINRepository/Util"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const jwkGoogleURL = "https://www.googleapis.com/oauth2/v3/certs"
const clientIDGoogle = "808175066760-mgsgdt1o4f8n2c84uh5pge7dik9iovvk.apps.googleusercontent.com"

var JWTSecret = []byte("aklsdfjklasdjflkasdjfklajsdlkfjaweriojcnjqwoiuarvjnokijaernjvkjlasdnhjf")

type TokenClaims struct {
	jwt.StandardClaims
	Username string      `json:"username"`
	EMail    string      `json:"email"`
	ID       Database.ID `json:"id"`
}

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	finish := func(ctx echo.Context, user *Database.User) error {
		newCtx := context.WithValue(ctx.Request().Context(), "auth", user)
		ctx.SetRequest(ctx.Request().WithContext(newCtx))
		return next(ctx)
	}

	return func(ctx echo.Context) error {

		tokenCookie, err := ctx.Cookie("token")

		if err != nil || tokenCookie == nil {
			return finish(ctx, nil)
		}

		var tokenClaims = TokenClaims{}
		_, err = jwt.ParseWithClaims(tokenCookie.Value, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})
		if err != nil {
			ctx.SetCookie(&http.Cookie{Name: "token", MaxAge: -1})
			return echo.NewHTTPError(http.StatusForbidden, "failed to verify token")
		}

		var user Database.User
		if err := Util.DBFromContext(ctx.Request().Context()).Find(&user, tokenClaims.ID).Error; err != nil {
			ctx.SetCookie(&http.Cookie{Name: "token", RawExpires: "-1"})
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get authenticated user")
		}

		return finish(ctx, &user)
	}
}

func AuthenticateUser(c echo.Context, email string, username string) (string, error) {
	ctx := c.Request().Context()
	var user Database.User
	query := Util.DBFromContext(ctx).Where("user_email = ?", email).Find(&user)
	if err := query.Error; err != nil || query.RowsAffected < 1 {
		// no existing account found, try to create one
		user = Database.User{ID: Database.ID(Util.GetSnowflakeFromCTX(ctx).Generate().Int64()), Name: username, EMail: email}
		if err = Util.DBFromContext(ctx).Create(&user).Error; err != nil {
			log.Printf("Error Auth User: %e", err)
			return "", echo.NewHTTPError(http.StatusInternalServerError, "Unable to authenticate user or create new user account")
		}
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &TokenClaims{
		Username: user.Name,
		EMail:    user.EMail,
		ID:       user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		log.Println("Error Token Gen: %e", err)
		return "", echo.NewHTTPError(http.StatusInternalServerError, "Unable to create authentication token")
	}

	c.SetCookie(&http.Cookie{Name: "token", Value: tokenString})
	return tokenString, nil
}

func OAuth2Request(ctx echo.Context) error {
	redirect, rErr := url.Parse(ctx.QueryParam("redirect"))
	failureRedirect, fErr := url.Parse(ctx.QueryParam("failure-redirect"))
	RespondError := func(code int, err string) error {
		r := failureRedirect
		if fErr != nil {
			r = redirect
		} else if rErr == nil {
			q := r.Query()
			q.Add("redirect", redirect.String())
			r.RawQuery = q.Encode()
		}
		if r != nil {
			q := r.Query()
			q.Add("failure", err)
			r.RawQuery = q.Encode()
			return ctx.Redirect(http.StatusTemporaryRedirect, r.String())
		} else {
			return echo.NewHTTPError(code, err)
		}
	}

	credential := ctx.FormValue("credential")
	if credential == "" {
		return RespondError(http.StatusBadRequest, "no credential given")
	}

	csrfCookie, err := ctx.Cookie("g_csrf_token")
	if err != nil || csrfCookie == nil || csrfCookie.Value == "" {
		return RespondError(http.StatusBadRequest, "no csrf token in cookie")
	}
	csrfBody := ctx.FormValue("g_csrf_token")
	if csrfBody == "" {
		return RespondError(http.StatusBadRequest, "no csrf token in post body")
	}
	if csrfCookie.Value != csrfBody {
		return RespondError(http.StatusBadRequest, "failed to verify double submit cookie")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(credential, claims, func(t *jwt.Token) (interface{}, error) {
		keyID, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("expecting JWT header to have string kid")
		}

		// TODO: Cache
		var response, err = http.Get(jwkGoogleURL)
		if err != nil {
			return nil, errors.New("unable to get google-auth keys")
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, errors.New("unable to read google-auth keys")
		}
		keys, err := jwk.Parse(body)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse google-auth keys: %v", err)
		}

		if key, ok := keys.LookupKeyID(keyID); ok {
			var raw interface{}
			return raw, key.Raw(&raw)
		}

		return nil, fmt.Errorf("unable to find key '%q'", keyID)
	})
	if err != nil {
		log.Printf("Error auth: %v", err)
		return RespondError(http.StatusBadRequest, "unable to verify token")
	}

	var claimMap = token.Claims.(jwt.MapClaims)
	if claimMap["aud"] != clientIDGoogle {
		return RespondError(http.StatusBadRequest, "no valid client id in token")
	}

	if !(claimMap["iss"] == "accounts.google.com" || claimMap["iss"] == "https://accounts.google.com") {
		return RespondError(http.StatusBadRequest, "no accepted iss in token")
	}

	tokenString, err := AuthenticateUser(ctx, claimMap["email"].(string), claimMap["name"].(string))
	if err != nil {
		return RespondError(http.StatusBadRequest, err.Error())
	}

	if redirect != nil {
		q := redirect.Query()
		q.Del("failure")
		redirect.RawQuery = q.Encode()
		return ctx.Redirect(http.StatusTemporaryRedirect, redirect.String())
	} else {
		response := struct {
			token string
		}{
			token: tokenString,
		}
		return ctx.JSON(http.StatusOK, response)
	}
}
