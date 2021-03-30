# Kubeedge Router Manager Bench

A stress testing tool for kubedge router. 

## Usage

**build:**

`go build`

**All commands:**

```bash
Usage of ./kubeedge-router-manager-bench:
  -X string
        Specify request command to use, accept: GET/POST/PUT/DELETE. (default "POST")
  -b string
        MQTT broker to edge. Ex: tcp://127.0.0.1:1883 (default "tcp://127.0.0.1:1883")
  -c int
        Maximum concurrency for sending. (default 1)
  -concurrency int
        Maximum concurrency for sending. (default 1)
  -e string
        MQTT topic to kubeedge eventbus.
  -eventbus string
        MQTT topic to kubeedge eventbus.
  -m string
        Message to send.
  -mb string
        Binary message to send.
  -mqtt-broker string
        MQTT broker to edge. Ex: tcp://127.0.0.1:1883 (default "tcp://127.0.0.1:1883")
  -msg string
        Message to send.
  -msg-binary string
        Binary message to send.
  -n int
        Number of meseages to send. (default 1)
  -num int
        Number of meseages to send. (default 1)
  -r string
        Restful path to kubeedge cloudcore. Ex: 
            're' mode: http://127.0.0.1:9443/centos78-edge-0/default/a
            'er' mode: 127.0.0.1:8080/test-1 
  -rest string
        Restful path to kubeedge cloudcore. Ex: 
            're' mode: http://127.0.0.1:9443/centos78-edge-0/default/a
            'er' mode: 127.0.0.1:8080/test-1 
  -t string
        Type of banch. Accept: 're' or 'er', means 'rest -> eventbus' or 'eventbus -> rest'. (default "re")
  -timeout int
        The maximum time in seconds to wait for the next message (default 300)
  -type string
        Type of banch. Accept: 're' or 'er', means 'rest -> eventbus' or 'eventbus -> rest'. (default "re")
```

## Example



**rest -> eventbus:**

1. Create rule:

    ``` yaml
    apiVersion: rules.kubeedge.io/v1
    kind: Rule
    metadata:
    name: my-rule
    labels:
        description: test
    spec:
    source: "my-rest"
    sourceResource: {"path":"/a"}
    target: "my-eventbus"
    targetResource: {"topic":"test"}
    ```

2. Stress testing:

    ```bash
    # Run on Edge
    # $ ./kubeedge-router-manager-bench -t re [-b MQTT_BROKER_URL] -e EVENTBUS_TOPIC [-n EXPECT_NUM] [-timeout TIMEOUT_IN_SECONDS]
    $ ./kubeedge-router-manager-bench -t re -b tcp://127.0.0.1:1883 -e test -n 100

    # Run on Cloud
    # $ ./kubeedge-router-manager-bench -t re -r REST_URL -m MESSAGE [-mb BINARY_MESSAGE] [-n EXPECT_NUM] [-c CONCURRENCY_NUM]
    $ ./kubeedge-router-manager-bench -t re -r http://127.0.0.1:9443/centos78-edge-0/default/a -m "hello" -n 100
    ```

**eventbus -> rest:**

1. Create rule:

    ```yaml
    apiVersion: rules.kubeedge.io/v1
    kind: Rule
    metadata:
      name: my-rule-eventbus-rest
      labels:
        description: test
    spec:
      source: "my-eventbus-1"
      sourceResource: {"topic": "test-1","node_name": "centos78-edge-0"}
      target: "my-rest-1"
      targetResource: {"resource":"http://127.0.0.1:8080/test-1"}
    ```

2. Stress testing:

    ```bash
    # Run on cloud
    # $ ./kubeedge-router-manager-bench -t er -r REST_URL_TO_CREATE [-n EXPECT_NUM] [-timeout TIMEOUT_IN_SECONDS]
    $ ./kubeedge-router-manager-bench -t er -r 127.0.0.1:8080/test-1 -n 100

    # Run on edge
    # $ ./kubeedge-router-manager-bench -t er [-b MQTT_BROKER_URL] -e EVENTBUS_TOPIC -m MESSAGE [-mb BINARY_MESSAGE] [-n EXPECT_NUM] [-c CONCURRENCY_NUM]
    $ ./kubeedge-router-manager-bench -t er -b tcp://127.0.0.1:1883 -e "default/test-1" -m "hello" -n 100
    ```

## Result

+ kubeedge v1.6.0
  - env:
    - cloud: 4c16g
    - edge: 1c1g
  - result:
    | type | node  | msg size | num  | time            | remark  |
    | ---- | ----- | -------- | ---- | --------------- |-------- |
    | re   | cloud | 100      | 1000 | 372.519152ms    | send    |
    |      | edge  |          | 1000 | 746.333895ms    | receive |
    | er   | edge  | 100      | 1000 | 1.288281784s    | send    |
    |      | cloud |          | 1000 | 1.286574437s    | receive |
    | re   | cloud | 2.2M     | 1000 | 8.772269883s    | send    |
    |      | edge  |          | 1000 | 1m23.201584146s | receive |
    | er   | edge  | 2.2M     | 100  | 971.079555ms    | send    |
    |      | cloud |          | 100  | 4.786262004s    | receive |
