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
>Note: Not all endpoints require the AdminPassword header, below it will specify

##### Body
```json
{
    "field" : value
}
```
#### Responses

##### Body
 
```json
{
    "data" : {},
    "error" : {},
    "success" : false
}
```

> Note: the data and error could be any number of objects, they are stored as `interface{}` types in go

### GET's

#### ListGames

Used to get all games with specified fields for each game

##### Endpoint

`/listGames`

##### Request

###### Body

```json
{
    "fields": ["field1", "field2", "field3"]
}
```

>Note: Possible fields include `name`, `gameID`, `startNode`, `timeLimit`, `teams`, `status`, and `startTime`

##### Response

###### Body

```json
{
    "data" : [
        {
            "gameID":"gameID",
            "name":"gameName",
            "startNode":"startNodeID",
            "timeLimit":100000,
            "status":"in-progress",
            "startTime":1556987407,
            "teams": [
                {
                    "teamID":"teamID",
                    "name":"teamName",
                    "score":0
                }
            ]
        }
    ],
    "error" : null,
    "success" : true
}
```

#### GameInfo

Similar to `ListGames` except only fetches info for the provided `gameID`

##### Endpoint

`/gameInfo?gameID=<gameID>`

##### Request

###### Body

```json
{
    "fields":["field1", "field2", "field3"]
}
```
>Note: Possible fields include `name`, `gameID`, `startNode`, `timeLimit`, `teams`, `status`, and `startTime`

##### Response

###### Body

```json
{
    "data" : 
        {
            "gameID":"gameID",
            "name":"gameName",
            "startNode":"startNodeID",
            "timeLimit":100000,
            "status":"in-progress",
            "startTime":1556987407,
            "teams": [
                {
                    "teamID":"teamID",
                    "name":"teamName",
                    "score":0
                }
            ]
        },
    "error" : null,
    "success" : true
}
```

### POST's

#### CreateGame

Create a game with the provided information, returns the new game's `gameID`

##### Endpoint

`/createGame`

##### Request

###### Header

```
AdminpPassword: <string>
```

###### Body

```json
{
    "name": "game-name",
    "timeLimit": 1000,
    "teams": ["team1", "team2", "team3"]    
}
```

##### Response

###### Body

```json
{
    "data" : {
        "gameID": "gameID"
    },
    "error" : null,
    "success" : true
}
```

#### JoinGame

**WIP**

Returns a websocket for communicating game state and player actions

##### Endpoint

`/joinGame?gameID=<gameID>`

##### Request

###### Body

##### Response

###### Body

```json
{
    "data" : {},
    "error" : {},
    "success" : false
}
```

#### AdminPasswordCheck

##### Endpoint

`/adminCheck`

##### Request

###### Header

```
AdminpPassword: <string>
```

###### Body

```json
{}
```

##### Response

###### Body

```json
{
    "data" : null,
    "error" : null,
    "success" : true
}
```

### DELETE's

#### DeleteGame

##### Endpoint

`/deleteGame?gameID=<gameID>`

##### Request

###### Header

```
AdminpPassword: <string>
```

###### Body

```json
{}
```

##### Response

###### Body

```json
{
    "data" : null,
    "error" : null,
    "success" : true
}
```

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

`games [gameID1key, gameID2key]` // a set of gameIDs, easier when getting all games

`game:gameID gameID string`
`game:gameID name string`
`game:gameID teamIDs []intToJSONString` // I'll just use a json string
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