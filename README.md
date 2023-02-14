# DogeFuzz inputs

This project will use Vandal to collect data about the smart contracts that will be used in benchmark of Dogefuzz

# Executing

First, initialize the Vandal API

```
docker compose -f ./infra/docker-compose.yml up -d
```

And, to execute the project, run the following command:

```
go run ./cmd/inputs
```

# Smart Contracts

The contracts in scope are located in [Drive](https://drive.google.com/drive/folders/1407C0Nnkpf6dadKKmCjiRgsZc10kVO4j?usp=share_link)
