# openssl-goproxy
```
sudo sh -c 'LD_PRELOAD="./openssl-source/dist/lib/libssl.so ./openssl-source/dist/lib/libcrypto.so" ./app -listen :80'
curl --resolve megadomain.vnn.vn:80:127.0.0.1 http://megadomain.vnn.vn -v
```
