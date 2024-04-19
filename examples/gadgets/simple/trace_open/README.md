# `trace_open`

This example shows how to run the `trace_open` gadget and print the events it
captures to the terminal in json format.

### How to compile

```bash
$ go build .
```

### How to run

The compiled binary doesn't need any parameters, just run it with root permissions:

```bash
$ sudo ./trace_open
```

In another terminal, execute some processes (as root because others are filtered out)

```bash
$ sudo cat /dev/null
```

Those will printed in the gadget's terminal:

```bash
$ sudo ./trace_open
{"args":"/usr/bin/cat","args_count":2,"args_size":23,"comm":"cat","gid":0,"loginuid":1001,"mntns_id":4026531841,"pid":9999,"ppid":98133,"retval":0,"sessionid":3,"timestamp":17145168898503,"uid":0,"upper_layer":false}
```
