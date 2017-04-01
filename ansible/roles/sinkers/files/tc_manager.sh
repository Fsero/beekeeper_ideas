#!/bin/bash -x

clean_rules() {
    local iface=$1
    tc qdisc del dev $iface root    2> /dev/null > /dev/null
    tc qdisc del dev $iface ingress 2> /dev/null > /dev/null
}

start_traffic_shape(){
    local iface="$1"
    local maxbwidth_download="${2:-56}"
    local maxbwidth_upload="${3:-56}"
    local latency_added="${4:-50ms}"
    local burst_allowed="${5:-1540}"

    # clean existing down- and uplink qdiscs, hide errors
    tc qdisc del dev $iface root    2> /dev/null > /dev/null
    tc qdisc del dev $iface ingress 2> /dev/null > /dev/null

    ###### uplink

    # install root HTB, point default traffic to 1:20:

    tc qdisc add dev $iface root handle 1: htb default 20

    # shape everything at $UPLINK speed - this prevents huge queues in your
    # DSL modem which destroy latency:

    tc class add dev $iface parent 1: classid 1:1 htb rate ${maxbwidth_upload}kbit burst 6k

    # high prio class 1:10:

    tc class add dev $iface parent 1:1 classid 1:10 htb rate ${maxbwidth_upload}kbit \
    burst 6k prio 1

    # bulk & default class 1:20 - gets slightly less traffic,
    # and a lower priority:

    tc class add dev $iface parent 1:1 classid 1:20 htb rate $[9*$maxbwidth_upload/10]kbit \
    burst 6k prio 2

    # both get Stochastic Fairness:
    tc qdisc add dev $iface parent 1:10 handle 10: sfq perturb 10
    tc qdisc add dev $iface parent 1:20 handle 20: sfq perturb 10

    # TOS Minimum Delay (ssh, NOT scp) in 1:10:
    tc filter add dev $iface parent 1:0 protocol ip prio 10 u32 \
    match ip tos 0x10 0xff  flowid 1:10

    # ICMP (ip protocol 1) in the interactive class 1:10 so we
    # can do measurements & impress our friends:
    tc filter add dev $iface parent 1:0 protocol ip prio 10 u32 \
    match ip protocol 1 0xff flowid 1:10

    # To speed up downloads while an upload is going on, put ACK packets in
    # the interactive class:

    tc filter add dev $iface parent 1: protocol ip prio 10 u32 \
    match ip protocol 6 0xff \
    match u8 0x05 0x0f at 0 \
    match u16 0x0000 0xffc0 at 2 \
    match u8 0x10 0xff at 33 \
    flowid 1:10

    # rest is 'non-interactive' ie 'bulk' and ends up in 1:20


    ########## downlink #############
    # slow downloads down to somewhat less than the real speed  to prevent
    # queuing at our ISP. Tune to see how high you can set it.
    # ISPs tend to have *huge* queues to make sure big downloads are fast
    #
    # attach ingress policer:

    tc qdisc add dev $iface handle ffff: ingress

    # filter *everything* to it (0.0.0.0/0), drop everything that's
    # coming in too fast:

    tc filter add dev $iface parent ffff: protocol ip prio 50 u32 match ip src \
    0.0.0.0/0 police rate ${maxbwidth_download}kbit burst 10k drop flowid :1

}

get_bridge_docker_device(){
    local _outvar=$1
    local docker_network_name=$2
    local docker_network_id=$(docker network ls | grep $docker_network_name | cut -f1 -d' ')
    local result=""
    if [[ -z "$docker_network_id" ]]; then
        echo "unable to get bridge id from docker, refusing to continue, maybe it not exists?"
        result="NO_DEV"
        eval $_outvar=\$result
	return
    fi
    IFACE="br-$docker_network_id"
    ifconfig $IFACE &> /dev/null
    if [[ $? -ne 0 ]]; then
        echo "unable to contact bridged interface <$IFACE>, refusing to continue"
        result="NO_DEV"
        eval $_outvar=\$result
	return
    fi
    result="br-$docker_network_id"
    eval $_outvar=\$result
}

do_start() {
    local iface=""
    get_bridge_docker_device "iface" "$docker_network_name"
    if [[ "$iface" == "NO_DEV" ]]; then
      echo "unable to get interface, refusing to continue"
      exit 1
    fi
    clean_rules $iface
    start_traffic_shape $iface
}

do_stop() {
    local iface=""
    get_bridge_docker_device "iface" "$docker_network_name"
    if [[ "$iface" == "NO_DEV" ]]; then
      echo "unable to get interface, maybe it not exists"
      exit 0
    fi
    clean_rules $iface
}

do_status() {
    local number_of_rules=$(tc qdisc show | wc -l)
    if [[ $number_of_rules -gt 2 ]]; then
        echo "$0 is ENABLED"
    else
        echo "$0 is DISABLED"
    fi
}


case "$1" in
   start)
     do_start $2
     ;;
   stop)
     do_stop $2
     ;;
   restart)
     do_stop $2
     do_start $2
     ;;
   status)
     do_status
     ;;
   *)
     echo "Usage: /etc/init.d/$NAME start DOCKER_NETWORK_NAME | stop DOCKER_NETWORK_NAME | status"
     exit 1
     ;;
esac


exit 0
