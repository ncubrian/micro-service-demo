CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o transService main.go
scp transService root@192.168.99.11:/opt/service/
scp transService root@192.168.99.22:/opt/service/
scp transService root@192.168.99.33:/opt/service/
rm ./transService
