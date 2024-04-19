# Local Manager Operator

This example shows how to use the `LocalManager` operator to filter and enrich
events with container information on the local host.

### How to compile

```bash
$ go build .
```

### How to run

The compiled binary doesn't need any parameters, just run it with root permissions:

```bash
$ sudo ./local_manager
```

The example is configured to only capture events from `mycontainer`. Run the
following commands in a different terminal:

```bash
$ docker run --name mycontainer --rm -it busybox sh -c "cat /dev/null"
$ docker run --name foocontainer --rm -it busybox sh -c "cat /dev/null"
```

The gadget only captured the events from `mycontainer`, and you can see how
they're enriched with the container name:

TODO: why all runtime things aren't displayed by default?

```bash
$ sudo ./local_manager
{"comm":"runc:[2:INIT]","err":0,"fd":4,"flags":524288,"fname":"/proc/self/fd","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855219045348,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/etc/ld.so.cache","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220073203,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/glibc-hwcaps/x86-64-v3/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220082220,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/glibc-hwcaps/x86-64-v2/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220088361,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/tls/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220093601,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220097990,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220103861,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/tls/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220109411,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220114671,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220119019,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220123358,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220127686,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/glibc-hwcaps/x86-64-v3/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220135851,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/glibc-hwcaps/x86-64-v2/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220140480,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/tls/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220145890,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220151180,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220155859,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/tls/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220160217,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220164536,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220168814,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220173172,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220177640,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/glibc-hwcaps/x86-64-v3/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220184022,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/glibc-hwcaps/x86-64-v2/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220188391,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/tls/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220194412,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220198720,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220203018,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/tls/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220207286,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220213548,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220217846,"uid":0}
{"comm":"sh","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220222165,"uid":0}
{"comm":"sh","err":0,"fd":3,"flags":524288,"fname":"/lib/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220237203,"uid":0}
{"comm":"sh","err":0,"fd":3,"flags":524288,"fname":"/lib/libresolv.so.2","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220296806,"uid":0}
{"comm":"sh","err":0,"fd":3,"flags":524288,"fname":"/lib/libc.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220354585,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/etc/ld.so.cache","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220913913,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/glibc-hwcaps/x86-64-v3/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220920145,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/glibc-hwcaps/x86-64-v2/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220925716,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/tls/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220931346,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220936656,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220941866,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/tls/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220947587,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220953027,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220958648,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220964199,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64-linux-gnu/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220970030,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/glibc-hwcaps/x86-64-v3/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220974899,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/glibc-hwcaps/x86-64-v2/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220979267,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/tls/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220984327,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220988685,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220993053,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/tls/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855220997632,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221002271,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221006539,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221010837,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/usr/lib/x86_64-linux-gnu/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221015115,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/glibc-hwcaps/x86-64-v3/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221019443,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/glibc-hwcaps/x86-64-v2/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221023791,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/tls/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221028089,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221032398,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/tls/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221037347,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/tls/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221041725,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221046043,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221050312,"uid":0}
{"comm":"cat","err":2,"fd":0,"flags":524288,"fname":"/lib/x86_64/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221054620,"uid":0}
{"comm":"cat","err":0,"fd":3,"flags":524288,"fname":"/lib/libm.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221062705,"uid":0}
{"comm":"cat","err":0,"fd":3,"flags":524288,"fname":"/lib/libresolv.so.2","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221106678,"uid":0}
{"comm":"cat","err":0,"fd":3,"flags":524288,"fname":"/lib/libc.so.6","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221155060,"uid":0}
{"comm":"cat","err":0,"fd":3,"flags":0,"fname":"/dev/null","gid":0,"mntns_id":4026534471,"mode":0,"pid":136666,"runtime":{"containerName":"mycontainer"},"timestamp":26855221357312,"uid":0}
```
