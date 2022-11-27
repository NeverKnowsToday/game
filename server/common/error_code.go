package common

/*
	错误码按照模块化，功能 区分 。中间采用 - 连接
public:
	标准化输入输出		0000
	用户及手机验证码		0100
	数据库相关       	0200
	SDK相关       		0300
    区块链浏览器相关     0400

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

	// SDK相关
	E_ReadUserKeyPairFailed     = "0300-0001" //读取用户证书私钥失败
	E_SendTxFailed              = "0300-0002" //交易发送失败（连接区块链节点提交交易失败， 需要重新发送）
	E_InvalidTx                 = "0300-0003" //交易发送失败（智能合约校验未通过）
	E_EnrollFailed              = "0300-0004" //证书注册失败
	E_GetBlockHeight            = "0300-0005" //获取区块高度失败
	E_HandleTxStatus            = "0300-0006" //设置同步发交易时，接收交易信息失败
	E_QueryFailed               = "0300-0007" //查询交易失败
	E_InvokeFailed              = "0300-0008" //执行交易失败
	E_ChannelNotExisted         = "0300-0009" //通道不存在或没有和相应组织绑定
	E_ChainCodeNotExisted       = "0300-0010" //智能合约不存在
	E_PeerNotFind               = "0300-0011" //未找到实例化了智能合约的peer
	E_OrdererNotFind            = "0300-0012" //未找到可用的orderer
	E_FunctionCanNotBeNull      = "0300-0013" //function不能为空
	E_ChannelCanNotBeNull       = "0300-0014" //channel不能为空
	E_FunctionIsInvalid         = "0300-0015" //function名非法
	E_EndorsementPolicyFailure  = "0300-0016" //交易背书策略不满足
	E_ChainCodeIsNotInstantiate = "0300-0017" //智能合约未实例化
	E_TxIdCanNotBeNull          = "0300-0018" //交易id不能为空
	E_OrgIsNotInChannel         = "0300-0019" //组织不在通道中

	// 监听通道相关
	E_ChannelListeningCreateFailed = "0400-0001" //监听通道创建失败
	E_ChannelListeningUpdateFailed = "0400-0002" //监听通道更新失败
	E_ChannelListeningDeleteFailed = "0400-0003" //监听通道删除失败
	E_ChannelListeningGetFailed    = "0400-0004" //监听通道查询失败

	//请求处理失败
	E_IoReadError        = "0500-0001" //io读取失败
	E_JsonUnmarshalError = "0500-0002" //json解析失败
	E_ExistSensitiveRord = "0500-0003" //输入存在敏感词
)
