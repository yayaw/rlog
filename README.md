# rlog
用标准库log封装的简单日志系统

使用说明：
  日志会以为数据大小为单位，进行轮询，默认为10MB，每个文件会以日期为名称进行建立
  调用使用rlog.SetMaxFileSizeMB(10)
  
  支持是否进行终端输出，或输入文件，或同时输入
