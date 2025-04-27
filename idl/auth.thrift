namespace go auth

struct AuthRequest {
    1 : required string authKey
}
struct AuthResponse {
    1: required bool success
    2: required map<string,string> res
}
service AuthService {
    AuthResponse Auth(1:AuthRequest req)
}