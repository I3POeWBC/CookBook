package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {

	now := time.Now()
	var (
		err          error
		tokenString  string
		tokenIssued  *jwt.Token
		tokenChecked *jwt.Token
		secretIssued string = "Only_I_know_this_secret!"
		//secretChecked string
	)

	tokenIssued = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"game": "warior, bard, druid, mag, cleric",
		"exp": jwt.NumericDate{
			Time: now.Add(time.Hour),
		},
		"nbf": jwt.NumericDate{
			Time: now.Add(time.Second),
		},
		"iss": JwtSampleIssuer,
		"sub": JwtSampleSubject,
		"aud": JwtSampleAudience,
	})

	if tokenString, err = tokenIssued.SignedString([]byte(secretIssued)); err != nil {
		log.Fatalf("token.SignedString err: [%v]", err)
	} else {
		fmt.Printf("tokenString: [%s]\n", tokenString)
	}

	explain := []string{
		"Must be not before error",
		"Expect all correct",
	}
	for _, v := range explain {

		if tokenChecked, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			//fmt.Printf("", metoken.Method.)

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(secretIssued), nil
		}); err != nil {
			fmt.Printf("[%s]\n\tjwt.Parse err: [%v]\n", v, err)
		} else {
			if claims, ok := tokenChecked.Claims.(jwt.MapClaims); ok {
				fmt.Printf("[%s]\n\tClime: [%s] NotBefore: [%v]\n", v, claims["game"], claims["nbf"])
			} else {
				fmt.Println(err)
			}

		}

		time.Sleep(time.Second * 2)
	}

}
