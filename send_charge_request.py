import json
import requests
import stripe
import os
stripe.api_key = os.environ["STRIPE_TEST_KEY"]

token = stripe.Token.create(
  card={
    "number": '4242424242424242',
    "exp_month": 12,
    "exp_year": 2018,
    "cvc": '123'
  },
)

token_id = token.id

data = {
    "email" : "eran+stripechargetest@sandler.co.il",
    "amount": 6500,
    "description": "this is a test charge",
    "metadata" : {
        "key1" : "value1",
        "key2" : "value2"
    },
    "token" : token_id
}

r = requests.post("http://localhost:8000/v1/api/charge", data=json.dumps(data), headers={ "Content-Type" : "application/json" })
print r.status_code
print r.text
