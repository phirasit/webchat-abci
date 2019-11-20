# Webchat ABCI

Webchat ABCI is an abci (abstract blockchain interface) web chat application that keeps information on Tendermint blockchains.   
This software is released under the **MIT License**.

## Installation & Running
```
# clone the project
hg clone https://phirasit@bitbucket.org/phirasit/webchat-abci
cd webchat-abci
# build the project
make build
# start all docker nodes
docker-compose up
```

#### Docker Image
If you want to manually run the node, you have to mount in the 
chain configuration file to ``/root/.tendermint``
inside the container. The communication port is the default 
Tendermint port(26657). The file ``docker-compose.yml`` 
can be used as an example of how to use the image.

## Available API
From Tendermint RPC (https://tendermint.com/rpc/?go), important APIs
in this projects are  ``/abci_query`` and ``/broadcast_tx_*``
(other methods still can be called). 
The value ``data`` in ``abci_query`` and the value ``tx`` in ``broadcast_tx``
have to be in JSON with the following specification  

| Parameter | Type   | Description |
|-----------|--------|-------------|
| type      | String | the type of the operation |
| data      | Object | the data of the operation specified by the operation |
| (optional) nonce | String | random string in case of deliberately sending the same transaction (This is used to prevent blockchains from assuming it's the same transaction) |

Note: All API should be called via HTTP GET method  

## List of operations
### /abci_query
#### type: "get_message"
Get the message from the group.
###### Request data
| Parameter | Type   | Description |
|-----------|--------|-------------|
| user      | String | id of the user |
| group     | String | id of the group |
| (optonal) last_message | String | id of the string to search from. If not given, the last_message will be the last message in the group |
| (optional) limit | Int | The number of messages (default 20) (maximum 100) |

###### Response data
| Parameter    | Type   | Description |
|--------------|--------|-------------|
| num_messages | Int    | number of the messages |
| messages     | []Message | list of the messages |
| prev_message | String | id of the latest message that comes before the returned messages (can be used for pagination) |
| timestamp    | Int    | the timestamp of the request time |

#### type: "get_unread_message"
Get the message from the group that are still not read.
###### Request data
| Parameter | Type   | Description |
|-----------|--------|-------------|
| user      | String | id of the user |
| group     | String | id of the group |

###### Response data
| Parameter    | Type   | Description |
|--------------|--------|-------------|
| num_messages | Int    | number of the messages |
| messages     | []Message | list of the messages |
| timestamp    | Int    | the timestamp of the request time |

### /broadcast_tx_*
#### type: "create_new_group"
Create a new group.
###### Request data
| Parameter | Type   | Description |
|-----------|--------|-------------|
| user      | String | id of the user |
| group     | String | id of the group (must not exist) |

#### type: "join_group"
Join an existing group.
###### Request data
| Parameter | Type   | Description |
|-----------|--------|-------------|
| user      | String | id of the user |
| group     | String | id of the group |

#### type: "leave_group"
Leave an existing group.
###### Request data
| Parameter | Type   | Description |
|-----------|--------|-------------|
| user      | String | id of the user |
| group     | String | id of the group |

#### type: "read_message"
Acknowledge the last message you received.
###### Request data
| Parameter | Type   | Description |
|-----------|--------|-------------|
| user      | String | id of the user |
| group     | String | id of the group |
| timestamp | Int    | timestamp of the last message |

#### type: "send_message"
Send a message to the chat group
###### Request data
| Parameter | Type   | Description |
|-----------|--------|-------------|
| user      | String | id of the user |
| group     | String | id of the group |
| message   | String | message |
| time      | String | time of the message |


