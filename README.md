# M3u8文件下载视频
通过m3u8网络文件或者本地文件下载合成视频
- 支持自动Aes-128解密
- 支持网络m3u8文件
- 支持本地m3u8文件
- 自动添加前缀host
- 支持多文件下载
- 支持多线程下载


## 下载地址
[m3u8](http://git.kaidon.cn/don178/m3u8/releases)

## 使用方法
在配置文件 conf.yaml 中添加本地路径或者网络路径
```yaml
address:
   - 
    # 网络路径或者m3u8本地文件
    path: 
    # [非必填] 当m3u8文件内ts地址没用域名时
    # 1. 使用该profix作为域名
    # 2. 从path中提取 path[:最后一个/的位置]
    prefix: 
    # [非必填] 保存文件名称(默认从网络地址中提取名称)
    # 最终文件地址 = dir + fname + ext
    fname: 
   - 
    path: 
    prefix: 
    fname: 
```

添加完成后直接运行m3u8,无需指定配置文件路径.

```shell
./m3u8
```