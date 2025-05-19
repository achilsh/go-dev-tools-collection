package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	// "github.com/dgrijalva/jwt-go"

	// "github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v5"
)

func TestTokenGen(t *testing.T){
	// jwt.RegisteredClaims{}

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
	// 	"foo":"bar",
	// 	"nbf": time.Now().Unix(),
	// })

	{
		//call 1: //签名方法是：jwt.SigningMethodHS256
		token :=jwt.New(jwt.SigningMethodHS256)
		//key 用于签名的私钥； 而公钥用于 验证 token.
		key := []byte("a3a4")
		tstr, err := token.SignedString(key)
		assert.Equal(t, err, nil)
		t.Logf("token str: %v", tstr)

		//解析，验证，校验：
		decodedToken, err := jwt.Parse(tstr, func(token *jwt.Token) (interface{}, error) {
					// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
					return key, nil
				}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		assert.Equal(t, err, nil)
		
		if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok{
			t.Logf("claims:%v", claims)
		}
	}

	{
		ta := jwt.NewWithClaims(
			// 签名方法
			jwt.SigningMethodHS256,
			//声明应该在如下字段中： `exp: 过期时间`, `iat：签发时间`, `nbf： 生效时间之前`, `iss: 颁发者`, `sub: 主题` and `aud: 受众`.
			jwt.MapClaims{
				"iss":"adfa", //颁发者
				"sub":"isdfad", //主题
				// 自定义声明，只要不和已经注册的声明名冲突即可
				"aaa":1232,
			},
		)
		// privateKey is 私钥
		privateKey := []byte("adfadf")
		//
		token, err := ta.SignedString(privateKey)
		assert.Equal(t, err,nil)
		t.Logf("token: %v", token)

		// 解析，校验：
		{
			decodedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
						// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
						return privateKey, nil
					}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
			assert.Equal(t, err, nil)
			
			if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok{
				t.Logf("2, claims:%v", claims)
				assert.Equal(t, claims["iss"],"adfa" )
				assert.Equal(t, claims["sub"],"isdfad" )

				aValue := int(claims["aaa"].(float64))
				assert.Equal(t, aValue, 1232)
			}
		}
	}

	{
		//从文件中读取： jwt.ParseECPrivateKeyFromPEM()
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			log.Fatalf("无法生成私钥: %v", err)
		}
		t.Logf("private key: %+v", privateKey)
		publicKey := &privateKey.PublicKey
	
		// 使用自定义声明类型， 比如：CustomClaims
		type CustomClaims struct {
			Foo string `json:"foo"`
			// 嵌入
			 jwt.RegisteredClaims 
		}

		claims := CustomClaims {
			Foo:"bar",
			RegisteredClaims: jwt.RegisteredClaims {
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(2*time.Hour)),

			},
		}

		t1 := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
		s,err:= t1.SignedString(privateKey)
		assert.Equal(t, err,nil)
		t.Logf("encoded token: %s", s)

		//解析，校验：
		{
			decodeToken, err := jwt.ParseWithClaims(s, new(CustomClaims), func(token *jwt.Token) (any, error) {
				t.Logf("alg method: %+v", token.Method)
				return publicKey, nil
			})
			t.Logf("err: %+v", err)

			assert.Equal(t, err, nil)
			if claims, ok := decodeToken.Claims.(*CustomClaims); ok {
				t.Logf("foo: %v", claims.Foo)	
			}
		}
	}
}
