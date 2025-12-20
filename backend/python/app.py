from flask import Flask, jsonify
from flask_cors import CORS
import logging
from config.config import SERVER_PORT
from database.db import engine, Base
from database.redis import test_redis_connection

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)

# 创建Flask应用
app = Flask(__name__)

# 启用CORS
CORS(app)

# 创建数据库表
Base.metadata.create_all(bind=engine)

# 健康检查
@app.route('/health', methods=['GET'])
def health_check():
    """健康检查"""
    return jsonify({
        "status": "ok",
        "message": "DdOaListDownload Python 服务运行正常"
    })

# API分组
@app.route('/api/v1/legacy', methods=['GET'])
def legacy_api():
    """旧版API调用示例"""
    return jsonify({
        "status": "ok",
        "message": "旧版API调用成功",
        "data": {
            "version": "legacy",
            "service": "python"
        }
    })

if __name__ == '__main__':
    logger.info("启动 DdOaListDownload Python 服务")
    
    # 测试Redis连接
    try:
        redis_pong = test_redis_connection()
        logger.info(f"Redis连接成功，响应: {redis_pong}")
    except Exception as e:
        logger.error(f"Redis连接失败: {str(e)}")
    
    # 启动服务器
    server_addr = f"0.0.0.0:{SERVER_PORT}"
    logger.info(f"Python 服务器正在启动，监听地址: {server_addr}")
    app.run(host='0.0.0.0', port=int(SERVER_PORT), debug=True)
