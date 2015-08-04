# Software Aware Load Balancing

This is a simple Golang based demo which demonstrates how software can be intelligent enough to do its own load balancing, easing the work load on system administrations and engineers.

In an ideal world, developers would have the time to ensure their systems are robust enough to deal with remote endpoints being unavailable, either failing gracefully or re-trying at appropriate intervals.

## Architecture

You have a remote API which serves up nothing but its instance name - it's clearly the greatest API to hit the Internet. You want to be able to launch multiple instances of this API and have the incoming requests load balanced for speed and reliability. Here are the common ways to do this.

### Load Balancer

```
HAProxy -> [API, API, ...]
```

The problem with this approach is HAProxy is a single point of failure - if it goes down, nothing is available.

### Two Load Balancers

```
[HAProxy | HAProxy] -> [API, API, API, ...]
```

(I'm making up syntax/symbols as I go along here; "," means each item is independant; "|" means the items on each side are dependent on each other, like in a cluster)

Much better, but this means you need to cluster the HAProxy instances and have a floating IP. There are situation were the latter isn't possible, stability of the cluster isn't performant, or you simply don't want the complexity of managing a cluster (I know I don't.)

### Two Load Balancers, No Cluster

```
[HAProxy, HAProxy] -> [API, API, API, ...]
```

Now we have two HAProxy instances which can LB requests between our pool of API servers. The problem here, however, is you now have two public IPs (and thus two public DNS records) for each HAProxy instance, but the clients talking to the API are probably only configured to talk to one... well that's the problem.

## A Solution

Have your client's be aware of both endpoints and programmatically switch between them as their availability comes and goes (which hopefully won't happen):

```python
endpoints = ["api01", "api02"]
use_endpoint = None

for endpoint in endpoints:
  if check_endpoint_health(endpoint):
    use_endpoint = endpoint
    break

if not use_endpoint:
  headless_chicken_mode()

# ...
```

This is very basic code, of course, but the idea is simple:

1. Configure you client with multiple endpoints;
1. Have the client iterate over the endpoints either on start up or per-request, or some other algorithm;
1. If an endpoint isn't available, move on to the next one;
1. If none are available, wait a bit, don't just die like a drama queen;

Such a simple solution resulting in a very stable experience for the end user, not to mention an even simpler infrastructure for the engineering team to design and the administration team to maintain.

## This Repository

This repository contains Go code which can be used to launch a few local APIs and run a multi-endpoint aware client to talk to them. Try running multiple API instances (which just need a configuration file) and then launching a correctly configured client. What happens if you kill an API and use the client again? it should skip the "broken" API endpoint. #win

### Example Session

Let's start the API servers:

```
$ api configs/api_01.json
Starting API with the following configuration...
	BindPort: 0.0.0.0
	BindIP: 8081
	APIName: API01
2015/08/04 15:12:42.878093 Starting Goji on [::]:8081
```

```
$ api configs/api_02.json
Starting API with the following configuration...
	BindPort: 0.0.0.0
	BindIP: 8082
	APIName: API02
2015/08/04 15:16:43.398193 Starting Goji on [::]:8082
```

Let's flip their health statuses to ensure they're showing as "OK":

```
$ curl localhost:8081/health/flip
{"Result":true,"Message":"OK"}

$ curl localhost:8082/health/flip
{"Result":true,"Message":"OK"}
```

And are they OK?

```
$ curl localhost:8082/health
{"Result":true,"Message":"OK"}

$ curl localhost:8081/health
{"Result":true,"Message":"OK"}
```

But what does the client think:

```
$ client configs/client.json
Configuration loaded.
Configured Endpoint: 127.0.0.1, Port: 8081, Name: Remote API 01
Configured Endpoint: 127.0.0.1, Port: 8082, Name: Remote API 02
Using endpoint: Remote API 01
Response Message (Result): This is API: API01 (true)
```

And if we turn "off" `API01`?

```
$ curl localhost:8081/health/flip
{"Result":false,"Message":"OK"}

$ client configs/client.json
Configuration loaded.
Configured Endpoint: 127.0.0.1, Port: 8081, Name: Remote API 01
Configured Endpoint: 127.0.0.1, Port: 8082, Name: Remote API 02
Using endpoint: Remote API 02
Response Message (Result): This is API: API02 (true)
```

The client automatically detected that the first endpoint in its configuration is dead, so it used the second one. What if both go down?

```
$ curl localhost:8082/health/flip
{"Result":false,"Message":"OK"}

$ client configs/client.json
Configuration loaded.
Configured Endpoint: 127.0.0.1, Port: 8081, Name: Remote API 01
Configured Endpoint: 127.0.0.1, Port: 8082, Name: Remote API 02
No endpoints to work with. Exiting.
```

It correctly detects that no endpoints are available and exits gracefully. We can exit gracefully here because this is a demo, but in the real world, you would want your application to hang around and periodically check if the endpoints have come back again - that's the sensible thing to do.

## Author

Michael Crilly
