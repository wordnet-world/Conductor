# Conductor
The server that handles client requests and communicates with databases and kafka

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
`game_id = <number>`
`game:game_id field value`

`game:game_id game_id string`
`game:game_id name string`
`game:game_id teams string` // I'll just use a json string
`game:game_id time_limit string`
`game:game_id start_node string`
`game:game_id status string` // waiting, in-progress, complete
`game:game_id start_time string` // 0 if not in progress

Teams
`team:game_id:team_id field value`

`team:game_id:team_id team_id string`
`team:game_id:team_id name string`
`team:game_id:team_id score int`

