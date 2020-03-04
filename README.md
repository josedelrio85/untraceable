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

kubectl -n bysidecar-pre apply -f ci/job-definition.yml

## Lista jobs

kubectl -n bysidecar-pre get jobs

## Describe job

kubectl -n bysidecar-pre describe jobs/untraceable-job

## Logs job

kubectl -n bysidecar-pre logs jobs/untraceable-job

## Delete jobs

kubectl -n bysidecar-pre delete job/untraceable-job