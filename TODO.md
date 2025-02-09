# TODO

## NextsNet (network) / NexusChat (app)

- [ ] TCPMessage impl & improvements
  - [x] Change the impl to support ctx cancellation (currently blocking in default)
  - [x] For the task above, `reader.ReadBytes` run in separate goroutine and listen for `messageHeader`s
  - [x] Write tests for the new impl (non blocking, goroutine reader)
  - [ ] Implement partial data processing.

- [ ] TCPConnection implementation `struct { conn: net.Conn }`

- [ ] Make a throw-away implementation of client-server or p2p client-client communication using TCPMessage and TCPConnection packages
- [ ] Create a `apps/experiments/tcp_comm_01` app and use TCPMessage and TCPConn packages
      and build a PoC TCP messaging app

Later:
- [ ] also use `log_prompt` package to implement interactive connection and communication between nodes trough TCPConnection and TCPMessage
