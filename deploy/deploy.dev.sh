#!/bin/bash

dir=""
server=""

if [ $1 ]
then
    dir=($1)
else
    echo "请输入要发布的目录"
    exit 1
fi

if [ $2 ]
then
    server=($2)
else
    echo "请输入要发布的服务名"
    exit 1
fi


path="/home/ecs-user/deploy/$dir"
targetPath="/home/ecs-user/webroot/$dir"
logPath="/home/ecs-user/log"
supervisorConfPath="/home/ecs-user/.local/etc/supervisor/conf.d"

if [ ! -d "$path" ]; then
	mkdir -p "$path"
fi

if [ ! -d "$targetPath" ]; then
	mkdir -p "$targetPath"
fi

if [ ! -d "$targetPath/configs" ]; then
	mkdir -p "$targetPath/configs"
fi

if [ ! -d "$logPath" ]; then
	mkdir -p "$logPath"
fi

#复制目录过去
dirs=("bin" "configs" "deploy")
for i in ${!dirs[@]}; do
	cp -rf "${path}/${dirs[i]}" "${targetPath}"
	if [ $? -ne 0 ]; then
		echo "error to copy agent to target";
		exit 1
	fi
done

# copy supervisor 配置文件
fileName="${path}/deploy/supervisor_dev/${server}.conf"
targetFile="${supervisorConfPath}/${server}.conf"

if [ -f "$fileName" ]; then
    if [ ! -f "${targetFile}" ];then
        cp -f "$fileName" "$targetFile"
        #更新 supervisor 配置
        #系统会自动启动进程
        supervisorctl update "${server}"
    else
        cp -f "$fileName" "$targetFile"
        supervisorctl restart "${server}"
    fi
else
    echo "supervisor config file not found: $fileName"
    exit 1
fi

# 等待进程重启
for k in {1..5}
do
    v=`supervisorctl status "${server}" | grep "RUNNING" | wc -l`
    if [ $v -eq "0" ]; then
        echo "${server} not running yet ${k}";
        sleep 10
    else
        break
    fi
done

# 判断进程状态
for k in {1..5}
do
    sleep 1
    v=`supervisorctl status "${process_name}" | grep "RUNNING" | wc -l`
    if [ $v -eq "0" ]; then
        echo "error status with ${process_name}";
        exit 1
    else
        echo "check status ${k} ok"
    fi
done

echo "ok"
exit 0
