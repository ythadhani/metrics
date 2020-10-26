# metrics

This repository contains:

* A Go application that generates metrics by probing external URLs. It is instrumented using the Prometheus Golang Client.
* Dockerfile to build a Docker image containing the Go binary.
* Kubernetes deployment specification and other related files to deploy the docker image in a Kubernetes cluster.
* A project report containing all the steps to be followed along with relevant screenshots.

### Getting Started

* Building the Go program to run locally:

    ```sh
    make build
    ```
  
* Running the Go binary:

    ```sh
    ./metrics
    ```
  
* View Prometheus scrape targets:
 
     ```sh
     http://localhost:9090/metrics
     ```
   
* Building a Docker image:

     ```sh
     make image
     ```
  
* Deploying on a Kubernetes cluster:

     ```sh
     make deploy
     ```
  
* Delete deployment from cluster:

     ```sh
     make undeploy
     ```