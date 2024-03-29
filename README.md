[![GoDoc](https://godoc.org/github.com/ncraft/nfs-api?status.svg)](http://godoc.org/github.com/ncraft/nfs-api)

# API to manage NFS shares

## Description

Provides a tiny REST API via unix domain socket which can be used to add NFS shares (`/etc/exports` on most linux systems).

The changes to provide remote access via TCP, instead of unix domain sockets, should be really small but that's beyond the scope of this project. Probably, to implement scenarios with remote access it would be more appropriate to call this service from a wep app which runs on the same node but does not require privileges to manipulate `/etc/exports`, thus splitting privileged and unprivileged parts over different services.


## Usage

Example JSON request body to add an NFS share for IPs `192.168.1.110` and `192.168.1.112` to directory `/var/nfs/pictures`:
```json
{
  "sharePath": "/var/nfs/pictures",
  "exportOptions": {
    "clients": [
      "192.168.1.110",
      "192.168.1.112"
    ],
    "rw": true,
    "sync": true,
    "noSubtreeCheck": true
  },
  "mkDir": true,
  "dirOwnerUid": 33,
  "dirOwnerGid": 33
}
```

Use `curl` or a simple Go [program](https://gist.github.com/teknoraver/5ffacb8757330715bcbcc90e6d46ac74) to post HTTP requests against the unix domain socket:
```bash
$ ./unix-socket-client -d '{"sharePath": "/var/nfs/pictures", "exportOptions": {"clients": ["192.168.1.110", "192.168.1.112"], "rw": true, "sync": true, "noSubtreeCheck": true}, "mkDir":true, "dirOwnerUid":33, "dirOwnerGid":33}' /path/to/socket /shares

Unix HTTP client
{"status":200,"message":"SHARED ADDED"}
```
