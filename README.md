# drone-plugin-helm-release

> drone plugin 参数模版

## 内置参数

- `PLUGIN_DEBUG` - 是否开启debug日志
- `PLUGIN_PAUSE` - 是否调试
- `PLUGIN_PROXY` - 代理

## helm参数

- `PLUGIN_USERNAME` - 用户名
- `PLUGIN_PASSWORD` - 密码
- `PLUGIN_TOKEN`	- token (与密码互斥)
- `PLUGIN_HUB`	- repo地址(需要以http://开头) 
- `PLUGIN_CONTEXT` - 路径，默认.
- `PLUGIN_MULTI` - 是否为多个chart，当前目录下有多个charts
- `PLUGIN_FORCE` - 是否强制发布
- `PLUGIN_EXTHUB` - 依赖的chart的hub地址
- `PLUGIN_EXCLUDE` - 排除的chart, 仅多个charts生效

## usage

```yaml
  - name: helm release
    image: ysicing/drone-plugin-helm-release
    privileged: true
    pull: always
    settings:
      debug: true
      username:
        from_secret: tcr-username
      password:
        from_secret: tcr-password
      hub:
        from_secret: tcr-hub
      context: ./stable
      multi: true
      exthub: 
      - https://charts.bitnami.com/bitnami
```
