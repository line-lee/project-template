package errorcode

const (
	UrlParse               = sdk*10000 + 1
	JsonMarshal            = sdk*10000 + 2
	HttpNewRequest         = sdk*10000 + 3
	HttpClientDo           = sdk*10000 + 4
	HttpResponseNil        = sdk*10000 + 5
	HttpResponseBodyNil    = sdk*10000 + 6
	HttpResponseStatusCode = sdk*10000 + 7
	IoReadAll              = sdk*10000 + 8
	JsonUnmarshal          = sdk*10000 + 9
	OSOpen                 = sdk*10000 + 10
	StrconvParseFloat      = sdk*10000 + 11
	ResponseBodyClose      = sdk*10000 + 12
	GinRunErr              = sdk*10000 + 13
	GrpcNewClient          = sdk*10000 + 14
	GrpcServiceUnknown     = sdk*10000 + 15
	GinShouldBindErr       = sdk*10000 + 16
	GrpcRequest            = sdk*10000 + 17
	StrconvParseInt        = sdk*10000 + 18
	AESEncrypt             = sdk*10000 + 19
	AESDecrypt             = sdk*10000 + 20
	MD5XError              = sdk*10000 + 21
	WebsocketUpgrade       = sdk*10000 + 22
	WebsocketClose         = sdk*10000 + 22
)

func SdkCode() {
	ErrorCode[UrlParse] = "url parse err"
	ErrorCode[JsonMarshal] = "json marshal err"
	ErrorCode[HttpNewRequest] = "http new request err"
	ErrorCode[HttpClientDo] = "http client do err"
	ErrorCode[HttpResponseNil] = "http response nil"
	ErrorCode[HttpResponseBodyNil] = "http response body nil"
	ErrorCode[HttpResponseStatusCode] = "http response status code [%v], body[%v]"
	ErrorCode[IoReadAll] = "io read all err"
	ErrorCode[JsonUnmarshal] = "json unmarshal err"
	ErrorCode[OSOpen] = "[%v]os open err"
	ErrorCode[StrconvParseFloat] = "strconv parse float err"
	ErrorCode[ResponseBodyClose] = "response body close err"
	ErrorCode[GinRunErr] = "gin run err"
	ErrorCode[GrpcNewClient] = "grpc new client err"
	ErrorCode[GrpcServiceUnknown] = "grpc service unknown"
	ErrorCode[GinShouldBindErr] = "gin should bind err"
	ErrorCode[GrpcRequest] = "grpc request err"
	ErrorCode[StrconvParseInt] = "strconv parse int err"
	ErrorCode[AESEncrypt] = "aes encrypt err"
	ErrorCode[AESDecrypt] = "aes decrypt err"
	ErrorCode[MD5XError] = "md5x error"
	ErrorCode[WebsocketUpgrade] = "websocket upgrade err"
	ErrorCode[WebsocketClose] = "websocket close err"
}
