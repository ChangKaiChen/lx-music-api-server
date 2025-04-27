namespace go wy

struct WyRequest {
    1: required string songId
    2: required string quality
}
struct WyResponse {
    1: required i16 code
    2: required string msg
    3: string data
    4: Extra extra
}
struct Extra {
    1: required bool cache
    2: required Quality quality
    3: required Expire expire
}
struct Quality {
    1: required string target
    2: required string result
}
struct Expire {
    1: required i64 time
    2: required bool canExpire
}
service WyService {
    WyResponse WyMusicUrl(1: WyRequest req)
}