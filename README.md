# M3u8下载

## 下载地址
[m3u8](http://git.kaidon.cn/attachments/f80462f5-20d6-41fb-81cc-2a5d0a5bcc19)

## 使用方法
在配置文件 conf.yaml 中添加本地路径或者网络路径
```yaml
address:
   - 
    # 网络路径或者m3u8本地文件路径
    path: 
    # 当m3u8文件内的地址不是完整地址时,会在地址前添加prefix组成完成地址
    # 如果m3u8文件内时完整地址,该项无效,可以为空.
    prefix: 
  - 
    # 网络路径或者m3u8本地文件路径
    path:
    # 当m3u8文件内的地址不是完整地址时,会在地址前添加prefix组成完成地址
    # 如果m3u8文件内时完整地址,该项无效,可以为空.
    prefix: 
```

添加完成后直接运行m3u8,无需指定配置文件路径.

```shell
./m3u8
```