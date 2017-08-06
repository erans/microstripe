package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/token"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	stripeKey = kingpin.Flag("stripekey", "Stripe Key").Required().OverrideDefaultFromEnvar("MICROSTRIPE_STRIPE_KEY").Short('s').String()
	port      = kingpin.Flag("port", "Port to listen").Default("8100").Short('p').String()
)

type chargeParams struct {
	Email       string            `json:"email"`  // Optional if customerId is supplied
	Amount      string            `json:"amount"` // Like in the Stripe API - whole number * 100.0 so $2.5 is 250
	Currency    stripe.Currency   `json:"currency"`
	Description string            `json:"description"` // Description to be set on the charge
	Metadata    map[string]string `json:"metadata"`    // Optional: Meta data to be included on the charge
	Token       string            `json:"token"`
}

func chargeRequest(c echo.Context) error {
	var err error

	params := &chargeParams{}
	jsonDecoder := json.NewDecoder(c.Request().Body)
	if err = jsonDecoder.Decode(params); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid parameters. Reason: %s", err))
	}

	if params.Currency == "" {
		params.Currency = "usd"
	}

	var stripeToken *stripe.Token
	if stripeToken, err = token.Get(params.Token, nil); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch token")
	}

	// Try to use the sent email, otherwise try to retrieve it from the token
	var email string
	if params.Email != "" {
		email = params.Email
	} else {
		email = stripeToken.Email
	}

	var amount int
	if amount, err = strconv.Atoi(params.Amount); err != nil {
		return c.String(http.StatusBadRequest, "Invalid amount value")
	}

	chargeParams := &stripe.ChargeParams{
		Params:   stripe.Params{Meta: params.Metadata},
		Amount:   uint64(amount),
		Currency: params.Currency,
		Email:    email,
		Desc:     params.Description,
	}
	chargeParams.SetSource(params.Token)

	var ch *stripe.Charge
	if ch, err = charge.New(chargeParams); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create a new charge.")
	}

	return c.JSON(http.StatusOK, ch)
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	e := echo.New()
	stripe.Key = *stripeKey

	e.Use(middleware.Logger())

	e.POST("/v1/api/charge", chargeRequest)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
