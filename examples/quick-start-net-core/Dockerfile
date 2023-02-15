FROM mcr.microsoft.com/dotnet/sdk:7.0.102 AS builder

WORKDIR /usr/src/app/
COPY ./src ./

RUN dotnet restore
RUN dotnet publish ./quick-start-net-core.csproj -c Release -o /tracetest

FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS runtime
WORKDIR /usr/src/app/
COPY --from=builder /tracetest .

EXPOSE 80
ENTRYPOINT ["dotnet", "./quick-start-net-core.dll"]

