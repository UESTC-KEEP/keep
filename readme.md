# 配置文件说明
- 系统配置文件：/etc/keepedge
  -  config/edgeagent.yml
  ```go 
  database:
    drivername: ""
    aliasname: ""
    datasource: /var/lib/keepedge/edgeagent.db
  modules:
    healthzagent:
      enable: true
      hostinfostat: null
      cpu: null
      cpuusage: 0
      mem: null
      diskpartitionstat: null
      diskiocountersstat: null
      netiocountersstat: null
      defaultedgehealthinterval: 10
      defaultmqttcachequeuesize: 10
    logsagent:
      enable: true
      loglevel: 6
      logtime: 2021-11-16T09:26:32.543688913+08:00
      logfilename: ""
      loginfo: ""
      logfiles:
      - /var/log/keepedge/keep_edgeagent_logs.log
      - /var/log/test.log
    edgepublisher:
      enable: true
      httpserver: http://192.168.1.140
      port: 20000
      serveport: 20350
      heartbeat: 15
      edgemsgqueens:
      - keep_log_topic
      - keep_data_topic
      tlscafile: ""
      tlscertfile: ""
      tlsprivatekeyfile: ""
  ```
  - 日志配置文件：logger_conf.json
    - 内容：
  ```go
      {
      "File": {
      "filename": "/var/logs/keepedge/keep_edgecore_logs.logs",
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
