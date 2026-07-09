# Metadata-Ingestion-Service

### Architecture

```
                  REST API

                     │

             Metadata Service

                     │

       +-------------+--------------+

       |                            |

Power BI Connector           Tableau Connector

       |                            |

       +-------------+--------------+

                     │

              Worker Pool

                     │

                 OpenSearch
```