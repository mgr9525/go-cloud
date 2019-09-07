#!/bin/sh

## service name
## 服务所在目录
SERVICE_DIR=/home/user/go/bin/
## 服务名称
SERVICE_NAME=main
PID=$SERVICE_NAME\.pid

cd $SERVICE_DIR

case "$1" in

    start)
        ##nohup &  以守护进程启动
        #nohup ./$SERVICE_NAME >/dev/null 2>&1 &
        nohup ./$SERVICE_NAME >/dev/null 2>&1 &
        echo $! > $SERVICE_DIR/$PID
        echo "=== start $SERVICE_NAME"
        ;;

    stop)
        if [ ! -e "$SERVICE_DIR/$PID" ];then
            echo "*** PID file not found"
            return
        fi
        PIDS=`cat $SERVICE_DIR/$PID`
        if [ "$PIDS" = "" ];then
            echo "*** PID file is emtpy"
            return
        fi

        kill $PIDS
        rm -rf $SERVICE_DIR/$PID
        echo "=== stop $SERVICE_NAME:$PIDS"

        ## 停止5秒
        sleep 1
        ##
        PIDS=`ps -aux  | awk '{print $2}' | grep "$PIDS"`
        ## ubuntu dash == upto =
        if [ "$PIDS" = "" ]; then
            echo "=== $SERVICE_NAME process not exists or stop success"
        else
            echo "=== $SERVICE_NAME process pid is:$PIDS"
            echo "=== begin kill $SERVICE_NAME process, pid is:$PIDS"
            kill -9 $PIDS
            sleep 1
        fi
        ;;

    restart)
        $0 stop
        sleep 2
        $0 start
        echo "=== restart $SERVICE_NAME"
        ;;

    *)
        ## restart
        $0 stop
        sleep 2
        $0 start
        ;;

esac
exit 0

