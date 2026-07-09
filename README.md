# Metadata-Ingestion-Service

### V1 Architecture

```
                         Metadata Ingestion Service

                   +-------------------------------+
                   |            cmd/app            |
                   +---------------+---------------+
                                   |
                                   |
                         Ingestion Service
                                   |
                  +----------------+----------------+
                  |                                 |
           Connector Manager                 Worker Pool
                  |                                 |
        +---------+---------+                Metadata Channel
        |         |         |                       |
     PowerBI   Tableau    MLflow                    |
                  |                                 |
                  +---------------+-----------------+
                                  |
                           Metadata Processor
                                  |
                           OpenSearch Sink
```