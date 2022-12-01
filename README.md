# game

##一、 目录结构
game
|---build 编译二进制文件和镜像文件
|---deploy 部署服务文件
|---server 服务文件
       |---router   路由
       |---handler  业务逻辑 （处理投骰子逻辑）
       |---main     服务主函数
       |---...

##二、 表结构

### user表
| Field      | Type         | Remark|
| ----- | ----- | ------ |
| created_at | string     | 创建时间  |
| updated_at | string     | 更新时间|
| name       | string | 用户名|
| password   | string | 加密后的用户密码|
| mail       | string | 邮箱|
| phone_num  | string | 电话号码|
| company    | string | 公司|
| type       | int    | 用户类型(0:超级管理员， 1:资源管理员，2:网络管理员)|
| is_valid   | bool   | 是否可用|

### token表

| Field      | Type         |  Remark|
| ----- | ----- | ------ |
| token      | string       | token串|
| name       | string       | 用户名|
| origin     | string       | 组织|
| created_at | string       | 创建时间|
| updated_at | string       | 更新时间|

### invoke 表
| Field      | Type         | Remark|
| ----- | ----- | ------ |
| created_at | string     | 创建时间  |
| updated_at | string     | 更新时间|
| name       | string | 用户名|
| CurrentPos   | int |  当前用户所在地图的位置｜
| IsWin       | Bool | 是否获胜|
| Room       | int | 游戏房间号码|


##三、 接口

**Path：** /v1/game/invoke

**Method：** GET

**接口描述：**
<p>投骰子</p>

### 请求参数

**Headers**

| 参数名称  | 参数值  |  是否必须 | 示例  | 备注  |
| ------------ | ------------ | ------------ | ------------ | ------------ |
| Content-Type  |  application/json | 是  |   |   |
| user  |   | 是  |   |  user (Only:undefined) |
| token  |   | 是  |   |  token (Only:undefined) |




* 响应参数

  ```json
  {
    "err_code": "0000-0000",
    "err_cue": "成功",
    "err_msg": "",
    "data": {
      "invoke": {
        "created_at": "2021-01-13T12:38:33Z",
        "updated_at": "2021-01-13T12:38:33Z",
        "name": "root",
        "CurrentPos":6,
        "IsWin": "false",
        "Room": "1"
      }
    }
  }
  ```


##备注：

在数据库初始化时，创建两个用户root 密码88888888,test 密码666666
登陆两个用户进行后，调用投骰子接口, （便于测试房间号固定都为1）

业务逻辑在server/handler/invoke.go







