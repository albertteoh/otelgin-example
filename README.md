# Gin Example with OTEL Instrumentation

## Running

```shell
# In one terminal session
$ cd inventory-frontend && go run main.go

# In another terminal session
$ cd inventory-backend && go run main.go
```


## Verification

1. Run [Jaeger All In One](https://www.jaegertracing.io/docs/latest/getting-started/#all-in-one).
2. From a terminal, hit the inventory-frontend: `curl localhost:8080/inventory`. 
   1. Note: please ignore the curl response, it's a bug that I haven't looked into yet.
3. Open http://localhost:16686/ in a browser.
4. Check the services `inventory-frontend` and `inventory-backend` are present.
5. Choose one of these services and click `Find Traces`.
6. You should see both services as spans in a single trace.

