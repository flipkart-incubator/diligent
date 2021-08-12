# What is Diligent?

_diligent | adjective. : characterized by steady, earnest, and energetic effort : painstaking_

Diligent is a tool we created at Flipkart for generating workloads for our SQL databases that enables us to answer specific questions about the performance of a database.

There are several benchmark definitions and benchmarking tools out there such as YCSB, TPCC, HammerDB and more. While these are useful, we had some specific needs at Flipkart for which we needed a good tool:

- Ability to simulate the workload of a prototypical application which uses secondary indexes and transactions, and answer questions such as: "What is the difference in read latency when we lookup by primary key vs a unique secondary key?"
- Ability to scale the load generator horizontally to generate more load and simulate the workload of a horizontally scaled application with many nodes
- Ability to generate load for large datasets (order of TBs)
- Ability observe the throughput and latency graphs - not just see the summarised stats for a run

Diligent was created to address these needs and more.

Currently Diligent works with MySQL compatible databases only. We have plans to support PostgreSQL compatible databases in the near future.

To get started with using Diligent, please refer to the documentation in the accompanying [Wiki](https://github.com/flipkart-incubator/diligent/wiki)

Hope you find Diligent useful!

