FORMAT: 1A
HOST: http://127.0.0.1:8889

# Leadschain API

## Getting started with Leadschain APIs

Welcome to the Leadschain API reference! 

Leadschain is distributed, high-loaded, blockchain database for tracking advertisers 
clicks, transitions and conversions.

The Leadschain API is based on REST. 
This documentation lists and describes the resources you can used to manipulate objects on the Leadschain Blockchain. 

## Authentication

Leadschain API uses [Public-key authentication](https://en.wikipedia.org/wiki/Public-key_cryptography) method, 
that requires to obtain public/pair for account identification. First you create account with Accounts resource.
After you can get access to post new transition and conversion in Leadschain.

## Server Responses

+ 200 OK - The request was successful.
+ 202 Accepted - The request's data was successfully posted to the blockchain queue.
+ 400 Bad Request - The request could not be understood or was missing required parameters (check msg field in response for details).
+ 401 Unauthorized - Authentication failed or user does not have permissions for the requested operation (check msg field in response for details).
+ 404 Not Found - Resource was not found.
+ 405 Method Not Allowed - Requested method is not supported for the specified resource.
+ 429 Too Many Requests - Exceeded Leadschain API limits.

## Accounts [/v1/accounts]

This resource is intended for create accounts.
Every users have their unique accounts in the system.

### Create a new account [POST]

+ Request (application/json)
    
+ Response 202 (application/json)
    + Attributes
        + code: 202 (number)
        + msg: Accepted (string)
        + data (Account)


## Transitions [/v1/transitions{?limit}{?offset}]

This resource is intended for listing and sending about transitions.

Transition is a first "touch" of client with the advertiser's product. There is no conversion happend, but the user transitions
should be captured in blockchain by advertiser.

### View a transitions list [GET]

+ Parameters
    + limit: 100 (number, optional)
    If a limit count is given, no more than that many rows will be returned. Limit can range between 1 and 500 items.
    + offset: 0 (number, optional)
    Offset says to skip that many rows before beginning to return rows.

+ Response 200 (application/json)
    + Attributes
        + code: 200 (number)
        + msg: OK (string)
        + data (array[TransitionGet])
        Transitions list

### Create a new transition [POST]

+ Request (application/json)
    + Attributes
        + account_id: 5acacd9b6d9bf091f214ad7b (required)
        Requester account identifier in blockchain
        + private_key: 6PSXoObyVM1slemJ+GfAluUzIbU9pNf7CX5J36O3iW8= (required)
        Requester private key in blockchain
        + public_key: BLnQWwtB2SEjisrmHLLAXU2drEaZZSVeFFuoWEwplMJwpEStOAzeZv0+SP/q4etJcaISoDOBnwvc9Pztuz9LUVw= (required)
        Requester public key in blockchain
        + data (TransitionPost)
        Transition object
    
+ Response 202 (application/json)
    + Attributes
        + code: 202 (number)
        + msg: Accepted (string)
        + data: 5acb5aa66d9bf0c526678d12 (string)
        Transition ID

## Transitions | Details [/v1/transitions/{id}]

This resource is intended for viewing details about transitions.

### View a transition details [GET]

+ Parameters
    + id (string)
    ID of the Transition in the form of an string

+ Response 200 (application/json)
    + Attributes
        + code: 200 (number)
        + msg: OK (string)
        + data (TransitionGet)
        Transition details

## Transitions | Search [/v1/transitions{?query}{?limit}{?offset}]

### Search transitions [GET]

This resource is intended for searching transitions with MongoDB query language.

See more at: https://docs.mongodb.com/manual/reference/method/db.collection.find/

+ Parameters
    + query: { status: { $in: [ "A", "D" ] } } (string)
    MongoDB search query language
    + limit: 100 (number, optional)
    If a limit count is given, no more than that many rows will be returned. Limit can range between 1 and 500 items.
    + offset: 0 (number, optional)
    Offset says to skip that many rows before beginning to return rows.

+ Response 200 (application/json)
    + Attributes
        + code: 200 (number)
        + msg: OK (string)
        + data (array[TransitionGet])
        Transitions filtered list

## Conversions [/v1/conversions{?limit}{?offset}]

This resource is intended for listing, sending and view details about conversions.

Conversion is "fact" of executing the initial agreement of advertiser with affiliate.
There may be action, when user filled form data on advertiser's site, installed app or some more.
The conversion may not be initial approved by advertiser and then it status will setted as "PENDING".
When the advertiser made a decision about user, then new conversion will created with updated status ("APPROVED" or "REJECTED").

### View a conversions list [GET]

+ Parameters
    + limit: 100 (number, optional)
    If a limit count is given, no more than that many rows will be returned. Limit can range between 1 and 500 items.
    + offset: 0 (number, optional)
    Offset says to skip that many rows before beginning to return rows.

+ Response 200 (application/json)
    + Attributes
        + code: 200 (number)
        + msg: OK (string)
        + data (array[ConversionGet])
        Conversions list

### Create a new conversion [POST]

+ Attributes
    + Statuses (Status)

+ Request (application/json)
    + Attributes
        + account_id: 5acacd9b6d9bf091f214ad7b (required)
        Requester account identifier in blockchain
        + private_key: 6PSXoObyVM1slemJ+GfAluUzIbU9pNf7CX5J36O3iW8= (required)
        Requester private key in blockchain
        + public_key: BLnQWwtB2SEjisrmHLLAXU2drEaZZSVeFFuoWEwplMJwpEStOAzeZv0+SP/q4etJcaISoDOBnwvc9Pztuz9LUVw= (required)
        Requester public key in blockchain
        + data (ConversionPost)
        Transition object

+ Response 202 (application/json)
    + Attributes
        + code: 202 (number)
        + msg: Accepted (string)
        + data: 5acb5aa66d9bf0c526678d12 (string)
        Conversion ID

## Conversions | Details [/v1/conversions/{id}]

This resource is intended for viewing details about conversions.

### View a conversion details [GET]

+ Parameters
    + id (string)
    ID of the Conversion in the form of an string

+ Response 200 (application/json)
    + Attributes
        + code: 200 (number)
        + msg: OK (string)
        + data (ConversionGet)
        Conversion details

## Conversions | Search [/v1/conversions{?query}{?limit}{?offset}]

### Search Conversions [GET]

This resource is intended for searching conversions with MongoDB query language.

See more at: https://docs.mongodb.com/manual/reference/method/db.collection.find/

+ Parameters
    + query: { status: { $in: [ "A", "D" ] } } (string)
    MongoDB search query language
    + limit: 100 (number, optional)
    If a limit count is given, no more than that many rows will be returned. Limit can range between 1 and 500 items.
    + offset: 0 (number, optional)
    Offset says to skip that many rows before beginning to return rows.

+ Response 200 (application/json)
    + Attributes
        + code: 200 (number)
        + msg: OK (string)
        + data (array[ConversionGet])
        Conversion filtered list

# Data Structures

## Status (object)
+ PENDING
+ APPROVED
+ REJECTED

## Account (object)
+ account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique user's identifier in blockchain
+ private_key: 6PSXoObyVM1slemJ+GfAluUzIbU9pNf7CX5J36O3iW8= (string)
Private key of user WARNING! Private Key should be kept in SAFE place!
+ public_key: BLnQWwtB2SEjisrmHLLAXU2drEaZZSVeFFuoWEwplMJwpEStOAzeZv0+SP/q4etJcaISoDOBnwvc9Pztuz9LUVw= (string)
Public key of user

## TransitionGet (object)

+ _id: 3acaad9v49bf591f212ad7b (string)
Unique transition identifier in blockchain
+ advertiser_account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique advertiser account identifier in blockchain
+ affiliate_account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique affiliate account identifier in blockchain (i.e. CPA Network)
+ click_id: test click (string)
Unique click identifier from external system
+ stream_id: test stream (string)
Stream identifier (synonym to platform id)
+ offer_id: test offer (string)
Offer identifier in affiliate network
+ created_at: 2512351252135 (number)
Unix time (seconds) datetime of transition
+ expires_in: 12312312379 (number)
Date of transition expiration in seconds from Unix time epoch

## TransitionPost (object)

+ advertiser_account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique advertiser account identifier in blockchain
+ affiliate_account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique affiliate account identifier in blockchain (i.e. CPA Network)
+ click_id: test click (string)
Unique click identifier from external system
+ stream_id: test stream (string)
Stream identifier (synonym to platform id)
+ offer_id: test offer (string)
Offer identifier in affiliate network
+ expires_in: 12312312379 (number)
Date of transition expiration in seconds from Unix time epoch

## ConversionGet (object)

+ _id: 3acaad9v49bf591f212ad7b (string)
Unique conversion identifier in blockchain
+ advertiser_account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique advertiser account identifier in blockchain
+ affiliate_account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique affiliate account identifier in blockchain (i.e. CPA Network)
+ click_id: test click (string)
Unique click identifier from external system
+ stream_id: test stream (string)
Stream identifier (synonym to platform id)
+ offer_id: test offer (string)
Offer identifier in affiliate network
+ client_id: test client (string)
Client identifier in advertiser CRM
+ goal_id: 0 (string)
Goal identifier discussed with affiliate network
+ created_at: 2512351252135 (number)
Unix time (seconds) datetime of conversion
+ comment: test comment (string)
Optional comment to conversion
+ status: PENDING (string, required)
Status of Conversion

## ConversionPost (object)

+ advertiser_account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique advertiser account identifier in blockchain
+ affiliate_account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique affiliate account identifier in blockchain (i.e. CPA Network)
+ click_id: test click (string)
Unique click identifier from external system
+ stream_id: test stream (string)
Stream identifier (synonym to platform id)
+ offer_id: test offer (string)
Offer identifier in affiliate network
+ client_id: test client (string)
Client identifier in advertiser CRM
+ goal_id: 0 (string)
Goal identifier discussed with affiliate network
+ comment: test comment (string)
Optional comment to conversion
+ status: PENDING (string, required)
Status of Conversion