## service user discussion items

---

### Error Handling:

The service efficiently handles errors at majorly the controller layer i.e the route handler. Errors are propagated & thrown at every layer like service, repo level or in future any middleware are majorly handled in controller layer. A very ideal case for middleware library/service will be the regex validator of path param or the input json validator both of which are made part of controller layer as of now to save time. The api's return pre-determined enum of errors so that clients build handlers at their end for specific functionality like retries and proper error message to user.

---

### Testing:

The service controller test is written to demonstrate nice and clean unit testing approach. There are more tests that needs to added for each handler function in the controller unit tests.

The repository is having an additional integration test to validate actual data CRUD simulations using test containers.

---

### Scalability Considerations:

##### caching layer

    The service is expected to handle 30k requests per minute with a latency of < 100ms. To reduce request latency, I have already put in place redis cache layer on the endpoints. As of now the controller is taking care of CRUD from the Redis cache with a default TTL using write-through strategy or proactive strategy where in after every update/write operation I am updating the cache.

##### database connection pooling

    I have included max open connections and max idle connections in gorm configuration to manage and respect number of open connections that mysql database would handle simultanously. The parameters defined for max and idle connection I would determine depending on the kind of infra we would provision for the database.

##### database indexes

    For the dummy service that I created and the simplistic approach I have taken to query the database, gorm would go ahead and create an Index on the ID column which is the primary key of the table besides a few others like email and its default deleted_at. As the functionalities grow and there is need of querying data with more columns, to improve efficiency of queries I would add more indexes as required on the gorm model definition. This will ensure that the required indices are created when gorm autoMigrates the model/table.

##### database read replicas

    I would expect a service user to be more optimally required to query specific user data. For this, I would prioritize setting up read replicas. Write operations will happen on master database and read will be distributed accoss multiple read replicas. This I would do to reduce stress on database.

##### autoscaling

    Assuming that our service is deployed on a k8s cluster, depending on the iac tooling where we define our hpa strategy, we specify a min and max replica for the service. We keep the strategy of hpa to scale up ig cpu utilization is 80% or memory utilization is 80%.

---

### Security:

The input(s) provided to exposed endpoints is validated before any kind of processing begins. This is taken care of at controller layer. Ideally I would write a middleware validator service as a gin Handler function injected at route level before the request arrives at controller. I have used regex pattern to validate path param and validator library from go-playground for input json validation. gorm decorators take care of validation at repo level for CRUD operations.

---

### Roadmap of things to be done:

##### CI & CD

I would write gitactions workflow file for testing, lint validation and building & deploying service to k8s cluster.

##### config/secrets management

I would utilize AWS secrets replacing the env variable injection for Database, Redis and other configurations.
Few things like default TTL on redis cache, sql max open connection, max idle connections can be from env variables or I will incorporate something like a viper to manage all this.

##### middleware services

I would seperate input validations in a seperate gin HandlerFunc middleware service to cut controller noise.

##### observability & health endpoints

I would integrate with observability tools like New Relic or data dog using their sdks. I will add healtz route to service to monitor it as well.

##### authentication and authorization checks

I would not want my service to take care of auth checks besides the functionality of say checking a few security headers that a Gateway service would add after it completes authentication and authorization checks. The checking of security headers etc would be done in a middleware service.

---
