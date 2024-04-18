docker network create docker_network

docker run --network=docker_network -e "ACCEPT_EULA=Y" -e "MSSQL_SA_PASSWORD=StrongP4ssword!" -p 1533:1433 --name sqlcontainer --hostname sqlcontainer -d mcr.microsoft.com/mssql/server:2022-latest

docker build -t e3-server .

docker run --network=docker_network -p 8080:8080 --name e3-server e3-server