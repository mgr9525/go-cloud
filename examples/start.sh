#!/bin/sh

## service name
## 服务所在目录
SERVICE_DIR=/home/user/go/bin/
## 服务名称
SERVICE_NAME=main



cd $SERVICE_DIR
PID_FILE=$SERVICE_NAME\.pid

case "$1" in

    start)
        ##nohup &  以守护进程启动
        #nohup ./$SERVICE_NAME >/dev/null 2>&1 &
        nohup ./$SERVICE_NAME >/dev/null 2>&1 &
        echo $! > $SERVICE_DIR/$PID_FILE
        echo "=== start $SERVICE_NAME"
        ;;

    stop)
        PID=`cat $SERVICE_DIR/$PID_FILE`

        if [ -z "$PID" ];then
            return
        fi

        PIDS=`ps -aux  | awk '{print $2}' | grep "$PID"`
        if [ -z "$PIDS" ]; then
            echo "=== $SERVICE_NAME process not exists or stop success"
        else
            echo "=== stop $SERVICE_NAME:$PID"
            kill $PID
        fi

        rm -rf $SERVICE_DIR/$PID_FILE

        ## 停止5秒
        sleep 2
        ##
        PIDS=`ps -aux  | awk '{print $2}' | grep "$PID"`

        if [ -z "$PIDS" ]; then
            echo "=== $SERVICE_NAME process not exists or stop success"
        else
            echo "=== $SERVICE_NAME process pid is:$PIDS"
            echo "=== begin kill $SERVICE_NAME process, pid is:$PIDS"
            kill -9 $PIDS
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

