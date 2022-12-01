package common

/*
	错误码按照模块化，功能 区分 。中间采用 - 连接
public:
	标准化输入输出		0000
	用户及手机验证码		0100
	数据库相关       	0200
	SDK相关       		0300
    业务相关     0400

*/

const (
	// 标准化输出		0000
	E_Success            = "0000-0000" //成功
	E_ParseBodyFailed    = "0000-0001" //解析请求body失败
	E_ParamFormatInvalid = "0000-0002" //传入参数错误，解析失败
	E_InvalidInput       = "0000-0003" //输入参数不合法
	E_MkdirFailed        = "0000-0004" //创建目录失败
	E_CreateFileFailed   = "0000-0005" //打开文件失败
	E_ReadFileFailed     = "0000-0006" //读取文件失败
	E_WriteFileFailed    = "0000-0007" //写入文件失败
	E_Base64DecodeFailed = "0000-0008" //base64解码失败

	// 用户及手机验证码	0100
	E_UserNameInvalid         = "0100-0001" //用户名无效
	E_UserTokenInvalid        = "0100-0002" //当前token无效
	E_UserNameOrPasswordError = "0100-0003" //用户名或密码错误
	E_UserHasNoPermission     = "0100-0004" //用户没有权限
	E_UserPasswordInvalid     = "0100-0005" //用户密码无效
	E_UserNameIsExist         = "0100-0006" //用户名已存在
	E_UserStateInvalid        = "0100-0008" //用户审核中,请通过后再登录
	E_PermissionDenied        = "0100-0009" //当前用户无操作权限
	E_UserTokenGetFailed      = "0100-0010" //获取token失败
	E_UserGetFailed           = "0100-0011" //查询用户失败
	E_InvalidOperation        = "0100-0012" //非法操作

	// 数据库 			0200
	E_InsertFailed = "0200-0001" //插入数据库失败
	E_GetFailed    = "0200-0002" //查询数据库失败
	E_DeleteFailed = "0200-0003" //删除数据库失败
	E_UpdateFailed = "0200-0004" //更新数据库失败

	// 业务相关
	E_NotExistUser     			= "0300-0001" //不存在该用户


	//请求处理失败
	E_IoReadError        = "0500-0001" //io读取失败
	E_JsonUnmarshalError = "0500-0002" //json解析失败
	E_ExistSensitiveRord = "0500-0003" //输入存在敏感词
)
