# Receipt Processor

## Introduction
This is a coding challenge from Fetch-Rewards. See the requirements from here: https://github.com/fetch-rewards/receipt-processor-challenge.

## Overview
This is a backend service written in Go which processes receipt awards points. It provides two API endpoints:
1. **Process Receipt**: Submits a receipt for processing and returns an ID.
2. **Get Points**: Retrieves the points awarded for a given receipt ID.

---
## Build and Run
### Using Go
0. Make sure you have **Go** installed on your machine.
1. Clone this repo to your local machine and navigate to the root directory.
2. Download dependencies.
```bash
go mod download
```

3. Build the application.
```bash
go build -o main
```

4. Run the application.
```bash
./main
```

5. Access the Application.
Once the application is running, you can access it at http://localhost:8080

### Using Docker
0. Make sure you have **Docker** installed on your machine.
1. Clone this repo to your local machine and navigate to the root directory.
2. Build the Docker image.
```bash
docker build -t reciept-processor .
```

3. Run the Docker container.
```bash
docker run -p 8080:8080 receipt-processor
```


---
## Unit Tests
There are one unit test for services that can be run by sending the following command in your project terminal.

`go test ./services`

---
## API Documentation
### 1. Process Receipt
- **URL:** `/receipts/process`
- **Method:** `POST`
- **Request:** JSON object representing the receipt
- **Response:** JSON object with the ID of the receipt

#### Request
The request body should be a JSON object with the following properties:

| Property | Type | Required | Description |
| -------- | ---- | -------- | ----------- |
| retailer | string | Yes | The name of the retailer. |
| purchaseDate | string | Yes | The date of the purchase in (yyyy-mm-dd). |
| purchaseTime | string | Yes | The time of the purchase in 24-hour format. |
| items | array | Yes | An array of purchased items. |
| total | string | Yes | The total amount of the purchase Dollar. |

#### Example Request

``` json 
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```

#### Response
| Property | Type | Description |
| -------- | ---- | ----------- |
| id | string | id of the receipt that was processed. |

#### Example Response

```json
{
  "id": "5cc04679-9360-4f23-adf6-342d6c45d5b8"
}
```

#### Status 

| Status Code | Description |
| ----------- | ----------- |
| 200 | Receipt processed successfully. |
| 400 | Invalid request body (receipt data). |
| 500 | Server error during processing. |


### 2. Get Points
- **URL:** `/receipts/{id}/points`
- **Method:** `GET`
- **Response:** JSON object with points awarded.

#### Request
Include the targeted receipt 'id' as parameter in the API endpoint.

#### Example Request

`http://localhost:8080/receipts/5cc04679-9360-4f23-adf6-342d6c45d5b8/points`

#### Response
| Property | Type | Description |
| -------- | ---- | ----------- |
| points | int | Point rewarded for the receipt. |

#### Example Response

```json
{
  "points": 28
}
```

#### Status 

| Status Code | Description |
| ----------- | ----------- |
| 200 | Points retrieved successfully. |
| 404 | Receipt ID not found. |


