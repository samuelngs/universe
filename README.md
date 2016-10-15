# universe
Distributed SSH authentication system

### Example

#### Server

```
$  go run main.go
2016/10/14 20:56:03 subscribe => {"message":"Starting server","topic":"server-start"}
2016/10/14 20:56:03 subscribe => {"message":"Listening on 127.0.0.1:2222","topic":"server-started"}
2016/10/14 20:56:07 RemoteAddr: 127.0.0.1:45608
2016/10/14 20:56:07 logging => {"message":"New connection from 127.0.0.1:45608 (SSH-2.0-OpenSSH_7.2p2 Ubuntu-4ubuntu2.1)","topic":"connect"}
2016/10/14 20:56:07 logging => {"message":"Pty initialized","topic":"channel"}
2016/10/14 20:56:07 logging => {"message":"Pty request","topic":"channel"}
```

#### Client

```
$ ssh localhost -p 2222

>>> uname -a
Linux 4.4.0-38-generic #57-Ubuntu SMP Tue Sep 6 15:42:33 UTC 2016 x86_64 x86_64 x86_64 GNU/Linux
```
