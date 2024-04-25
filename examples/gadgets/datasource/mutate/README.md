# Mutating fields from a `datasource`

This example shows how an operator can add and mutate fields from a datasource.

### How to compile

```bash
$ go build .
```

### How to run

The compiled binary doesn't need any parameters, just run it with root permissions:

```bash
$ sudo ./mutate
```

In another terminal, open some files

```bash
$ cat /dev/null
```

Those will printed in the gadget's terminal:

```bash
$ sudo ./mutate
...
PID      UID      GID      MNTNS… E… FD      F… MODE    COMM    FNAME          IS_ROOT
1097     108      117      402653  0 7       52 0       systemd /proc/meminfo  0
1097     108      117      402653  0 7       52 0       systemd /proc/meminfo  0
4720     1001     1001     402653  0 175     52 0       ThreadP /proc/4720/sta 0
28941    0        0        402653  0 3       52 0       cat     /etc/ld.so.cac 1
28941    0        0        402653  0 3       52 0       cat     /lib/x86_64-li 1
28941    0        0        402653  0 3       52 0       cat     /usr/lib/local 1
28941    0        0        402653  0 3       0  0       cat     /dev/null      1
...
```
