# go_api_systemd
Go API server which to do start/stop/restart/status of services in linux

## go run *.go memcached start && systemctl status memcached


curl http://localhost:8080/service \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"serviceName":"db","action":"status"}'



curl http://localhost:8080/service \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"serviceName":"memcached","action":"status"}'

curl http://localhost:8080/service \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"serviceName":"memcached","action":"start"}'

curl http://localhost:8080/service \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"serviceName":"memcached","action":"stop"}'


to do:
- api endpoint for action on single service
- api endpoint for action on multiple services






////////
curl http://localhost:8080/multiservice \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '[{"serviceName":"memcached","action":"stop"},{"serviceName":"memcachedd","action":"status"}]'