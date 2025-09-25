
帮我编写一个发布本服务的Jenkins pipeline流水线，Jenkins在我的内网中，服务器在阿里云/华为云的机房中，你需要帮我编写这个部署流水线。
Jenkins pipeline文件： deply\prd\jenkins-pipline.groovy
dockerfile文件：deply\prd\Dockerfile

go程序调用： /root/.g/versions/1.23.7/bin
python程序： /root/anaconda3/envs/jenkins/bin/python

流程说明：
1. 在本地Jenkins中将拉取的代码使用go程序编译完成 git地址： https://gitee.com/xxxxx/xxxxx-xxxxx.git
2. 然后需要将编译好的程序和config文件夹、deply文件夹、web文件夹以及编译好的服务（服务名：nano-banana-prompt-api）上传到线上服务器中 /data/wwwroot/nano-banana-prompt-api 目录下
3. 再在线上服务中调用docker命令，基于当前新文件构建一个新的镜像
4. 然后结束旧的服务，启动新的服务


（不需要下面的蓝绿部署的逻辑；也不需要你编译、上传docker镜像，镜像到线上直接生成；）
参考pipeline：

pipeline {
    agent any
    
    environment {
        DOCKER_DIR="build" // Dockerfile 所在的根目录
        SERVICE_NAME = "am-go-api"
        PYTHON_SUB_SERVICE = "api"
        DOCKER_REGISTRY = "192.168.0.1:8018/product-test"
        DOCKER_IMAGE = "${DOCKER_REGISTRY}/${SERVICE_NAME}"
        JENKINS_SERVER = "aaaa@192.168.0.2"
        JENKINS_GO =  "/root/.g/versions/1.23.7/bin/go"
        DOCKER_SERVER = "aaaa@192.168.0.2"
        NGINX_SERVER = "aaaa@192.168.0.2"
        CURRENT_VERSION = "blue"
        NEW_VERSION = "green"
        BLUE_PORT = "8001"
        GREEN_PORT = "8002"
        DOCKER_SRV_PORT="12010"
        // 清理旧Docker镜像 配置
        REGISTRY = "192.168.0.3:8018"
        REPO = "product-test/${SERVICE_NAME}"
        KEEP_COUNT = 3
        // nginx 域名
        NGINX_DOMAIN = "http://yxxx.xxxx.xx"
        NGINX_PREFIX_PATH="go-api"
    }
    
    stages {
        stage('Checkout') {
            steps {
                git(
                    url: 'https://gitee.com/xxxx/xxxxxxxxx.git',
                    credentialsId: 'giteeUsername', // 填写你的 SSH 凭据 ID
                    branch: 'master'
                )
            }
        }
        
        stage('获取当前活动版本') {
            steps {
                script {
                    // 更健壮的SSH命令，处理可能的错误情况
                    def activeVersion = sh(
                        script: """
                            ssh -o StrictHostKeyChecking=no ${NGINX_SERVER} '/usr/local/bin/get_service_version ${SERVICE_NAME}'
                        """,
                        returnStdout: true
                    ).trim().toLowerCase()
                    
                    // 添加调试输出
                    echo "从Nginx服务器获取的原始版本信息: ${activeVersion}"
                    
                    // 验证获取的值是否有效
                    if (activeVersion != "blue" && activeVersion != "green") {
                        echo "警告：无法确定当前活动版本，默认使用blue"
                        activeVersion = "blue"
                    }
                    
                    // 确定当前和新版本
                    CURRENT_VERSION = activeVersion
                    NEW_VERSION = activeVersion == "blue" ? "green" : "blue"
                    
                    echo "当前版本: ${CURRENT_VERSION}, 将部署到: ${NEW_VERSION}"
                    
                    // 将版本信息存储为环境变量供后续阶段使用
                    env.CURRENT_VERSION = CURRENT_VERSION
                    env.NEW_VERSION = NEW_VERSION
                }
            }
        }

        // 需要修改
        stage('编译代码') {
            steps {
                echo "current path: ${pwd()}"
                echo "JENKINS_GO: ${JENKINS_GO}"

                sh "${JENKINS_GO} env -w GO111MODULE=on"
                sh "${JENKINS_GO} env -w  GOPROXY=https://goproxy.cn,direct"
                sh "${JENKINS_GO} env | grep GOPROXY"
                sh "${JENKINS_GO} mod tidy"
                sh "${JENKINS_GO} build -o main ./cmd/main.go"

                sh "echo '编译代码完成: ${pwd()}/main'"
            }
        }

        stage('构建Docker镜像') {
            steps {
                // sh 'go test ./...'
                sh "echo 'Image name:  ${DOCKER_IMAGE}:${BUILD_NUMBER}. Building...'"
                sh "rm -rf .git"
                sh "rm -rf README.*"
                sh "docker build -t ${DOCKER_IMAGE}:${BUILD_NUMBER}  ."
                sh "docker tag ${DOCKER_IMAGE}:${BUILD_NUMBER} ${DOCKER_IMAGE}:latest"
                sh "docker push ${DOCKER_IMAGE}:${BUILD_NUMBER}"
                sh "docker push ${DOCKER_IMAGE}:latest"
            }
        }

        stage('部署新版本') {
            steps {
                script {
                    def port = NEW_VERSION == "blue" ? BLUE_PORT : GREEN_PORT
                    
                    // 停止并移除旧的新版本容器（如果存在）
                    sh "ssh ${DOCKER_SERVER} 'docker stop ${SERVICE_NAME}-${NEW_VERSION} || true'"
                    sh "ssh ${DOCKER_SERVER} 'docker rm ${SERVICE_NAME}-${NEW_VERSION} || true'"
                    
                    // 启动新版本容器
                    sh """
                        ssh ${DOCKER_SERVER} 'docker run -d --name ${SERVICE_NAME}-${NEW_VERSION} \
                        -p ${port}:${DOCKER_SRV_PORT} \
                        --restart unless-stopped \
                        -e APP_ENV=production \
                        ${DOCKER_IMAGE}:${BUILD_NUMBER}'
                    """
                }
            }
        }

        stage('健康检查') {
            steps {
                script {
                    def port = NEW_VERSION == "blue" ? BLUE_PORT : GREEN_PORT
                    
                    // 等待容器启动
                    sh "sleep 10"
                    
                    // 健康检查
                    def maxRetries = 10
                    def retries = 0
                    def healthy = false
                    
                    while (!healthy && retries < maxRetries) {
                        try {
                            def health = sh(
                                script: "ssh ${DOCKER_SERVER} 'curl -s -o /dev/null -w \"%{http_code}\" http://localhost:${port}/heartbeat'",
                                returnStdout: true
                            ).trim()
                            
                            if (health == "200") {
                                healthy = true
                                echo "Service is healthy!"
                            } else {
                                retries++
                                echo "Health check failed (${health}), retrying... ${retries}/${maxRetries}"
                                sh "sleep 5"
                            }
                        } catch (Exception e) {
                            retries++
                            echo "Health check failed with exception, retrying... ${retries}/${maxRetries}"
                            sh "sleep 5"
                        }
                    }
                    
                    if (!healthy) {
                        error "Service failed health check after ${maxRetries} attempts"
                    }
                }
            }
        }

        stage('将流量切换到新版本中') {
            steps {
                script {
                    // 切换Nginx配置
                    sh """
                        ssh ${NGINX_SERVER} 'ln -sf /etc/nginx/versions/${SERVICE_NAME}/${NEW_VERSION}.conf /etc/nginx/conf.d/includes/${SERVICE_NAME}.conf'
                    """
                    
                    // 重载Nginx配置
                    sh "ssh ${NGINX_SERVER} '/usr/local/bin/tengine-ctl -t && /usr/local/bin/tengine-ctl reload'"
                    
                    echo "Traffic switched to ${NEW_VERSION} version"
                }
            }
        }

        stage('验证新版本') {
            steps {
                script {
                    // 验证新版本是否正常工作
                    sh "sleep 5"
                    
                    def statusCode = sh(
                        script: "curl -s -o /dev/null -w \"%{http_code}\" ${NGINX_DOMAIN}/${NGINX_PREFIX_PATH}/heartbeat",
                        returnStdout: true
                    ).trim()
                    
                    if (statusCode != "200") {
                        error "New version verification failed with status code: ${statusCode}"
                    }
                    
                    echo "New version verified successfully"
                }
            }
        }

        stage('清理旧版本容器') {
            steps {
                script {
                    // 停止并移除旧版本容器
                    sh "ssh ${DOCKER_SERVER} 'docker stop ${SERVICE_NAME}-${CURRENT_VERSION} || true'"
                    sh "ssh ${DOCKER_SERVER} 'docker rm ${SERVICE_NAME}-${CURRENT_VERSION} || true'"
                    
                    echo "Old version cleaned up"
                }
            }
        }   
        
        stage('清理旧Docker镜像') {
            steps {
                script {
                    // 获取所有镜像tag并按创建时间排序（最新的在前面）
                    def tags = sh(
                        script: """
                            docker images --format '{{.Tag}} {{.CreatedAt}}' ${REGISTRY}/${REPO} | \
                            sort -k2 -r | awk '{print \$1}'
                        """,
                        returnStdout: true
                    ).trim().split('\n')
                    
                    // 计算需要删除的数量
                    int total = tags.size()
                    int toDelete = total - KEEP_COUNT.toInteger()
                    
                    if (toDelete <= 0) {
                        echo "不需要删除 - 只有 ${total} 个tag，保留 ${KEEP_COUNT} 个"
                        return
                    }
                    
                    echo "当前有 ${total} 个tag，将保留最新的 ${KEEP_COUNT} 个，删除 ${toDelete} 个"
                    
                    // 删除旧tag
                    int count = 0
                    for (int i = KEEP_COUNT.toInteger(); i < tags.size(); i++) {
                        def tag = tags[i]
                        if (tag.trim()) {
                            echo "正在删除 ${REGISTRY}/${REPO}:${tag}"
                            try {
                                sh "docker rmi ${REGISTRY}/${REPO}:${tag}"
                                count++
                            } catch (e) {
                                echo "删除 ${REGISTRY}/${REPO}:${tag} 失败: ${e}"
                            }
                        }
                    }
                    
                    echo "删除完成，共删除了 ${count} 个旧tag"
                }
            }
        }

    }

    post {
        success {
            dingtalk (
                robot: '9c6bb3cc-645f-44fd-b592-843757c20792',
                type: 'TEXT',
                text: [
                    "${SERVICE_NAME} 服务部署成功！",
                    "容器名称: ${SERVICE_NAME}-${NEW_VERSION}",
                ],
                at: [
                    'all'
                ]
            )
        }

        failure {
            script {
                // 如果部署失败，回滚到旧版本
                sh """
                    ssh ${NGINX_SERVER} 'ln -sf /etc/nginx/versions/${SERVICE_NAME}/${CURRENT_VERSION}.conf /etc/nginx/conf.d/includes/${SERVICE_NAME}.conf'
                """
                sh "ssh ${NGINX_SERVER} '/usr/local/bin/tengine-ctl -t && /usr/local/bin/tengine-ctl reload'"
                
                // 清理失败的新版本
                sh "ssh ${DOCKER_SERVER} 'docker stop ${SERVICE_NAME}-${NEW_VERSION} || true'"
                sh "ssh ${DOCKER_SERVER} 'docker rm ${SERVICE_NAME}-${NEW_VERSION} || true'"
                
                echo "Deployment failed, rolled back to ${CURRENT_VERSION} version"
            }
        }
        always {
            // 清理构建环境
            sh "docker rmi ${DOCKER_IMAGE}:${BUILD_NUMBER} || true"
        }
    }
}

