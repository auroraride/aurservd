node {
    stage('Development') {
        if (TAG == 'latest' || TAG == 'dev') {
            echo '开始部署Development环境'
            sshagent (credentials: ['Jenkins']) {
                sh "ssh -o StrictHostKeyChecking=no root@39.106.77.239 '${deploy('next-api', TAG, 'https://next-api.auroraride.com/maintain/update')}'"
            }
        }
    }
    stage('Production') {
        if (TAG == 'latest' || TAG == 'prod') {
            echo '开始部署Production环境'
            sshagent (credentials: ['Jenkins']) {
                sh "ssh -o StrictHostKeyChecking=no root@39.106.77.239 '${deploy('api', TAG, 'https://api.auroraride.com/maintain/update')}'"
            }
        }
    }
}

def deploy(path, tag, url) {
    def str = """
        docker pull registry-vpc.cn-beijing.aliyuncs.com/liasica/aurservd:$tag
        curl $url
        docker stop ${path}
        docker rm -f ${path}
        mkdir -p /var/www/${path}.auroraride.com/runtime
        docker run -itd --user 0 --name ${path} --restart=always \
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
