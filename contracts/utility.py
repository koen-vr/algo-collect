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

@Subroutine(TealType.bytes)
def get_asset_unit_name(arg):

    out = ScratchVar(TealType.bytes)
    tmp = ScratchVar(TealType.bytes)
    num = ScratchVar(TealType.uint64)

    return Seq([
        out.store(Bytes("Nr#")),
        tmp.store(uint_to_bytes(arg)),
        For(num.store(Int(5) - Len(tmp.load())), num.load() > Int(0), num.store(num.load() - Int(1))).Do(
            out.store(Concat(
                out.load(), 
                Bytes("0")
            )),
        ),
        out.store(Concat(
            out.load(),
            uint_to_bytes(arg)
        )),
        out.load()
    ])