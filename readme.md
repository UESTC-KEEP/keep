# 配置文件说明
- 系统配置文件：/etc/keepedge
  - 日志配置文件：logger_conf.json
    - 内容：
  ```go
      {
      "File": {
      "filename": "/var/log/keepedge/keep_edgecore_logs.log",
      "level": "TRAC",
      "daily": true,
      "maxlines": 1000000,
      "maxsize": 1,
      "maxdays": -1,
      "append": true,
      "permit": "0660"
      }
      }
  ```

- 日志文件存放：/var/log/keepedge