# Conductor
The server that handles client requests and communicates with databases and kafka

## API Specification

### Basic Message Structure

#### Requests
##### Headers
```
AdminpPassword: <string>
Content-Type: application/json
```
##### Body
```json
{
    "field" : {<value>}
}
```
#### Responses

##### Body
```json
{
    "data" : {},  // Data and error are interface types in golang, so could be string or map or array, etc
    "error" : {},
    "success" : false
}
```

### GET's

### POST's

#### CreateGame

##### Endpoint

##### Request

###### Headers
###### Body

##### Response

###### Body

#### JoinGame

##### Endpoint

##### Request

###### Headers
###### Body

##### Response

###### Body

#### AdminPasswordCheck

##### Endpoint

##### Request

###### Headers
###### Body

##### Response

###### Body

### DELETE's

## Redis Data Model

There are a handful of things we need to keep track of, especially since we strive to keep the application stateless.

There are two categories of data we want to store in Redis, game metadata and the edge nodes for a team's graph

### Game Meta Data

#### Games

We need to store:
- An ID or lobby name
- Team ids (rather than team names, teams will be under a hash)
- status: waiting, in-progress, complete
- time duration
- start node id
- start time if in-progress

We could potentially use a zset and unix time to show either most recent games or oldest games first
We can also make unique id's by having a count variable in redis and incrementing it, one for each type maybe

#### Teams

We don't really need to keep track of individual players if we use the Kafka well, we can simply track team progress

We need to store:
- Name
- Score (might have game thread increment this stuff)

#### Actual Redis structures
UUID variables
`game:id string` // technically can be int, look at redis STRING type
`team:id string`

Games
`gameID = <number>`
`game:gameID field value`

`game:gameID gameID string`
`game:gameID name string`
`game:gameID teams string` // I'll just use a json string
`game:gameID timeLimit string`
`game:gameID startNode string`
`game:gameID status string` // waiting, in-progress, complete
`game:gameID startTime string` // 0 if not in progress

Teams
`teamID = <number>`
`team:teamID field value`

`team:teamID teamID string`
`team:teamID name string`
`team:teamID score int`
`team:teamID memberCount int`
