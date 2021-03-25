# Kubeedge Router Manager Bench

## Usage

**build**

`go build`

**rest -> eventbus:**

```bash
# Run on Edge
# $ ./kubeedge-router-manager-bench -t re [-b MQTT_URL] -e EVENTBUS_TOPIC [-n EXPECT_NUM]
$ ./kubeedge-router-manager-bench -t re -b tcp://127.0.0.1:1883 -e test -n 100

# Run on Cloud
# $ ./kubeedge-router-manager-bench -t re -r REST_URL -m MESSAGE [-n EXPECT_NUM] [-c CONCURRENCY_NUM]
$ ./kubeedge-router-manager-bench -t re -r http://127.0.0.1:9443/centos78-edge-0/default/a -m "hello" -n 100 -c 3
```

**eventbus -> rest:**

```bash
# Run on cloud
# $ ./kubeedge-router-manager-bench -t er -r REST_URL_TO_CREATE [-n EXPECT_NUM]
$ ./kubeedge-router-manager-bench -t er -r 127.0.0.1:8080/test-1 -n 100

# Run on edge
# $ ./kubeedge-router-manager-bench -t er [-b MQTT_URL] -e EVENTBUS_TOPIC [-n EXPECT_NUM]
$ ./kubeedge-router-manager-bench -t er -b tcp://127.0.0.1:1883 -e "default/test-1" -m "hello" -n 100 -c 3
```

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
  -type string
    	Type of banch. Accept: 're' or 'er', means 'rest -> eventbus' or 'eventbus -> rest'. (default "re")
```