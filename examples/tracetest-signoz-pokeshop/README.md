# Tracetest + SigNoz + Pokeshop

<!-- > [Read the detailed recipe for setting up Tracetest + SigNoz + Pokeshop in our documentation.]() -->

This examples' objective is to show how you can:

1. Configure Tracetest to ingest traces.
2. Configure SigNoz and use it as a trace data store.
3. Configure SigNoz to query traces.

## Steps

1. [Install the tracetest CLI.](https://docs.tracetest.io/installing/)
2. Run `tracetest configure --endpoint http://localhost:11633` on a terminal.
3. Run the project by using docker-compose: `docker-compose up -d` (Linux) or `docker compose up -d` (Mac)
4. Test if it works by running: `tracetest run test -f tests/test.yaml`. This would trigger a test that will send trace spans to both the SigNoz instance and Tracetest instance that are running on your machine. View the test on `http://localhost:11633`.
5. View traces in SigNoz on `http://localhost:3301`. Use this query:

    ```yaml
    http://localhost:3301/trace?selected={%22operation%22:[%22POST%20/pokemon/import%22]}&filterToFetchData=[%22duration%22,%22status%22,%22serviceName%22,%22operation%22]&spanAggregateCurrentPage=1&selectedTags=[]&&isFilterExclude={%22operation%22:false}&userSelectedFilter={%22serviceName%22:[%22pokeshop%22,%22tracetest%22],%22status%22:[%22error%22,%22ok%22],%22operation%22:[%22POST%20/pokemon/import%22]}&spanAggregateCurrentPage=1&spanAggregateOrder=&spanAggregateCurrentPageSize=10&spanAggregateOrderParam=
    ```
