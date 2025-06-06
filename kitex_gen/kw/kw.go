// Code generated by thriftgo (0.4.1). DO NOT EDIT.

package kw

import (
	"context"
	"fmt"
)

type KwRequest struct {
	SongId  string `thrift:"songId,1,required" frugal:"1,required,string" json:"songId"`
	Quality string `thrift:"quality,2,required" frugal:"2,required,string" json:"quality"`
}

func NewKwRequest() *KwRequest {
	return &KwRequest{}
}

func (p *KwRequest) InitDefault() {
}

func (p *KwRequest) GetSongId() (v string) {
	return p.SongId
}

func (p *KwRequest) GetQuality() (v string) {
	return p.Quality
}
func (p *KwRequest) SetSongId(val string) {
	p.SongId = val
}
func (p *KwRequest) SetQuality(val string) {
	p.Quality = val
}

func (p *KwRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("KwRequest(%+v)", *p)
}

var fieldIDToName_KwRequest = map[int16]string{
	1: "songId",
	2: "quality",
}

type KwResponse struct {
	Code  int16  `thrift:"code,1,required" frugal:"1,required,i16" json:"code"`
	Msg   string `thrift:"msg,2,required" frugal:"2,required,string" json:"msg"`
	Data  string `thrift:"data,3" frugal:"3,default,string" json:"data"`
	Extra *Extra `thrift:"extra,4" frugal:"4,default,Extra" json:"extra"`
}

func NewKwResponse() *KwResponse {
	return &KwResponse{}
}

func (p *KwResponse) InitDefault() {
}

func (p *KwResponse) GetCode() (v int16) {
	return p.Code
}

func (p *KwResponse) GetMsg() (v string) {
	return p.Msg
}

func (p *KwResponse) GetData() (v string) {
	return p.Data
}

var KwResponse_Extra_DEFAULT *Extra

func (p *KwResponse) GetExtra() (v *Extra) {
	if !p.IsSetExtra() {
		return KwResponse_Extra_DEFAULT
	}
	return p.Extra
}
func (p *KwResponse) SetCode(val int16) {
	p.Code = val
}
func (p *KwResponse) SetMsg(val string) {
	p.Msg = val
}
func (p *KwResponse) SetData(val string) {
	p.Data = val
}
func (p *KwResponse) SetExtra(val *Extra) {
	p.Extra = val
}

func (p *KwResponse) IsSetExtra() bool {
	return p.Extra != nil
}

func (p *KwResponse) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("KwResponse(%+v)", *p)
}

var fieldIDToName_KwResponse = map[int16]string{
	1: "code",
	2: "msg",
	3: "data",
	4: "extra",
}

type Extra struct {
	Cache   bool     `thrift:"cache,1,required" frugal:"1,required,bool" json:"cache"`
	Quality *Quality `thrift:"quality,2,required" frugal:"2,required,Quality" json:"quality"`
	Expire  *Expire  `thrift:"expire,3,required" frugal:"3,required,Expire" json:"expire"`
}

func NewExtra() *Extra {
	return &Extra{}
}

func (p *Extra) InitDefault() {
}

func (p *Extra) GetCache() (v bool) {
	return p.Cache
}

var Extra_Quality_DEFAULT *Quality

func (p *Extra) GetQuality() (v *Quality) {
	if !p.IsSetQuality() {
		return Extra_Quality_DEFAULT
	}
	return p.Quality
}

var Extra_Expire_DEFAULT *Expire

func (p *Extra) GetExpire() (v *Expire) {
	if !p.IsSetExpire() {
		return Extra_Expire_DEFAULT
	}
	return p.Expire
}
func (p *Extra) SetCache(val bool) {
	p.Cache = val
}
func (p *Extra) SetQuality(val *Quality) {
	p.Quality = val
}
func (p *Extra) SetExpire(val *Expire) {
	p.Expire = val
}

func (p *Extra) IsSetQuality() bool {
	return p.Quality != nil
}

func (p *Extra) IsSetExpire() bool {
	return p.Expire != nil
}

func (p *Extra) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Extra(%+v)", *p)
}

