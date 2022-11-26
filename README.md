# M3u8文件下载视频
通过m3u8网络文件或者本地文件下载合成视频
- 支持自动Aes解密
- 支持网络m3u8文件
- 支持本地m3u8文件
- 自动添加前缀host
- 支持多文件下载


## 下载地址
[m3u8](http://git.kaidon.cn/attachments/f80462f5-20d6-41fb-81cc-2a5d0a5bcc19)

## 使用方法
在配置文件 conf.yaml 中添加本地路径或者网络路径
```yaml
address:
   - 
    # 网络路径或者m3u8本地文件
    path: 
    # [非必填] 当m3u8文件内没有host时,在地址前添加prefix 
    prefix: 
    # [非必填] 保存文件名称(默认从网络地址中提取名称)
    # 最终文件地址 = dir + fname + ext
    fname: 
```

添加完成后直接运行m3u8,无需指定配置文件路径.

```shell
./m3u8
```