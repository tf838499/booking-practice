MASTER1_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis-master1)
MASTER2_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis-master2)
MASTER3_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis-master3)
SLAVE1_1_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis-slave11)
SLAVE1_2_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis-slave12)
SLAVE2_1_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis-slave21)
SLAVE2_2_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis-slave22)
SLAVE3_1_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis-slave31)
SLAVE3_2_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis-slave32)
echo "Master1 IP: $MASTER1_IP"
echo "Master2 IP: $MASTER2_IP"
echo "Master3 IP: $MASTER3_IP"
echo "Slave1-1 IP: $SLAVE1_1_IP"
echo "Slave1-2 IP: $SLAVE1_2_IP"
echo "Slave2-1 IP: $SLAVE2_1_IP"
echo "Slave2-2 IP: $SLAVE2_2_IP"
echo "Slave3-1 IP: $SLAVE3_1_IP"
echo "Slave3-2 IP: $SLAVE3_2_IP"

docker exec -it redis-master1 redis-cli --cluster create \
  $MASTER1_IP:6379 \
  $MASTER2_IP:6379 \
  $MASTER3_IP:6379 \
  $SLAVE1_1_IP:6379 \
  $SLAVE1_2_IP:6379 \
  $SLAVE2_1_IP:6379 \
  $SLAVE2_2_IP:6379 \
  $SLAVE3_1_IP:6379 \
  $SLAVE3_2_IP:6379 \
  --cluster-replicas 2