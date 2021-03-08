# Balancer

A simple load balancer implemented for learning purposes

## Example of usage

In order to start the balancer on port 3003 and distribute load between backends ht<span>tp://</span>localhost:50001 and ht<span>tp://</span>localhost:50002, issue the following command.

```
./balancer --port=3003 --backends="http://localhost:50001,http://localhost:50002"
```