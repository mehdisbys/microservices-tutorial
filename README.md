##### Proposed Solution

Dependencies :

- docker
- docker-compose
- golangci
- go `1.12` or above

Install : 

- fetch the branch locally : `git fetch origin solution-test:solution_test`
- checkout to the branch : `git checkout solution-test`
- `make deps`
- `make test` (this will create a local redis instance for test purposes)
- `make all`
- `make deploy`
 
What was deployed ?

By entering the command `make deploy` the following things happened :

- A private and public network were created
- Containers for redis, nsq, gateway, driver-location and zombie driver services were placed inside the private network
- Gateway container was also placed on the public network
- Gateway container listens on localhost:9000

Only the Gateway is exposed all other services and dependencies are private


How to test :

- make one or multiple curl requests to the gateway :
    ``` 
       curl -X PATCH \
       http://localhost:9000/drivers/6/locations \
       -H 'Content-Type: application/json' \
       -d '{
       "latitude": 28.864193,
       "longitude": 9.550498
     }'
     ```
     
- you should get a 200 back

- logs will be created with a traceID to be able to follow the request across the queue and the driver-location services

- make a another request to the gateway to check whether a driver is a zombie :

```
curl -X GET http://localhost:9000/drivers/6?minutes=15
```

- a response similar to this one should be returned :

```
{"id":6,"zombie":false}
```

### Solution description

#### Approach

- All three services do not know about each other, they only the hold the name of the service they have a dependency on. This is essential
for added flexibility as services often get assigned network addresses dynamically.
The only thing that services are sharing are message formats and payloads. 

- Inside each service there is an emphasis on using interfaces, this allows us to formally define dependencies and could allow dependencies change 
(e.g using dynamoDB instead of Redis) with minimal code update. It also makes testing much easier.

- Each service comes with its tests for the most essential scenarios.

- Fail early : if there is missing configuration at startup time the service will fail and log an error message, this is preferred to failing while handling a request.

- Pass on traceIDs to track a request going through different services

- Structured logging : logging in JSON format enables to do log analysis / querying and much more with appropriate tools


#### Gateway Service

It dynamically creates endpoint based on `config.yaml`, the following scenarios have tests :

- It can create asynchronous handlers (forward message to a queue)
- It can create synchronous handlers (proxy the request to another service and return response)
- Dynamically created asynchronous handler pushes to the topic specified in topic
- Dynamically created synchronous handler proxies request and returns response and status code
- Forwards messages or proxies requests only for matching requests


In a nutshell this service listen for requests and publishes them to the queue or proxies them to another service when appropriate.
It also assigns a traceID if none where passed in the request, the way to pass it in the request is to add the header `X-Trace-Id` with a uuid-like value.

It adds url parameters and traceID inside the message for downstream use and logging

##### How to improve it

- Add an authentication layer
- Add metrics
- Add an endpoint to retrieve all the supported different dynamically created endpoints
- Add circuit breakers when the queue or the zombie service are not responding correctly
- Add integration tests with a live queue
- Add integration tests with a zombie and driver service


### Driver Service

Driver service's responbility is to store and allow to retrieve drivers locations.

This service has two handlers :
- a queue handlers that listen on the same queue as the gateway service publishes messages to and stores them to a Redis database. 
- a HTTP handler that enables to retrieves locations for a given driver and for the last x minutes.
- it is designed so that queue or database implementation can easily be switched

The following scenarios have tests :

Integration test with a redis instance : 
 - store and fetch locations
 - ensure a timestamp is assigned when saving to db
 - ensure locations are sorted by date
 - ensure only the locations within the `minutes` argument are returned
 
- It can unmarshal and save a message to a mock db
- It returns an error when a driverID is missing in message
- It returns an error when the message envelope is malformed 
- It returns an error when the message inside the envelope is invalid

##### How to improve it

- Add ability to process messages in batches as opposed to individual message (slice of locations)
- Add healthcheck
- Add metrics (count locations saved, time to fetch locations, etc..)
- Add endpoint to retrieve locations for multiple drivers at once
- Generate an event when location is saved for other services uses and for datawarehouse / BI


### Zombie Service

Zombie service's responsibility is to fetch locations for a driver and determine whether he is a zombie in function of two parameters : distance and time

This service will query the driver service which is therefore a dependency.

- The chosen way to calculate a distance between two coordinates is the Haversine method. 
- The code is written in a way that allows to easily swap it for another method of calculation through the `DistanceEstimator` interface, 
it can even be a third party dependency (e.g google maps).
- Locations can easily be fetched from somewhere else than the driver service which does not need to have a http interface
- Locations fetching, distance calculation and zombie determination are decoupled in the code

The following scenarios have tests : 
- tests for Haversine method
- tests to determine if driver is a zombie

#### How to improve it 
- Add healthcheck
- Add metrics (time to serve, non-success http statuses etc..)
- Add ability to determine if drivers are zombies in batch
- Generate an event each time a driver is changing state (zombie/ not zombie) for other services use


### Future Considerations
There are many ways to make this project better and while I have enjoyed developing this solution I thought to stop myself once an acceptable and clean solution was reached.

Given more time these are the things I would do :

- Add a client for driver service so that zombie service can just use it without having to know about what endpoint to request
- Deploy on kubernetes and make the services highly available and scalable
- Integration testing of all three services, which probably would become its own service