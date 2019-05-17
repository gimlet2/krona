# krona
This project targets very specific setup: Kubernetes 1.6 with Kubeless and Scheduled Functions. This version of k8s has some issues with CronJobs. And that is temporar adhoc solution.

## Build:
`CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o krona .`
### then:
`docker build .`

