package codecs 

type BaseCodecs interface {
    Encode(in any)([]byte, error)
    Decode(data []byte, dst any) (error)
}
