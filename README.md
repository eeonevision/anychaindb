# Leadschain
Leadschain - secure, fast and open blockchain platform for providing transparancy in interaction between advertisers and webmasters.
Leadschain uses the modern technologies not only from blockchain world. 
Let's take a look at them:
  * [MongoDB] - the high-performable database with full-text search
  * [Tendermint] - the heart of blockchain platform
  * [Docker] - all components of platform wrapps in containers for fast deploy and ease to use
  * [MsgPack] - transport between messages in platform
  * [Golang] - fast and beautiful language

## Introduction to idea
### Problem
At the current time, in the affiliate networks market (CPA) webmasters are faced with such an unpleasant event as a “shave”. Shave-fraudulent operations of the advertiser, aimed at reducing payments to webmasters from the above clients. In other words, when a webmaster has led one number of users to the advertiser, some users may be intentionally “lost” or “not processed” by the advertiser.

### Proposed solution
In this system, each user's transition to the advertiser's page will be sent in a distributed database (blockchain), which will be publicly available to all users of the network. Each advertiser can connect to this system for free and use it as a guarantor of honesty and transparency. 
That gives us some features as:
1. Transparency of actions provides public access to header data on transactions by any users of the network, but keeps the anonymity of the advertiser's offer and the identity of the webmaster
2. The reliability of the information is supported by the use of the Consensus algorithm, with which all participants of the blockchain network have the right to make decisions about the veracity or falsity of the incoming information automatically, based on the programmed algorithm
3. This solution posted in opensource and each participant will be able to verify its operation and use in their developments and systems
4. Actions of the advertiser and its employees (affecting to conversions and payments) become more transparent
5. The advertiser is better informed about the unfair actions of its employees

## Architecture
Below picture represents the general concept of interaction between participants of the system.

![Base concept](../docs/architecture/concept.png)

Let's breakdown Blockchain component:

![Nodes communication](../docs/architecture/nodes_cm.png)

The advertiser creates an account in our distributed storage and expands the blockchain, or it can deploy the full node on its local machine and integrate with Afifliate Network (CPA).
Now each user's transition to the advertiser's page from the webmaster will be sent to the blockchain. The fact of clicking will be confirmed by all network participants using consensus algorithm.
If the conversion status is changed and the payment is made, the data will also be sent to the blockchain
Now any user of the network will be able to see transactions in real time. To ensure anonymity, additional information about the transaction (name of the offer, advertiser and affiliate) will be available through the interface of the chosen CPA-network from the affiliate account.

## Installation
Will be later...

## Additional docs
  * [Leadschain REST API] - REST API for Leadschain client
  * [Tendermint Docs] - Tendermint documentation

License
----
Apache 2.0

   [MongoDB]: <https://www.mongodb.com/>
   [Tendermint]: <https://github.com/tendermint/tendermint>
   [Docker]: <https://www.docker.com/>
   [MsgPack]: <https://msgpack.org/>
   [Golang]: <https://golang.org/>
   [Leadschain REST API]: <https://leadschain1.docs.apiary.io/>
   [Tendermint Docs]: <http://tendermint.readthedocs.io/en/master/introduction.html>