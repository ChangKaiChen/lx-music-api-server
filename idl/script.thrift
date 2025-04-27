namespace go script

struct ScriptRequest {
    1: required string url
    2: required string key
    3: string checkUpdate
}
struct ScriptResponse {
    1: required bool success
    2: string content
    3: Res res
}
struct Res {
    1: i16 code
    2: string msg
    3: map<string, string> data
    4: string content
}
service ScriptService {
    ScriptResponse Script(1: ScriptRequest req)
}