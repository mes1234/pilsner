# pilsner

Stream processing with filters

[![Go](https://github.com/mes1234/pilsner/actions/workflows/go.yml/badge.svg)](https://github.com/mes1234/pilsner/actions/workflows/go.yml)

# Pub sub mode (PS mode)

In pub sub mode:

* publisher will send data to memoryStream
* consumer will receive new messages until it is active (online consumer Ocon)
    * active consumer is one which has active connection with server
* ACK is based on policy:
    * one time no ACK
    * n times with ACK delay based:
        * linear
        * power
        * random
* Data is in memoryStream as long as:
    * Expiration condition is met
    * There is no ACK for active consumer and policy is still executed
* Data is in memory only

# Pub sub wih authentication

* auth is done based on JWT token

# Pub sub with authorization

* authorization is done based on JWT token id
* memoryStream writers and readers can be managed via console


# Pub sub replay mode (PSR mode)

In Pub sub replay:

* All pub sub apply
* Exceptions:
    * Data is never deleted
    * Replay consumer (RCon) 
      * Is a consumer which can be defined as replay same as online
      * Will get all item from beginning of memoryStream
      
# Preserved streams and consumers (only Replay)
  * Data is never deleted
  * Data is saved in persistent storage:
    * File
    * Document DB
    * SQL DB
  * only Replay consumers state is preserved 

# Smart consumer
  * Stream has corresponding `.proto` file 
  * New entity is defined - Filter
    * Filter can be defined as way to postprocess item based on `.proto` file 
    * Filter will allow to pass or not pass item forward 
    * Filter can be follow map reduce pattern