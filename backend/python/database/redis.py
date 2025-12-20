import redis
from config.config import REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, REDIS_DB

# 创建Redis客户端
redis_client = redis.Redis(
    host=REDIS_HOST,
    port=REDIS_PORT,
    password=REDIS_PASSWORD,
    db=REDIS_DB,
    decode_responses=True
)

# 测试连接
def test_redis_connection():
    """测试Redis连接"""
    try:
        pong = redis_client.ping()
        return pong
    except Exception as e:
        raise Exception(f"Redis连接失败: {str(e)}")

# 获取Redis客户端
def get_redis():
    """获取Redis客户端"""
    return redis_client
