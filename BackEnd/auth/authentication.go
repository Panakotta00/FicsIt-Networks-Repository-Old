package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const jwkGoogleURL = "https://www.googleapis.com/oauth2/v3/certs"
const clientIDGoogle = "808175066760-mgsgdt1o4f8n2c84uh5pge7dik9iovvk.apps.googleusercontent.com"

type JWK struct {
	KID string `json:"kid"`
	Use string `json:"use"`
	KTY string `json:"kty"`
	N   string `json:"n"`
	Alg string `json:"alg"`
	E   string `json:"e"`
}

type JWKList struct {
	Keys []JWK `json:"keys"`
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
			log.Println(r.String())
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

	if redirect != nil {
		q := redirect.Query()
		q.Del("failure")
		redirect.RawQuery = q.Encode()
		return ctx.Redirect(http.StatusTemporaryRedirect, redirect.String())
	} else {
		return ctx.NoContent(http.StatusOK)
	}
}
