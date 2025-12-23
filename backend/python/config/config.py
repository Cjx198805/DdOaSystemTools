import os
from dotenv import load_dotenv

# 加载环境变量
load_dotenv()

# 服务器配置
SERVER_PORT = os.getenv('SERVER_PORT', '8081')

# MySQL配置
MYSQL_HOST = os.getenv('MYSQL_HOST', 'localhost')
MYSQL_PORT = os.getenv('MYSQL_PORT', '3306')
MYSQL_USERNAME = os.getenv('MYSQL_USERNAME', 'root')
MYSQL_PASSWORD = os.getenv('MYSQL_PASSWORD', 'root')
MYSQL_DBNAME = os.getenv('MYSQL_DBNAME', 'dd_oa_download')
MYSQL_CHARSET = os.getenv('MYSQL_CHARSET', 'utf8mb4')

# Redis配置
REDIS_HOST = os.getenv('REDIS_HOST', 'localhost')
REDIS_PORT = os.getenv('REDIS_PORT', '6379')
REDIS_PASSWORD = os.getenv('REDIS_PASSWORD', '')
REDIS_DB = int(os.getenv('REDIS_DB', '0'))

# 钉钉配置
DINGTALK_APPKEY = os.getenv('DINGTALK_APPKEY', '')
DINGTALK_APPSECRET = os.getenv('DINGTALK_APPSECRET', '')
