# Untraceable Job

This component is responsible of send sms (using LLeidanet API services) to unreachable clients. To achieve this, is needed to retrieve those clients, analyze if an sms was already sended in a period of 30 days, and it this criteria is satisfied, send it using LLeidanet API service.

## How to run the component

This component has been developed using the following Go! version:

```bash
go version go1.13.4 windows/amd64
```

```bash
go run cmd/main.go
```

## How to run the tests

```bash
go test ./...
```

## Build and run dockerfile

You will need a database to use this component, not included yet in dockerfile

```bash
docker image build -t untraceable:[version] .
docker container run -it -d --name untraceable -e DB_PORT -e DB_HOST -e DB_USER -e DB_PASS -e DB_NAME untraceable:[version]
```

## Kube Jobs

You can't rerun a job, you have to delete it and create it again

## Apply

```bash
kubectl -n bysidecar-pre apply -f ci/job-definition.yml
```

## Lista jobs/cronjob

```bash
kubectl -n bysidecar-pre get jobs | cronjob
```

## Describe job

```bash
kubectl -n bysidecar-pre describe job/untraceable-job
kubectl -n bysidecar-pre describe cronjob/untraceable-cronjob
```

## Logs job

A job creates one or more pods and ensures that a specified number of them successfully terminate.

All you need is to view logs for a pod that was created for the job.

```bash
kubectl -n bysidecar-pre get pods
# view the pod(s) generated by the job | cronjob
kubectl -n bysidecar-pre logs [pod_name]
```

## Delete jobs

```bash
kubectl -n bysidecar-pre delete job/untraceable-job
kubectl -n bysidecar-pre delete cronjob/untraceable-cronjob
```
