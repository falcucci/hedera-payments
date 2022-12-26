### Hedera Coin Payments API

This repository contains a simple API for integrating Hedera Coin payments into your application. With this API, you can easily send and receive payments using Hedera's native cryptocurrency, HBAR.

### Prerequisites

Before you can use this API, you will need to sign up for a Hedera account and obtain API keys. These keys will allow you to authenticate your API requests and access the payment functionality.

### Getting Started

To get started with the Hedera Coin Payments API, you will need to clone this repository and install the required dependencies.

```bash
git clone https://github.com/falcucci/hedera-coin-payments-api.git
cd hedera-coin-payments-api
npm install

Next, you will need to configure the API by setting your Hedera API keys in the .env file.

```bash
# .env
HEDERA_PUBLIC_KEY=<your-public-key>
HEDERA_PRIVATE_KEY=<your-private-key>
```

Once you have configured the API, you can start the server by running the following command:

```bash
npm start
```
The API will be running at http://localhost:3000.

### API Endpoints

The Hedera Coin Payments API provides the following endpoints:

POST /payment: Send a payment to a specific address.
GET /balance/:accountId: Check the balance of a specific account.
GET /transactions/:accountId: Retrieve a list of transactions for a specific account.
Examples

Here are some examples of how to use the API:

```javascript
const axios = require('axios');

async function sendPayment() {
  const response = await axios.post('http://localhost:3000/payment', {
    to: '3f9a07d83c604dba400d13df4d3456',
    amount: 100,
  });
  console.log(response.data);
}

async function checkBalance() {
  const response = await axios.get('http://localhost:3000/balance/3f9a07d83c604dba400d13df4d3456');
  console.log(response.data);
}

async function getTransactions() {
  const response = await axios.get('http://localhost:3000/transactions/3f9a07d83c604dba400d13df4d3456');
  console.log(response.data);
}
```

### Contributing

If you would like to contribute to the development of this API, please fork the repository and submit a pull request.

### License

This project is licensed under the MIT License - see the LICENSE file for details.
