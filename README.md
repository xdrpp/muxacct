# muxacct

This is a simple command-line tool for packing and unpacking Stellar
`MuxedAccount` structures in the strkey format proposed by [SEP-0023].

[SEP-0023]: https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0023.md

# Building

To build this you need goxdr, which you can install with:

    go get github.com/xdrpp/goxdr/cmd/goxdr

Then just run make.

# Running

To create a `MuxedAccount`:

    $ ./muxacct mux 100 GCHWOCH5OIMSAXBJFS4AWJR6SG3DSXY4IRV3AX32N6G2VY4WABUEIGAV
    MAAAAAAAAAAAAZEPM4EP24QZEBOCSLFYBMTD5ENWHFPRYRDLWBPXU34NVLRZMADIISM34

To unpack one:

    $ ./muxacct demux MAAAAAAAAAAAAZEPM4EP24QZEBOCSLFYBMTD5ENWHFPRYRDLWBPXU34NVLRZMADIISM34
    100 GCHWOCH5OIMSAXBJFS4AWJR6SG3DSXY4IRV3AX32N6G2VY4WABUEIGAV
