node {
    stage('Development') {
        if (TAG == 'latest' || TAG == 'dev') {
            echo '开始部署Development环境'
            sshagent (credentials: ['Jenkins']) {
                sh "ssh -o StrictHostKeyChecking=no root@39.106.77.239 '${deploy('next-api', TAG)}'"
            }
        } else {
            echo '不需要部署Development环境'
        }
        echo '完成Development环境部署'
    }
    stage('Production') {
        timeout (time: 1, unit: 'HOURS' )  {
            input 'Deploy to Production?'
        }
        if (TAG == 'latest' || TAG == 'prod') {
            echo '开始部署Production环境'
            sshagent (credentials: ['Jenkins']) {
                sh "ssh -o StrictHostKeyChecking=no root@39.106.77.239 '${deploy('api', TAG)}'"
            }
        } else {
            echo '不需要部署Production环境'
        }
        echo '已结束Production环境部署'
    }
}

def deploy(path, tag) {
    def str = """
        docker pull registry-vpc.cn-beijing.aliyuncs.com/liasica/aurservd:$tag
        docker rm -f ${path}
        mkdir -p /var/www/${path}.auroraride.com/runtime
        docker run -itd --name ${path} --restart=always \
        --network host \
        -v /var/www/${path}.auroraride.com/config:/app/config \
        -v /var/www/${path}.auroraride.com/runtime:/app/runtime \
        -v /var/www/${path}.auroraride.com/public:/app/public \
        registry-vpc.cn-beijing.aliyuncs.com/liasica/aurservd:$tag
        docker image prune -f
        docker container prune -f
        docker volume prune -f
    """
    return str
}