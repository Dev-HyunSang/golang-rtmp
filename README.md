# golang-rtmp

## What is the RTMP?
Real-Time Messaging Protocol (RTMP) is a communication protocol for streaming audio, video, and data over the Internet.   
Originally developed as a proprietary protocol by Macromedia for streaming between Flash Player and the Flash Communication Server, Adobe (which acquired Macromedia) has released an incomplete version of the specification of the protocol for public use.

**The RTMP protocol has multiple variations:**

- RTMP proper, the "plain" protocol which works on top of Transmission Control Protocol (TCP) and uses port number 1935 by default.
- RTMPS, which is RTMP over a Transport Layer Security (TLS/SSL) connection.
- RTMPE, which is RTMP encrypted using Adobe's own security mechanism. While the details of the implementation are proprietary, the mechanism uses industry standard cryptographic primitives.[1]
- RTMPT, which is encapsulated within HTTP requests to traverse firewalls. RTMPT is frequently found utilizing cleartext requests on TCP ports 80 and 443 to bypass most corporate traffic filtering. The encapsulated session may carry plain RTMP, RTMPS, or RTMPE packets within.
- RTMFP, which is RTMP over User Datagram Protocol (UDP) instead of TCP, replacing RTMP Chunk Stream. The Secure Real-Time Media Flow Protocol suite has been developed by Adobe Systems and enables end‐users to connect and communicate directly with each other (P2P).

## Getting Started
### Installation
Use the go-rtmp package.

```shell
$ go get github.com/yutopp/go-rtmp
```