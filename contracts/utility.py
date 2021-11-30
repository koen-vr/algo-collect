from pyteal import *

@Subroutine(TealType.bytes)
def uint_to_bytes(arg):

    num = ScratchVar(TealType.uint64)
    digit = ScratchVar(TealType.uint64)
    string = ScratchVar(TealType.bytes)

    return If(
        arg == Int(0),
        Bytes("0"),
        Seq([
            string.store(Bytes("")),
            For(num.store(arg), num.load() > Int(0), num.store(num.load() / Int(10))).Do(
                Seq([
                    digit.store(num.load() % Int(10)),
                    string.store(
                        Concat(
                            Substring(
                                Bytes('0123456789'),
                                digit.load(),
                                digit.load() + Int(1)
                            ),
                            string.load()
                        )
                    )
                ])

            ),
            string.load()
        ])
    )