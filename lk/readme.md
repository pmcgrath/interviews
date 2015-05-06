# This is far from ready
Should be able to get what is completed running if you have golang on a machine with the following - I have not tested on a clean machine
```bash
# This will get and install the user connections service - Will pull dependencies like echo and uuid into your gopath 
go get github.com/pmcgrath/interviews/lk/ucsvc

# Get and install user connections client - has very specific steps demoing service usage
go get github.com/pmcgrath/interviews/lk/ucsvc_client

# Start service - listens on port 8090 by default
ucsvc

# In other teminal run client - will make calls on the ucsvc service
ucsvc_client
```



# Assumptions
- Just to acknowledge separating the services seems a little odd, normally a service would be responsible for its own storage
- Will use a uuid as the type for the user Id
- Assuming the service usage profile is more reads than writes
- Will use json over http for both services for interop with other stacks - Could use http://golang.org/pkg/net/rpc/ or http://golang.org/pkg/net/rpc/jsonrpc/ for better perf (Would need to do perf measurements based on expected usage)
- Will only support application/json media type
- Does not require TLS support
- Will not support adding a user with connections - This is achieved by adding a user and then adding user connections as seperate service calls
- Will not support creating multiple connections for a user at one time, requirements seem to indicate you create a single connection between 2 users, could allow this but does not seem to be a requirement now
- Will not use a circuitbreaker https://github.com/afex/hystrix-go


# Questions
- Are the user connections unidirectional or bidirectional ?
- Is it the responsibility of the user connections service to assign a user Id or will the client already have assigned an Id ?
- Is it okay to use curl to demo using both services ? Would you prefer a CLI app ?
- Does the user name need to be unique ? If so does the client optimistically submit a user name when creating a user and the service ensures it is unique
- Do I need to deal with authentication\authorisation - Requirements do not indicate if the client is authenticated, seem to indicate the user connections service client can create users and connections for any user
- Can I confirm that the requirements indicate the persistence service may be used in future to support the storage of other services data ? I know you indicate that it need not do so at this time
- Do you want to include conditional gets and to leverage caching headers ? REST features - This assumes json/http
- Do you want me to use a circuit breaker for downstream service calls - user connections service to persistence service - See https://github.com/afex/hystrix-go


# Pending phases
- Get the user connection service up and running - Partially complete
- Add authentication\authorisation using the http Authorization header
- Create user connection interface so we can mock and test the handlers
- Consider the persistence service - This sounds like it needs to mimic memcached or redis interface - simple get\set\del with a key - will have two collections - users and userconnections
- Sort godeps so dependencies are vendored rather than being based on go get
- Raise observed panics within the echo package router - Investagate\raise github issue

