node {
    stage('Development') {
        echo '开始部署Development环境'
        def deploy = '''
            docker pull registry-vpc.cn-beijing.aliyuncs.com/liasica/aurservd
            docker rm -f next-api
            mkdir -p /var/www/next-api.auroraride.com/runtime
            docker run -itd --name next-api --restart=always --network host -v /var/www/next-api.auroraride.com/config:/app/config -v /var/www/next-api.auroraride.com/runtime:/app/runtime registry-vpc.cn-beijing.aliyuncs.com/liasica/aurservd
        '''
        sshagent (credentials: ['Jenkins']) {
            sh "ssh -o StrictHostKeyChecking=no root@39.106.77.239 '${deploy}'"
        }
        echo '完成Development环境部署'
    }
    stage('Production') {
        input 'Deploy to Production?'
        echo "开始部署Production环境"
        echo "已终止部署"
    }
}