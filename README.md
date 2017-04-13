# qframe-collector-gelf
GELF collector for the qframe framework.


```
$ docker run -ti --rm --log-driver gelf --log-opt gelf-address=udp://<IP_OF_THE_GELF_HOST>:12201\
                      --log-opt gelf-compression-type=none debian:latest hostname
```
