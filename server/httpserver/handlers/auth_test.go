package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator"
	"generator/accounts"
	"generator/accounts/auth"
	"generator/internal/logger/zaplog"
	pkgauth "generator/pkg/auth"
	"generator/pkg/ethsig"
	"generator/server/emails"
	"generator/server/httpserver/handlers"
	"generator/storage/testdb"
)

const cookieName = "generator_auth"

func TestAuthHandlers(t *testing.T) {
	testdb.Run(t, func(ctx context.Context, t *testing.T, db generator.DB) {
		log := zaplog.NewLog()
		service := auth.NewService(db.Accounts(), accounts.NewService(db.Accounts()), pkgauth.JWTSigner{Secret: []byte("test secret")}, &emails.Service{}, log)
		cookie := pkgauth.NewSessionCookie(pkgauth.SessionSettings{Name: cookieName, Path: "/"})
		controller := handlers.NewAuth(log, service, cookie, nil)

		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		sign := func(message string) string {
			signature, err := crypto.Sign(ethsig.SignHash([]byte(message)), key)
			require.NoError(t, err)
			signature[crypto.RecoveryIDOffset] += 27
			return hexutil.Encode(signature)
		}
		call := func(handler http.HandlerFunc, body interface{}) *httptest.ResponseRecorder {
			data, err := json.Marshal(body)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			handler(recorder, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(data))))
			return recorder
		}
		// assertError checks the status and that the body holds exactly one error object.
		assertError := func(t *testing.T, recorder *httptest.ResponseRecorder, status int) {
			require.Equal(t, status, recorder.Code, recorder.Body.String())
			decoder := json.NewDecoder(recorder.Body)
			var response struct {
				Error string `json:"error"`
			}
			require.NoError(t, decoder.Decode(&response))
			assert.NotEmpty(t, response.Error)
			assert.Equal(t, io.EOF, decoder.Decode(&response), "the handler must stop after the first error")
			assert.Empty(t, cookiesOf(t, recorder))
		}

		registration := sign(accounts.DefaultMessageForRegistration)
		require.Equal(t, http.StatusOK, call(controller.MetamaskRegister, registration).Code)

		nonce, err := service.Nonce(ctx, crypto.PubkeyToAddress(key.PublicKey))
		require.NoError(t, err)
		signature := sign(nonce)

		type login struct {
			Nonce     string `json:"nonce"`
			Signature string `json:"signature"`
		}
		for _, test := range []struct {
			name    string
			handler http.HandlerFunc
			body    interface{}
			status  int
		}{
			{"register: not hex", controller.MetamaskRegister, "0xzz", http.StatusBadRequest},
			{"register: short signature", controller.MetamaskRegister, registration[:len(registration)-2], http.StatusUnauthorized},
			{"register: wallet is registered already", controller.MetamaskRegister, registration, http.StatusConflict},
			{"login: not hex", controller.MetamaskLogin, login{nonce, "0xzz"}, http.StatusBadRequest},
			{"login: short signature", controller.MetamaskLogin, login{nonce, signature[:len(signature)-2]}, http.StatusUnauthorized},
			{"login: another nonce", controller.MetamaskLogin, login{"0x01", sign("0x01")}, http.StatusUnauthorized},
		} {
			test := test
			t.Run(test.name, func(t *testing.T) {
				assertError(t, call(test.handler, test.body), test.status)
			})
		}

		t.Run("login sets the cookie", func(t *testing.T) {
			recorder := call(controller.MetamaskLogin, login{nonce, signature})
			require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
			cookies := cookiesOf(t, recorder)
			require.Len(t, cookies, 1)
			assert.Equal(t, cookieName, cookies[0].Name)
			assert.True(t, cookies[0].HttpOnly)
			_, err := service.Authorize(ctx, cookies[0].Value)
			assert.NoError(t, err)
		})

		t.Run("logout removes the cookie", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", nil)
			request.Header.Set("X-Forwarded-Proto", "https")
			controller.Logout(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)
			cookies := cookiesOf(t, recorder)
			require.Len(t, cookies, 1)
			assert.Equal(t, cookieName, cookies[0].Name)
			assert.Equal(t, "/", cookies[0].Path)
			assert.Empty(t, cookies[0].Value)
			assert.Negative(t, cookies[0].MaxAge)
			assert.True(t, cookies[0].Secure)
		})
	})
}

// cookiesOf returns the cookies the handler set, closing the body of the recorded response.
func cookiesOf(t *testing.T, recorder *httptest.ResponseRecorder) []*http.Cookie {
	response := recorder.Result()
	require.NoError(t, response.Body.Close())
	return response.Cookies()
}
