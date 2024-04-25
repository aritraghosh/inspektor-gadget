# `trace_dns`

This example shows how to run the `trace_dns` gadget and print the events it
captures to the terminal in json format.

### How to compile

```bash
$ go build .
```

### How to run

The compiled binary doesn't need any parameters, just run it with root permissions:

```bash
$ sudo ./trace_dns
```

In another terminal, perform some DNS request within a container. Requests on
the host aren't traced by the example.

```bash
$ docker run --name c3 --rm -it busybox sh -c "nslookup inspektor-gadget.io"
Server:         190.248.0.7
Address:        190.248.0.7:53

Non-authoritative answer:
Name:   inspektor-gadget.io
Address: 172.67.166.105
Name:   inspektor-gadget.io
Address: 104.21.11.160

Non-authoritative answer:
Name:   inspektor-gadget.io
Address: 2606:4700:3030::6815:ba0
Name:   inspektor-gadget.io
Address: 2606:4700:3037::ac43:a669
```

Those will printed in the gadget's terminal:

TODO: check why this gadget is lacking so much information compared to the built-in one.

```bash
$ sudo ./trace_dns
{"anaddrcount":0,"ancount":0,"gid":0,"id":38647,"latency_ns":0,"mntns_id":4026533917,"name":"\u0010inspektor-gadget\u0002io","netns":4026534560,"pid":52845,"pkt_type":4,"qr":0,"qtype":1,"rcode":0,"runtime":{"containerName":"c3"},"task":"nslookup","tid":52845,"timestamp":10931213449822,"uid":0}
{"anaddrcount":0,"ancount":0,"gid":0,"id":51941,"latency_ns":0,"mntns_id":4026533917,"name":"\u0010inspektor-gadget\u0002io","netns":4026534560,"pid":52845,"pkt_type":4,"qr":0,"qtype":28,"rcode":0,"runtime":{"containerName":"c3"},"task":"nslookup","tid":52845,"timestamp":10931213495608,"uid":0}
{"anaddrcount":1,"ancount":2,"gid":0,"id":38647,"latency_ns":131283419,"mntns_id":4026533917,"name":"\u0010inspektor-gadget\u0002io","netns":4026534560,"pid":52845,"pkt_type":0,"qr":1,"qtype":1,"rcode":0,"runtime":{"containerName":"c3"},"task":"nslookup","tid":52845,"timestamp":10931344733241,"uid":0}
{"anaddrcount":1,"ancount":2,"gid":0,"id":51941,"latency_ns":896357186,"mntns_id":4026533917,"name":"\u0010inspektor-gadget\u0002io","netns":4026534560,"pid":52845,"pkt_type":0,"qr":1,"qtype":28,"rcode":0,"runtime":{"containerName":"c3"},"task":"nslookup","tid":52845,"timestamp":10932109852794,"uid":0}
```

> ⚠️ The DNS name isn't shown in the right format. See
> https://github.com/inspektor-gadget/inspektor-gadget/issues/2316.
