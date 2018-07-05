FORMAT: 1A
HOST: http://127.0.0.1:8889

# AnychainDB API

## Getting started with AnychainDB APIs

Welcome to the AnychainDB API reference! 

AnychainDB is distributed, high-loaded, blockchain database.

The AnychainDB API is based on REST. 
This documentation lists and describes the resources you can used to manipulate objects on the AnychainDB Blockchain. 

## Authentication

AnychainDB API uses [Elliptic Curve Digital Signature Algorithm](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm) method, 
that requires to obtain public/private keys for account identification.
First you create account with Accounts resource.
After you can get access to post new data (payload) to AnychainDB.

## Server Responses

+ 200 OK - The request was successful.
+ 202 Accepted - The request's data was successfully posted to the blockchain queue.
+ 400 Bad Request - The request could not be understood or was missing required parameters (check msg field in response for details).
+ 401 Unauthorized - Authentication failed or user does not have permissions for the requested operation (check msg field in response for details).
+ 404 Not Found - Resource was not found.
+ 405 Method Not Allowed - Requested method is not supported for the specified resource.
+ 429 Too Many Requests - Exceeded AnychainDB API limits.

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

## Payloads [/v1/payloads{?limit}{?offset}]

This resource is intended for listing, sending and view details about transaction data (payload).

Payload contains from two main parts: private and public data. You can place any data to this fields.
*Private data* encrypted by ECDH algorithm with public key of receiver. For encryption of private data AnychainDB uses [Elliptic-curve Diffie–Hellman](https://en.wikipedia.org/wiki/Elliptic-curve_Diffie%E2%80%93Hellman) protocol, that also uses ECDSA but for assymetric data encryption.
*Public data* keeps any open data without encryption.
ReceiverAccountID fields contains information about receiver's account, that decrypts private data with it private key.

### View a payloads list [GET]

+ Parameters
    + limit: 100 (number, optional)
    If a limit count is given, no more than that many rows will be returned. Limit can range between 1 and 500 items.
    + offset: 0 (number, optional)
    Offset says to skip that many rows before beginning to return rows.

+ Response 200 (application/json)
    + Attributes
        + code: 200 (number)
        + msg: OK (string)
        + data (array[PayloadGet])
        Payloads list

### Create a new payload [POST]

+ Request (application/json)
    + Attributes
        + account_id: 5acacd9b6d9bf091f214ad7b (required)
        Requester account identifier in blockchain
        + private_key: 6PSXoObyVM1slemJ+GfAluUzIbU9pNf7CX5J36O3iW8= (required)
        Requester private key in blockchain
        + public_key: BLnQWwtB2SEjisrmHLLAXU2drEaZZSVeFFuoWEwplMJwpEStOAzeZv0+SP/q4etJcaISoDOBnwvc9Pztuz9LUVw= (required)
        Requester public key in blockchain
        + data (PayloadPost)
        Payload object

+ Response 202 (application/json)
    + Attributes
        + code: 202 (number)
        + msg: Accepted (string)
        + data: 5acb5aa66d9bf0c526678d12 (string)
        Payload ID

## Payloads | Details [/v1/payloads/{id}]

This resource is intended for viewing details about transaction data (payload).

### View a payload details [GET]

+ Parameters
    + id (string)
    ID of the Payload in the form of an string

+ Response 200 (application/json)
    + Attributes
        + code: 200 (number)
        + msg: OK (string)
        + data (PayloadGet)
        Payload details

## Payloads | Search [/v1/payloads{?query}{?limit}{?offset}]

### Search Payloads [GET]

This resource is intended for searching payloads using MongoDB query language.

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
        + data (array[PayloadGet])
        Payload filtered list

# Data Structures

## Account (object)
+ account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique user's identifier in blockchain
+ private_key: 6PSXoObyVM1slemJ+GfAluUzIbU9pNf7CX5J36O3iW8= (string)
Private key of user WARNING! Private Key should be kept in SAFE place!
+ public_key: BLnQWwtB2SEjisrmHLLAXU2drEaZZSVeFFuoWEwplMJwpEStOAzeZv0+SP/q4etJcaISoDOBnwvc9Pztuz9LUVw= (string)
Public key of user

## PayloadGet (object)

+ _id: 3acaad9v49bf591f212ad7b (string)
Unique payload identifier in blockchain
+ sender_account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique sender account identifier in blockchain (i.e. Advertiser)
+ receiver_account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique receiver account identifier in blockchain (i.e. CPA Network)
+ public_data: test click (string)
Public data available to all
+ private_data: test stream (string)
Private data encrypted with public key of receiver
+ created_at: 2512351252135 (number)
Unix time (seconds) datetime of conversion

## PayloadPost (object)

+ _id: 3acaad9v49bf591f212ad7b (string)
Unique payload identifier in blockchain
+ receiver_account_id: 5acacd9b6d9bf091f214ad7b (string)
Unique receiver account identifier in blockchain (i.e. CPA Network)
+ public_data: test click (string)
Public data available to all
+ private_data: test stream (string)
Private data encrypted with public key of receiver
+ created_at: 2512351252135 (number)
Unix time (seconds) datetime of conversion