var fieldIDToName_Extra = map[int16]string{
	1: "cache",
	2: "quality",
	3: "expire",
}

type Quality struct {
	Target  string `thrift:"target,1,required" frugal:"1,required,string" json:"target"`
	Result_ string `thrift:"result,2,required" frugal:"2,required,string" json:"result"`
}

func NewQuality() *Quality {
	return &Quality{}
}

func (p *Quality) InitDefault() {
}

func (p *Quality) GetTarget() (v string) {
	return p.Target
}

func (p *Quality) GetResult_() (v string) {
	return p.Result_
}
func (p *Quality) SetTarget(val string) {
	p.Target = val
}
func (p *Quality) SetResult_(val string) {
	p.Result_ = val
}

func (p *Quality) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Quality(%+v)", *p)
}

var fieldIDToName_Quality = map[int16]string{
	1: "target",
	2: "result",
}

type Expire struct {
	Time      int64 `thrift:"time,1,required" frugal:"1,required,i64" json:"time"`
	CanExpire bool  `thrift:"canExpire,2,required" frugal:"2,required,bool" json:"canExpire"`
}

func NewExpire() *Expire {
	return &Expire{}
}

func (p *Expire) InitDefault() {
}

func (p *Expire) GetTime() (v int64) {
	return p.Time
}

func (p *Expire) GetCanExpire() (v bool) {
	return p.CanExpire
}
func (p *Expire) SetTime(val int64) {
	p.Time = val
}
func (p *Expire) SetCanExpire(val bool) {
	p.CanExpire = val
}

func (p *Expire) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Expire(%+v)", *p)
}

var fieldIDToName_Expire = map[int16]string{
	1: "time",
	2: "canExpire",
}

type KwService interface {
	KwMusicUrl(ctx context.Context, req *KwRequest) (r *KwResponse, err error)
}

type KwServiceKwMusicUrlArgs struct {
	Req *KwRequest `thrift:"req,1" frugal:"1,default,KwRequest" json:"req"`
}

func NewKwServiceKwMusicUrlArgs() *KwServiceKwMusicUrlArgs {
	return &KwServiceKwMusicUrlArgs{}
}

func (p *KwServiceKwMusicUrlArgs) InitDefault() {
}

var KwServiceKwMusicUrlArgs_Req_DEFAULT *KwRequest

func (p *KwServiceKwMusicUrlArgs) GetReq() (v *KwRequest) {
	if !p.IsSetReq() {
		return KwServiceKwMusicUrlArgs_Req_DEFAULT
	}
	return p.Req
}
func (p *KwServiceKwMusicUrlArgs) SetReq(val *KwRequest) {
	p.Req = val
}

func (p *KwServiceKwMusicUrlArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *KwServiceKwMusicUrlArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("KwServiceKwMusicUrlArgs(%+v)", *p)
}

var fieldIDToName_KwServiceKwMusicUrlArgs = map[int16]string{
	1: "req",
}

type KwServiceKwMusicUrlResult struct {
	Success *KwResponse `thrift:"success,0,optional" frugal:"0,optional,KwResponse" json:"success,omitempty"`
}

func NewKwServiceKwMusicUrlResult() *KwServiceKwMusicUrlResult {
	return &KwServiceKwMusicUrlResult{}
}

func (p *KwServiceKwMusicUrlResult) InitDefault() {
}

var KwServiceKwMusicUrlResult_Success_DEFAULT *KwResponse

func (p *KwServiceKwMusicUrlResult) GetSuccess() (v *KwResponse) {
	if !p.IsSetSuccess() {
		return KwServiceKwMusicUrlResult_Success_DEFAULT
	}
	return p.Success
}
func (p *KwServiceKwMusicUrlResult) SetSuccess(x interface{}) {
	p.Success = x.(*KwResponse)
}

func (p *KwServiceKwMusicUrlResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *KwServiceKwMusicUrlResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("KwServiceKwMusicUrlResult(%+v)", *p)
}

var fieldIDToName_KwServiceKwMusicUrlResult = map[int16]string{
	0: "success",
}
