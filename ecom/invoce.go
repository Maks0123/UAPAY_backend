package ecom

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	createInvoiceUrl string = "https://api.demo.uapay.ua/api/invoicer/invoices/create"
)

/*
{
    "description": "Some new book 12",
	"amount": 400077
}
*/

// Create Invoce with key

func CreateInvoce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var invoiceData map[string]interface{}
	json.NewDecoder(r.Body).Decode(&invoiceData)

	// json.NewEncoder(w).Encode(invoiceData)
	// json.NewEncoder(w).Encode(invoiceData["description"])

	// geting values amount, description from posts
	var description = invoiceData["description"]
	var amount = invoiceData["amount"]

	//currentTime := time.Now().Unix()
	currentTime := time.Now()
	var externalId string = currentTime.String()

	iat := getUnixTime()

	// Token jwt Standard Claim Object

	type params struct {
		SessionId  string `json:"sessionId"`
		SystemType string `json:"systemType"`
	}

	type data struct {
		ExternalId  string  `json:"externalId"`
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		RedirectUrl string  `json:"redirectUrl"`
		Type        string  `json:"type"`
		ExtraInfo   string  `json:"extraInfo"`
	}

	type InvoceToken struct {
		Params params `json:"params"`
		Data   data   `json:"data"`
		Iat    int64  `json:"iat"`
		Token  string `json:"token"`
		jwt.StandardClaims
	}

	var InvoceTokenClaim = InvoceToken{
		Params: params{
			SessionId:  sessionId,
			SystemType: "ECOM",
		},
		Data: data{
			ExternalId:  externalId,
			Description: description.(string),
			Amount:      amount.(float64),
			RedirectUrl: "https://uapay.ua",
			Type:        "PAY",
			ExtraInfo:   "{\"phoneFrom\":\"380971112233\",\"phoneTo\":\"380631112233\",\"cardToId\":\"216f8390-9abc-428d-89d6-7be50183afb5\"}",
		},
		Iat:            iat,
		StandardClaims: jwt.StandardClaims{
			// Enter expiration in milisecond
			// ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}

	// Create a new claim with HS256 algorithm and token claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, InvoceTokenClaim)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		log.Fatal(err)
	}

	invoiceMessage := map[string]interface{}{
		"params": map[string]string{
			"sessionId":  sessionId,
			"systemType": "ECOM",
		},
		"data": map[string]interface{}{
			"externalId":  externalId,
			"description": description,
			"amount":      amount,
			"redirectUrl": "https://uapay.ua",
			"type":        "PAY",
			"extraInfo":   "{\"phoneFrom\":\"380971112233\",\"phoneTo\":\"380631112233\",\"cardToId\":\"216f8390-9abc-428d-89d6-7be50183afb5\"}",
		},
		"iat":   iat,
		"token": tokenString,
	}

	bytesRepresentationInvoice, err := json.Marshal(invoiceMessage)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(createInvoiceUrl, "application/json", bytes.NewBuffer(bytesRepresentationInvoice))
	if err != nil {
		log.Fatalln(err)
	}

	var resultInvoice map[string]interface{}

	// Decode result
	json.NewDecoder(resp.Body).Decode(&resultInvoice)

	var dataInvoceToken = resultInvoice["data"].(map[string]interface{})

	//json.NewEncoder(w).Encode(resultInvoice)
	// json.NewEncoder(w).Encode(dataInvoceToken["token"])

	var invoceTokenString = dataInvoceToken["token"].(string)

	type InvoiceRespClaims struct {
		Id               string `json:"id"`
		PaymentPageUrl   string `json:"paymentPageUrl"`
		PaymentPageUrlQR string `json:"paymentPageUrlQR"`
		Iat              string `json:"iat"`
		jwt.StandardClaims
	}

	decodeToken, err := jwt.ParseWithClaims(
		invoceTokenString,
		&InvoiceRespClaims{},
		func(decodeToken *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	claims := decodeToken.Claims.(*InvoiceRespClaims)
	json.NewEncoder(w).Encode(claims.PaymentPageUrl)

}

// create invoce without key
func DemoCreateInvoce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var invoiceData map[string]interface{}
	json.NewDecoder(r.Body).Decode(&invoiceData)

	json.NewEncoder(w).Encode(invoiceData)
	json.NewEncoder(w).Encode(invoiceData["description"])

	// geting values amount, description from posts
	var description = invoiceData["description"]
	var amount = invoiceData["amount"]

	//currentTime := time.Now().Unix()
	currentTime := time.Now()
	var externalId string = currentTime.String()

	iat := getUnixTime()
	invoiceMessage := map[string]interface{}{
		"params": map[string]string{
			"sessionId":  sessionId,
			"systemType": "ECOM",
		},
		"data": map[string]interface{}{
			"externalId":  externalId,
			"description": description,
			"amount":      amount,
			"redirectUrl": "https://uapay.ua",
			"type":        "PAY",
			"extraInfo":   "{\"phoneFrom\":\"380971112233\",\"phoneTo\":\"380631112233\",\"cardToId\":\"216f8390-9abc-428d-89d6-7be50183afb5\"}",
		},
		"iat":   iat,
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXJhbXMiOnsic2Vzc2lvbklkIjoiZTQ2Zjk5YWQtNWZhNi00Njg2LWE0ZTMtYjdiODhhZjZjM2VhIiwic3lzdGVtVHlwZSI6IkVDT00ifSwiZGF0YSI6eyJleHRlcm5hbElkIjoiMTUwMDM4MzA3NTAwIiwicmV1c2FiaWxpdHkiOmZhbHNlLCJkZXNjcmlwdGlvbiI6ItGC0LXRgdGC0L7QstGL0Lkg0L_Qu9Cw0YJp0LYiLCJhbW91bnQiOjEwMCwicmVkaXJlY3RVcmwiOiJodHRwczovL3VhcGF5LnVhIiwidHlwZSI6IlBBWSIsImV4dHJhSW5mbyI6IntcInBob25lRnJvbVwiOlwiMzgwOTcxMTEyMjMzXCIsXCJwaG9uZVRvXCI6XCIzODA2MzExMTIyMzNcIixcImNhcmRUb0lkXCI6XCIyMTZmODM5MC05YWJjLTQyOGQtODlkNi03YmU1MDE4M2FmYjVcIn0ifX0.5M4zgtmEqfMViuCBigILlzKRGSY6VrmKw-g9CtY7KP8",
	}

	bytesRepresentationInvoice, err := json.Marshal(invoiceMessage)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(createInvoiceUrl, "application/json", bytes.NewBuffer(bytesRepresentationInvoice))
	if err != nil {
		log.Fatalln(err)
	}

	var resultInvoice map[string]interface{}

	// Decode result
	json.NewDecoder(resp.Body).Decode(&resultInvoice)

	// map resultInvoice data
	var paymentPageUrl = resultInvoice["data"].(map[string]interface{})

	//log.Println(resultInvoice)
	//log.Println(paymentPageUrl["paymentPageUrl"])

	//json.NewEncoder(w).Encode(resultInvoice)
	json.NewEncoder(w).Encode(paymentPageUrl["paymentPageUrl"])

}